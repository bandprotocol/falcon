package evm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

var _ Client = &client{}

// Client is the interface that handles interactions with the EVM chain.
type Client interface {
	Connect() error
	Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error)
}

// Client is the struct that handles interactions with the EVM chain.
type client struct {
	ChainName    string
	RpcEndpoints []string
	Log          *zap.Logger

	SelectedEndpoint string
	Client           *ethclient.Client

	KeyStore *KeyStore
	Keys     []Key
}

// NewClient creates a new EVM client from config file and load keys.
func NewClient(chainName string, rpcEndpoints []string, log *zap.Logger) *client {
	return &client{
		ChainName:    chainName,
		RpcEndpoints: rpcEndpoints,
		Log:          log,
	}
}

// Connect connects to the EVM chain.
func (c *client) Connect() error {
	res, err := getClientWithMaxHeight(c.RpcEndpoints)
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
	if c.Client == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	callMsg := ethereum.CallMsg{
		To:   &gethAddr,
		Data: data,
	}

	res, err := c.Client.CallContract(ctx, callMsg, nil)
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
func getClientWithMaxHeight(rpcEndpoints []string) (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(rpcEndpoints))

	for _, endpoint := range rpcEndpoints {
		go func(endpoint string) {
			client, err := ethclient.Dial(endpoint)
			if err != nil {
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			block, err := client.BlockByNumber(context.Background(), nil)
			if err != nil {
				ch <- ClientConnectionResult{endpoint, client, 0}
				return
			}

			ch <- ClientConnectionResult{endpoint, client, block.NumberU64()}
		}(endpoint)
	}

	var result ClientConnectionResult
	for i := 0; i < len(rpcEndpoints); i++ {
		r := <-ch
		if r.Client != nil && r.BlockHeight > result.BlockHeight {
			result = r
			if result.Client != nil {
				result.Client.Close()
			}
		} else if r.Client != nil {
			r.Client.Close()
		}
	}

	if result.Client == nil {
		return ClientConnectionResult{}, fmt.Errorf("failed to connect to EVM chain")
	}

	return result, nil
}
