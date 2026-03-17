package types

import "math/big"

// Tunnel defines the tunnel information on target chain.
type Tunnel struct {
	ID            uint64 `json:"-"`
	TargetAddress string `json:"-"`
	IsActive      bool   `json:"is_active"`
	// pointer is used to distinguish between zero value and unknown value.
	LatestSequence  *uint64  `json:"latest_sequence"`
	Balance         *big.Int `json:"balance"`
	SupportContract bool     `json:"-"`
}

// NewTunnel creates a new tunnel object.
func NewTunnel(
	id uint64,
	targetAddress string,
	isActive bool,
	latestSequence *uint64,
	balance *big.Int,
	supportContract bool,
) *Tunnel {
	return &Tunnel{
		ID:              id,
		TargetAddress:   targetAddress,
		IsActive:        isActive,
		LatestSequence:  latestSequence,
		Balance:         balance,
		SupportContract: supportContract,
	}
}
