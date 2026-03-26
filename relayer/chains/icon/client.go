package icon

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/logger"
	iconclient "github.com/icon-project/goloop/client"
	"github.com/icon-project/goloop/server/jsonrpc"
	v3 "github.com/icon-project/goloop/server/v3"
)

// IconClients holds Icon RPC clients and the selected endpoint.
type IconClients = chains.ClientPool[iconclient.ClientV3]

// NewIconClients creates and returns a new IconClients instance with no endpoints.
func NewIconClients() IconClients {
	return chains.NewClientPool[iconclient.ClientV3]()
}

var _ Client = (*client)(nil)

type Client interface {
	Connect() error
	CheckAndConnect() error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	BroadcastTx(txParams v3.TransactionParam) (string, error)
	GetTx(txHash string) (*iconclient.TransactionResult, error)
	GetContractPrices(contract string, symbols []string) ([]ContractOutput, error)
	GetBalance(account string) (*big.Int, error)
}

// Client is the struct that handles interactions with the Icon chain.
type client struct {
	ChainName string
	Endpoints []string

	Log logger.Logger

	clients IconClients
	alert   alert.Alert
}

// NewClient creates a new Icon client from config file and load keys.
func NewClient(chainName string, cfg *IconChainProviderConfig, log logger.Logger, alert alert.Alert) *client {
	return &client{
		ChainName: chainName,
		Endpoints: cfg.Endpoints,
		Log:       log.With("chain_name", chainName),
		alert:     alert,
		clients:   NewIconClients(),
	}
}

// Connect connects to the EVM chain.
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
			client := iconclient.NewClientV3(endpoint)
			c.clients.SetClient(endpoint, client)
		}(endpoint)
	}

	wg.Wait()
	res, err := c.getClientWithMaxHeight()
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
			err := c.Connect()
			if err != nil {
				c.Log.Error("Liveliness check: unable to reconnect to any endpoints", err)
			}
		}
	}
}

// getClientWithMaxHeight connects to the endpoint that has the highest block height.
func (c *client) getClientWithMaxHeight() (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.Endpoints))

	for _, endpoint := range c.Endpoints {
		go func(endpoint string) {
			client, ok := c.clients.GetClient(endpoint)

			if !ok {
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			block, err := client.GetLastBlock()
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
				"block_number", block.Height,
			)
			alert.HandleReset(
				c.alert,
				alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
					WithChainName(c.ChainName).
					WithEndpoint(endpoint),
			)

			ch <- ClientConnectionResult{endpoint, client, uint64(block.Height)}
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
			fmt.Sprintf("failed to connect to icon chain on all endpoints: %s", c.Endpoints),
		)
		return ClientConnectionResult{}, fmt.Errorf("[IconClient] failed to connect to icon chain")
	}

	alert.HandleReset(c.alert, alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName))

	return result, nil
}

// checkAndConnect checks if the client is connected to the EVM chain, if not connect it.
func (c *client) CheckAndConnect() error {
	if _, err := c.clients.GetSelectedClient(); err != nil {
		return c.Connect()
	}

	return nil
}

func (c *client) BroadcastTx(txParams v3.TransactionParam) (string, error) {
	c.Log.Debug(
		"Broadcasting tx",
		"endpoint", c.clients.GetSelectedEndpoint(),
		"tx_params", txParams,
	)

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return "", fmt.Errorf("[IconClient] failed to get client: %w", err)
	}

	var result jsonrpc.HexBytes
	if _, err := client.Do("icx_sendTransaction", txParams, &result); err != nil {
		c.Log.Error(
			"Failed to broadcast tx",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"tx_hash", string(result),
			err,
		)

		return "", fmt.Errorf("[IconClient] failed to broadcast tx with error %s", err.Error())
	}

	return string(result), nil
}

func (c *client) GetTx(txHash string) (*iconclient.TransactionResult, error) {
	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[IconClient] failed to get client: %w", err)
	}

	txResult, err := client.GetTransactionResult(&v3.TransactionHashParam{Hash: jsonrpc.HexBytes(txHash)})
	if err != nil {
		c.Log.Error("Failed to get transaction result", "endpoint", c.clients.GetSelectedEndpoint(), "txHash", txHash, err)
		return nil, fmt.Errorf("[IconClient] failed to get transaction result: %w", err)
	}
	return txResult, nil

}

// GetContractPrices retrieves contract prices for the given symbols.
func (c *client) GetContractPrices(contract string, symbols []string) ([]ContractOutput, error) {
	message := v3.CallParam{
		ToAddress: jsonrpc.Address(contract),
		DataType:  "call",
		Data: ContractData{
			Method: "getReferenceDataBulk",
			Params: GetRefDataBulkParams{
				Bases:  symbols,
				Quotes: make([]string, len(symbols)), // Fill with "USD"
			},
		},
	}

	for i := range message.Data.(ContractData).Params.(GetRefDataBulkParams).Quotes {
		message.Data.(ContractData).Params.(GetRefDataBulkParams).Quotes[i] = "USD"
	}

	latestPrices := make([]ContractOutput, len(symbols))
	success := false

	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[IconClient] failed to get client: %w", err)
	}

	var contractPrices []ContractOutput
	if _, err := client.Do("icx_call", message, &contractPrices); err != nil {
		c.Log.Error("Failed to get contract prices", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[IconClient] failed to get contract prices: %w", err)
	}

	if len(contractPrices) != len(latestPrices) {
		c.Log.Error("Mismatch in symbol count", "expected", len(latestPrices), "got", len(contractPrices))
		return nil, fmt.Errorf("expected %d symbols, got %d", len(latestPrices), len(contractPrices))
	}

	for i := range contractPrices {
		if latestPrices[i].LastUpdateBase.Value() < contractPrices[i].LastUpdateBase.Value() {
			latestPrices[i] = contractPrices[i]
		}
	}
	success = true

	if !success {
		return nil, fmt.Errorf("cannot get contract price from endpoint")
	}

	return latestPrices, nil
}

func (c *client) GetBalance(account string) (*big.Int, error) {
	client, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("[EVMClient] failed to get client: %w", err)
	}

	res, err := client.GetBalance(&v3.AddressParam{Address: jsonrpc.Address(account)})
	if err != nil {
		c.Log.Error(
			"Failed to query balance",
			"endpoint", c.clients.GetSelectedEndpoint(),
			"account", account,
			err,
		)
		return nil, fmt.Errorf("[IconClient] failed to query balance: %w", err)
	}

	return res.BigInt()
}
