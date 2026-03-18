package chains

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/bytes"

	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
)

// SelectSigning selects the signing from the packet.
func SelectSigning(packet *bandtypes.Packet) (*bandtypes.Signing, error) {
	// get signing from packet; prefer to use signing from
	// current group than incoming group
	if packet.CurrentGroupSigning != nil {
		return packet.CurrentGroupSigning, nil
	}

	if packet.IncomingGroupSigning != nil {
		return packet.IncomingGroupSigning, nil
	}

	return nil, fmt.Errorf("missing signing")
}

// ExtractEVMSignature extracts the EVM signature from the signing.
// If the signing is nil, it returns empty byte slices.
func ExtractEVMSignature(evmSignature *bandtypes.EVMSignature) (bytes.HexBytes, bytes.HexBytes) {
	rAddress := []byte{}
	signature := []byte{}
	if evmSignature != nil {
		rAddress = evmSignature.RAddress
		signature = evmSignature.Signature
	}

	return rAddress, signature
}
