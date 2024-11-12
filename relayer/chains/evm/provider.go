package evm

import (
	"context"
	"fmt"
	"math"
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
	"github.com/bandprotocol/falcon/relayer/chains/evm/gas"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

var _ chains.ChainProvider = (*EVMChainProvider)(nil)

// EVMChainProvider is the struct that handles interactions with the EVM chain.
type EVMChainProvider struct {
	Config    *EVMChainProviderConfig
	ChainName string

	// TODO: add lock object.
	Client   Client
	GasModel gas.GasModel

	FreeSenders FreeSenders

	TunnelRouterAddress gethcommon.Address
	TunnelRouterABI     abi.ABI

	Log *zap.Logger

	KeyStore *keyStore.KeyStore
}

// NewEVMChainProvider creates a new EVM chain provider.
func NewEVMChainProvider(
	chainName string,
	client Client,
	gasModel gas.GasModel,
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

	// create free senders
	freeSenders, err := LoadFreeSenders(homePath, chainName, keyStore)
	if err != nil {
		log.Error("ChainProvider: cannot create a sender",
			zap.Error(err),
			zap.String("chain_name", chainName),
		)
		return nil, fmt.Errorf("[EVMProvider] failed to create a sender: %w", err)
	}

	return &EVMChainProvider{
		Config:              cfg,
		ChainName:           chainName,
		Client:              client,
		GasModel:            gasModel,
		FreeSenders:         freeSenders,
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

	retryCount := 0
	for retryCount < cp.Config.MaxRetry {
		cp.Log.Info(
			"relaying a message",
			zap.String("chain_name", cp.ChainName),
			zap.Uint64("tunnel_id", packet.TunnelID),
			zap.Uint64("sequence", packet.Sequence),
			zap.Int("retry_count", retryCount),
		)

		txHash, err := cp.handleRelay(ctx, packet, retryCount)
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
			"submitting a message",
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
			result, err := cp.checkConfirmedTx(ctx, txHash, createdAt)
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
			case TX_STATUS_FAILED, TX_STATUS_TIMEOUT:
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

		retryCount += 1
	}

	return fmt.Errorf("[EVMProvider] failed to relay packet after %d retries", cp.Config.MaxRetry)
}

// handleRelay handles the relay message from the source chain to the destination chain.
func (cp *EVMChainProvider) handleRelay(
	ctx context.Context,
	packet *bandtypes.Packet,
	retryCount int,
) (txHash string, err error) {
	calldata, err := cp.createCalldata(packet)
	if err != nil {
		return "", fmt.Errorf("failed to create calldata: %w", err)
	}

	var selectedSender *Sender
	var selectedKeyName string

	if len(cp.FreeSenders) == 0 {
		return "", fmt.Errorf("no key available to relay packet")
	}

	// use available sender
	for selectedSender == nil {
		for keyName, sender := range cp.FreeSenders {
			if !sender.IsExecuting {
				cp.FreeSenders[selectedKeyName].Mutex.Lock()
				sender.IsExecuting = true
				selectedKeyName = keyName
				selectedSender = sender
				break
			}
		}
	}
	defer func() {
		cp.FreeSenders[selectedKeyName].Mutex.Unlock()
	}()

	cp.Log.Debug(
		fmt.Sprintf("Relaying packet using address: %v", selectedSender.Address),
		zap.String("evm_sender_address", selectedSender.Address.String()),
		zap.String("chain_name", cp.ChainName),
		zap.Uint64("tunnel_id", packet.TunnelID),
		zap.Uint64("sequence", packet.Sequence),
	)

	tx, err := cp.newRelayTx(ctx, calldata, selectedSender.Address, retryCount)
	if err != nil {
		return "", fmt.Errorf("failed to create an evm transaction: %w", err)
	}

	signedTx, err := cp.signTx(tx, selectedSender)
	if err != nil {
		return "", fmt.Errorf("failed to sign an evm transaction: %w", err)
	}

	txHash, err = cp.Client.BroadcastTx(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast an evm transaction: %w", err)
	}

	cp.FreeSenders[selectedKeyName].IsExecuting = false

	return txHash, nil
}

// checkConfirmedTx checks the confirmed transaction status.
func (cp *EVMChainProvider) checkConfirmedTx(
	ctx context.Context,
	txHash string,
	createdAt time.Time,
) (*ConfirmTxResult, error) {
	failResult := NewConfirmTxResult(
		txHash,
		TX_STATUS_UNMINED,
		decimal.NullDecimal{},
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
		if time.Now().Unix() > createdAt.Add(cp.Config.WaitingTxDuration).Unix() {
			return failResult.WithStatus(TX_STATUS_TIMEOUT), nil
		}
		return failResult.WithStatus(TX_STATUS_UNMINED), nil
	}

	// calculate gas used and effective gas price
	gasUsed := decimal.NewNullDecimal(decimal.New(int64(receipt.GasUsed), 0))
	effGasPrice, err := cp.Client.GetEffectiveGasPrice(ctx, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to get effective gas price: %w", err)
	}

	return NewConfirmTxResult(txHash, TX_STATUS_SUCCESS, gasUsed, effGasPrice), nil
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
	retryCount int,
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

	// calculate new fee info
	feeInfo := cp.GasModel.
		GetGas(ctx).
		Bump(math.Pow(cp.Config.GasMultiplier, float64(retryCount))).
		Param()

	var tx *gethtypes.Transaction
	switch cp.GasModel.GasType() {
	case gas.GasTypeLegacy:
		tx = gethtypes.NewTx(&gethtypes.LegacyTx{
			Nonce:    nonce,
			To:       &cp.TunnelRouterAddress,
			Value:    decimal.NewFromInt(0).BigInt(),
			Data:     data,
			Gas:      gasLimit,
			GasPrice: big.NewInt(int64(feeInfo.GasPrice)),
		})

	case gas.GasTypeEIP1559:
		tx = gethtypes.NewTx(&gethtypes.DynamicFeeTx{
			ChainID:   big.NewInt(int64(cp.Config.ChainID)),
			Nonce:     nonce,
			To:        &cp.TunnelRouterAddress,
			Value:     decimal.NewFromInt(0).BigInt(),
			Data:      data,
			Gas:       gasLimit,
			GasFeeCap: big.NewInt(int64(feeInfo.MaxPriorityFee + feeInfo.MaxBaseFee)),
			GasTipCap: big.NewInt(int64(feeInfo.MaxPriorityFee)),
		})

	default:
		return nil, fmt.Errorf("unsupported gas type: %v", cp.GasModel.GasType())
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
	switch cp.GasModel.GasType() {
	case gas.GasTypeLegacy:
		signer = gethtypes.NewEIP155Signer(big.NewInt(int64(cp.Config.ChainID)))
	case gas.GasTypeEIP1559:
		signer = gethtypes.NewLondonSigner(big.NewInt(int64(cp.Config.ChainID)))
	default:
		return nil, fmt.Errorf("unsupported gas type: %v", cp.GasModel.GasType())
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
	gethaddr := cp.FreeSenders[keyName].Address

	return cp.Client.GetBalance(ctx, gethaddr)
}
