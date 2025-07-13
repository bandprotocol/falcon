package types

// SignalPrice stores information about signal ID and it price.
type SignalPrice struct {
	SignalID string `json:"signal_id"`
	Price    uint64 `json:"price"`
}

// Packet stores information about the packet generated from the tunnel.
type Packet struct {
	TunnelID             uint64        `json:"tunnel_id"`
	Sequence             uint64        `json:"sequence"`
	SignalPrices         []SignalPrice `json:"signal_prices"`
	CurrentGroupSigning  *Signing      `json:"current_group_signing"`
	IncomingGroupSigning *Signing      `json:"incoming_group_signing"`
}

// NewPacket creates a new Packet instance.
func NewPacket(
	tunnelID uint64,
	sequence uint64,
	signalPrices []SignalPrice,
	currentGroupSigning *Signing,
	incomingGroupSigning *Signing,
) *Packet {
	return &Packet{
		TunnelID:             tunnelID,
		Sequence:             sequence,
		SignalPrices:         signalPrices,
		CurrentGroupSigning:  currentGroupSigning,
		IncomingGroupSigning: incomingGroupSigning,
	}
}

// GetCurrentGroupSigningID returns the signing ID of the current group.
func (p *Packet) GetCurrentGroupSigningID() uint64 {
	if p.CurrentGroupSigning == nil {
		return 0
	}

	return p.CurrentGroupSigning.ID
}

// GetIncomingGroupSigningID returns the signing ID of the incoming group.
func (p *Packet) GetIncomingGroupSigningID() uint64 {
	if p.IncomingGroupSigning == nil {
		return 0
	}

	return p.IncomingGroupSigning.ID
}
