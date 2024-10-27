package types

// Packet stores an information of the packet that is generated from the tunnel.
type Packet struct {
	TunnelID              uint64
	Sequence              uint64
	SigningID             uint64
	TargetChainID         string
	TargetContractAddress string
}

// NewPacket creates a new packet instance.
func NewPacket(
	tunnelID uint64,
	sequence uint64,
	signingID uint64,
) *Packet {
	return &Packet{
		TunnelID:  tunnelID,
		Sequence:  sequence,
		SigningID: signingID,
	}
}
