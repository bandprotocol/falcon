package tunnel

import (
	"github.com/cosmos/gogoproto/proto"
)

// events
const (
	EventTypeProducePacketSuccess = "produce_packet_success"
	EventTypeTriggerTunnel        = "trigger_tunnel"

	AttributeKeyTunnelID = "tunnel_id"
)

// isTssRouteType checks if the route type is TSSRoute
func IsTssRouteType(routeType string) bool {
	tssRouteType := proto.MessageName(&TSSRoute{})
	return routeType == tssRouteType || routeType == "/"+tssRouteType
}
