package types

// Packet defines the packet information.
type Packet struct {
	TunnelID  uint64 `json:"tunnel_id"`
	Sequence  uint64 `json:"sequence"`
	SigningID uint64 `json:"signing_id"`
}

// NewPacket creates a new packet object.
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
