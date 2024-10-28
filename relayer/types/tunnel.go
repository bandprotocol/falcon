package types

import chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"

type BandChainInfo struct {
	LatestSequence uint64 `json:"latest_sequence"`
	IsActive       bool   `json:"is_active"`
}

// Tunnel defines the tunnel information.
type Tunnel struct {
	ID              uint64              `json:"id"`
	TargetChain     string              `json:"target_chain"`
	TargetAddress   string              `json:"target_address"`
	TunnelChainInfo *chainstypes.Tunnel `json:"tunnel_chain_info"`
	BandChainInfo   *BandChainInfo      `json:"band_chain_info"`
}

// NewBandChainInfo creates a new band chain info object
func NewBandChainInfo(
	latestSequence uint64,
	isActive bool,
) *BandChainInfo {
	return &BandChainInfo{
		LatestSequence: latestSequence,
		IsActive:       isActive,
	}
}

// NewTunnel creates a new tunnel object.
func NewTunnel(
	id uint64,
	targetChain string,
	targetAddress string,
	bandChainInfo *BandChainInfo,
	tunnelChainInfo *chainstypes.Tunnel,
) *Tunnel {
	return &Tunnel{
		ID:              id,
		TargetChain:     targetChain,
		TargetAddress:   targetAddress,
		BandChainInfo:   bandChainInfo,
		TunnelChainInfo: tunnelChainInfo,
	}
}
