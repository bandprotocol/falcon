package evm

import (
	"context"
	"fmt"
	"math/big"
	"path"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	keyStore "github.com/ethereum/go-ethereum/accounts/keystore"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

var _ chains.ChainProvider = (*EVMChainProvider)(nil)

// EVMChainProvider is the struct that handles interactions with the EVM chain.
type EVMChainProvider struct {
	Config    *EVMChainProviderConfig
	ChainName string

	Client  Client
	GasType GasType

	KeyInfo     KeyInfo
	FreeSenders chan *Sender

	TunnelRouterAddress gethcommon.Address
	TunnelRouterABI     abi.ABI

	Log *zap.Logger

	KeyStore *keyStore.KeyStore
}

// NewEVMChainProvider creates a new EVM chain provider.
func NewEVMChainProvider(
	chainName string,
	client Client,
	cfg *EVMChainProviderConfig,
	log *zap.Logger,
	homePath string,
) (*EVMChainProvider, error) {
	// load abis here
	abi, err := abi.JSON(strings.NewReader(gasPriceTunnelRouterABI))
	if err != nil {
		log.Error("ChainProvider: failed to load abi",
			zap.Error(err),
			zap.String("chain_name", chainName),
		)
		return nil, ErrLoadAbi(err)
	}

	addr, err := HexToAddress(cfg.TunnelRouterAddress)
	if err != nil {
		log.Error("ChainProvider: cannot convert tunnel router address",
			zap.Error(err),
			zap.String("chain_name", chainName),
		)
		return nil, ErrInvalidAddress(err)
	}

	keyStoreDir := path.Join(homePath, keyDir, chainName, privateKeyDir)
	keyStore := keyStore.NewKeyStore(keyStoreDir, keyStore.StandardScryptN, keyStore.StandardScryptP)

	keyInfo, err := LoadKeyInfo(homePath, chainName)
	if err != nil {
		return nil, err
	}

	return &EVMChainProvider{
		Config:              cfg,
		ChainName:           chainName,
		Client:              client,
		GasType:             cfg.GasType,
		KeyInfo:             keyInfo,
		TunnelRouterAddress: addr,
		TunnelRouterABI:     abi,
		Log:                 log.With(zap.String("chain_name", chainName)),
		KeyStore:            keyStore,
	}, nil
}

// Connect connects to the EVM chain.
func (cp *EVMChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(ctx); err != nil {
		return err
	}

	go cp.Client.StartLivelinessCheck(ctx, cp.Config.LivelinessCheckingInterval)

	return nil
}

// QueryTunnelInfo queries the tunnel info from the tunnel router contract.
func (cp *EVMChainProvider) QueryTunnelInfo(
	ctx context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*chainstypes.Tunnel, error) {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		cp.Log.Error("Connect client error", zap.Error(err))
		return nil, ErrClientConnection(err)
	}

	addr, err := HexToAddress(tunnelDestinationAddr)
	if err != nil {
		cp.Log.Error("Invalid address", zap.Error(err), zap.String("address", tunnelDestinationAddr))
		return nil, ErrInvalidAddress(err)
	}

	info, err := cp.queryTunnelInfo(ctx, tunnelID, addr)
	if err != nil {
		cp.Log.Error(
			"Query contract error",
			zap.Error(err),
			zap.Uint64("tunnel_id", tunnelID),
			zap.String("address", tunnelDestinationAddr),
		)

		return nil, ErrQueryData(err)
	}

	return &chainstypes.Tunnel{
		ID:             tunnelID,
		TargetAddress:  tunnelDestinationAddr,
		IsActive:       info.IsActive,
		LatestSequence: info.LatestSequence,
		Balance:        info.Balance,
	}, nil
}

