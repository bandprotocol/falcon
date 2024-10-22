package types

import chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"

// Tunnel defines the tunnel information.
type Tunnel struct {
	ID              uint64              `json:"id"`
	TargetChain     string              `json:"target_chain"`
	TargetAddress   string              `json:"target_address"`
	TunnelChainInfo *chainstypes.Tunnel `json:"tunnel_chain_info"`
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
