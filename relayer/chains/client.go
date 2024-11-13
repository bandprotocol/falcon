package chains

import "math/big"

// Client defines the interface for the target chain client
type Client interface {
	// GetNonce returns the nonce of the given address
	GetNonce(address string) (uint64, error)

	// BroadcastTx broadcasts the given raw transaction
	BroadcastTx(rawTx string) (string, error)

	// GetBalances returns the balances of the given accounts
	GetBalances(accounts []string) ([]*big.Int, error)

	// GetTunnelNonce returns the nonce of the given tunnel
	GetTunnelNonce(targetAddress string, tunnelID uint64) (uint64, error)
}
