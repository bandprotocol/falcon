package tunnel

import (
	"github.com/cosmos/gogoproto/proto"

	bandtsstypes "github.com/bandprotocol/falcon/internal/bandchain/bandtss"
	feedstypes "github.com/bandprotocol/falcon/internal/bandchain/feeds"
)

// PacketReceiptI defines an interface for confirming the delivery of a packet to its destination via the specified route.
type PacketReceiptI interface {
	proto.Message
}

// RouteI defines a routing path to deliver data to the destination.
type RouteI interface {
	proto.Message
}

// NewTSSRoute return a new TSSRoute instance.
func NewTSSRoute(
	destinationChainID string,
	destinationContractAddress string,
	encoder feedstypes.Encoder,
) TSSRoute {
	return TSSRoute{
		DestinationChainID:         destinationChainID,
		DestinationContractAddress: destinationContractAddress,
		Encoder:                    encoder,
	}
}

// NewTSSPacketReceipt creates a new TSSPacketReceipt instance.
func NewTSSPacketReceipt(signingID bandtsstypes.SigningID) *TSSPacketReceipt {
	return &TSSPacketReceipt{
		SigningID: signingID,
	}
}

// NewIBCRoute creates a new IBCRoute instance.
func NewIBCRoute(channelID string) *IBCRoute {
	return &IBCRoute{
		ChannelID: channelID,
	}
}

// NewIBCPacketReceipt creates a new IBCPacketReceipt instance.
func NewIBCPacketReceipt(sequence uint64) *IBCPacketReceipt {
	return &IBCPacketReceipt{
		Sequence: sequence,
	}
}
