package xrpl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Peersyst/xrpl-go/xrpl/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/logger"
)

func TestBroadcastTxWaitsForQueuedTransaction(t *testing.T) {
	var submitRequests atomic.Int32
	var txRequests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		switch request.Method {
		case "submit":
			submitRequests.Add(1)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"engine_result":         queuedTransactionResult,
					"engine_result_message": "Held until escalated fee drops",
					"tx_json": map[string]any{
						"hash": "HASH",
						"Fee":  "10",
					},
				},
			})
		case "tx":
			if txRequests.Add(1) == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"result": map[string]any{"error": "txnNotFound"},
				})
				return
			}

			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"hash":         "HASH",
					"ledger_index": 123,
					"meta": map[string]any{
						"TransactionResult": successfulTransactionResult,
					},
					"validated": true,
				},
			})
		default:
			http.Error(w, "unexpected method", http.StatusBadRequest)
		}
	}))
	defer server.Close()

	rpcConfig, err := rpc.NewClientConfig(server.URL)
	require.NoError(t, err)
	rpcClient := rpc.NewClient(rpcConfig)

	client := &client{
		ChainName:         "xrpl-test",
		TxPollingInterval: time.Millisecond,
		Log:               logger.NewZapLogWrapper(zap.NewNop().Sugar()),
		clients:           NewXRPLClients(),
	}
	client.clients.SetClient(server.URL, rpcClient)
	client.clients.SetSelectedEndpoint(server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := client.BroadcastTx(ctx, "signed-blob")
	require.NoError(t, err)
	require.Equal(t, "HASH", result.TxHash)
	require.Equal(t, "10", result.Fee)
	require.Equal(t, 123, result.LedgerIndex.Int())
	require.EqualValues(t, 1, submitRequests.Load())
	require.EqualValues(t, 2, txRequests.Load())
}
