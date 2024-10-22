package evm

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

var _ Client = &client{}

// Client is the interface that handles interactions with the EVM chain.
type Client interface {
	Connect(ctx context.Context) error
	Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error)
}

// Client is the struct that handles interactions with the EVM chain.
type client struct {
	ChainName    string
	Endpoints    []string
	QueryTimeout time.Duration
	Log          *zap.Logger

	SelectedEndpoint string
	Client           *ethclient.Client

	KeyStore *KeyStore
	Keys     []Key
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
	if c.Client != nil {
		return nil
	}

	res, err := c.getClientWithMaxHeight(ctx)
	if err != nil {
		return err
	}

	c.SelectedEndpoint = res.Endpoint
	c.Client = res.Client
	c.Log.Info("Connected to EVM chain", zap.String("endpoint", c.SelectedEndpoint))
	return nil
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

	res, err := c.Client.CallContract(newCtx, callMsg, nil)
	if err != nil {
		c.Log.Error("Failed to query EVM chain", zap.Error(err))
		return nil, err
	}

	return res, nil
}

// ClientConnectionResult is the struct that contains the result of connecting to the specific endpoint.
type ClientConnectionResult struct {
	Endpoint    string
	Client      *ethclient.Client
	BlockHeight uint64
}

// getClientWithMaxHeight connects to the endpoint that has the highest block height.
func (c *client) getClientWithMaxHeight(ctx context.Context) (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.Endpoints))

	for _, endpoint := range c.Endpoints {
		go func(endpoint string) {
			client, err := ethclient.Dial(endpoint)
			if err != nil {
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			newCtx, cancel := context.WithTimeout(ctx, c.QueryTimeout)
			defer cancel()

			block, err := client.BlockByNumber(newCtx, nil)
			if err != nil {
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
		return ClientConnectionResult{}, fmt.Errorf("failed to connect to EVM chain")
	}

	return result, nil
}

// checkConnection checks if the client is connected to the EVM chain, if not connect it.
func (c *client) checkAndConnect(ctx context.Context) error {
	return c.Connect(ctx)
}
