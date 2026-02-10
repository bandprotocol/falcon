package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/logger"
)

// EVMClients holds Ethereum RPC clients and the selected endpoint.
type EVMClients struct {
	mu               sync.RWMutex
	selectedEndpoint string                       // Currently selected endpoint
	clients          map[string]*ethclient.Client // Endpoint to client map
}

// NewEVMClients creates and returns a new EVMClients instance with no endpoints.
func NewEVMClients() EVMClients {
	return EVMClients{
		clients: make(map[string]*ethclient.Client),
	}
}

// GetClient returns the ethclient.Client for a given endpoint, and a boolean indicating if it exists.
func (ec *EVMClients) GetClient(endpoint string) (*ethclient.Client, bool) {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	client, exists := ec.clients[endpoint]
	return client, exists
}

// SetClient sets the ethclient.Client for a given endpoint in the clients map.
func (ec *EVMClients) SetClient(endpoint string, client *ethclient.Client) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.clients[endpoint] = client
}

// SetSelectedEndpoint sets the currently selected endpoint.
func (ec *EVMClients) SetSelectedEndpoint(endpoint string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.selectedEndpoint = endpoint
}

// GetSelectedEndpoint returns the currently selected endpoint.
func (ec *EVMClients) GetSelectedEndpoint() string {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	return ec.selectedEndpoint
}

// GetSelectedClient returns the ethclient.Client for the selected endpoint.
// Returns an error if no endpoint is selected or if the selected client does not exist.
func (ec *EVMClients) GetSelectedClient() (*ethclient.Client, error) {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	if ec.selectedEndpoint == "" {
		return nil, fmt.Errorf("no selected endpoint")
	}

	selectedClient, exists := ec.clients[ec.selectedEndpoint]
	if !exists {
		return nil, fmt.Errorf("selected endpoint client not found: %s", ec.selectedEndpoint)
	}

	return selectedClient, nil
}

var _ Client = &client{}

// Client is the interface that handles interactions with the EVM chain.
type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	NonceAt(ctx context.Context, address gethcommon.Address) (uint64, error)
	GetBlockHeight(ctx context.Context) (uint64, error)
	GetHeaderBlock(ctx context.Context, height *big.Int) (*gethtypes.Header, error)
	GetTxReceipt(ctx context.Context, txHash string) (*TxReceipt, error)
	Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	EstimateGasPrice(ctx context.Context) (*big.Int, error)
	EstimateBaseFee(ctx context.Context) (*big.Int, error)
	EstimateGasTipCap(ctx context.Context) (*big.Int, error)
	BroadcastTx(ctx context.Context, tx *gethtypes.Transaction) (string, error)
	GetBalance(ctx context.Context, gethAddr gethcommon.Address, blockNumber *big.Int) (*big.Int, error)
}

// Client is the struct that handles interactions with the EVM chain.
type client struct {
	ChainName      string
	Endpoints      []string
	QueryTimeout   time.Duration
	ExecuteTimeout time.Duration

	Log logger.Logger

	clients EVMClients
	alert   alert.Alert
}

// NewClient creates a new EVM client from config file and load keys.
func NewClient(chainName string, cfg *EVMChainProviderConfig, log logger.Logger, alert alert.Alert) *client {
	return &client{
		ChainName:      chainName,
		Endpoints:      cfg.Endpoints,
		QueryTimeout:   cfg.QueryTimeout,
		ExecuteTimeout: cfg.ExecuteTimeout,
		Log:            log.With("chain_name", chainName),
		alert:          alert,
		clients:        NewEVMClients(),
	}
}

// Connect connects to the EVM chain.
func (c *client) Connect(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, endpoint := range c.Endpoints {
		_, ok := c.clients.GetClient(endpoint)
		if ok {
			continue
		}

		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			client, err := ethclient.Dial(endpoint)
			if err != nil {
				c.Log.Warn(
					"Failed to connect to EVM chain",
					"endpoint", endpoint,
					err,
				)
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
			c.clients.SetClient(endpoint, client)
		}(endpoint)
	}

	wg.Wait()
	res, err := c.getClientWithMaxHeight(ctx)
	if err != nil {
		c.Log.Error("Failed to connect to EVM chain", err)
		return err
	}

	// only log when new endpoint is used
	if c.clients.GetSelectedEndpoint() != res.Endpoint {
		c.Log.Info("Connected to EVM chain", "endpoint", res.Endpoint)
	}

	c.clients.SetSelectedEndpoint(res.Endpoint)

	return nil
}

