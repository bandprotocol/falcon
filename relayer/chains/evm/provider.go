package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/alert"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/db"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ chains.ChainProvider = (*EVMChainProvider)(nil)

// EVMChainProvider is the struct that handles interactions with the EVM chain.
type EVMChainProvider struct {
	Config    *EVMChainProviderConfig
	ChainName string

	Client  Client
	GasType GasType

	FreeSigners chan wallet.Signer

	TunnelRouterAddress gethcommon.Address
	TunnelRouterABI     abi.ABI

	Log logger.Logger

	Wallet wallet.Wallet
	DB     db.Database

	Alert alert.Alert
}

// NewEVMChainProvider creates a new EVM chain provider.
func NewEVMChainProvider(
	chainName string,
	client Client,
	cfg *EVMChainProviderConfig,
	log logger.Logger,
	wallet wallet.Wallet,
	alert alert.Alert,
) (*EVMChainProvider, error) {
	// load abis here
	abi, err := abi.JSON(strings.NewReader(gasPriceTunnelRouterABI))
	if err != nil {
		log.Error("ChainProvider: failed to load abi",
			"chain_name", chainName,
			err,
		)
		return nil, fmt.Errorf("[EVMProvider] failed to load abi: %w", err)
	}

	addr, err := HexToAddress(cfg.TunnelRouterAddress)
	if err != nil {
		log.Error("ChainProvider: cannot convert tunnel router address",
			"chain_name", chainName,
			err,
		)
		return nil, fmt.Errorf("[EVMProvider] incorrect address: %w", err)
	}

	return &EVMChainProvider{
		Config:              cfg,
		ChainName:           chainName,
		Client:              client,
		GasType:             cfg.GasType,
		TunnelRouterAddress: addr,
		TunnelRouterABI:     abi,
		Log:                 log.With("chain_name", chainName),
		Wallet:              wallet,
		Alert:               alert,
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

// SetDatabase assigns the given database instance to the EVMChainProvider.
func (cp *EVMChainProvider) SetDatabase(database db.Database) {
	cp.DB = database
}

// QueryTunnelInfo queries the tunnel info from the tunnel router contract.
func (cp *EVMChainProvider) QueryTunnelInfo(
	ctx context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		cp.Log.Error("Connect client error", err)
		return nil, fmt.Errorf("[EVMProvider] failed to connect client: %w", err)
	}

	addr, err := HexToAddress(tunnelDestinationAddr)
	if err != nil {
		cp.Log.Error("Invalid address", "address", tunnelDestinationAddr, err)
		return nil, fmt.Errorf("[EVMProvider] invalid address: %w", err)
	}

	info, err := cp.queryTunnelInfo(ctx, tunnelID, addr)
	if err != nil {
		cp.Log.Error(
			"Query contract error",
			"tunnel_id", tunnelID,
			"address", tunnelDestinationAddr,
			err,
		)

		return nil, fmt.Errorf("[EVMProvider] failed to query contract: %w", err)
	}

	return &types.Tunnel{
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
		cp.Log.Error("Connect client error", err)
		return fmt.Errorf("[EVMProvider] failed to connect client: %w", err)
	}

	// get a free signer
	cp.Log.Debug("Waiting for a free signer...")
	freeSigner := <-cp.FreeSigners
	defer func() { cp.FreeSigners <- freeSigner }()

	log := cp.Log.With(
		"tunnel_id", packet.TunnelID,
		"sequence", packet.Sequence,
		"signer_address", freeSigner.GetAddress(),
	)

	// get gas information
	gasInfo, err := cp.EstimateGasFee(ctx)
	if err != nil {
		cp.Log.Error("Failed to estimate gas fee", err)
		alert.HandleAlert(
			cp.Alert,
			alert.NewTopic(alert.EstimateGasFeeErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
			err.Error(),
		)
		return fmt.Errorf("[EVMProvider] failed to estimate gas fee: %w", err)
	}
	alert.HandleReset(cp.Alert, alert.NewTopic(alert.EstimateGasFeeErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName))

	var lastErr error
	var bumpGasErr error
	for retryCount := 1; retryCount <= cp.Config.MaxRetry; retryCount++ {
		log.Info("Relaying a message", "retry_count", retryCount)

		// create and submit a transaction; if failed, retry, no need to bump gas.
		signedTx, err := cp.createAndSignRelayTx(ctx, packet, freeSigner, gasInfo)
		if err != nil {
			lastErr = fmt.Errorf("create and sign tx error: %v", err)
			log.Error("CreateAndSignTx error", "retry_count", retryCount, err)
			continue
		}

		balance, err := cp.Client.GetBalance(ctx, gethcommon.HexToAddress(freeSigner.GetAddress()), nil)
		if err != nil {
			log.Error("GetBalance error", err)
		}

		// submit the transaction, if failed, bump gas and retry
		txHash, err := cp.Client.BroadcastTx(ctx, signedTx)
		if err != nil {
			lastErr = fmt.Errorf("broadcast tx error: %v", err)
			log.Error("HandleRelay error", "retry_count", retryCount, err)
			// bump gas and retry
			gasInfo, bumpGasErr = cp.BumpAndBoundGas(ctx, gasInfo, cp.Config.GasMultiplier)
			if bumpGasErr != nil {
				log.Error("Cannot bump gas", "retry_count", retryCount, bumpGasErr)
			}

			continue
		}

		createdAt := time.Now()

		log.Info(
			"Submitted a message; checking transaction status",
			"tx_hash", txHash,
			"retry_count", retryCount,
		)

		if err := cp.saveUnconfirmedTransaction(txHash, types.TX_STATUS_PENDING, packet, freeSigner.GetAddress()); err != nil {
			log.Error("SaveTransaction error", "retry_count", retryCount, err)
		}

		txResult := cp.WaitForConfirmedTx(ctx, txHash, log)

		cp.handleMetrics(packet.TunnelID, createdAt, txResult)
		cp.handleSaveTransaction(ctx, freeSigner.GetAddress(), balance, packet, txResult, retryCount, log)

		if txResult.Status == types.TX_STATUS_SUCCESS {
			log.Info(
				"Packet is successfully relayed",
				"tx_hash", txHash,
				"retry_count", retryCount,
			)
			alert.HandleReset(
				cp.Alert,
				alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
			)
			return nil
		}

		lastErr = fmt.Errorf("%s", txResult.FailureReason)
		log.Error(
			"Failed to relaying a packet with status and error",
			"status", txResult.Status.String(),
			"tx_hash", txHash,
			"retry_count", retryCount,
			lastErr,
		)

		// bump gas and retry
		gasInfo, bumpGasErr = cp.BumpAndBoundGas(ctx, gasInfo, cp.Config.GasMultiplier)
		if bumpGasErr != nil {
			log.Error("Cannot bump gas", "retry_count", retryCount, bumpGasErr)
		}
	}

	if bumpGasErr != nil {
		// add bump gas error detail to the latest error
		lastErr = fmt.Errorf("%v; %v", lastErr, bumpGasErr)
	}

	alert.HandleAlert(
		cp.Alert,
		alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		lastErr.Error(),
	)

	return fmt.Errorf("[EVMProvider] failed to relay packet after %d retries", cp.Config.MaxRetry)
}

// createAndSignRelayTx creates and signs the relay transaction.
func (cp *EVMChainProvider) createAndSignRelayTx(
	ctx context.Context,
	packet *bandtypes.Packet,
	signer wallet.Signer,
	gasInfo GasInfo,
) (*gethtypes.Transaction, error) {
	calldata, err := cp.CreateCalldata(packet)
	if err != nil {
		return nil, fmt.Errorf("failed to create calldata: %w", err)
	}

	tx, err := cp.NewRelayTx(ctx, calldata, signer, gasInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create an evm transaction: %w", err)
	}

	signedTx, err := cp.signTx(tx, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to sign an evm transaction: %w", err)
	}

	return signedTx, nil
}

// WaitForConfirmedTx polls the transaction until it reaches a terminal state.
// It NEVER returns an error. Instead, it always returns a TxResult where:
//   - Status == TX_STATUS_SUCCESS or TX_STATUS_FAILED when confirmed.
//   - Status == TX_STATUS_TIMEOUT if it did not reach the required confirmations
//     within WaitingTxDuration (or the context was canceled); in this case,
//     the result’s FailureReason field is populated with details.
//
// The function sleeps for CheckingTxInterval between polls.
func (cp *EVMChainProvider) WaitForConfirmedTx(
	ctx context.Context,
	txHash string,
	log logger.Logger,
) TxResult {
	createdAt := time.Now()
	var lastErr error
	for time.Since(createdAt) <= cp.Config.WaitingTxDuration {
		result, err := cp.CheckConfirmedTx(ctx, txHash)
		if err != nil {
			lastErr = err
			log.Debug(
				"Failed to check tx status",
				"tx_hash", txHash,
				err,
			)
		}

		switch result.Status {
		case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
			return result
		case types.TX_STATUS_PENDING:
			log.Debug(
				"Waiting for tx to be mined",
				"tx_hash", txHash,
			)
			time.Sleep(cp.Config.CheckingTxInterval)
		}
	}

	failureReason := fmt.Sprintf("timed out waiting %s for tx %s to reach %d confirmations",
		cp.Config.WaitingTxDuration, txHash, cp.Config.BlockConfirmation)

	if lastErr != nil {
		failureReason = fmt.Sprintf("%s: %v", failureReason, lastErr)
	}

	return NewTxResult(
		txHash,
		types.TX_STATUS_TIMEOUT,
		decimal.NullDecimal{},
		decimal.NullDecimal{},
		nil,
		failureReason,
	)
}

// handleMetrics increments tx count and, for success/failed, records processing time (ms) and gas used.
func (cp *EVMChainProvider) handleMetrics(tunnelID uint64, createdAt time.Time, txResult TxResult) {
	// increment the transactions count metric
	relayermetrics.IncTxsCount(tunnelID, cp.ChainName, txResult.Status.String())

	switch txResult.Status {
	case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
		// track transaction processing time (ms)
		relayermetrics.ObserveTxProcessTime(
			tunnelID,
			cp.ChainName,
			txResult.Status.String(),
			time.Since(createdAt).Milliseconds(),
		)

		// track gas used for the relayed transaction
		relayermetrics.ObserveGasUsed(
			tunnelID,
			cp.ChainName,
			txResult.Status.String(),
			txResult.GasUsed.Decimal.InexactFloat64(),
		)
	}
}

// handleSaveTransaction saves the transaction to the database based on its status.
func (cp *EVMChainProvider) handleSaveTransaction(ctx context.Context,
	signerAddress string,
	oldBalance *big.Int,
	packet *bandtypes.Packet,
	txResult TxResult,
	retryCount int,
	log logger.Logger,
) {
	switch txResult.Status {
	case types.TX_STATUS_SUCCESS, types.TX_STATUS_FAILED:
		if err := cp.saveConfirmedTransaction(ctx, signerAddress, oldBalance, packet, txResult); err != nil {
			log.Error("SaveTransaction error", "retry_count", retryCount, err)
		}
	default:
		if err := cp.saveUnconfirmedTransaction(txResult.TxHash, txResult.Status, packet, signerAddress); err != nil {
			log.Error("SaveTransaction error", "retry_count", retryCount, err)
		}
	}
}

// CheckConfirmedTx checks the confirmed transaction status.
func (cp *EVMChainProvider) CheckConfirmedTx(
	ctx context.Context,
	txHash string,
) (TxResult, error) {
	receipt, err := cp.Client.GetTxReceipt(ctx, txHash)
	if err != nil {
		err = fmt.Errorf("failed to get tx receipt: %w", err)
		return NewTxResult(
			txHash,
			types.TX_STATUS_PENDING,
			decimal.NullDecimal{},
			decimal.NullDecimal{},
			nil,
			err.Error(),
		), err
	}

	// calculate gas used and effective gas price
	gasUsed := decimal.NewNullDecimal(decimal.New(int64(receipt.GasUsed), 0))
	gasPrice := decimal.NewNullDecimal(decimal.New(int64(receipt.EffectiveGasPrice.Uint64()), 0))

	if receipt.Status == gethtypes.ReceiptStatusFailed {
		return NewTxResult(
			txHash,
			types.TX_STATUS_FAILED,
			gasUsed,
			gasPrice,
			receipt.BlockNumber,
			"transaction reverted on-chain",
		), nil
	}

	latestBlock, err := cp.Client.GetBlockHeight(ctx)
	if err != nil {
		err = fmt.Errorf("failed to get latest block height: %w", err)
		return NewTxResult(
			txHash,
			types.TX_STATUS_PENDING,
			decimal.NullDecimal{},
			decimal.NullDecimal{},
			nil,
			err.Error(),
		), err
	}

	// if tx block is not confirmed and waiting too long return status with timeout
	if receipt.BlockNumber.Uint64() > latestBlock-cp.Config.BlockConfirmation {
		return NewTxResult(
			txHash,
			types.TX_STATUS_PENDING,
			decimal.NullDecimal{},
			decimal.NullDecimal{},
			nil,
			"",
		), nil
	}

	return NewTxResult(txHash, types.TX_STATUS_SUCCESS, gasUsed, gasPrice, receipt.BlockNumber, ""), nil
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

		maxBaseFee := big.NewInt(int64(cp.Config.MaxBaseFee))
		newBaseFee := gasInfo.GasBaseFee
		if maxBaseFee.Cmp(big.NewInt(0)) > 0 && newBaseFee.Cmp(maxBaseFee) > 0 {
			newBaseFee = maxBaseFee
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

// NewRelayTx creates a new relay transaction.
func (cp *EVMChainProvider) NewRelayTx(
	ctx context.Context,
	data []byte,
	signer wallet.Signer,
	gasInfo GasInfo,
) (*gethtypes.Transaction, error) {
	addr := gethcommon.HexToAddress(signer.GetAddress())
	nonce, err := cp.Client.NonceAt(ctx, addr)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:      addr,
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
		return nil, fmt.Errorf("unsupported gas type: %v", cp.GasType)
	}

	return tx, nil
}

// CreateCalldata creates the calldata for the relay transaction.
func (cp *EVMChainProvider) CreateCalldata(packet *bandtypes.Packet) ([]byte, error) {
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

// signTx signs the transaction with the signer.
func (cp *EVMChainProvider) signTx(
	tx *gethtypes.Transaction,
	signer wallet.Signer,
) (*gethtypes.Transaction, error) {
	var (
		rlpEncoded []byte
		err        error
		gethSigner gethtypes.Signer
	)

	chainID := big.NewInt(int64(cp.Config.ChainID))

	switch cp.GasType {
	case GasTypeLegacy:
		rlpEncoded, err = rlp.EncodeToBytes(
			[]interface{}{
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				chainID, uint(0), uint(0),
			},
		)
		if err != nil {
			return nil, err
		}

		gethSigner = gethtypes.NewEIP155Signer(chainID)
	case GasTypeEIP1559:
		rlpEncoded, err = rlp.EncodeToBytes(
			[]interface{}{
				chainID,
				tx.Nonce(),
				tx.GasTipCap(),
				tx.GasFeeCap(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				tx.AccessList(),
			},
		)
		if err != nil {
			return nil, err
		}

		rlpEncoded = append([]byte{tx.Type()}, rlpEncoded...)
		gethSigner = gethtypes.NewLondonSigner(chainID)

	default:
		return nil, fmt.Errorf("unsupported gas type: %v", cp.GasType)
	}

	signature, err := signer.Sign(rlpEncoded)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(gethSigner, signature)
}

// QueryBalance queries balance of specific account address.
func (cp *EVMChainProvider) QueryBalance(
	ctx context.Context,
	keyName string,
) (*big.Int, error) {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		cp.Log.Error(
			"Connect client error",
			"chain_name", cp.ChainName,
			err,
		)
		return nil, fmt.Errorf("[EVMProvider] failed to connect client: %w", err)
	}

	signer, ok := cp.Wallet.GetSigner(keyName)
	if !ok {
		cp.Log.Error("Key name does not exist", "key_name", keyName)
		return nil, fmt.Errorf("key name does not exist: %s", keyName)
	}

	address, err := HexToAddress(signer.GetAddress())
	if err != nil {
		return nil, err
	}

	return cp.Client.GetBalance(ctx, address, nil)
}

// GetChainName retrieves the chain name from the chain provider.
func (cp *EVMChainProvider) GetChainName() string {
	return cp.ChainName
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

func (cp *EVMChainProvider) saveUnconfirmedTransaction(
	txHash string,
	txStatus types.TxStatus,
	packet *bandtypes.Packet,
	sender string,
) error {
	// db was disabled
	if cp.DB == nil {
		return nil
	}

	var signalPrices []db.SignalPrice
	for _, p := range packet.SignalPrices {
		signalPrices = append(signalPrices, *db.NewSignalPrice(p.SignalID, p.Price))
	}

	tx := db.NewUnconfirmedTransaction(
		txHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeEVM,
		sender,
		txStatus,
		signalPrices,
	)

	if err := cp.DB.AddOrUpdateTransaction(tx); err != nil {
		return fmt.Errorf("failed to save transaction to database: %w", err)
	}

	return nil
}

// saveConfirmedTransaction stores the transaction result and related metadata (e.g. gas, status, balance delta) to the database if enabled.
func (cp *EVMChainProvider) saveConfirmedTransaction(
	ctx context.Context,
	signerAddress string,
	oldBalance *big.Int,
	packet *bandtypes.Packet,
	txResult TxResult,
) error {
	// db was disabled
	if cp.DB == nil {
		return nil
	}

	var signalPrices []db.SignalPrice
	for _, p := range packet.SignalPrices {
		signalPrices = append(signalPrices, *db.NewSignalPrice(p.SignalID, p.Price))
	}

	var blockTimestamp time.Time
	balanceDelta := decimal.NullDecimal{}

	block, err := cp.Client.GetBlock(ctx, txResult.BlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get block: %w", err)
	}

	blockTimestamp = time.Unix(int64(block.Time()), 0).UTC()

	// Compute new balance
	// Note: this may be incorrect if other transactions affected the user's balance during this period.
	if oldBalance != nil {
		newBalance, err := cp.Client.GetBalance(ctx, gethcommon.HexToAddress(signerAddress), txResult.BlockNumber)
		if err != nil {
			return fmt.Errorf("failed to get balance: %w", err)
		}
		diff := new(big.Int).Sub(newBalance, oldBalance)
		balanceDelta = decimal.NewNullDecimal(decimal.NewFromBigInt(diff, 0))
	}

	tx := db.NewConfirmedTransaction(
		txResult.TxHash,
		packet.TunnelID,
		packet.Sequence,
		cp.ChainName,
		types.ChainTypeEVM,
		signerAddress,
		txResult.Status,
		txResult.GasUsed,
		txResult.EffectiveGasPrice,
		balanceDelta,
		signalPrices,
		blockTimestamp,
	)

	if err := cp.DB.AddOrUpdateTransaction(tx); err != nil {
		return fmt.Errorf("failed to save transaction to database: %w", err)
	}

	return nil
}
