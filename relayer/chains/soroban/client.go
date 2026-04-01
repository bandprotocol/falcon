package soroban

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/logger"
)

type Client interface {
	Connect(ctx context.Context) error
	CheckAndConnect(ctx context.Context) error
	StartLivelinessCheck(ctx context.Context, interval time.Duration)
	GetAccountSequenceNumber(account string) (uint64, error)
	GetBalance(account string) (*big.Int, error)
	GetLatestLedger() (uint64, *time.Time, error)
	BroadcastTx(txBlob string) (TxResult, error)
	GetLedgerCloseTime(ledgerIndex uint64) (*time.Time, error)
	GetEndpoint() string
}

type client struct {
	ChainName        string
	Endpoints        []string
	HorizonEndpoint  string
	SelectedEndpoint string

	Log   logger.Logger
	alert alert.Alert
}

type TxResult struct {
	TxHash      string
	LedgerIndex uint64
}

func NewClient(chainName string, cfg *SorobanChainProviderConfig, log logger.Logger, alert alert.Alert) Client {
	return &client{
		ChainName:       chainName,
		Endpoints:       cfg.Endpoints,
		HorizonEndpoint: cfg.HorizonEndpoint,
		Log:             log.With("chain_name", chainName),
		alert:           alert,
	}
}

func (c *client) Connect(ctx context.Context) error {
	// Simple endpoint selection based on reaching getLatestLedger
	var bestEndpoint string
	var highestSequence uint64

	for _, endpoint := range c.Endpoints {
		reqBody := `{"jsonrpc": "2.0", "id": 1, "method": "getLatestLedger"}`
		resp, err := http.Post(endpoint, "application/json", strings.NewReader(reqBody)) // #nosec G107
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var result struct {
				Result struct {
					Sequence uint64 `json:"sequence"`
				} `json:"result"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				continue
			}
			if result.Result.Sequence > highestSequence {
				highestSequence = result.Result.Sequence
				bestEndpoint = endpoint
			}
		}
	}

	if bestEndpoint == "" {
		return fmt.Errorf("could not connect to any soroban endpoint")
	}

	c.SelectedEndpoint = bestEndpoint
	return nil
}

func (c *client) CheckAndConnect(ctx context.Context) error {
	if c.SelectedEndpoint == "" {
		return c.Connect(ctx)
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

func (c *client) GetAccountSequenceNumber(account string) (uint64, error) {
	url := fmt.Sprintf("%s/accounts/%s", strings.TrimRight(c.HorizonEndpoint, "/"), account)
	resp, err := http.Get(url) // #nosec G107
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("failed to fetch sequence, status: %d", resp.StatusCode)
	}

	var result struct {
		Sequence string `json:"sequence"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	var seq uint64
	if _, err := fmt.Sscanf(result.Sequence, "%d", &seq); err != nil {
		return 0, err
	}
	return seq, nil
}

func (c *client) GetBalance(account string) (*big.Int, error) {
	url := fmt.Sprintf("%s/accounts/%s", strings.TrimRight(c.HorizonEndpoint, "/"), account)
	resp, err := http.Get(url) // #nosec G107
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch balance, status: %d", resp.StatusCode)
	}

	var result struct {
		Balances []struct {
			Balance   string `json:"balance"`
			AssetType string `json:"asset_type"`
		} `json:"balances"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, bal := range result.Balances {
		if bal.AssetType == "native" {
			// Convert "10.1234567" to drops/stroops (x 10^7)
			parts := strings.Split(bal.Balance, ".")
			base := parts[0]
			frac := ""
			if len(parts) > 1 {
				frac = parts[1]
			}
			if len(frac) > 7 {
				frac = frac[:7]
			}
			for len(frac) < 7 {
				frac += "0"
			}
			bigBal := new(big.Int)
			bigBal.SetString(base+frac, 10)
			return bigBal, nil
		}
	}

	return big.NewInt(0), nil
}

func (c *client) GetLatestLedger() (uint64, *time.Time, error) {
	reqBody := `{"jsonrpc": "2.0", "id": 1, "method": "getLatestLedger"}`
	resp, err := http.Post(c.SelectedEndpoint, "application/json", strings.NewReader(reqBody)) // #nosec G107
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			Sequence uint64 `json:"sequence"`
			// Note: may not include close time directly from RPC
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, nil, err
	}
	return result.Result.Sequence, nil, nil
}

func (c *client) BroadcastTx(txBlob string) (TxResult, error) {
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "sendTransaction",
		"params":  map[string]interface{}{"transaction": txBlob},
	}
	b, _ := json.Marshal(reqBody)
	resp, err := http.Post(c.SelectedEndpoint, "application/json", bytes.NewReader(b)) // #nosec G107
	if err != nil {
		return TxResult{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			Status         string `json:"status"`
			Hash           string `json:"hash"`
			LatestLedger   uint64 `json:"latestLedger"`
			ErrorResultXdr string `json:"errorResultXdr"`
		} `json:"result"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return TxResult{}, err
	}

	if result.Error != nil {
		return TxResult{}, fmt.Errorf("rpc error: %s", result.Error.Message)
	}

	if result.Result.Status == "ERROR" {
		return TxResult{TxHash: result.Result.Hash}, fmt.Errorf("transaction error: %s", result.Result.ErrorResultXdr)
	}

	return TxResult{
		TxHash:      result.Result.Hash,
		LedgerIndex: result.Result.LatestLedger,
	}, nil
}

func (c *client) GetEndpoint() string {
	return c.SelectedEndpoint
}

func (c *client) GetLedgerCloseTime(ledgerIndex uint64) (*time.Time, error) {
	// From Horizon: /ledgers/{id}
	url := fmt.Sprintf("%s/ledgers/%d", strings.TrimRight(c.HorizonEndpoint, "/"), ledgerIndex)
	resp, err := http.Get(url) // #nosec G107
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch ledger, status: %d", resp.StatusCode)
	}

	var result struct {
		ClosedAt string `json:"closed_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, result.ClosedAt)
	if err != nil {
		// Fallback to RFC3339
		t, err = time.Parse(time.RFC3339, result.ClosedAt)
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}