// StartLivelinessCheck starts the liveliness check for the EVM chain.
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

// NonceAt retrieves the current account nonce for the given address
// at the latest known block.
func (c *client) NonceAt(ctx context.Context, address gethcommon.Address) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return 0, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	nonce, err := client.NonceAt(newCtx, address, nil)
	if err != nil {
		c.Log.Error(
			"Failed to get nonce",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"evm_address", address.Hex(),
			err,
		)
		return 0, fmt.Errorf("[EVMClient] failed to get nonce: %w", err)
	}

	return nonce, nil
}

// GetBlockHeight returns the current block height of the EVM chain on the selected endpoint.
func (c *client) GetBlockHeight(ctx context.Context) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return 0, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	blockHeight, err := client.BlockNumber(newCtx)
	if err != nil {
		c.Log.Error("Failed to get block height", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return 0, fmt.Errorf("[EVMClient] failed to get block height: %w", err)
	}

	return blockHeight, nil
}

// GetHeaderBlock returns the block header at the given height.
func (c *client) GetHeaderBlock(ctx context.Context, height *big.Int) (*gethtypes.Header, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	header, err := client.HeaderByNumber(newCtx, height)
	if err != nil {
		c.Log.Error(
			"Failed to get header block by height",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"height", height.String(),
			err,
		)
		return nil, fmt.Errorf("[EVMClient] failed to get header block by height: %w", err)
	}

	return header, nil
}

// GetTxReceipt returns the transaction receipt of the given transaction hash.
func (c *client) GetTxReceipt(ctx context.Context, txHash string) (*TxReceipt, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	var receipt *TxReceipt
	err = client.Client().CallContext(newCtx, &receipt, "eth_getTransactionReceipt", txHash)
	if err == nil && receipt == nil {
		// it's normal to not have receipt for pending tx
		err = ethereum.NotFound
	}

	if err != nil {
		c.Log.Debug(
			"Failed to get tx receipt",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"tx_hash", txHash,
			err,
		)
		return nil, fmt.Errorf("[EVMClient] failed to get tx receipt: %w", err)
	}
	return receipt, nil
}

// GetTxByHash returns the transaction of the given transaction hash.
func (c *client) GetTxByHash(ctx context.Context, txHash string) (*gethtypes.Transaction, bool, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, false, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	tx, isPending, err := client.TransactionByHash(newCtx, gethcommon.HexToHash(txHash))
	if err != nil {
		c.Log.Error(
			"Failed to get tx by hash",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"tx_hash", txHash,
			err,
		)
		return nil, false, fmt.Errorf("[EVMClient] failed to get tx by hash: %w", err)
	}

	return tx, isPending, nil
}

// Query queries the EVM chain, if never connected before, it will try to connect to the available one.
func (c *client) Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	var result string
	err = client.Client().CallContext(newCtx, &result, "eth_call", map[string]interface{}{
		"to":   gethAddr.Hex(),
		"data": fmt.Sprintf("0x%x", data),
	}, "latest")
	if err != nil {
		c.Log.Error(
			"Failed to query contract",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"evm_address", gethAddr.Hex(),
			err,
		)
		return nil, fmt.Errorf("[EVMClient] failed to query: %w", err)
	}

	// Convert hex result to bytes
	if len(result) < 2 || result[:2] != "0x" {
		return nil, fmt.Errorf("[EVMClient] invalid hex result: %s", result)
	}

	return gethcommon.FromHex(result), nil
}

// EstimateGas estimates the gas amount being used to submit the given message.
func (c *client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return 0, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	gas, err := client.EstimateGas(newCtx, msg)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"evm_address", msg.To.Hex(),
			err,
		)
		return 0, fmt.Errorf("[EVMClient] failed to estimate gas: %w", err)
	}

	return gas, nil
}

// EstimateGasPrice estimates the current gas price on the EVM chain.
func (c *client) EstimateGasPrice(ctx context.Context) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(newCtx)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas price",
			"endpoint", c.clients.GetSelectedEndpoint(),
			err,
		)
		return nil, err
	}

	return gasPrice, nil
}

