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
	executeFn := func(c *ethclient.Client) (any, error) { return nil, nil }
	if _, err := tryExecute(c, executeFn); err != nil {
		c.Log.Error("Failed to connect to EVM chain", zap.Error(err))
		return fmt.Errorf("failed to connect to EVM chain")
	}

	return nil
}

// Query queries the EVM chain, if never connected before, it will try to connect to the available one.
func (c *client) Query(ctx context.Context, gethAddr gethcommon.Address, data []byte) ([]byte, error) {
	callMsg := ethereum.CallMsg{
		To:   &gethAddr,
		Data: data,
	}

	executeFn := func(c *ethclient.Client) ([]byte, error) {
		return c.CallContract(ctx, callMsg, nil)
	}

	res, err := tryExecute(c, executeFn)
	if err != nil {
		c.Log.Error("Failed to query EVM chain", zap.Error(err))
		return nil, err
	}

	return res, nil
}

// tryExecute tries to execute the given function with the client identified in Client object.
// If fail to execute due to the rpc connection, it will try to connect to another rpc endpoint
// until it succeeds or every endpoint is selected.
//
// Note: cannot use generic type with a method of non-generic object.
func tryExecute[T any](c *client, executeFn func(client *ethclient.Client) (T, error)) (T, error) {
	var res T

	// try to execute with the current client
	if c.Client != nil {
		if res, err := executeFn(c.Client); err == nil {
			return res, nil
		}
	}

	// if not success, try to execute with another client
	selectedEndpoint := c.SelectedEndpoint
	c.SelectedEndpoint = ""
	for _, endpoint := range c.RpcEndpoints {
		if endpoint == selectedEndpoint {
			continue
		}

		client, err := ethclient.Dial(endpoint)
		if err != nil {
			c.Log.Error("Failed to connect to EVM chain", zap.Error(err))
		}

		res, err = executeFn(client)
		// TODO: check if the error is due to the contract or not, if so, return with the error
		// else continue to the next endpoint.
		if err != nil {
			c.Log.Error("Failed to execute function", zap.Error(err))
			continue
		}

		c.Client = client
		c.SelectedEndpoint = endpoint
		c.Log.Info("Connected to EVM chain", zap.String("endpoint", endpoint))
		return res, nil
	}

	return res, fmt.Errorf("failed to execute function")
}
