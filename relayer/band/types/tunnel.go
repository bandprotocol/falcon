package types

// Tunnel stores an information of the tunnel.
type Tunnel struct {
	ID             uint64
	LatestSequence uint64
	TargetAddress  string
	TargetChainID  string
	IsActive       bool
}

// NewTunnel creates a new tunnel instance.
func NewTunnel(
	id uint64,
	latestSequence uint64,
	targetAddress string,
	targetChainID string,
	isActive bool,
) *Tunnel {
	return &Tunnel{
		ID:             id,
		LatestSequence: latestSequence,
		TargetAddress:  targetAddress,
		TargetChainID:  targetChainID,
		IsActive:       isActive,
	}
}
