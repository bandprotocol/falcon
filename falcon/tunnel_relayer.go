package falcon

import (
	"time"

	"github.com/bandprotocol/falcon/falcon/band"
	"github.com/bandprotocol/falcon/falcon/chains"
)

// TunnelRelayer is a relayer that listens to the tunnel and relays the packet
type TunnelRelayer struct {
	TunnelID               uint64
	ContractAddress        string
	CheckingPacketInterval time.Duration

	BandClient        band.Client
	TargetChainClient chains.Client
}

// NewTunnelRelayer creates a new TunnelRelayer
func NewTunnelRelayer(
	tunnelID uint64,
	contractAddress string,
	checkingPacketInterval time.Duration,
	bandClient band.Client,
	targetChainClient chains.Client,
) *TunnelRelayer {
	return &TunnelRelayer{
		TunnelID:               tunnelID,
		ContractAddress:        contractAddress,
		CheckingPacketInterval: checkingPacketInterval,

		BandClient:        bandClient,
		TargetChainClient: targetChainClient,
	}
}
