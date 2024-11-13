package types

import (
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

// Tunnel defines the tunnel information.
type Tunnel struct {
	ID              uint64              `json:"id"`
	BandChainInfo   *bandtypes.Tunnel   `json:"band_chain_info"`
	TunnelChainInfo *chainstypes.Tunnel `json:"tunnel_chain_info"`
}

// NewTunnel creates a new tunnel object.
func NewTunnel(
	bandChainInfo *bandtypes.Tunnel,
	tunnelChainInfo *chainstypes.Tunnel,
) *Tunnel {
	return &Tunnel{
		ID:              bandChainInfo.ID,
		BandChainInfo:   bandChainInfo,
		TunnelChainInfo: tunnelChainInfo,
	}
}
