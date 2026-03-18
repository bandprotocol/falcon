package types

import "math/big"

// Tunnel defines the tunnel information on target chain.
type Tunnel struct {
	ID            uint64 `json:"-"`
	TargetAddress string `json:"-"`
	IsActive      bool   `json:"is_active"`
	// LatestSequence is the last sequence confirmed on the target chain.
	// A non-nil value means the chain provides an authoritative on-chain sequence
	// (e.g. EVM contract). A nil value means the chain does not track sequence
	// on-chain (e.g. XRPL): the relayer skips the current BandChain latest on
	// cold start and only relays genuinely new packets, avoiding duplicates and
	// stale sends.
	LatestSequence *uint64  `json:"latest_sequence"`
	Balance        *big.Int `json:"balance"`
}

// NewTunnel creates a new tunnel object.
func NewTunnel(
	id uint64,
	targetAddress string,
	isActive bool,
	latestSequence *uint64,
	balance *big.Int,
) *Tunnel {
	return &Tunnel{
		ID:             id,
		TargetAddress:  targetAddress,
		IsActive:       isActive,
		LatestSequence: latestSequence,
		Balance:        balance,
	}
}
