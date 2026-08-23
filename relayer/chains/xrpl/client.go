package xrpl

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	xrplaccount "github.com/Peersyst/xrpl-go/xrpl/queries/account"
	"github.com/Peersyst/xrpl-go/xrpl/queries/common"
	"github.com/Peersyst/xrpl-go/xrpl/queries/ledger"
	requests "github.com/Peersyst/xrpl-go/xrpl/queries/transactions"
	"github.com/Peersyst/xrpl-go/xrpl/rpc"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/logger"
)

const RippleEpochOffset = 946684800

const (
	queuedTransactionResult     = "terQUEUED"
	successfulTransactionResult = "tesSUCCESS"
	transactionNotFoundError    = "txnNotFound"
	defaultTxPollingInterval    = time.Second
)

// idleConnTimeout is shorter than the typical load-balancer idle timeout
// (e.g. AWS ALB default 60s) to prevent "connection reset by peer" errors
// from stale keep-alive connections being reused after the LB closes them.
const idleConnTimeout = 30 * time.Second

// XRPLClients holds XRPL RPC clients and the selected endpoint.
type XRPLClients = chains.ClientPool[rpc.Client]

// NewXRPLClients creates and returns a new XRPLClients instance with no endpoints.
func NewXRPLClients() XRPLClients {
	return chains.NewClientPool[rpc.Client]()
}

// Client is the interface that handles interactions with the XRPL chain.
type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	GetAccountSequenceNumber(account string) (uint32, error)
	GetBalance(account string) (*big.Int, error)
	Autofill(tx *transaction.FlatTransaction) error
	BroadcastTx(ctx context.Context, txBlob string) (TxResult, error)
	GetLedgerCloseTime(ledgerIndex common.LedgerIndex) (*time.Time, error)
}

var _ Client = (*client)(nil)

// client is the concrete implementation that handles XRPL JSON-RPC interactions.
type client struct {
	ChainName         string
	Endpoints         []string
	BlockConfirmation uint32
	TxPollingInterval time.Duration

	Log   logger.Logger
	alert alert.Alert

	clients XRPLClients
}

type TxResult struct {
	TxHash      string
	Fee         string
	LedgerIndex common.LedgerIndex
}

// NewClient creates a new XRPL client from config.
func NewClient(chainName string, cfg *XRPLChainProviderConfig, log logger.Logger, alert alert.Alert) Client {
	return &client{
		ChainName:         chainName,
		Endpoints:         cfg.Endpoints,
		BlockConfirmation: 5,
		TxPollingInterval: cfg.NonceInterval,
		Log:               log.With("chain_name", chainName),
		alert:             alert,
		clients:           NewXRPLClients(),
	}
}

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *rpc.Client
	LedgerIndex uint32
}

// Connect connects to all endpoints and selects the first eligible one by config order
// within the confirmed ledger range.
func (c *client) Connect(_ context.Context) error {
	var wg sync.WaitGroup
	for _, endpoint := range c.Endpoints {
		_, ok := c.clients.GetClient(endpoint)
		if ok {
			continue
		}

		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			transport := http.DefaultTransport.(*http.Transport).Clone()
			transport.IdleConnTimeout = idleConnTimeout
			opts := []rpc.ConfigOpt{
				rpc.WithHTTPClient(&http.Client{Transport: transport}),
			}
			cfg, err := rpc.NewClientConfig(endpoint, opts...)
			if err != nil {
				c.Log.Warn("XRPL endpoint config error", "endpoint", endpoint, err)
				alert.HandleAlert(
					c.alert,
					alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
						WithChainName(c.ChainName).
						WithEndpoint(endpoint),
					err.Error(),
				)
				return
			}

			alert.HandleReset(
				c.alert,
				alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
					WithChainName(c.ChainName).
					WithEndpoint(endpoint),
			)
			client := rpc.NewClient(cfg)
			c.clients.SetClient(endpoint, client)
		}(endpoint)
	}

	wg.Wait()
	res, err := c.getClientWithMaxHeight()
	if err != nil {
		c.Log.Error("Failed to connect to XRPL chain", err)
		return fmt.Errorf("failed to connect to XRPL chain: %w", err)
	}

	// only log when new endpoint is used
	if c.clients.GetSelectedEndpoint() != res.Endpoint {
		c.Log.Info("Connected to XRPL chain", "endpoint", res.Endpoint)
	}

	c.clients.SetSelectedEndpoint(res.Endpoint)

	return nil
}

