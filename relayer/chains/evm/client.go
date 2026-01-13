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

type EVMClients struct {
	mu      sync.RWMutex
	clients map[string]*ethclient.Client
}

func NewEVMClients() EVMClients {
	return EVMClients{
		clients: make(map[string]*ethclient.Client),
	}
}

func (ec *EVMClients) GetClient(endpoint string) (*ethclient.Client, bool) {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	client, exists := ec.clients[endpoint]
	return client, exists
}

func (ec *EVMClients) SetClient(endpoint string, client *ethclient.Client) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.clients[endpoint] = client
}

var _ Client = &client{}

// Client is the interface that handles interactions with the EVM chain.
type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	NonceAt(ctx context.Context, address gethcommon.Address) (uint64, error)
	GetBlockHeight(ctx context.Context) (uint64, error)
	GetBlock(ctx context.Context, height *big.Int) (*gethtypes.Block, error)
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

	selectedEndpoint string
	selectedClient   *ethclient.Client
	clients          EVMClients
	alert            alert.Alert
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
	for idx, endpoint := range c.Endpoints {
		_, ok := c.clients.GetClient(endpoint)
		if ok {
			continue
		}

		wg.Add(1)
		go func(idx int, endpoint string) {
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
		}(idx, endpoint)
	}

	wg.Wait()
	res, err := c.getClientWithMaxHeight(ctx)

	if err != nil {
		c.selectedEndpoint = ""
		c.selectedClient = nil
		c.Log.Error("Failed to connect to EVM chain", err)
		return err
	}

	// only log when new endpoint is used
	if c.selectedEndpoint != res.Endpoint {
		c.Log.Info("Connected to EVM chain", "endpoint", res.Endpoint)
	}

	c.selectedEndpoint = res.Endpoint
	c.selectedClient = res.Client

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

	nonce, err := c.selectedClient.NonceAt(newCtx, address, nil)
	if err != nil {
		c.Log.Error(
			"Failed to get nonce",
			"endpoint", c.selectedEndpoint,
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

	blockHeight, err := c.selectedClient.BlockNumber(newCtx)
	if err != nil {
		c.Log.Error("Failed to get block height", "endpoint", c.selectedEndpoint, err)
		return 0, fmt.Errorf("[EVMClient] failed to get block height: %w", err)
	}

	return blockHeight, nil
}

// GetBlock returns the blocks of the given block height
func (c *client) GetBlock(ctx context.Context, height *big.Int) (*gethtypes.Block, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	block, err := c.selectedClient.BlockByNumber(newCtx, height)
	if err != nil {
		c.Log.Error(
			"Failed to get block by height",
			"endpoint", c.selectedEndpoint,
			"height", height.String(),
			err,
		)
		return nil, fmt.Errorf("[EVMClient] failed to get block by height: %w", err)
	}

	return block, nil
}

// GetTxReceipt returns the transaction receipt of the given transaction hash.
func (c *client) GetTxReceipt(ctx context.Context, txHash string) (*TxReceipt, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	var receipt *TxReceipt
	err := c.selectedClient.Client().CallContext(newCtx, &receipt, "eth_getTransactionReceipt", txHash)
	if err == nil && receipt == nil {
		// it's normal to not have receipt for pending tx
		err = ethereum.NotFound
	}

	if err != nil {
		c.Log.Debug(
			"Failed to get tx receipt",
			"endpoint", c.selectedEndpoint,
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

	tx, isPending, err := c.selectedClient.TransactionByHash(newCtx, gethcommon.HexToHash(txHash))
	if err != nil {
		c.Log.Error(
			"Failed to get tx by hash",
			"endpoint", c.selectedEndpoint,
			"tx_hash", txHash,
			err,
		)
		return nil, false, fmt.Errorf("[EVMClient] failed to get tx by hash: %w", err)
	}

	return tx, isPending, nil
}

// Query queries the EVM chain, if never connected before, it will try to connect to the available one.
func (c *client) Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error) {
	callMsg := ethereum.CallMsg{
		To:   &gethAddr,
		Data: data,
	}

	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	res, err := c.selectedClient.CallContract(newCtx, callMsg, nil)
	if err != nil {
		c.Log.Error(
			"Failed to query contract",
			"endpoint", c.selectedEndpoint,
			"evm_address", gethAddr.Hex(),
			err,
		)
		return nil, fmt.Errorf("[EVMClient] failed to query: %w", err)
	}

	return res, nil
}

// EstimateGas estimates the gas amount being used to submit the given message.
func (c *client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	gas, err := c.selectedClient.EstimateGas(newCtx, msg)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas",
			"endpoint", c.selectedEndpoint,
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

	gasPrice, err := c.selectedClient.SuggestGasPrice(newCtx)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas price",
			"endpoint", c.selectedEndpoint,
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

	latestHeader, err := c.selectedClient.HeaderByNumber(newCtx, nil)
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

	gasTipCap, err := c.selectedClient.SuggestGasTipCap(newCtx)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas tip cap",
			"endpoint", c.selectedEndpoint,
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
		"endpoint", c.selectedEndpoint,
		"tx_hash", tx.Hash().Hex(),
		"to", tx.To().Hex(),
		"gas_fee_cap", tx.GasFeeCap().String(),
		"gas_price", tx.GasPrice().String(),
		"gas_tip_cap", tx.GasTipCap().String(),
		"gas", tx.Gas(),
	)

	newCtx, cancel := context.WithTimeout(ctx, c.ExecuteTimeout)
	defer cancel()

	if err := c.selectedClient.SendTransaction(newCtx, tx); err != nil {
		c.Log.Error(
			"Failed to broadcast tx",
			"endpoint", c.selectedEndpoint,
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

	for idx, endpoint := range c.Endpoints {
		go func(idx int, endpoint string) {
			client, ok := c.clients.GetClient(endpoint)

			if !ok {
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
			defer cancel()

			syncProgress, err := client.SyncProgress(newCtx)
			if err != nil {
				c.Log.Warn(
					"Failed to get sync progress",
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

			if syncProgress != nil {
				c.Log.Warn(
					"Skipping client because it is not fully synced",
					"endpoint", endpoint,
				)
				ch <- ClientConnectionResult{endpoint, nil, 0}
				alert.HandleAlert(
					c.alert,
					alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
						WithChainName(c.ChainName).
						WithEndpoint(endpoint),
					"Skipping client because it is not fully synced",
				)
				return
			}

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
		}(idx, endpoint)
	}

	var result ClientConnectionResult
	for i := 0; i < len(c.Endpoints); i++ {
		r := <-ch
		if r.Client != nil {
			if r.BlockHeight > result.BlockHeight || (r.Endpoint == c.selectedEndpoint && r.BlockHeight == result.BlockHeight) {
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
	if c.selectedClient != nil {
		return nil
	}

	return c.Connect(ctx)
}

// GetBalance get the balance of specific account the EVM chain.
func (c *client) GetBalance(ctx context.Context, gethAddr gethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	res, err := c.selectedClient.BalanceAt(newCtx, gethAddr, blockNumber)
	if err != nil {
		c.Log.Error(
			"Failed to query balance",
			"endpoint", c.selectedEndpoint,
			"evm_address", gethAddr.Hex(),
			err,
		)
		return nil, fmt.Errorf("[EVMClient] failed to query balance: %w", err)
	}

	return res, nil
}