// RelayPacket relays the packet from the source chain to the destination chain.
func (cp *EVMChainProvider) RelayPacket(ctx context.Context, packet *bandtypes.Packet) error {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		cp.Log.Error("Connect client error", zap.Error(err))
		return ErrClientConnection(err)
	}

	// get a free sender
	cp.Log.Debug("Waiting for a free sender...")
	sender := <-cp.FreeSenders
	defer func() { cp.FreeSenders <- sender }()

	log := cp.Log.With(
		zap.Uint64("tunnel_id", packet.TunnelID),
		zap.Uint64("sequence", packet.Sequence),
		zap.String("sender_address", sender.Address.String()),
	)

	// get gas information
	gasInfo, err := cp.EstimateGasFee(ctx)
	if err != nil {
		cp.Log.Error("Failed to estimate gas fee", zap.Error(err))
		return fmt.Errorf("failed to estimate gas fee: %w", err)
	}

	retryCount := 1
	for retryCount <= cp.Config.MaxRetry {
		log.Info("Relaying a message", zap.Int("retry_count", retryCount))

		// create and submit a transaction; if failed, retry, no need to bump gas.
		signedTx, err := cp.createAndSignRelayTx(ctx, packet, sender, gasInfo)
		if err != nil {
			log.Error("CreateAndSignTx error", zap.Error(err), zap.Int("retry_count", retryCount))
			retryCount += 1
			continue
		}

		// submit the transaction, if failed, bump gas and retry
		txHash, err := cp.Client.BroadcastTx(ctx, signedTx)
		if err != nil {
			log.Error("HandleRelay error", zap.Error(err), zap.Int("retry_count", retryCount))
			// bump gas and retry
			gasInfo, err = cp.BumpAndBoundGas(ctx, gasInfo, cp.Config.GasMultiplier)
			if err != nil {
				log.Error("Cannot bump gas", zap.Error(err), zap.Int("retry_count", retryCount))
			}

			retryCount += 1
			continue
		}

		createdAt := time.Now()

		log.Info(
			"Submitted a message; checking transaction status",
			zap.String("tx_hash", txHash),
			zap.Int("retry_count", retryCount),
		)

		var checkTxErr error
		var txStatus TxStatus
	checkTxLogic:
		for time.Since(createdAt) < cp.Config.WaitingTxDuration {
			result, err := cp.checkConfirmedTx(ctx, txHash)
			if err != nil {
				log.Debug(
					"Failed to check tx status",
					zap.Error(err),
					zap.String("tx_hash", txHash),
					zap.Int("retry_count", retryCount),
				)

				checkTxErr = err
				txStatus = TX_STATUS_UNDEFINED
				time.Sleep(cp.Config.CheckingTxInterval)
				continue
			}

			checkTxErr = nil
			txStatus = result.Status
			switch result.Status {
			case TX_STATUS_SUCCESS:
				log.Info(
					"Packet is successfully relayed",
					zap.String("tx_hash", txHash),
					zap.Int("retry_count", retryCount),
				)
				return nil
			case TX_STATUS_FAILED:
				retryCount += 1
				break checkTxLogic
			case TX_STATUS_UNMINED:
				log.Debug(
					"Waiting for tx to be mined",
					zap.Error(err),
					zap.String("tx_hash", txHash),
					zap.Int("retry_count", retryCount),
				)

				time.Sleep(cp.Config.CheckingTxInterval)
			}
		}

		log.Error(
			"Failed to relaying a packet with status and error",
			zap.Error(checkTxErr),
			zap.String("status", txStatus.String()),
			zap.String("tx_hash", txHash),
			zap.Int("retry_count", retryCount),
		)

		// bump gas and retry
		gasInfo, err = cp.BumpAndBoundGas(ctx, gasInfo, cp.Config.GasMultiplier)
		if err != nil {
			log.Error("Cannot bump gas", zap.Error(err), zap.Int("retry_count", retryCount))
		}

		retryCount += 1
	}

	return ErrRelayPacketRetries(cp.Config.MaxRetry)
}

// createAndSignRelayTx creates and signs the relay transaction.
func (cp *EVMChainProvider) createAndSignRelayTx(
	ctx context.Context,
	packet *bandtypes.Packet,
	sender *Sender,
	gasInfo GasInfo,
) (*gethtypes.Transaction, error) {
	calldata, err := cp.createCalldata(packet)
	if err != nil {
		return nil, fmt.Errorf("failed to create calldata: %w", err)
	}

	tx, err := cp.newRelayTx(ctx, calldata, sender.Address, gasInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create an evm transaction: %w", err)
	}

	signedTx, err := cp.signTx(tx, sender)
	if err != nil {
		return nil, fmt.Errorf("failed to sign an evm transaction: %w", err)
	}

	return signedTx, nil
}

// checkConfirmedTx checks the confirmed transaction status.
func (cp *EVMChainProvider) checkConfirmedTx(
	ctx context.Context,
	txHash string,
) (*ConfirmTxResult, error) {
	failResult := NewConfirmTxResult(
		txHash,
		TX_STATUS_UNMINED,
		decimal.NullDecimal{},
		cp.GasType,
	)

	receipt, err := cp.Client.GetTxReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx receipt: %w", err)
	}

	if receipt.Status == gethtypes.ReceiptStatusFailed {
		return failResult.WithStatus(TX_STATUS_FAILED), nil
	}

	latestBlock, err := cp.Client.GetBlockHeight(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block height: %w", err)
	}

	// if tx block is not confirmed and waiting too long return status with timeout
	if receipt.BlockNumber.Uint64() > latestBlock-cp.Config.BlockConfirmation {
		return failResult.WithStatus(TX_STATUS_UNMINED), nil
	}

	// calculate gas used and effective gas price
	gasUsed := decimal.NewNullDecimal(decimal.New(int64(receipt.GasUsed), 0))
	return NewConfirmTxResult(txHash, TX_STATUS_SUCCESS, gasUsed, cp.GasType), nil
}

