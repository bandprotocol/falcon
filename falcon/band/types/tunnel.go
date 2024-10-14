package types

// Tunnel stores an information of the tunnel.
type Tunnel struct {
	ID             uint64
	LatestSequence uint64
	TargetAddress  string
	TargetChainID  string
}

// NewTunnel creates a new tunnel instance.
func NewTunnel(
	id uint64,
	latestSequence uint64,
	targetAddress string,
	targetChainID string,
) *Tunnel {
	return &Tunnel{
		ID:             id,
		LatestSequence: latestSequence,
		TargetAddress:  targetAddress,
		TargetChainID:  targetChainID,
	}
}
