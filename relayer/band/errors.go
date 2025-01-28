package band

import (
	"fmt"

	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
)

var (
	ErrBandChainNotConnect = fmt.Errorf("cannot connect to Bandchain")

	ErrUnsupportedRouteType = func(route string) error {
		return fmt.Errorf("unsupported route type: %s", route)
	}

	ErrUnsupportedPacketContentType = func(packetReceipt tunneltypes.PacketReceiptI) error {
		return fmt.Errorf("unsupported packet content type: %T", packetReceipt)
	}
)
