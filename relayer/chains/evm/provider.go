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
		return nil, fmt.Errorf("[EVMProvider] failed to load abi: %w", err)
	}

	addr, err := HexToAddress(cfg.TunnelRouterAddress)
	if err != nil {
		log.Error("ChainProvider: cannot convert tunnel router address",
			zap.Error(err),
			zap.String("chain_name", chainName),
		)
		return nil, fmt.Errorf("[EVMProvider] incorrect address: %w", err)
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
		Log:                 log,
		KeyStore:            keyStore,
	}, nil
}

// Connect connects to the EVM chain.
func (cp *EVMChainProvider) Init(ctx context.Context) error {
	// TODO: implement loading private key from store

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
		cp.Log.Error(
			"connect client error",
			zap.Error(err),
			zap.String("chain_name", cp.ChainName),
		)
		return nil, fmt.Errorf("[EVMProvider] failed to connect client: %w", err)
	}

	addr, err := HexToAddress(tunnelDestinationAddr)
	if err != nil {
		cp.Log.Error(
			"invalid address",
			zap.Error(err),
			zap.String("chain_name", cp.ChainName),
			zap.String("address", tunnelDestinationAddr),
		)
		return nil, fmt.Errorf("[EVMProvider] invalid address: %w", err)
	}

	info, err := cp.queryTunnelInfo(ctx, tunnelID, addr)
	if err != nil {
		cp.Log.Error(
			"query contract error",
			zap.Error(err),
			zap.String("chain_name", cp.ChainName),
			zap.Uint64("tunnel_id", tunnelID),
			zap.String("address", tunnelDestinationAddr),
		)

		return nil, fmt.Errorf("[EVMProvider] failed to query contract: %w", err)
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
func (cp *EVMChainProvider) RelayPacket(
	ctx context.Context,
	packet *bandtypes.Packet,
) error {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		cp.Log.Error(
			"connect client error",
			zap.Error(err),
			zap.String("chain_name", cp.ChainName),
		)
		return fmt.Errorf("[EVMProvider] failed to connect client: %w", err)
	}

	gasInfo, err := cp.EstimateGas(ctx)
	if err != nil {
		cp.Log.Error(
			"failed to estimate gas",
			zap.Error(err),
			zap.String("chain_name", cp.ChainName),
		)
		return fmt.Errorf("failed to estimate gas: %w", err)
	}

	retryCount := 0
	for retryCount < cp.Config.MaxRetry {
		cp.Log.Info(
			"relaying a message",
			zap.String("chain_name", cp.ChainName),
			zap.Uint64("tunnel_id", packet.TunnelID),
			zap.Uint64("sequence", packet.Sequence),
			zap.Int("retry_count", retryCount),
		)

		// create and submit a transaction; if failed, retry, no need to bump gas.
		txHash, err := cp.handleRelay(ctx, packet, gasInfo)
		if err != nil {
			cp.Log.Error(
				"HandleRelay error",
				zap.Error(err),
				zap.String("chain_name", cp.ChainName),
				zap.Uint64("tunnel_id", packet.TunnelID),
				zap.Uint64("sequence", packet.Sequence),
				zap.Int("retry_count", retryCount),
			)
			retryCount += 1
			continue
		}
		createdAt := time.Now()

		cp.Log.Info(
			"submitted a message; checking transaction status",
			zap.String("chain_name", cp.ChainName),
			zap.Uint64("tunnel_id", packet.TunnelID),
			zap.Uint64("sequence", packet.Sequence),
			zap.String("tx_hash", txHash),
			zap.Int("retry_count", retryCount),
		)

		var checkTxErr error
		var txStatus TxStatus
	checkTxLogic:
		for time.Since(createdAt) < cp.Config.WaitingTxDuration {
			result, err := cp.checkConfirmedTx(ctx, txHash)
			if err != nil {
				cp.Log.Debug(
					"Failed to check tx status",
					zap.Error(err),
					zap.String("chain_name", cp.ChainName),
					zap.Uint64("tunnel_id", packet.TunnelID),
					zap.Uint64("sequence", packet.Sequence),
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
				cp.Log.Info(
					"Packet is successfully relayed",
					zap.String("chain_name", cp.ChainName),
					zap.Uint64("tunnel_id", packet.TunnelID),
					zap.Uint64("sequence", packet.Sequence),
					zap.String("tx_hash", txHash),
					zap.Int("retry_count", retryCount),
				)
				return nil
			case TX_STATUS_FAILED:
				retryCount += 1
				break checkTxLogic
			case TX_STATUS_UNMINED:
				cp.Log.Debug(
					"Waiting for tx to be mined",
					zap.Error(err),
					zap.String("chain_name", cp.ChainName),
					zap.Uint64("tunnel_id", packet.TunnelID),
					zap.Uint64("sequence", packet.Sequence),
					zap.String("tx_hash", txHash),
					zap.Int("retry_count", retryCount),
				)

				time.Sleep(cp.Config.CheckingTxInterval)
			}
		}

		cp.Log.Error(
			"Failed to relaying a packet with status and error",
			zap.Error(checkTxErr),
			zap.String("status", txStatus.String()),
			zap.String("chain_name", cp.ChainName),
			zap.Uint64("tunnel_id", packet.TunnelID),
			zap.Uint64("sequence", packet.Sequence),
			zap.String("tx_hash", txHash),
			zap.Int("retry_count", retryCount),
		)

		// bump gas and retry
		gasInfo, err = cp.BumpAndBoundGas(ctx, gasInfo, cp.Config.GasMultiplier)
		if err != nil {
			cp.Log.Error(
				"cannot bump gas",
				zap.Error(err),
				zap.String("chain_name", cp.ChainName),
				zap.Uint64("tunnel_id", packet.TunnelID),
				zap.Uint64("sequence", packet.Sequence),
				zap.Int("retry_count", retryCount),
			)
		}

		retryCount += 1
	}

	return fmt.Errorf("[EVMProvider] failed to relay packet after %d retries", cp.Config.MaxRetry)
}

// handleRelay handles the relay message from the source chain to the destination chain.
func (cp *EVMChainProvider) handleRelay(
	ctx context.Context,
	packet *bandtypes.Packet,
	gasInfo GasInfo,
) (txHash string, err error) {
	calldata, err := cp.createCalldata(packet)
	if err != nil {
		return "", fmt.Errorf("failed to create calldata: %w", err)
	}

	if len(cp.FreeSenders) == 0 {
		return "", fmt.Errorf("no key available to relay packet")
	}

	sender := <-cp.FreeSenders
	defer func() { cp.FreeSenders <- sender }()

	cp.Log.Debug(
		fmt.Sprintf("Relaying packet using address: %v", sender.Address),
		zap.String("evm_sender_address", sender.Address.String()),
		zap.String("chain_name", cp.ChainName),
		zap.Uint64("tunnel_id", packet.TunnelID),
		zap.Uint64("sequence", packet.Sequence),
	)

	tx, err := cp.newRelayTx(ctx, calldata, sender.Address, gasInfo)
	if err != nil {
		return "", fmt.Errorf("failed to create an evm transaction: %w", err)
	}

	signedTx, err := cp.signTx(tx, sender)
	if err != nil {
		return "", fmt.Errorf("failed to sign an evm transaction: %w", err)
	}

	txHash, err = cp.Client.BroadcastTx(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast an evm transaction: %w", err)
	}

	return txHash, nil
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
		decimal.NullDecimal{},
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

	effectiveGas, err := cp.GetEffectiveGas(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to get effective gas price: %w", err)
	}

	return NewConfirmTxResult(txHash, TX_STATUS_SUCCESS, gasUsed, cp.GasType, effectiveGas), nil
}

func (cp *EVMChainProvider) GetEffectiveGas(
	ctx context.Context,
	receipt *gethtypes.Receipt,
) (decimal.NullDecimal, error) {
	switch cp.GasType {
	case GasTypeLegacy:
		return cp.Client.GetEffectiveGasPrice(ctx, receipt)
	case GasTypeEIP1559:
		return cp.Client.GetEffectiveGasTipValue(ctx, receipt)
	default:
		return decimal.NullDecimal{}, fmt.Errorf("unsupported gas type: %v", cp.GasType)
	}
}

// EstimateGas estimates the gas for the transaction.
func (cp *EVMChainProvider) EstimateGas(ctx context.Context) (GasInfo, error) {
	switch cp.GasType {
	case GasTypeLegacy:
		gasPrice, err := cp.Client.EstimateGasPrice(ctx)
		if err != nil {
			return GasInfo{}, err
		}

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

		return cp.BumpAndBoundGas(ctx, NewGasEIP1559Info(priorityFee, baseFee), 1.0)
	default:
		return GasInfo{}, fmt.Errorf("unsupported gas type: %v", cp.GasType)
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

		newBaseFee := gasInfo.GasBaseFee
		if newBaseFee.Cmp(big.NewInt(int64(cp.Config.MaxBaseFee))) > 0 {
			newBaseFee = big.NewInt(int64(cp.Config.MaxBaseFee))
		}

		return NewGasEIP1559Info(newPriorityFee, newBaseFee), nil
	default:
		return GasInfo{}, fmt.Errorf("unsupported gas type: %v", cp.GasType)
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
		return nil, fmt.Errorf("failed to pack calldata: %w", err)
	}

	b, err := cp.Client.Query(ctx, cp.TunnelRouterAddress, calldata)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}

	var output TunnelInfoOutputRaw
	if err := cp.TunnelRouterABI.UnpackIntoInterface(&output, "tunnelInfo", b); err != nil {
		return nil, fmt.Errorf("failed to unpack data: %w", err)
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
	nonce, err := cp.Client.GetNonce(ctx, sender)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From: sender,
		To:   &cp.TunnelRouterAddress,
		Data: data,
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
		gasFeeCap := new(big.Int).Add(gasInfo.GasBaseFee, gasInfo.GasPriorityFee)

		tx = gethtypes.NewTx(&gethtypes.DynamicFeeTx{
			ChainID:   big.NewInt(int64(cp.Config.ChainID)),
			Nonce:     nonce,
			To:        &cp.TunnelRouterAddress,
			Value:     decimal.NewFromInt(0).BigInt(),
			Data:      data,
			Gas:       gasLimit,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasInfo.GasPriorityFee,
		})

	default:
		return nil, fmt.Errorf("unsupported gas type: %v", cp.GasType)
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

	if signing.EVMSignature == nil {
		return nil, fmt.Errorf("evm signature is unavailable")
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
		return nil, fmt.Errorf("unsupported gas type: %v", cp.GasType)
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
			"connect client error",
			zap.Error(err),
			zap.String("chain_name", cp.ChainName),
		)
		return nil, fmt.Errorf("[EVMProvider] failed to connect client: %w", err)
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
		return nil, fmt.Errorf("failed to pack calldata: %w", err)
	}

	b, err := cp.Client.Query(ctx, cp.TunnelRouterAddress, calldata)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}

	var output *big.Int
	if err := cp.TunnelRouterABI.UnpackIntoInterface(&output, "gasFee", b); err != nil {
		return nil, fmt.Errorf("failed to unpack data: %w", err)
	}

	return output, nil
}
