package types

import (
	feedstypes "github.com/bandprotocol/falcon/internal/bandchain/feeds"
)

// IBCRoute defines the IBC route for the tunnel module
var _ RouteI = &IBCRoute{}

// NewIBCRoute creates a new IBCRoute instance.
func NewIBCRoute(channelID string) *IBCRoute {
	return &IBCRoute{
		ChannelID: channelID,
	}
}

// ValidateBasic validates the IBCRoute
func (r *IBCRoute) ValidateBasic() error {
	return nil
}

// NewIBCPacketReceipt creates a new IBCPacketReceipt instance.
func NewIBCPacketReceipt(sequence uint64) *IBCPacketReceipt {
	return &IBCPacketReceipt{
		Sequence: sequence,
	}
}

// NewTunnelPricesPacketData creates a new TunnelPricesPacketData instance.
func NewTunnelPricesPacketData(
	tunnelID uint64,
	sequence uint64,
	prices []feedstypes.Price,
	createdAt int64,
) TunnelPricesPacketData {
	return TunnelPricesPacketData{
		TunnelID:  tunnelID,
		Sequence:  sequence,
		Prices:    prices,
		CreatedAt: createdAt,
	}
}
