package xrpl

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	xrplaccount "github.com/Peersyst/xrpl-go/xrpl/queries/account"
	"github.com/Peersyst/xrpl-go/xrpl/queries/common"
	"github.com/Peersyst/xrpl-go/xrpl/queries/utility"
	"github.com/Peersyst/xrpl-go/xrpl/rpc"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/logger"
)

// Client handles XRPL JSON-RPC interactions.
type Client struct {
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

// NewClient creates a new XRPL client from config.
func NewClient(chainName string, cfg *XRPLChainProviderConfig, log logger.Logger, alert alert.Alert) *Client {
	return &Client{
		ChainName:      chainName,
		Endpoints:      cfg.Endpoints,
		QueryTimeout:   cfg.QueryTimeout,
		ExecuteTimeout: cfg.ExecuteTimeout,
		Log:            log.With("chain_name", chainName),
		alert:          alert,
	}
}

func (c *Client) getSelectedEndpoint() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.selectedEndpoint
}

func (c *Client) setSelectedEndpoint(endpoint string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.selectedEndpoint = endpoint
}

func (c *Client) getRPCClient() (*rpc.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.rpcClient == nil {
		return nil, fmt.Errorf("xrpl rpc client not initialized")
	}
	return c.rpcClient, nil
}

func (c *Client) setRPCClient(client *rpc.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rpcClient = client
}

// Connect selects a responsive endpoint by pinging the server.
func (c *Client) Connect(ctx context.Context) error {
	return c.ping(ctx)
}

// Ping checks connectivity to the XRPL endpoint.
func (c *Client) ping(ctx context.Context) error {
	endpoints := make([]string, 0, len(c.Endpoints))
	if selected := c.getSelectedEndpoint(); selected != "" {
		endpoints = append(endpoints, selected)
	}
	for _, endpoint := range c.Endpoints {
		if endpoint == c.getSelectedEndpoint() {
			continue
		}
		endpoints = append(endpoints, endpoint)
	}

	var lastErr error
	for _, endpoint := range endpoints {
		if err := ctx.Err(); err != nil {
			return err
		}

		timeout := c.QueryTimeout
		if c.ExecuteTimeout > timeout {
			timeout = c.ExecuteTimeout
		}
		var opts []rpc.ConfigOpt
		if timeout > 0 {
			opts = append(opts, rpc.WithTimeout(timeout))
		}
		cfg, err := rpc.NewClientConfig(endpoint, opts...)
		if err != nil {
			lastErr = err
			c.Log.Warn("XRPL endpoint error", "endpoint", endpoint, err)
			continue
		}

		client := rpc.NewClient(cfg)
		_, err = client.Ping(&utility.PingRequest{})
		if err == nil {
			if c.getSelectedEndpoint() != endpoint {
				c.Log.Info("Connected to XRPL endpoint", "endpoint", endpoint)
			}
			c.setSelectedEndpoint(endpoint)
			c.setRPCClient(client)
			return nil
		}

		lastErr = err
		c.Log.Warn("XRPL endpoint error", "endpoint", endpoint, err)
	}

	return lastErr
}

// GetAccountSequenceNumber fetches the sequence for the given account.
func (c *Client) GetAccountSequenceNumber(ctx context.Context, account string) (uint64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
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

	return uint64(result.AccountData.Sequence), nil
}

// GetBalance fetches the XRP balance for the given account (drops).
func (c *Client) GetBalance(ctx context.Context, account string) (*big.Int, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	client, err := c.getRPCClient()
	if err != nil {
		return nil, err
	}
	result, err := client.GetAccountInfo(&xrplaccount.InfoRequest{
		Account:     types.Address(account),
		LedgerIndex: common.Validated,
		Strict:      true,
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
func (c *Client) Autofill(tx *transaction.FlatTransaction) error {
	client, err := c.getRPCClient()
	if err != nil {
		return err
	}
	return client.Autofill(tx)
}

// BroadcastTx submits a signed tx blob and returns its hash.
func (c *Client) BroadcastTx(ctx context.Context, txBlob string) (string, error) {
	client, err := c.getRPCClient()
	if err != nil {
		return "", err
	}

	result, err := client.SubmitTxBlob(txBlob, false)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(result.EngineResult, "tes") {
		return "", fmt.Errorf(
			"failed to broadcast with engine result %s: %s",
			result.EngineResult,
			result.EngineResultMessage,
		)
	}

	txHash, ok := result.Tx["hash"].(string)
	if !ok || txHash == "" {
		return "", fmt.Errorf("missing tx hash in submit response")
	}

	return txHash, nil
}
