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

// XRPLClients holds XRPL RPC clients and the selected endpoint.
type XRPLClients struct {
	mu               sync.RWMutex
	selectedEndpoint string                 // Currently selected endpoint
	clients          map[string]*rpc.Client // Endpoint to client map
}

// NewXRPLClients creates and returns a new XRPLClients instance with no endpoints.
func NewXRPLClients() XRPLClients {
	return XRPLClients{
		clients: make(map[string]*rpc.Client),
	}
}

// GetClient returns the rpc.Client for a given endpoint, and a boolean indicating if it exists.
func (xc *XRPLClients) GetClient(endpoint string) (*rpc.Client, bool) {
	xc.mu.RLock()
	defer xc.mu.RUnlock()

	client, exists := xc.clients[endpoint]
	return client, exists
}

// SetClient sets the rpc.Client for a given endpoint in the clients map.
func (xc *XRPLClients) SetClient(endpoint string, client *rpc.Client) {
	xc.mu.Lock()
	defer xc.mu.Unlock()

	xc.clients[endpoint] = client
}

// SetSelectedEndpoint sets the currently selected endpoint.
func (xc *XRPLClients) SetSelectedEndpoint(endpoint string) {
	xc.mu.Lock()
	defer xc.mu.Unlock()

	xc.selectedEndpoint = endpoint
}

// GetSelectedEndpoint returns the currently selected endpoint.
func (xc *XRPLClients) GetSelectedEndpoint() string {
	xc.mu.RLock()
	defer xc.mu.RUnlock()

	return xc.selectedEndpoint
}

// GetSelectedClient returns the rpc.Client for the selected endpoint.
// Returns an error if no endpoint is selected or if the selected client does not exist.
func (xc *XRPLClients) GetSelectedClient() (*rpc.Client, error) {
	xc.mu.RLock()
	defer xc.mu.RUnlock()

	if xc.selectedEndpoint == "" {
		return nil, fmt.Errorf("no selected endpoint")
	}

	selectedClient, exists := xc.clients[xc.selectedEndpoint]
	if !exists {
		return nil, fmt.Errorf("selected endpoint client not found: %s", xc.selectedEndpoint)
	}

	return selectedClient, nil
}

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

	clients XRPLClients
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
		clients:        NewXRPLClients(),
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
	var wg sync.WaitGroup
	for _, endpoint := range c.Endpoints {
		_, ok := c.clients.GetClient(endpoint)
		if ok {
			continue
		}

		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			timeout := c.QueryTimeout
			opts := []rpc.ConfigOpt{rpc.WithTimeout(timeout)}
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
		return err
	}

	// only log when new endpoint is used
	if c.clients.GetSelectedEndpoint() != res.Endpoint {
		c.Log.Info("Connected to XRPL chain", "endpoint", res.Endpoint)
	}

	c.clients.SetSelectedEndpoint(res.Endpoint)

	return nil
}

// getClientWithMaxHeight connects to the endpoint that has the highest ledger index.
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

	var result ClientConnectionResult
	for range c.Endpoints {
		r := <-ch
		if r.Client != nil {
			if r.LedgerIndex > result.LedgerIndex || (r.Endpoint == c.clients.GetSelectedEndpoint() && r.LedgerIndex == result.LedgerIndex) {
				result = r
			}
		}
	}

	if result.Client == nil {
		alert.HandleAlert(
			c.alert,
			alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName),
			fmt.Sprintf("failed to connect to XRPL chain on all endpoints: %s", c.Endpoints),
		)
		return ClientConnectionResult{}, fmt.Errorf("[XRPLClient] failed to connect to XRPL chain on all endpoints")
	}

	alert.HandleReset(c.alert, alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName))

	return result, nil
}

// CheckAndConnect checks if the client is connected to the XRPL chain, if not connect it.
func (c *client) CheckAndConnect() error {
	if _, err := c.clients.GetSelectedClient(); err != nil {
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
	client, err := c.clients.GetSelectedClient()
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
	client, err := c.clients.GetSelectedClient()
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

	return b, nil
}

// Autofill completes a transaction with missing Sequence, Fee, and LastLedgerSequence fields.
func (c *client) Autofill(tx *transaction.FlatTransaction) error {
	client, err := c.clients.GetSelectedClient()
	if err != nil {
		return err
	}
	return client.Autofill(tx)
}

// BroadcastTx submits a signed tx blob and returns its hash.
func (c *client) BroadcastTx(ctx context.Context, txBlob string) (TxResult, error) {
	client, err := c.clients.GetSelectedClient()
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

	if result.EngineResult != "tesSUCCESS" {
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
