package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var _ Client = &client{}

// Client is the interface that handles interactions with the EVM chain.
type Client interface {
	Connect(ctx context.Context) error
	GetNonce(ctx context.Context, address gethcommon.Address) (uint64, error)
	GetBlockHeight(ctx context.Context) (uint64, error)
	GetTxReceipt(ctx context.Context, txHash string) (*gethtypes.Receipt, error)
	GetEffectiveGasPrice(
		ctx context.Context,
		receipt *gethtypes.Receipt,
	) (decimal.NullDecimal, error)
	GetEffectiveGasTipValue(
		ctx context.Context,
		receipt *gethtypes.Receipt,
	) (decimal.NullDecimal, error)
	Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	BroadcastTx(ctx context.Context, tx *gethtypes.Transaction) (string, error)
}

// Client is the struct that handles interactions with the EVM chain.
type client struct {
	ChainName    string
	Endpoints    []string
	QueryTimeout time.Duration

	Log *zap.Logger

	selectedEndpoint string
	client           *ethclient.Client
}

// NewClient creates a new EVM client from config file and load keys.
func NewClient(chainName string, cfg *EVMChainProviderConfig, log *zap.Logger) *client {
	return &client{
		ChainName:    chainName,
		Endpoints:    cfg.Endpoints,
		QueryTimeout: cfg.QueryTimeout,
		Log:          log,
	}
}

// Connect connects to the EVM chain.
func (c *client) Connect(ctx context.Context) error {
	if c.client != nil {
		return nil
	}

	res, err := c.getClientWithMaxHeight(ctx)
	if err != nil {
		return err
	}

	c.selectedEndpoint = res.Endpoint
	c.client = res.Client
	c.Log.Info(
		"Connected to EVM chain",
		zap.String("chain_name", c.ChainName),
		zap.String("endpoint", c.selectedEndpoint),
	)
	return nil
}

// GetNonce returns the current nonce of the given address.
func (c *client) GetNonce(ctx context.Context, address gethcommon.Address) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	nonce, err := c.client.NonceAt(newCtx, address, nil)
	if err != nil {
		c.Log.Error(
			"Failed to get nonce",
			zap.Error(err),
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("evm_address", address.Hex()),
		)
		return 0, fmt.Errorf("[EVMClient] failed to get nonce: %w", err)
	}

	return nonce, nil
}

// GetBlockHeight returns the current block height of the EVM chain on the selected endpoint.
func (c *client) GetBlockHeight(ctx context.Context) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	block, err := c.client.BlockByNumber(newCtx, nil)
	if err != nil {
		c.Log.Error(
			"Failed to get block height",
			zap.Error(err),
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
		)
		return 0, fmt.Errorf("[EVMClient] failed to get block height: %w", err)
	}

	return block.NumberU64(), nil
}

// GetBlock returns the block information of the given height.
func (c *client) GetBlock(ctx context.Context, height uint64) (*gethtypes.Block, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	block, err := c.client.BlockByNumber(newCtx, new(big.Int).SetUint64(height))
	if err != nil {
		c.Log.Error(
			"Failed to get block information",
			zap.Error(err),
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
		)
		return nil, fmt.Errorf("[EVMClient] failed to get block information: %w", err)
	}

	return block, nil
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
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("tx_hash", txHash),
		)
		return nil, fmt.Errorf("[EVMClient] failed to get tx receipt: %w", err)
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
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("tx_hash", txHash),
		)
		return nil, false, fmt.Errorf("[EVMClient] failed to get tx by hash: %w", err)
	}

	return tx, isPending, nil
}

