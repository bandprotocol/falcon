package secret

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	std "github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/logger"
)

type SecretClients = chains.ClientPool[sdkclient.Context]

// Client defines the Cosmos RPC operations required by the Secret chain provider.
type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)

	BroadcastTx(txBlob []byte) (string, error)
	GetBalance(ctx context.Context, address string) (*big.Int, error)
	GetTx(ctx context.Context, txHash string) (*typesTxResult, error)
	GetBlockByHeight(ctx context.Context, height *big.Int) (*typesBlockResult, error)
	GetAccount(ctx context.Context, sender string) (accountNumber uint64, sequence uint64, err error)
}

// typesTxResult mirrors the subset of TxResult fields we need in provider.
type typesTxResult struct {
	StatusCode uint32
	GasUsed    int64
	Height     int64
	Log        string
}

// typesBlockResult mirrors the subset of block fields we need in provider.
type typesBlockResult struct {
	Time time.Time
}

var _ Client = (*client)(nil)

type client struct {
	chainName string
	endpoints []string
	denom     string

	log     logger.Logger
	alert   alert.Alert
	clients SecretClients
}

func NewClient(chainName string, cpc *SecretChainProviderConfig, log logger.Logger, alert alert.Alert) *client {
	return &client{
		chainName: chainName,
		endpoints: cpc.Endpoints,
		denom:     cpc.Denom,
		log:       log.With("chain_name", chainName),
		alert:     alert,
		clients:   chains.NewClientPool[sdkclient.Context](),
	}
}

func makeEncodingConfig() struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          sdkclient.TxConfig
} {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := authtx.NewTxConfig(marshaler, authtx.DefaultSignModes)

	std.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)

	return struct {
		InterfaceRegistry codectypes.InterfaceRegistry
		Marshaler         codec.Codec
		TxConfig          sdkclient.TxConfig
	}{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
	}
}

func newRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = timeout

	// Websocket endpoint is required by cometbft client for some queries.
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}
	return rpcClient, nil
}

func (c *client) Connect(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, endpoint := range c.endpoints {
		_, ok := c.clients.GetClient(endpoint)
		if ok {
			continue
		}

		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()

			enc := makeEncodingConfig()

			// NOTE: band-feeder uses a larger timeout; we keep it conservative here.
			rpcClient, err := newRPCClient(endpoint, 3*time.Second)
			if err != nil {
				c.log.Error("Failed to create cometbft rpc client", "endpoint", endpoint, "err", err)
				return
			}

			ctxCli := sdkclient.Context{}.
				WithAccountRetriever(authtypes.AccountRetriever{}).
				WithBroadcastMode(flags.BroadcastSync).
				WithCodec(enc.Marshaler).
				WithInterfaceRegistry(enc.InterfaceRegistry).
				WithTxConfig(enc.TxConfig).
				WithClient(rpcClient).
				WithNodeURI(endpoint)

			c.clients.SetClient(endpoint, &ctxCli)
		}(endpoint)
	}

	wg.Wait()

	res, err := c.getClientWithMaxHeight(ctx)
	if err != nil {
		c.log.Error("Failed to connect to secret chain", err)
		return err
	}

	// only log when new endpoint is used
	if c.clients.GetSelectedEndpoint() != res.Endpoint {
		c.log.Info("Connected to secret chain", "endpoint", res.Endpoint)
	}

	c.clients.SetSelectedEndpoint(res.Endpoint)

	return nil
}

func (c *client) StartLivelinessCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = c.Connect(ctx)
		}
	}
}

func (c *client) CheckAndConnect(ctx context.Context) error {
	if _, err := c.clients.GetSelectedClient(); err != nil {
		return c.Connect(ctx)
	}
	return nil
}

func (c *client) getSelectedClient() (*sdkclient.Context, error) {
	cli, err := c.clients.GetSelectedClient()
	if err != nil {
		return nil, fmt.Errorf("no selected rpc endpoint: %w", err)
	}
	return cli, nil
}

func (c *client) BroadcastTx(txBlob []byte) (string, error) {
	cli, err := c.getSelectedClient()
	if err != nil {
		return "", err
	}

	res, err := cli.BroadcastTx(txBlob)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", fmt.Errorf("transaction failed with code %d: %s", res.Code, res.RawLog)
	}
	return res.TxHash, nil
}