// getClientWithMaxHeight selects an endpoint by config order within the confirmed
// ledger range. It collects ledger indices from all endpoints in parallel, then
// picks the first endpoint (in config order) whose ledger index is within
// [maxLedger - BlockConfirmation, maxLedger].
func (c *client) getClientWithMaxHeight() (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.Endpoints))

	for _, endpoint := range c.Endpoints {
		go func(endpoint string) {
			client, ok := c.clients.GetClient(endpoint)
			if !ok {
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			ledgerIndex, err := client.GetLedgerIndex()
			if err != nil {
				c.Log.Warn("Failed to get ledger index", "endpoint", endpoint, "err", err)
				ch <- ClientConnectionResult{endpoint, nil, 0}
				alert.HandleAlert(
					c.alert,
					alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
						WithChainName(c.ChainName).
						WithEndpoint(endpoint),
					err.Error(),
				)
				return
			}

			alert.HandleReset(
				c.alert,
				alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
					WithChainName(c.ChainName).
					WithEndpoint(endpoint),
			)

			ch <- ClientConnectionResult{endpoint, client, uint32(ledgerIndex)}
		}(endpoint)
	}

	// Collect all results into a map and track the maximum ledger index.
	resultMap := make(map[string]ClientConnectionResult, len(c.Endpoints))
	var maxLedger uint32
	for range c.Endpoints {
		r := <-ch
		if r.Client != nil {
			resultMap[r.Endpoint] = r
			if r.LedgerIndex > maxLedger {
				maxLedger = r.LedgerIndex
			}
		}
	}

	// Determine the minimum acceptable ledger index based on BlockConfirmation.
	var minLedger uint32
	if maxLedger >= c.BlockConfirmation {
		minLedger = maxLedger - c.BlockConfirmation
	}

	// Pick the first endpoint in config order that is within the confirmed range.
	var result ClientConnectionResult
	for _, endpoint := range c.Endpoints {
		if r, ok := resultMap[endpoint]; ok && r.LedgerIndex >= minLedger {
			result = r
			break
		}
	}

	if result.Client == nil {
		alert.HandleAlert(
			c.alert,
			alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName),
			fmt.Sprintf("failed to connect to XRPL chain on all endpoints: %s", c.Endpoints),
		)
		return ClientConnectionResult{}, fmt.Errorf("failed to connect to XRPL chain on all endpoints")
	}

	alert.HandleReset(c.alert, alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName))

	return result, nil
}

// CheckAndConnect checks if the client is connected to the XRPL chain, if not connect it.
func (c *client) CheckAndConnect(ctx context.Context) error {
	if _, err := c.clients.GetSelectedClient(); err != nil {
		return c.Connect(ctx)
	}

	return nil
}

// StartLivelinessCheck starts the liveliness check for the XRPL chain.
func (c *client) StartLivelinessCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.Log.Info("Stopping liveliness check")
			return
		case <-ticker.C:
			err := c.Connect(ctx)
			if err != nil {
				c.Log.Error("Liveliness check: unable to reconnect to any endpoints", err)
			}
		}
	}
}

