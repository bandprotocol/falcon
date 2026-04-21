package soroban

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stellar/go-stellar-sdk/clients/horizonclient"
	hProtocol "github.com/stellar/go-stellar-sdk/protocols/horizon"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/logger"
)

// HorizonClients holds Horizon RPC clients and the selected endpoint.
type HorizonClients = chains.ClientPool[horizonclient.Client]

// NewHorizonClients creates and returns a new SorobanClients instance with no endpoints.
func NewHorizonClients() HorizonClients {
	return chains.NewClientPool[horizonclient.Client]()
}

// ClientConnectionResult contains the result of connecting to a specific Horizon endpoint.
type ClientConnectionResult struct {
	Endpoint       string
	Client         *horizonclient.Client
	LedgerSequence uint64
}

type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	GetAccountSequenceNumber(account string) (int64, error)
	GetBalance(account string) (*big.Int, error)
	BroadcastTx(txBlob string) (TxResult, error)
	GetLedgerCloseTime(ledgerIndex uint64) (*time.Time, error)
	GetTransactionStatus(txHash string) (hProtocol.Transaction, error)
}

type client struct {
	ChainName        string
	HorizonEndpoints []string
	QueryTimeout     time.Duration

	Log   logger.Logger
	alert alert.Alert

	clients HorizonClients
}

func NewClient(chainName string, cfg *SorobanChainProviderConfig, log logger.Logger, alert alert.Alert) Client {
	return &client{
		ChainName:        chainName,
		HorizonEndpoints: cfg.HorizonEndpoints,
		QueryTimeout:     cfg.QueryTimeout,
		Log:              log.With("chain_name", chainName),
		alert:            alert,
		clients:          NewHorizonClients(),
	}
}

// Connect connects to all Horizon endpoints in parallel and selects the one with the highest ledger.
func (c *client) Connect(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, endpoint := range c.HorizonEndpoints {
		_, ok := c.clients.GetClient(endpoint)
		if ok {
			continue
		}

		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			hc := &horizonclient.Client{
				HorizonURL: strings.TrimRight(endpoint, "/"),
				HTTP:       &http.Client{Timeout: c.QueryTimeout},
			}
			if _, err := hc.Root(); err != nil {
				c.Log.Warn(
					"Failed to connect to Soroban chain",
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
			c.clients.SetClient(endpoint, hc)
		}(endpoint)
	}

	wg.Wait()
	res, err := c.getClientWithMaxLedger(ctx)
	if err != nil {
		c.Log.Error("Failed to connect to Soroban chain", err)
		return err
	}

	// only log when new endpoint is used
	if c.clients.GetSelectedEndpoint() != res.Endpoint {
		c.Log.Info("Connected to Soroban chain", "endpoint", res.Endpoint)
	}

	c.clients.SetSelectedEndpoint(res.Endpoint)
	return nil
}

// CheckAndConnect checks if the client is connected; if not, it connects.
func (c *client) CheckAndConnect(ctx context.Context) error {
	if _, err := c.clients.GetSelectedClient(); err != nil {
		return c.Connect(ctx)
	}
	return nil
}

// StartLivelinessCheck starts the liveliness check for the Soroban chain.
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

// getClientWithMaxLedger returns the connected client with the highest ingested ledger sequence.
func (c *client) getClientWithMaxLedger(ctx context.Context) (ClientConnectionResult, error) {
	ch := make(chan ClientConnectionResult, len(c.HorizonEndpoints))

	for _, endpoint := range c.HorizonEndpoints {
		go func(endpoint string) {
			hc, ok := c.clients.GetClient(endpoint)
			if !ok {
				ch <- ClientConnectionResult{endpoint, nil, 0}
				return
			}

			root, err := hc.Root()
			if err != nil {
				c.Log.Warn(
					"Failed to get latest ledger",
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

			ledgerSeq := uint64(root.IngestSequence)
			c.Log.Debug(
				"Get latest ledger of the given client",
				"endpoint", endpoint,
				"ledger_sequence", ledgerSeq,
			)
			alert.HandleReset(
				c.alert,
				alert.NewTopic(alert.ConnectSingleChainClientErrorMsg).
					WithChainName(c.ChainName).
					WithEndpoint(endpoint),
			)

			ch <- ClientConnectionResult{endpoint, hc, ledgerSeq}
		}(endpoint)
	}

	var result ClientConnectionResult
	for i := 0; i < len(c.HorizonEndpoints); i++ {
		r := <-ch
		if r.Client != nil {
			if r.LedgerSequence > result.LedgerSequence ||
				(r.Endpoint == c.clients.GetSelectedEndpoint() && r.LedgerSequence == result.LedgerSequence) {
				result = r
			}
		}
	}

	if result.Client == nil {
		alert.HandleAlert(
			c.alert,
			alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName),
			fmt.Sprintf("failed to connect to Soroban chain on all endpoints: %s", c.HorizonEndpoints),
		)
		return ClientConnectionResult{}, fmt.Errorf("failed to connect to Soroban chain")
	}

	alert.HandleReset(c.alert, alert.NewTopic(alert.ConnectAllChainClientErrorMsg).WithChainName(c.ChainName))
	return result, nil
}

func (c *client) GetAccountSequenceNumber(account string) (int64, error) {
	hc, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return 0, fmt.Errorf("failed to get client: %w", err)
	}

	acc, err := hc.AccountDetail(horizonclient.AccountRequest{AccountID: account})
	if err != nil {
		return 0, fmt.Errorf("failed to fetch account: %w", err)
	}

	seq, err := acc.GetSequenceNumber()
	if err != nil {
		return 0, fmt.Errorf("failed to get sequence number: %w", err)
	}

	if seq < 0 {
		return 0, fmt.Errorf("negative sequence number: %d", seq)
	}

	return seq, nil
}

func (c *client) GetBalance(account string) (*big.Int, error) {
	hc, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	acc, err := hc.AccountDetail(horizonclient.AccountRequest{AccountID: account})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	native, err := acc.GetNativeBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get native balance: %w", err)
	}

	// Convert "10.1234567" to stroops (x 10^7) using exact decimal arithmetic.
	d, err := decimal.NewFromString(native)
	if err != nil {
		return nil, fmt.Errorf("failed to parse balance %q: %w", native, err)
	}
	stroops := d.Mul(decimal.NewFromInt(1e7)).BigInt()
	return stroops, nil
}

func (c *client) BroadcastTx(txBlob string) (TxResult, error) {
	hc, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return TxResult{}, fmt.Errorf("failed to get client: %w", err)
	}

	resp, err := hc.SubmitTransactionXDR(txBlob)
	if err != nil {
		return TxResult{}, fmt.Errorf("failed to submit transaction: %w", err)
	}

	return TxResult{TxHash: resp.Hash}, nil
}

func (c *client) GetLedgerCloseTime(ledgerIndex uint64) (*time.Time, error) {
	hc, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	ledger, err := hc.LedgerDetail(uint32(ledgerIndex))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ledger %d: %w", ledgerIndex, err)
	}

	t := ledger.ClosedAt
	return &t, nil
}

func (c *client) GetTransactionStatus(txHash string) (hProtocol.Transaction, error) {
	hc, err := c.clients.GetSelectedClient()
	if err != nil {
		c.Log.Error("Failed to get client", "endpoint", c.clients.GetSelectedEndpoint(), err)
		return hProtocol.Transaction{}, fmt.Errorf("failed to get client: %w", err)
	}

	return hc.TransactionDetail(txHash)
}