// EstimateGasFee estimates the gas for the transaction.
func (cp *EVMChainProvider) EstimateGasFee(ctx context.Context) (GasInfo, error) {
	switch cp.GasType {
	case GasTypeLegacy:
		gasPrice, err := cp.Client.EstimateGasPrice(ctx)
		if err != nil {
			return GasInfo{}, err
		}

		// bound gas fee
		return cp.BumpAndBoundGas(ctx, NewGasLegacyInfo(gasPrice), 1.0)
	case GasTypeEIP1559:
		priorityFee, err := cp.Client.EstimateGasTipCap(ctx)
		if err != nil {
			return GasInfo{}, err
		}

		baseFee, err := cp.Client.EstimateBaseFee(ctx)
		if err != nil {
			return GasInfo{}, err
		}

		// bound gas fee
		return cp.BumpAndBoundGas(ctx, NewGasEIP1559Info(priorityFee, baseFee), 1.0)
	default:
		return GasInfo{}, ErrUnsupportedGasType(cp.GasType)
	}
}

// BumpAndBoundGas bumps the gas price.
func (cp *EVMChainProvider) BumpAndBoundGas(
	ctx context.Context,
	gasInfo GasInfo,
	multiplier float64,
) (newGasInfo GasInfo, err error) {
	switch gasInfo.Type {
	case GasTypeLegacy:
		// calculate new gas price and compare with the cap being setup in the configuration.
		// if the cap is not set in the configuration, should query from the contract.
		newGasPrice := MultiplyBigIntWithFloat64(gasInfo.GasPrice, multiplier)

		maxGasPrice := big.NewInt(int64(cp.Config.MaxGasPrice))
		if maxGasPrice.Cmp(big.NewInt(0)) <= 0 {
			maxGasPrice, err = cp.queryRelayerGasFee(ctx)
			if err != nil {
				return GasInfo{}, err
			}
		}

		if newGasPrice.Cmp(maxGasPrice) > 0 {
			newGasPrice = maxGasPrice
		}

		return NewGasLegacyInfo(newGasPrice), nil
	case GasTypeEIP1559:
		// calculate new priority fee and compare with the cap being setup in the configuration.
		// if the cap is not set in the configuration, should query from the contract.
		newPriorityFee := MultiplyBigIntWithFloat64(gasInfo.GasPriorityFee, multiplier)

		maxPriorityFee := big.NewInt(int64(cp.Config.MaxPriorityFee))
		if maxPriorityFee.Cmp(big.NewInt(0)) <= 0 {
			maxPriorityFee, err = cp.queryRelayerGasFee(ctx)
			if err != nil {
				return GasInfo{}, err
			}
		}

		if newPriorityFee.Cmp(maxPriorityFee) > 0 {
			newPriorityFee = maxPriorityFee
		}

		maxBaseFee := big.NewInt(int64(cp.Config.MaxBaseFee))
		newBaseFee := gasInfo.GasBaseFee
		if maxBaseFee.Cmp(big.NewInt(0)) > 0 && newBaseFee.Cmp(maxBaseFee) > 0 {
			newBaseFee = maxBaseFee
		}

		return NewGasEIP1559Info(newPriorityFee, newBaseFee), nil
	default:
		return GasInfo{}, ErrUnsupportedGasType(cp.GasType)
	}
}

// queryTunnelInfo queries the target contract information.
func (cp *EVMChainProvider) queryTunnelInfo(
	ctx context.Context,
	tunnelID uint64,
	addr gethcommon.Address,
) (*TunnelInfoOutput, error) {
	calldata, err := cp.TunnelRouterABI.Pack("tunnelInfo", tunnelID, addr)
	if err != nil {
		return nil, ErrPackCalldata(err)
	}

	b, err := cp.Client.Query(ctx, cp.TunnelRouterAddress, calldata)
	if err != nil {
		return nil, ErrQueryData(err)
	}

	var output TunnelInfoOutputRaw
	if err := cp.TunnelRouterABI.UnpackIntoInterface(&output, "tunnelInfo", b); err != nil {
		return nil, ErrUnpackData(err)
	}

	return &output.Info, nil
}

