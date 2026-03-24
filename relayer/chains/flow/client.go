package flow

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/onflow/flow-go-sdk"
	flowhttp "github.com/onflow/flow-go-sdk/access/http"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/logger"
)

// FlowClients holds Flow HTTP clients and the selected endpoint.
type FlowClients = chains.ClientPool[flowhttp.Client]

// NewFlowClients creates and returns a new FlowClients instance with no endpoints.
func NewFlowClients() FlowClients {
	return chains.NewClientPool[flowhttp.Client]()
}

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *flowhttp.Client
	BlockHeight uint64
}

// Client defines the interface for interacting with the Flow blockchain.
type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	GetAccount(ctx context.Context, address string) (*flow.Account, error)
	GetLatestBlockID(ctx context.Context) (string, error)
	BroadcastTx(ctx context.Context, txBlob []byte) (string, error)
	GetTxResult(ctx context.Context, txHash string) (*flow.TransactionResult, error)
	GetBalance(ctx context.Context, address string) (*big.Int, error)
	GetBlockTimestamp(ctx context.Context, txHash string) (*time.Time, error)
}

var _ Client = (*client)(nil)

// client is the concrete implementation that handles Flow HTTP interactions.
type client struct {
	ChainName      string
	Endpoints      []string
	QueryTimeout   time.Duration
	ExecuteTimeout time.Duration

	Log   logger.Logger
	alert alert.Alert

	clients FlowClients
}

// NewClient creates a new Flow client from config.
func NewClient(chainName string, cfg *FlowChainProviderConfig, log logger.Logger, a alert.Alert) Client {
	return &client{
		ChainName:      chainName,
		Endpoints:      cfg.Endpoints,
		QueryTimeout:   cfg.QueryTimeout,
		ExecuteTimeout: cfg.ExecuteTimeout,
		Log:            log.With("chain_name", chainName),
		alert:          a,
		clients:        NewFlowClients(),
	}
}

// Connect connects to all endpoints and selects the one with the highest block height.
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
			fc, err := flowhttp.NewClient(endpoint)
			if err != nil {
				c.Log.Warn("Flow endpoint error", "endpoint", endpoint, err)
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
			c.clients.SetClient(endpoint, fc)
		}(endpoint)
	}

	wg.Wait()
	res, err := c.getClientWithMaxHeight()
	if err != nil {
		c.Log.Error("Failed to connect to Flow chain", err)
		return err
	}

	// only log when new endpoint is used
	if c.clients.GetSelectedEndpoint() != res.Endpoint {
		c.Log.Info("Connected to Flow chain", "endpoint", res.Endpoint)
	}

	c.clients.SetSelectedEndpoint(res.Endpoint)

	return nil
}

// getClientWithMaxHeight selects the endpoint with the highest sealed block height.
func (c *client) getClientWithMaxHeight() (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.Endpoints))

	for _, endpoint := range c.Endpoints {
		go func(endpoint string) {
			fc, ok := c.clients.GetClient(endpoint)
			if !ok {
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			queryCtx, queryCancel := context.WithTimeout(context.Background(), c.QueryTimeout)
			defer queryCancel()
			bh, err := fc.GetLatestBlockHeader(queryCtx, true)
			if err != nil {
				c.Log.Warn("Failed to get block height", "endpoint", endpoint, err)
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

			c.Log.Debug("Get height of the given client", "endpoint", endpoint, "block_height", bh.Height)
			alert.HandleReset(
				c.alert,
				alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
					WithChainName(c.ChainName).
					WithEndpoint(endpoint),
			)

			ch <- ClientConnectionResult{endpoint, fc, bh.Height}
		}(endpoint)
	}

	var result ClientConnectionResult
	for range c.Endpoints {
		r := <-ch
		if r.Client != nil {
			if r.BlockHeight > result.BlockHeight ||
				(r.Endpoint == c.clients.GetSelectedEndpoint() && r.BlockHeight == result.BlockHeight) {
				result = r
			}
		}
	}

	if result.Client == nil {
		alert.HandleAlert(
			c.alert,
			alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName),
			fmt.Sprintf("failed to connect to Flow chain on all endpoints: %s", c.Endpoints),
		)
		return ClientConnectionResult{}, fmt.Errorf("[FlowClient] failed to connect to any endpoint")
	}

	alert.HandleReset(c.alert, alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName))

	return result, nil
}

