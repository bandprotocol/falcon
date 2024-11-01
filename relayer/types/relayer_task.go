package types

import bandtypes "github.com/bandprotocol/falcon/relayer/band/types"

// RelayerTask defines a task for the relayer.
type RelayerTask struct {
	Tunnel *bandtypes.Tunnel
	Packet *bandtypes.Packet
}

// NewTssMessage creates a new TSS message object.
func NewRelayerTask(
	tunnel *bandtypes.Tunnel,
	packet *bandtypes.Packet,
) RelayerTask {
	return RelayerTask{
		Tunnel: tunnel,
		Packet: packet,
	}
}
