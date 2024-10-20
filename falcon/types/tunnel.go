package types

import chainstypes "github.com/bandprotocol/falcon/falcon/chains/types"

// Tunnel defines the tunnel information.
type Tunnel struct {
	ID              uint64              `json:"id"`
	TargetChain     string              `json:"target_chain"`
	TargetAddress   string              `json:"target_address"`
	TunnelChainInfo *chainstypes.Tunnel `json:"chain_info,omitempty"`
}

// NewTunnel creates a new tunnel object.
func NewTunnel(
	id uint64,
	targetChain string,
	targetAddress string,
	tunnelChainInfo *chainstypes.Tunnel,
) *Tunnel {
	return &Tunnel{
		ID:              id,
		TargetChain:     targetChain,
		TargetAddress:   targetAddress,
		TunnelChainInfo: tunnelChainInfo,
	}
}