func (c *client) GetAccount(ctx context.Context, sender string) (accountNumber uint64, sequence uint64, err error) {
	cli, err := c.getSelectedClient()
	if err != nil {
		return 0, 0, err
	}

	queryClient := authtypes.NewQueryClient(*cli)
	req := &authtypes.QueryAccountRequest{Address: sender}
	resp, err := queryClient.Account(ctx, req)
	if err != nil {
		return 0, 0, err
	}

	var account sdk.AccountI
	err = cli.InterfaceRegistry.UnpackAny(resp.Account, &account)
	if err != nil {
		return 0, 0, err
	}

	return account.GetAccountNumber(), account.GetSequence(), nil
}

func (c *client) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	cli, err := c.getSelectedClient()
	if err != nil {
		return nil, err
	}

	queryClient := banktypes.NewQueryClient(*cli)
	req := banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   c.denom,
	}

	resp, err := queryClient.Balance(ctx, &req)
	if err != nil {
		return nil, err
	}

	amount := resp.GetBalance().Amount
	return amount.BigInt(), nil
}

func (c *client) GetTx(ctx context.Context, txHash string) (*typesTxResult, error) {
	cli, err := c.getSelectedClient()
	if err != nil {
		return nil, err
	}

	node, err := cli.GetNode()
	if err != nil {
		return nil, err
	}

	txHashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		return nil, err
	}

	resultTx, err := node.Tx(ctx, txHashBytes, true)
	if err != nil {
		return nil, err
	}

	code := resultTx.TxResult.Code
	gasUsed := resultTx.TxResult.GasUsed
	height := resultTx.Height
	log := resultTx.TxResult.Log

	// normalize logs: some nodes return empty log on success
	if code == 0 {
		log = ""
	}

	return &typesTxResult{
		StatusCode: code,
		GasUsed:    gasUsed,
		Height:     height,
		Log:        log,
	}, nil
}

func (c *client) GetBlockByHeight(ctx context.Context, height *big.Int) (*typesBlockResult, error) {
	cli, err := c.getSelectedClient()
	if err != nil {
		return nil, err
	}

	node, err := cli.GetNode()
	if err != nil {
		return nil, err
	}

	h := height.Int64()
	resBlock, err := node.Block(ctx, &h)
	if err != nil {
		return nil, err
	}

	if resBlock.Block == nil {
		return nil, fmt.Errorf("block not found at height %d", h)
	}

	// Block timestamp is in UTC.
	return &typesBlockResult{Time: resBlock.Block.Time.UTC()}, nil
}

// getClientWithMaxHeight connects to the endpoint that has the highest block height.
func (c *client) getClientWithMaxHeight(ctx context.Context) (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.endpoints))

	for _, endpoint := range c.endpoints {
		go func(endpoint string) {
			cli, ok := c.clients.GetClient(endpoint)

			if !ok {
				ch <- ClientConnectionResult{Endpoint: endpoint, Client: nil, BlockHeight: 0}
				return
			}

			node, err := cli.GetNode()
			if err != nil {
				c.log.Warn(
					"Failed to get node from client",
					"endpoint", endpoint,
					"err", err,
				)
				ch <- ClientConnectionResult{Endpoint: endpoint, Client: nil, BlockHeight: 0}
				alert.HandleAlert(
					c.alert,
					alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
						WithChainName(c.chainName).
						WithEndpoint(endpoint),
					err.Error(),
				)
				return
			}

			status, err := node.Status(ctx)
			if err != nil {
				c.log.Warn(
					"Failed to get status from node",
					"endpoint", endpoint,
					"err", err,
				)
				ch <- ClientConnectionResult{Endpoint: endpoint, Client: nil, BlockHeight: 0}
				alert.HandleAlert(
					c.alert,
					alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
						WithChainName(c.chainName).
						WithEndpoint(endpoint),
					err.Error(),
				)
				return
			}

			c.log.Debug(
				"Get height of the given client",
				"endpoint", endpoint,
				"block_number", status.SyncInfo.LatestBlockHeight,
			)
			alert.HandleReset(
				c.alert,
				alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
					WithChainName(c.chainName).
					WithEndpoint(endpoint),
			)

			ch <- ClientConnectionResult{Endpoint: endpoint, Client: cli, BlockHeight: status.SyncInfo.LatestBlockHeight}
		}(endpoint)
	}

	var result ClientConnectionResult
	for i := 0; i < len(c.endpoints); i++ {
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
			alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.chainName),
			fmt.Sprintf("failed to connect to secret chain on all endpoints: %s", c.endpoints),
		)
		return ClientConnectionResult{}, fmt.Errorf("failed to connect to secret chain")
	}

	alert.HandleReset(c.alert, alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.chainName))

	return result, nil
}
