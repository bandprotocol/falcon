package evm

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

var _ Client = &client{}

// Client is the interface that handles interactions with the EVM chain.
type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	PendingNonceAt(ctx context.Context, address gethcommon.Address) (uint64, error)
	GetBlockHeight(ctx context.Context) (uint64, error)
	GetTxReceipt(ctx context.Context, txHash string) (*gethtypes.Receipt, error)
	Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	EstimateGasPrice(ctx context.Context) (*big.Int, error)
	EstimateBaseFee(ctx context.Context) (*big.Int, error)
	EstimateGasTipCap(ctx context.Context) (*big.Int, error)
	BroadcastTx(ctx context.Context, tx *gethtypes.Transaction) (string, error)
	GetBalance(ctx context.Context, gethAddr gethcommon.Address) (*big.Int, error)
}

// Client is the struct that handles interactions with the EVM chain.
type client struct {
	ChainName      string
	Endpoints      []string
	QueryTimeout   time.Duration
	ExecuteTimeout time.Duration

	Log *zap.Logger

	selectedEndpoint string
	client           *ethclient.Client
}

// NewClient creates a new EVM client from config file and load keys.
func NewClient(chainName string, cfg *EVMChainProviderConfig, log *zap.Logger) *client {
	return &client{
		ChainName:      chainName,
		Endpoints:      cfg.Endpoints,
		QueryTimeout:   cfg.QueryTimeout,
		ExecuteTimeout: cfg.ExecuteTimeout,
		Log:            log.With(zap.String("chain_name", chainName)),
	}
}

// Connect connects to the EVM chain.
func (c *client) Connect(ctx context.Context) error {
	res, err := c.getClientWithMaxHeight(ctx)
	if err != nil {
		c.Log.Error("Failed to connect to EVM chain", zap.Error(err))
		return err
	}

	c.selectedEndpoint = res.Endpoint
	c.client = res.Client
	c.Log.Info("Connected to EVM chain", zap.String("endpoint", c.selectedEndpoint))

	return nil
}

// StartLivelinessCheck starts the liveliness check for the EVM chain.
func (c *client) StartLivelinessCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			c.Log.Info("Stopping liveliness check")

			ticker.Stop()

			return
		case <-ticker.C:
			err := c.Connect(ctx)
			if err != nil {
				c.Log.Error("Liveliness check: unable to reconnect to any endpoints", zap.Error(err))
			}
		}
	}
}

// PendingNonceAt returns the current pending nonce of the given address.
func (c *client) PendingNonceAt(ctx context.Context, address gethcommon.Address) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	nonce, err := c.client.PendingNonceAt(newCtx, address)
	if err != nil {
		c.Log.Error(
			"Failed to get pending nonce",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("evm_address", address.Hex()),
		)
		return 0, ErrGetNonce(err)
	}

	return nonce, nil
}

// GetBlockHeight returns the current block height of the EVM chain on the selected endpoint.
func (c *client) GetBlockHeight(ctx context.Context) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	blockHeight, err := c.client.BlockNumber(newCtx)
	if err != nil {
		c.Log.Error("Failed to get block height", zap.Error(err), zap.String("endpoint", c.selectedEndpoint))
		return 0, ErrGetBlockHeight(err)
	}

	return blockHeight, nil
}

// GetTxReceipt returns the transaction receipt of the given transaction hash.
func (c *client) GetTxReceipt(ctx context.Context, txHash string) (*gethtypes.Receipt, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	receipt, err := c.client.TransactionReceipt(newCtx, gethcommon.HexToHash(txHash))
	if err != nil {
		// tend to be debug log, as it's normal to not have receipt for pending tx
		c.Log.Debug(
			"Failed to get tx receipt",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("tx_hash", txHash),
		)
		return nil, ErrGetTxReceipt(err)
	}

	return receipt, nil
}

// GetTxByHash returns the transaction of the given transaction hash.
func (c *client) GetTxByHash(ctx context.Context, txHash string) (*gethtypes.Transaction, bool, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	tx, isPending, err := c.client.TransactionByHash(newCtx, gethcommon.HexToHash(txHash))
	if err != nil {
		c.Log.Error(
			"Failed to get tx by hash",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("tx_hash", txHash),
		)
		return nil, false, ErrGetTxByHash(err)
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

	res, err := c.client.CallContract(newCtx, callMsg, nil)
	if err != nil {
		c.Log.Error(
			"Failed to query contract",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("evm_address", gethAddr.Hex()),
		)
		return nil, ErrQuery(err)
	}

	return res, nil
}

// EstimateGas estimates the gas amount being used to submit the given message.
func (c *client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	gas, err := c.client.EstimateGas(newCtx, msg)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("evm_address", msg.To.Hex()),
		)
		return 0, ErrEstimateGas(err, true)
	}

	// NOTE: Add 20% buffer to the estimated gas.
	gas = gas * 120 / 100

	return gas, nil
}