// GetEffectiveGasTipValue returns the effective gas tip of the given transaction receipt.
func (c *client) GetEffectiveGasTipValue(
	ctx context.Context,
	receipt *gethtypes.Receipt,
) (decimal.NullDecimal, error) {
	tx, isPending, err := c.GetTxByHash(ctx, receipt.TxHash.String())
	if err != nil {
		return decimal.NullDecimal{}, err
	}

	if isPending {
		c.Log.Debug(
			"tx is pending",
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("tx_hash", receipt.TxHash.String()),
		)
		return decimal.NullDecimal{}, fmt.Errorf("[EVMClient] tx is pending")
	}

	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	header, err := c.GetBlock(newCtx, receipt.BlockNumber.Uint64())
	if err != nil {
		return decimal.NullDecimal{}, err
	}

	baseFee := header.BaseFee()
	priorityFee := tx.EffectiveGasTipValue(baseFee)
	return decimal.NewNullDecimal(
		decimal.New(new(big.Int).Add(baseFee, priorityFee).Int64(), 0),
	), nil
}

// GetEffectiveGasPrice returns the effective gas price of the given transaction receipt.
func (c *client) GetEffectiveGasPrice(
	ctx context.Context,
	receipt *gethtypes.Receipt,
) (decimal.NullDecimal, error) {
	tx, isPending, err := c.GetTxByHash(ctx, receipt.TxHash.String())
	if err != nil {
		return decimal.NullDecimal{}, err
	}

	if isPending {
		c.Log.Debug(
			"tx is pending",
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("tx_hash", receipt.TxHash.String()),
		)
		return decimal.NullDecimal{}, fmt.Errorf("[EVMClient] tx is pending")
	}

	return decimal.NewNullDecimal(
		decimal.New(int64(tx.GasPrice().Uint64()), 0),
	), nil
}

// Query queries the EVM chain, if never connected before, it will try to connect to the available one.
func (c *client) Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error) {
	if err := c.checkAndConnect(ctx); err != nil {
		return nil, err
	}

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
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("evm_address", gethAddr.Hex()),
		)
		return nil, fmt.Errorf("[EVMClient] failed to query: %w", err)
	}

	return res, nil
}

// EstimateGas estimates the gas of the given message.
func (c *client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
	defer cancel()

	gas, err := c.client.EstimateGas(newCtx, msg)
	if err != nil {
		c.Log.Error(
			"Failed to estimate gas",
			zap.Error(err),
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("evm_address", msg.To.Hex()),
		)
		return 0, fmt.Errorf("[EVMClient] failed to estimate gas: %w", err)
	}

	return gas, nil
}

// BroadcastTx sends the transaction to the EVM chain.
func (c *client) BroadcastTx(ctx context.Context, tx *gethtypes.Transaction) (string, error) {
	c.Log.Debug(
		"Broadcasting tx",
		zap.String("chain_name", c.ChainName),
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
			zap.String("chain_name", c.ChainName),
			zap.String("endpoint", c.selectedEndpoint),
			zap.String("tx_hash", tx.Hash().Hex()),
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
			client, err := ethclient.Dial(endpoint)
			if err != nil {
				c.Log.Debug(
					"Failed to connect to EVM chain",
					zap.Error(err),
					zap.String("endpoint", endpoint),
					zap.String("chain_name", c.ChainName),
				)
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
			defer cancel()

			block, err := client.BlockByNumber(newCtx, nil)
			if err != nil {
				c.Log.Debug(
					"Failed to get block height",
					zap.Error(err),
					zap.String("endpoint", endpoint),
					zap.String("chain_name", c.ChainName),
				)
				ch <- ClientConnectionResult{endpoint, client, 0}
				return
			}

			ch <- ClientConnectionResult{endpoint, client, block.NumberU64()}
		}(endpoint)
	}

	var result ClientConnectionResult
	for i := 0; i < len(c.Endpoints); i++ {
		r := <-ch
		if r.Client != nil {
			if r.BlockHeight >= result.BlockHeight {
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
		c.Log.Error(
			"failed to connect to EVM chain",
			zap.String("chain_name", c.ChainName),
		)
		return ClientConnectionResult{}, fmt.Errorf("[EVMClient] failed to connect to EVM chain")
	}

	return result, nil
}

// checkAndConnect checks if the client is connected to the EVM chain, if not connect it.
func (c *client) checkAndConnect(ctx context.Context) error {
	return c.Connect(ctx)
}
