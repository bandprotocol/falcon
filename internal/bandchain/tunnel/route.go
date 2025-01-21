package types

import "github.com/cosmos/gogoproto/proto"

// RouteI defines a routing path to deliver data to the destination.
type RouteI interface {
	proto.Message

	ValidateBasic() error
}
