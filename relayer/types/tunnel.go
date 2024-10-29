package types

import (
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

// Tunnel defines the tunnel information.
type Tunnel struct {
	BandChainInfo   *bandtypes.Tunnel   `json:"band_chain_info"`
	TunnelChainInfo *chainstypes.Tunnel `json:"tunnel_chain_info"`
}

// NewTunnel creates a new tunnel object.
func NewTunnel(
	bandChainInfo *bandtypes.Tunnel,
	tunnelChainInfo *chainstypes.Tunnel,
) *Tunnel {
	return &Tunnel{
		BandChainInfo:   bandChainInfo,
		TunnelChainInfo: tunnelChainInfo,
	}
}
