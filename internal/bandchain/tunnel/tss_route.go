package types

import (
	bandtsstypes "github.com/bandprotocol/falcon/internal/bandchain/bandtss"
	feedstypes "github.com/bandprotocol/falcon/internal/bandchain/feeds"
)

var _ RouteI = &TSSRoute{}

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

// ValidateBasic performs basic validation of the TSSRoute fields.
func (r *TSSRoute) ValidateBasic() error {
	return nil
}

// NewTSSPacketReceipt creates a new TSSPacketReceipt instance.
func NewTSSPacketReceipt(signingID bandtsstypes.SigningID) *TSSPacketReceipt {
	return &TSSPacketReceipt{
		SigningID: signingID,
	}
}
