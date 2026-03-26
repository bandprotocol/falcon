package soroban

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/alert"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/db"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/wallet"
	"github.com/bandprotocol/falcon/relayer/wallet/soroban"
)

var _ chains.ChainProvider = (*SorobanChainProvider)(nil)

// SorobanChainProvider handles interactions with Soroban.
type SorobanChainProvider struct {
	Config    *SorobanChainProviderConfig
	ChainName string

	Client Client
	Log    logger.Logger
	DB     db.Database
	Alert  alert.Alert

	FreeSigners chan wallet.Signer
	Wallet      wallet.Wallet
}

// NewSorobanChainProvider creates a new Soroban chain provider.
func NewSorobanChainProvider(
	chainName string,
	client Client,
	cfg *SorobanChainProviderConfig,
	log logger.Logger,
	w wallet.Wallet,
	a alert.Alert,
) *SorobanChainProvider {
	return &SorobanChainProvider{
		Config:      cfg,
		ChainName:   chainName,
		Client:      client,
		Log:         log.With("chain_name", chainName),
		Alert:       a,
		FreeSigners: chains.LoadSigners(w),
		Wallet:      w,
	}
}

func (cp *SorobanChainProvider) Init(ctx context.Context) error {
	if err := cp.Client.Connect(ctx); err != nil {
		return err
	}
	go cp.Client.StartLivelinessCheck(ctx, cp.Config.LivelinessCheckingInterval)
	return nil
}

func (cp *SorobanChainProvider) SetDatabase(database db.Database) {
	cp.DB = database
}

func (cp *SorobanChainProvider) QueryTunnelInfo(
	_ context.Context,
	tunnelID uint64,
	tunnelDestinationAddr string,
) (*types.Tunnel, error) {
	// Soroban uses Skipable tunnels without sequence tracking similar to XRPL
	tunnel := types.NewTunnel(tunnelID, tunnelDestinationAddr, true, nil, nil)
	return tunnel, nil
}

