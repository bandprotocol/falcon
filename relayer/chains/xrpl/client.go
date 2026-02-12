package xrpl

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	xrplaccount "github.com/Peersyst/xrpl-go/xrpl/queries/account"
	"github.com/Peersyst/xrpl-go/xrpl/queries/common"
	requests "github.com/Peersyst/xrpl-go/xrpl/queries/transactions"
	"github.com/Peersyst/xrpl-go/xrpl/rpc"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/logger"
)

// Client is the interface that handles interactions with the XRPL chain.
type Client interface {
	Connect() error
	CheckAndConnect() error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	GetAccountSequenceNumber(ctx context.Context, account string) (uint32, error)
	GetBalance(ctx context.Context, account string) (*big.Int, error)
	Autofill(tx *transaction.FlatTransaction) error
	BroadcastTx(ctx context.Context, txBlob string) (TxResult, error)
}

var _ Client = (*client)(nil)

// client is the concrete implementation that handles XRPL JSON-RPC interactions.
type client struct {
	ChainName      string
	Endpoints      []string
	QueryTimeout   time.Duration
	ExecuteTimeout time.Duration

	Log   logger.Logger
	alert alert.Alert

	rpcClient *rpc.Client

	mu               sync.RWMutex
	selectedEndpoint string
}

type TxResult struct {
	TxHash string
	Fee    string
}

// NewClient creates a new XRPL client from config.
func NewClient(chainName string, cfg *XRPLChainProviderConfig, log logger.Logger, alert alert.Alert) Client {
	return &client{
		ChainName:      chainName,
		Endpoints:      cfg.Endpoints,
		QueryTimeout:   cfg.QueryTimeout,
		ExecuteTimeout: cfg.ExecuteTimeout,
		Log:            log.With("chain_name", chainName),
		alert:          alert,
	}
}

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *rpc.Client
	LedgerIndex uint32
}

// Connect selects a responsive endpoint with the highest ledger index.
func (c *client) Connect() error {
	res, err := c.getClientWithMaxHeight()
	if err != nil {
		c.Log.Error("Failed to connect to XRPL chain", err)
		return err
	}

	// only log when new endpoint is used
	if c.getSelectedEndpoint() != res.Endpoint {
		c.Log.Info("Connected to XRPL chain", "endpoint", res.Endpoint)
	}

	c.setSelectedEndpoint(res.Endpoint)
	c.setRPCClient(res.Client)

	return nil
}

// getClientWithMaxHeight connects to the endpoint that has the highest ledger index.
func (c *client) getClientWithMaxHeight() (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.Endpoints))

	for _, endpoint := range c.Endpoints {
		go func(endpoint string) {
			timeout := c.QueryTimeout
			opts := []rpc.ConfigOpt{rpc.WithTimeout(timeout)}
			cfg, err := rpc.NewClientConfig(endpoint, opts...)
			if err != nil {
				c.Log.Warn("XRPL endpoint config error", "endpoint", endpoint, err)
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			client := rpc.NewClient(cfg)
			ledgerIndex, err := client.GetLedgerIndex()
			if err != nil {
				c.Log.Warn("Failed to get ledger index", "endpoint", endpoint, "err", err)
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			ch <- ClientConnectionResult{endpoint, client, uint32(ledgerIndex)}
		}(endpoint)
	}

	var result ClientConnectionResult
	for range c.Endpoints {
		r := <-ch
		if r.Client != nil {
			if r.LedgerIndex > result.LedgerIndex || (r.Endpoint == c.getSelectedEndpoint() && r.LedgerIndex == result.LedgerIndex) {
				result = r
			}
		}
	}

	if result.Client == nil {
		return ClientConnectionResult{}, fmt.Errorf("[XRPLClient] failed to connect to XRPL chain on all endpoints")
	}

	return result, nil
}

// CheckAndConnect checks if the client is connected to the XRPL chain, if not connect it.
func (c *client) CheckAndConnect() error {
	if _, err := c.getRPCClient(); err != nil {
		return c.Connect()
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
			err := c.Connect()
			if err != nil {
				c.Log.Error("Liveliness check: unable to reconnect to any endpoints", err)
			}
		}
	}
}

// GetAccountSequenceNumber fetches the sequence for the given account.
func (c *client) GetAccountSequenceNumber(ctx context.Context, account string) (uint32, error) {
	client, err := c.getRPCClient()
	if err != nil {
		return 0, err
	}

	result, err := client.GetAccountInfo(&xrplaccount.InfoRequest{
		Account:     types.Address(account),
		LedgerIndex: common.Validated,
		Strict:      true,
	})
	if err != nil {
		return 0, err
	}

	return result.AccountData.Sequence, nil
}

// GetBalance fetches the XRP balance for the given account (drops).
func (c *client) GetBalance(ctx context.Context, account string) (*big.Int, error) {
	client, err := c.getRPCClient()
	if err != nil {
		return nil, err
	}

	result, err := client.GetAccountInfo(&xrplaccount.InfoRequest{
		Account: types.Address(account),
	})
	if err != nil {
		return nil, err
	}

	b := new(big.Int)
	b, ok := b.SetString(result.AccountData.Balance.String(), 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse balance of %s (%s)", account, result.AccountData.Balance.String())
	}

	fmt.Println(b)

	return b, nil
}

// Autofill completes a transaction with missing Sequence, Fee, and LastLedgerSequence fields.
func (c *client) Autofill(tx *transaction.FlatTransaction) error {
	client, err := c.getRPCClient()
	if err != nil {
		return err
	}
	return client.Autofill(tx)
}

// BroadcastTx submits a signed tx blob and returns its hash.
func (c *client) BroadcastTx(ctx context.Context, txBlob string) (TxResult, error) {
	client, err := c.getRPCClient()
	if err != nil {
		return TxResult{}, err
	}

	res, err := client.Request(&requests.SubmitRequest{
		TxBlob:   txBlob,
		FailHard: false,
	})
	if err != nil {
		return TxResult{}, err
	}

	var result requests.SubmitResponse
	if err := res.GetResult(&result); err != nil {
		return TxResult{}, err
	}

	txHash, ok := result.Tx["hash"].(string)
	if !ok || txHash == "" {
		return TxResult{}, fmt.Errorf("missing tx hash in submit response")
	}

	if result.EngineResultCode != 0 {
		return TxResult{
				TxHash: txHash,
			}, fmt.Errorf(
				"failed to broadcast with engine result %s: %s",
				result.EngineResult,
				result.EngineResultMessage,
			)
	}

	fee, ok := result.Tx["Fee"].(string)
	if !ok || fee == "" {
		return TxResult{
			TxHash: txHash,
		}, fmt.Errorf("missing fee in submit response")
	}

	return TxResult{
		TxHash: txHash,
		Fee:    fee,
	}, nil
}

func (c *client) getSelectedEndpoint() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.selectedEndpoint
}

func (c *client) setSelectedEndpoint(endpoint string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.selectedEndpoint = endpoint
}

func (c *client) getRPCClient() (*rpc.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.rpcClient == nil {
		return nil, fmt.Errorf("xrpl rpc client not initialized")
	}
	return c.rpcClient, nil
}

func (c *client) setRPCClient(client *rpc.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rpcClient = client
}