// newRelayTx creates a new relay transaction.
func (cp *EVMChainProvider) newRelayTx(
	ctx context.Context,
	data []byte,
	sender gethcommon.Address,
	gasInfo GasInfo,
) (*gethtypes.Transaction, error) {
	nonce, err := cp.Client.PendingNonceAt(ctx, sender)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:      sender,
		To:        &cp.TunnelRouterAddress,
		Data:      data,
		GasPrice:  gasInfo.GasPrice,
		GasFeeCap: gasInfo.GasFeeCap,
		GasTipCap: gasInfo.GasPriorityFee,
	}

	// calculate gas limit
	gasLimit := cp.Config.GasLimit
	if gasLimit == 0 {
		gasLimit, err = cp.Client.EstimateGas(ctx, callMsg)
		if err != nil {
			return nil, err
		}
	}

	// set fee info
	var tx *gethtypes.Transaction
	switch cp.GasType {
	case GasTypeLegacy:
		tx = gethtypes.NewTx(&gethtypes.LegacyTx{
			Nonce:    nonce,
			To:       &cp.TunnelRouterAddress,
			Value:    decimal.NewFromInt(0).BigInt(),
			Data:     data,
			Gas:      gasLimit,
			GasPrice: gasInfo.GasPrice,
		})

	case GasTypeEIP1559:
		tx = gethtypes.NewTx(&gethtypes.DynamicFeeTx{
			ChainID:   big.NewInt(int64(cp.Config.ChainID)),
			Nonce:     nonce,
			To:        &cp.TunnelRouterAddress,
			Value:     decimal.NewFromInt(0).BigInt(),
			Data:      data,
			Gas:       gasLimit,
			GasFeeCap: gasInfo.GasFeeCap,
			GasTipCap: gasInfo.GasPriorityFee,
		})

	default:
		return nil, ErrUnsupportedGasType(cp.GasType)
	}

	return tx, nil
}

// createCalldata creates the calldata for the relay transaction.
func (cp *EVMChainProvider) createCalldata(packet *bandtypes.Packet) ([]byte, error) {
	var signing *bandtypes.Signing

	// get signing from packet; prefer to use signing from
	// current group than incoming group
	if packet.CurrentGroupSigning != nil {
		signing = packet.CurrentGroupSigning
	} else if packet.IncomingGroupSigning != nil {
		signing = packet.IncomingGroupSigning
	} else {
		return nil, fmt.Errorf("missing signing")
	}

	rAddr, err := HexToAddress(signing.EVMSignature.RAddress.String())
	if err != nil {
		return nil, err
	}

	return cp.TunnelRouterABI.Pack(
		"relay",
		signing.Message.Bytes(),
		rAddr,
		new(big.Int).SetBytes(signing.EVMSignature.Signature),
	)
}

// signTx signs the transaction with the sender.
func (cp *EVMChainProvider) signTx(
	tx *gethtypes.Transaction,
	sender *Sender,
) (*gethtypes.Transaction, error) {
	var signer gethtypes.Signer
	switch cp.GasType {
	case GasTypeLegacy:
		signer = gethtypes.NewEIP155Signer(big.NewInt(int64(cp.Config.ChainID)))
	case GasTypeEIP1559:
		signer = gethtypes.NewLondonSigner(big.NewInt(int64(cp.Config.ChainID)))
	default:
		return nil, ErrUnsupportedGasType(cp.GasType)
	}

	return gethtypes.SignTx(tx, signer, sender.PrivateKey)
}

// QueryBalance queries balance of specific account address.
func (cp *EVMChainProvider) QueryBalance(
	ctx context.Context,
	keyName string,
) (*big.Int, error) {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		cp.Log.Error(
			"Connect client error",
			zap.Error(err),
			zap.String("chain_name", cp.ChainName),
		)
		return nil, ErrClientConnection(err)
	}

	address, err := HexToAddress(cp.KeyInfo[keyName])
	if err != nil {
		return nil, err
	}

	return cp.Client.GetBalance(ctx, address)
}

// queryRelayerGasFee queries the relayer gas fee being set on tunnel router.
func (cp *EVMChainProvider) queryRelayerGasFee(ctx context.Context) (*big.Int, error) {
	calldata, err := cp.TunnelRouterABI.Pack("gasFee")
	if err != nil {
		return nil, ErrPackCalldata(err)
	}

	b, err := cp.Client.Query(ctx, cp.TunnelRouterAddress, calldata)
	if err != nil {
		return nil, ErrQueryData(err)
	}

	var output *big.Int
	if err := cp.TunnelRouterABI.UnpackIntoInterface(&output, "gasFee", b); err != nil {
		return nil, ErrUnpackData(err)
	}

	return output, nil
}
