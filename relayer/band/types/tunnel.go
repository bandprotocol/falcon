package types

// Tunnel stores an information of the tunnel.
type Tunnel struct {
	ID             uint64 `json:"-"`
	LatestSequence uint64 `json:"latest_sequence"`
	TargetAddress  string `json:"target_address"`
	TargetChainID  string `json:"target_chain_id"`
	IsActive       bool   `json:"is_active"`
	Creator        string `json:"creator"`
}

// NewTunnel creates a new tunnel instance.
func NewTunnel(
	id uint64,
	latestSequence uint64,
	targetAddress string,
	targetChainID string,
	isActive bool,
	creator string,
) *Tunnel {
	return &Tunnel{
		ID:             id,
		LatestSequence: latestSequence,
		TargetAddress:  targetAddress,
		TargetChainID:  targetChainID,
		IsActive:       isActive,
		Creator:        creator,
	}
}
