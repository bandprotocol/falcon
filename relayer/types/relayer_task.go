package types

import bandtypes "github.com/bandprotocol/falcon/relayer/band/types"

// RelayerTask defines a task for the relayer.
type RelayerTask struct {
	Packet  *bandtypes.Packet
	Signing *bandtypes.Signing
}

// NewTssMessage creates a new TSS message object.
func NewRelayerTask(
	packet *bandtypes.Packet,
	signing *bandtypes.Signing,
) RelayerTask {
	return RelayerTask{
		Packet:  packet,
		Signing: signing,
	}
}
