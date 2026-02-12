package chains

import (
	"fmt"

	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
)

// SelectSigning selects the signing from the packet.
func SelectSigning(packet *bandtypes.Packet) (*bandtypes.Signing, error) {
	// get signing from packet; prefer to use signing from
	// current group than incoming group
	if packet.CurrentGroupSigning != nil {
		return packet.CurrentGroupSigning, nil
	} else if packet.IncomingGroupSigning != nil {
		return packet.IncomingGroupSigning, nil
	} else {
		return nil, fmt.Errorf("missing signing")
	}
}