// GetAccountSequenceNumber fetches the sequence for the given account.
func (c *client) GetAccountSequenceNumber(account string) (uint32, error) {
	client, err := c.clients.GetSelectedClient()
	if err != nil {
		return 0, fmt.Errorf("failed to get client: %w", err)
	}

	result, err := client.GetAccountInfo(&xrplaccount.InfoRequest{
		Account: types.Address(account),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get account info: %w", err)
	}

	return result.AccountData.Sequence, nil
}

// GetBalance fetches the XRP balance for the given account (drops).
func (c *client) GetBalance(account string) (*big.Int, error) {
	client, err := c.clients.GetSelectedClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	result, err := client.GetAccountInfo(&xrplaccount.InfoRequest{
		Account: types.Address(account),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	b := new(big.Int)
	b, ok := b.SetString(result.AccountData.Balance.String(), 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse balance of %s (%s)", account, result.AccountData.Balance.String())
	}

	return b, nil
}

// Autofill completes a transaction with missing Sequence, Fee, and LastLedgerSequence fields.
func (c *client) Autofill(tx *transaction.FlatTransaction) error {
	client, err := c.clients.GetSelectedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	if err := client.Autofill(tx); err != nil {
		return fmt.Errorf("failed to autofill tx: %w", err)
	}
	return nil
}

// BroadcastTx submits a signed tx blob and returns its validated result when it
// is queued. A queued transaction must not be rebuilt with a new sequence while
// its final result is still unknown.
func (c *client) BroadcastTx(ctx context.Context, txBlob string) (TxResult, error) {
	client, err := c.clients.GetSelectedClient()
	if err != nil {
		return TxResult{}, fmt.Errorf("failed to get client: %w", err)
	}

	res, err := client.Request(&requests.SubmitRequest{
		TxBlob:   txBlob,
		FailHard: false,
	})
	if err != nil {
		return TxResult{}, fmt.Errorf("failed to submit tx: %w", err)
	}

	var result requests.SubmitResponse
	if err := res.GetResult(&result); err != nil {
		return TxResult{}, fmt.Errorf("failed to get submit result: %w", err)
	}

	txHash, ok := result.Tx["hash"].(string)
	if !ok || txHash == "" {
		return TxResult{}, fmt.Errorf("missing tx hash in submit response")
	}

	fee, ok := result.Tx["Fee"].(string)
	if !ok || fee == "" {
		return TxResult{
			TxHash: txHash,
		}, fmt.Errorf("missing fee in submit response")
	}

	txResult := TxResult{
		TxHash:      txHash,
		Fee:         fee,
		LedgerIndex: result.ValidatedLedgerIndex,
	}

	if result.EngineResult == queuedTransactionResult {
		return c.waitForQueuedTransaction(ctx, txResult)
	}

	if result.EngineResult != successfulTransactionResult {
		return TxResult{
				TxHash:      txHash,
				Fee:         fee,
				LedgerIndex: result.ValidatedLedgerIndex,
			}, fmt.Errorf(
				"failed to broadcast with engine result %s: %s",
				result.EngineResult,
				result.EngineResultMessage,
			)
	}

	return txResult, nil
}

// waitForQueuedTransaction polls the original transaction hash until its
// validated result is known. Without LastLedgerSequence, txnNotFound and RPC
// errors cannot prove that a queued transaction will never validate, so they
// remain pending until the relayer context is canceled.
func (c *client) waitForQueuedTransaction(
	ctx context.Context,
	txResult TxResult,
) (TxResult, error) {
	pollingInterval := c.TxPollingInterval
	if pollingInterval <= 0 {
		pollingInterval = defaultTxPollingInterval
	}
	for {
		select {
		case <-ctx.Done():
			return txResult, fmt.Errorf("waiting for queued transaction %s: %w", txResult.TxHash, ctx.Err())
		default:
		}

		client, err := c.clients.GetSelectedClient()
		if err == nil {
			var res rpc.XRPLResponse
			res, err = client.Request(&requests.TxRequest{Transaction: txResult.TxHash})
			if err == nil {
				var txResponse requests.TxResponse
				if err = res.GetResult(&txResponse); err == nil && txResponse.Validated {
					txResult.LedgerIndex = txResponse.LedgerIndex
					if txResponse.Meta.TransactionResult != successfulTransactionResult {
						return txResult, fmt.Errorf(
							"queued transaction %s validated with engine result %s",
							txResult.TxHash,
							txResponse.Meta.TransactionResult,
						)
					}

					return txResult, nil
				}
			}
		}

		if err != nil && !isTransactionNotFound(err) {
			c.Log.Warn("Failed to query queued XRPL transaction", "tx_hash", txResult.TxHash, "err", err)
		}
		timer := time.NewTimer(pollingInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return txResult, fmt.Errorf("waiting for queued transaction %s: %w", txResult.TxHash, ctx.Err())
		case <-timer.C:
		}
	}
}

func isTransactionNotFound(err error) bool {
	var clientErr *rpc.ClientError
	return errors.As(err, &clientErr) && clientErr.ErrorString == transactionNotFoundError
}

// GetLedgerCloseTime fetches the close time of the ledger with the given index.
func (c *client) GetLedgerCloseTime(ledgerIndex common.LedgerIndex) (*time.Time, error) {
	client, err := c.clients.GetSelectedClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	ledger, err := client.GetLedger(&ledger.Request{LedgerIndex: ledgerIndex})
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger: %w", err)
	}

	closeTime := time.Unix(int64(ledger.Ledger.CloseTime)+RippleEpochOffset, 0)

	return &closeTime, nil
}