// CheckAndConnect connects if not already connected.
func (c *client) CheckAndConnect(ctx context.Context) error {
	if _, err := c.clients.GetSelectedClient(); err != nil {
		return c.Connect(ctx)
	}

	return nil
}

// StartLivelinessCheck periodically reconnects to verify endpoint health.
func (c *client) StartLivelinessCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.Log.Info("Stopping liveliness check")
			return
		case <-ticker.C:
			if err := c.Connect(ctx); err != nil {
				c.Log.Error("Liveliness check: unable to reconnect to any endpoints", err)
			}
		}
	}
}

// GetAccount returns the account for the given address.
func (c *client) GetAccount(ctx context.Context, address string) (*flow.Account, error) {
	fc, err := c.clients.GetSelectedClient()
	if err != nil {
		return nil, fmt.Errorf("[FlowClient] failed to get client: %w", err)
	}

	queryCtx, queryCancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer queryCancel()
	acc, err := fc.GetAccount(queryCtx, flow.HexToAddress(address))
	if err != nil {
		return nil, fmt.Errorf("[FlowClient] failed to get account %s: %w", address, err)
	}

	return acc, nil
}

// GetLatestBlockID returns the hex-encoded ID of the latest sealed block.
func (c *client) GetLatestBlockID(ctx context.Context) (string, error) {
	fc, err := c.clients.GetSelectedClient()
	if err != nil {
		return "", fmt.Errorf("[FlowClient] failed to get client: %w", err)
	}

	queryCtx, queryCancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer queryCancel()
	bh, err := fc.GetLatestBlockHeader(queryCtx, true)
	if err != nil {
		return "", fmt.Errorf("[FlowClient] failed to get latest block header: %w", err)
	}

	return bh.ID.Hex(), nil
}

// BroadcastTx decodes a flow.Transaction and broadcasts it to the network.
func (c *client) BroadcastTx(ctx context.Context, txBlob []byte) (string, error) {
	fc, err := c.clients.GetSelectedClient()
	if err != nil {
		return "", fmt.Errorf("[FlowClient] failed to get client: %w", err)
	}

	tx, err := flow.DecodeTransaction(txBlob)
	if err != nil {
		return "", fmt.Errorf("[FlowClient] failed to decode tx blob: %w", err)
	}

	execCtx, execCancel := context.WithTimeout(ctx, c.ExecuteTimeout)
	defer execCancel()
	if err := fc.SendTransaction(execCtx, *tx); err != nil {
		return "", fmt.Errorf("[FlowClient] failed to broadcast transaction: %w", err)
	}

	return tx.ID().String(), nil
}

// GetTxResult fetches the result of a transaction by its hash.
func (c *client) GetTxResult(ctx context.Context, txHash string) (*flow.TransactionResult, error) {
	fc, err := c.clients.GetSelectedClient()
	if err != nil {
		return nil, fmt.Errorf("[FlowClient] failed to get client: %w", err)
	}

	queryCtx, queryCancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer queryCancel()
	result, err := fc.GetTransactionResult(queryCtx, flow.HexToID(txHash))
	if err != nil {
		return nil, fmt.Errorf("[FlowClient] failed to get transaction result for %s: %w", txHash, err)
	}

	return result, nil
}

// GetBlockTimestamp fetches the block timestamp for the block containing the given transaction.
func (c *client) GetBlockTimestamp(ctx context.Context, txHash string) (*time.Time, error) {
	fc, err := c.clients.GetSelectedClient()
	if err != nil {
		return nil, fmt.Errorf("[FlowClient] failed to get client: %w", err)
	}

	queryCtx, queryCancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer queryCancel()

	result, err := fc.GetTransactionResult(queryCtx, flow.HexToID(txHash))
	if err != nil {
		return nil, fmt.Errorf("[FlowClient] failed to get transaction result for %s: %w", txHash, err)
	}

	bh, err := fc.GetBlockHeaderByID(queryCtx, result.BlockID)
	if err != nil {
		return nil, fmt.Errorf("[FlowClient] failed to get block header for tx %s: %w", txHash, err)
	}

	t := bh.Timestamp
	return &t, nil
}

// GetBalance returns the FLOW token balance of the given address in UFix64 units (as big.Int).
func (c *client) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	acc, err := c.GetAccount(ctx, address)
	if err != nil {
		return nil, err
	}

	return big.NewInt(int64(acc.Balance)), nil
}