func (cp *SorobanChainProvider) RelayPacket(ctx context.Context, packet *bandtypes.Packet) error {
	if err := cp.Client.CheckAndConnect(ctx); err != nil {
		return fmt.Errorf("[SorobanProvider] failed to connect client: %w", err)
	}

	var freeSigner wallet.Signer
	defer func() {
		if freeSigner != nil {
			cp.FreeSigners <- freeSigner
		}
	}()

SignerLoop:
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("[SorobanProvider] context canceled while waiting for signer: %w", ctx.Err())
		case s := <-cp.FreeSigners:
			freeSigner = s
			break SignerLoop
		}
	}

	log := cp.Log.With(
		"tunnel_id", packet.TunnelID,
		"sequence", packet.Sequence,
		"signer_address", freeSigner.GetAddress(),
	)

	var lastErr error
	for retryCount := 1; retryCount <= cp.Config.MaxRetry; retryCount++ {
		log.Info("Relaying a message", "retry_count", retryCount)

		sequence, err := cp.Client.GetAccountSequenceNumber(freeSigner.GetAddress())
		if err != nil {
			log.Error("Get account sequence number error", "retry_count", retryCount, err)
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}

		signing, err := chains.SelectSigning(packet)
		if err != nil {
			log.Error("Select signing error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		signerPayload := soroban.NewSignerPayload(
			freeSigner.GetAddress(),
			cp.Config.ContractAddress,
			cp.Config.Fee,
			sequence,
			cp.Config.NetworkPassphrase,
		)

		payloadBytes, err := json.Marshal(signerPayload)
		if err != nil {
			log.Error("Marshal signer payload error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}

		rAddress, signature := chains.ExtractEVMSignature(signing.EVMSignature)
		tssPayload := wallet.NewTssPayload(
			signing.Message,
			rAddress,
			signature,
		)

		result, err := freeSigner.Sign(payloadBytes, tssPayload)
		if err != nil {
			log.Error("Sign transaction error", "retry_count", retryCount, err)
			lastErr = err
			continue
		}
		txBlob := string(result)

		var balance *big.Int
		if cp.DB != nil {
			balance, err = cp.Client.GetBalance(freeSigner.GetAddress())
			if err != nil {
				log.Error("Failed to get balance", "retry_count", retryCount, err)
				alert.HandleAlert(cp.Alert, alert.NewTopic(alert.GetBalanceErrorMsg).
					WithTunnelID(packet.TunnelID).
					WithChainName(cp.ChainName), err.Error())
			}
		}

		txResult, err := cp.Client.BroadcastTx(txBlob)
		if err != nil {
			log.Error("Broadcast transaction error", "retry_count", retryCount, err)
			lastErr = err
			time.Sleep(2 * time.Second)

			cp.handleSaveTransaction(
				txResult,
				types.TX_STATUS_FAILED,
				freeSigner.GetAddress(),
				packet,
				balance,
				log,
				retryCount,
			)
			continue
		}

		log.Info("Packet is successfully relayed", "tx_hash", txResult.TxHash, "retry_count", retryCount)

		cp.handleSaveTransaction(
			txResult,
			types.TX_STATUS_SUCCESS,
			freeSigner.GetAddress(),
			packet,
			balance,
			log,
			retryCount,
		)

		relayermetrics.IncTxsCount(
			packet.TunnelID,
			cp.ChainName,
			types.ChainTypeSoroban.String(),
			types.TX_STATUS_SUCCESS.String(),
		)
		alert.HandleReset(
			cp.Alert,
			alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		)

		return nil
	}

	alert.HandleAlert(
		cp.Alert,
		alert.NewTopic(alert.RelayTxErrorMsg).WithTunnelID(packet.TunnelID).WithChainName(cp.ChainName),
		lastErr.Error(),
	)
	return fmt.Errorf("[SorobanProvider] failed to relay packet after %d attempts", cp.Config.MaxRetry)
}

func (cp *SorobanChainProvider) QueryBalance(ctx context.Context, address string) (*big.Int, error) {
	return cp.Client.GetBalance(address)
}

func (cp *SorobanChainProvider) GetChainName() string { return cp.ChainName }
func (cp *SorobanChainProvider) ChainType() types.ChainType { return types.ChainTypeSoroban }
func (cp *SorobanChainProvider) GetWallet() wallet.Wallet { return cp.Wallet }

func (cp *SorobanChainProvider) handleSaveTransaction(
	txResult TxResult,
	txStatus types.TxStatus,
	signerAddress string,
	packet *bandtypes.Packet,
	oldBalance *big.Int,
	log logger.Logger,
	retryCount int,
) {
	if cp.DB != nil {
		if txResult.TxHash == "" {
			return
		}

		var signalPrices []db.SignalPrice
		for _, p := range packet.SignalPrices {
			signalPrices = append(signalPrices, *db.NewSignalPrice(p.SignalID, p.Price))
		}

		fee := decimal.NullDecimal{}
		feeDecimal, err := decimal.NewFromString(cp.Config.Fee)
		if err == nil {
			fee = decimal.NewNullDecimal(feeDecimal)
		}

		balanceDelta := decimal.NullDecimal{}
		if oldBalance != nil {
			newBalance, err := cp.Client.GetBalance(signerAddress)
			if err == nil {
				diff := new(big.Int).Sub(newBalance, oldBalance)
				balanceDelta = decimal.NewNullDecimal(decimal.NewFromBigInt(diff, 7))
			}
		}

		var closeTime *time.Time
		if txResult.LedgerIndex != 0 {
			closeTime, _ = cp.Client.GetLedgerCloseTime(txResult.LedgerIndex)
		}

		tx := db.NewTransaction(
			txResult.TxHash,
			packet.TunnelID,
			packet.Sequence,
			cp.ChainName,
			types.ChainTypeSoroban,
			signerAddress,
			txStatus,
			decimal.NewNullDecimal(decimal.NewFromInt(1)),
			fee,
			balanceDelta,
			signalPrices,
			closeTime,
		)

		chains.HandleSaveTransaction(cp.DB, cp.Alert, tx, log)
	}
}
