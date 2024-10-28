package types

// SignalPrice stores information about signal ID and it price.
type SignalPrice struct {
	SignalID string `json:"signal_id"`
	Price    uint64 `json:"price"`
}

// Packet stores information about the packet generated from the tunnel.
type Packet struct {
	TunnelID     uint64        `json:"tunnel_id"`
	Sequence     uint64        `json:"sequence"`
	SignalPrices []SignalPrice `json:"signal_prices"`
	SigningInfo  *SigningInfo  `json:"signing_info"`
}

// NewPacket creates a new Packet instance.
func NewPacket(
	tunnelID uint64,
	sequence uint64,
	signalPrices []SignalPrice,
	signingInfo *SigningInfo,
) *Packet {
	return &Packet{
		TunnelID:     tunnelID,
		Sequence:     sequence,
		SignalPrices: signalPrices,
		SigningInfo:  signingInfo,
	}
}