// EstimateBaseFee estimates the current base fee on the EVM chain.
// ref: https://ethereum.stackexchange.com/questions/132333/how-can-we-calculate-next-base-fee
func (c *client) EstimateBaseFee(ctx context.Context) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	latestHeader, err := client.HeaderByNumber(newCtx, nil)
	if err != nil {
		return nil, err
	}

	estimatedBaseFee := MultiplyBigIntWithFloat64(latestHeader.BaseFee, 1.125)
	return estimatedBaseFee, nil
}

// EstimateGasTipCap estimates the current gas tip cap on the EVM chain.
func (c *client) EstimateGasTipCap(ctx context.Context) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	gasTipCap, err := client.SuggestGasTipCap(newCtx)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas tip cap",
			"endpoint", c.clients.GetSelectedEndpoint(),
			err,
		)
		return nil, err
	}

	return gasTipCap, nil
}

// BroadcastTx sends the transaction to the EVM chain.
func (c *client) BroadcastTx(ctx context.Context, tx *gethtypes.Transaction) (string, error) {
	c.Log.Debug(
		"Broadcasting tx",
		"endpoint", c.clients.GetSelectedEndpoint(),
		"tx_hash", tx.Hash().Hex(),
		"to", tx.To().Hex(),
		"gas_fee_cap", tx.GasFeeCap().String(),
		"gas_price", tx.GasPrice().String(),
		"gas_tip_cap", tx.GasTipCap().String(),
		"gas", tx.Gas(),
	)

	newCtx, cancel := context.WithTimeout(ctx, c.ExecuteTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return "", fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	if err := client.SendTransaction(newCtx, tx); err != nil {
		c.Log.Error(
			"Failed to broadcast tx",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"tx_hash", tx.Hash().Hex(),
			err,
		)

		return "", fmt.Errorf("[EVMClient] failed to broadcast tx with error %s", err.Error())
	}

	return tx.Hash().Hex(), nil
}

// getClientWithMaxHeight connects to the endpoint that has the highest block height.
func (c *client) getClientWithMaxHeight(ctx context.Context) (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.Endpoints))

	for _, endpoint := range c.Endpoints {
		go func(endpoint string) {
			client, ok := c.clients.GetClient(endpoint)

			if !ok {
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
			defer cancel()

			blockHeight, err := client.BlockNumber(newCtx)
			if err != nil {
				c.Log.Warn(
					"Failed to get block height",
					"endpoint", endpoint,
					err,
				)
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

			c.Log.Debug(
				"Get height of the given client",
				"endpoint", endpoint,
				"block_number", blockHeight,
			)
			alert.HandleReset(
				c.alert,
				alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
					WithChainName(c.ChainName).
					WithEndpoint(endpoint),
			)

			ch <- ClientConnectionResult{endpoint, client, blockHeight}
		}(endpoint)
	}

	var result ClientConnectionResult
	for i := 0; i < len(c.Endpoints); i++ {
		r := <-ch
		if r.Client != nil {
			if r.BlockHeight > result.BlockHeight || (r.Endpoint == c.clients.GetSelectedEndpoint() && r.BlockHeight == result.BlockHeight) {
				result = r
			}
		}
	}

	if result.Client == nil {
		alert.HandleAlert(
			c.alert,
			alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName),
			fmt.Sprintf("failed to connect to EVM chain on all endpoints: %s", c.Endpoints),
		)
		return ClientConnectionResult{}, fmt.Errorf("[EVMClient] failed to connect to EVM chain")
	}

	alert.HandleReset(c.alert, alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName))

	return result, nil
}

// checkAndConnect checks if the client is connected to the EVM chain, if not connect it.
func (c *client) CheckAndConnect(ctx context.Context) error {
	if _, err := c.clients.GetSelectedClient(); err != nil {
		return c.Connect(ctx)
	}

	return nil
}

// GetBalance get the balance of specific account the EVM chain.
func (c *client) GetBalance(ctx context.Context, gethAddr gethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	res, err := client.BalanceAt(newCtx, gethAddr, blockNumber)
	if err != nil {
		c.Log.Error(
			"Failed to query balance",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"evm_address", gethAddr.Hex(),
			err,
		)
		return nil, fmt.Errorf("[EVMClient] failed to query balance: %w", err)
	}

	return res, nil
}