// EstimateGasPrice estimates the current gas price on the EVM chain.
func (c *client) EstimateGasPrice(ctx context.Context) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	gasPrice, err := c.client.SuggestGasPrice(newCtx)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas price",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
		)
		return nil, err
	}

	return gasPrice, nil
}

// EstimateBaseFee estimates the current base fee on the EVM chain.
func (c *client) EstimateBaseFee(ctx context.Context) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	latestHeader, err := c.client.HeaderByNumber(newCtx, nil)
	if err != nil {
		return nil, err
	}

	return latestHeader.BaseFee, nil
}

// EstimateGasTipCap estimates the current gas tip cap on the EVM chain.
func (c *client) EstimateGasTipCap(ctx context.Context) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	gasTipCap, err := c.client.SuggestGasTipCap(newCtx)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas tip cap",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
		)
		return nil, err
	}

	return gasTipCap, nil
}

// BroadcastTx sends the transaction to the EVM chain.
func (c *client) BroadcastTx(ctx context.Context, tx *gethtypes.Transaction) (string, error) {
	c.Log.Debug(
		"Broadcasting tx",
		zap.String("endpoint", c.selectedEndpoint),
		zap.String("tx_hash", tx.Hash().Hex()),
		zap.String("to", tx.To().Hex()),
		zap.String("gas_fee_cap", tx.GasFeeCap().String()),
		zap.String("gas_price", tx.GasPrice().String()),
		zap.String("gas_tip_cap", tx.GasTipCap().String()),
		zap.Uint64("gas", tx.Gas()),
	)

	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	if err := c.client.SendTransaction(newCtx, tx); err != nil {
		c.Log.Error(
			"Failed to broadcast tx",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("tx_hash", tx.Hash().Hex()),
		)

		return "", ErrBroadcastTx(err.Error())
	}

	return tx.Hash().Hex(), nil
}

// getClientWithMaxHeight connects to the endpoint that has the highest block height.
func (c *client) getClientWithMaxHeight(ctx context.Context) (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.Endpoints))

	for _, endpoint := range c.Endpoints {
		go func(endpoint string) {
			client, err := ethclient.Dial(endpoint)
			if err != nil {
				c.Log.Debug(
					"Failed to connect to EVM chain",
					zap.Error(err),
					zap.String("endpoint", endpoint),
				)
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
			defer cancel()

			blockHeight, err := client.BlockNumber(newCtx)
			if err != nil {
				c.Log.Debug(
					"Failed to get block height",
					zap.Error(err),
					zap.String("endpoint", endpoint),
				)
				ch <- ClientConnectionResult{endpoint, client, 0}
				return
			}

			c.Log.Debug(
				"Get height of the given client",
				zap.Error(err),
				zap.String("endpoint", endpoint),
				zap.Uint64("block_number", blockHeight),
			)

			ch <- ClientConnectionResult{endpoint, client, blockHeight}
		}(endpoint)
	}

	var result ClientConnectionResult
	for i := 0; i < len(c.Endpoints); i++ {
		r := <-ch
		if r.Client != nil {
			if r.BlockHeight > result.BlockHeight {
				if result.Client != nil {
					result.Client.Close()
				}
				result = r
			} else {
				r.Client.Close()
			}
		}
	}

	if result.Client == nil {
		return ClientConnectionResult{}, ErrConnectEVMChain
	}

	return result, nil
}

// checkAndConnect checks if the client is connected to the EVM chain, if not connect it.
func (c *client) CheckAndConnect(ctx context.Context) error {
	if c.client != nil {
		return nil
	}

	return c.Connect(ctx)
}

// GetBalance get the balance of specific account the EVM chain.
func (c *client) GetBalance(ctx context.Context, gethAddr gethcommon.Address) (*big.Int, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	res, err := c.client.BalanceAt(newCtx, gethAddr, nil)
	if err != nil {
		c.Log.Error(
			"Failed to query balance",
			zap.Error(err),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("evm_address", gethAddr.Hex()),
		)
		return nil, ErrQueryBalance(err)
	}

	return res, nil
}
