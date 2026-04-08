package secret

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/logger"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	std "github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
)

type SecretClients = chains.ClientPool[sdkclient.Context]

// Client defines the Cosmos RPC operations required by the Secret chain provider.
type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)

	BroadcastTx(txBlob []byte) (string, error)
	GetBalance(address string) (*big.Int, error)
	GetTx(txHash string) (*typesTxResult, error)
	GetBlockByHeight(height *big.Int) (*typesBlockResult, error)
	GetAccount(sender string) (accountNumber uint64, sequence uint64, err error)
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

func (c *client) Connect(_ context.Context) error {
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

			ctx := sdkclient.Context{}.
				WithAccountRetriever(authtypes.AccountRetriever{}).
				WithBroadcastMode(flags.BroadcastSync).
				WithCodec(enc.Marshaler).
				WithInterfaceRegistry(enc.InterfaceRegistry).
				WithTxConfig(enc.TxConfig).
				WithClient(rpcClient).
				WithNodeURI(endpoint)

			c.clients.SetClient(endpoint, &ctx)
		}(endpoint)
	}

	wg.Wait()

	// Select the first connected endpoint (or just first endpoint if none exist yet).
	if c.clients.GetSelectedEndpoint() == "" && len(c.endpoints) > 0 {
		c.clients.SetSelectedEndpoint(c.endpoints[0])
	}

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

func (c *client) GetAccount(sender string) (accountNumber uint64, sequence uint64, err error) {
	cli, err := c.getSelectedClient()
	if err != nil {
		return 0, 0, err
	}

	acc, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return 0, 0, err
	}

	account, err := cli.AccountRetriever.GetAccount(*cli, acc)
	if err != nil {
		return 0, 0, err
	}

	return account.GetAccountNumber(), account.GetSequence(), nil
}

func (c *client) GetBalance(address string) (*big.Int, error) {
	cli, err := c.getSelectedClient()
	if err != nil {
		return nil, err
	}

	queryClient := banktypes.NewQueryClient(*cli)
	req := banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   c.denom,
	}

	resp, err := queryClient.Balance(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	amount := resp.GetBalance().Amount
	return amount.BigInt(), nil
}

func (c *client) GetTx(txHash string) (*typesTxResult, error) {
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

	resultTx, err := node.Tx(context.Background(), txHashBytes, true)
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

func (c *client) GetBlockByHeight(height *big.Int) (*typesBlockResult, error) {
	cli, err := c.getSelectedClient()
	if err != nil {
		return nil, err
	}

	node, err := cli.GetNode()
	if err != nil {
		return nil, err
	}

	h := height.Int64()
	resBlock, err := node.Block(context.Background(), &h)
	if err != nil {
		return nil, err
	}

	if resBlock.Block == nil {
		return nil, fmt.Errorf("block not found at height %d", h)
	}

	// Block timestamp is in UTC.
	return &typesBlockResult{Time: resBlock.Block.Header.Time.UTC()}, nil
}
