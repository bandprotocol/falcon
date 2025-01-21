package feeds

import "github.com/cosmos/gogoproto/proto"

// PriceStatus is a structure that defines the price status of a price.
type PriceStatus int32

const (
	// PRICE_STATUS_UNSPECIFIED is an unspecified price status.
	PRICE_STATUS_UNSPECIFIED PriceStatus = 0
	// PRICE_STATUS_UNKNOWN_SIGNAL_ID is an unknown signal id price status.
	PRICE_STATUS_UNKNOWN_SIGNAL_ID PriceStatus = 1
	// PRICE_STATUS_NOT_READY is a not ready price status.
	PRICE_STATUS_NOT_READY PriceStatus = 2
	// PRICE_STATUS_AVAILABLE is an available price status.
	PRICE_STATUS_AVAILABLE PriceStatus = 3
	// PRICE_STATUS_NOT_IN_CURRENT_FEEDS is a not in current feed price status.
	PRICE_STATUS_NOT_IN_CURRENT_FEEDS PriceStatus = 4
)

var PriceStatus_name = map[int32]string{
	0: "PRICE_STATUS_UNSPECIFIED",
	1: "PRICE_STATUS_UNKNOWN_SIGNAL_ID",
	2: "PRICE_STATUS_NOT_READY",
	3: "PRICE_STATUS_AVAILABLE",
	4: "PRICE_STATUS_NOT_IN_CURRENT_FEEDS",
}

var PriceStatus_value = map[string]int32{
	"PRICE_STATUS_UNSPECIFIED":          0,
	"PRICE_STATUS_UNKNOWN_SIGNAL_ID":    1,
	"PRICE_STATUS_NOT_READY":            2,
	"PRICE_STATUS_AVAILABLE":            3,
	"PRICE_STATUS_NOT_IN_CURRENT_FEEDS": 4,
}

func (x PriceStatus) String() string {
	return proto.EnumName(PriceStatus_name, int32(x))
}

// Price is a structure that defines the price of a signal id.
type Price struct {
	// status is the status of a the price.
	Status PriceStatus `protobuf:"varint,1,opt,name=status,proto3,enum=band.feeds.v1beta1.PriceStatus" json:"status,omitempty"`
	// signal_id is the signal id of the price.
	SignalID string `protobuf:"bytes,2,opt,name=signal_id,json=signalId,proto3"                     json:"signal_id,omitempty"`
	// price is the price of the signal id.
	Price uint64 `protobuf:"varint,3,opt,name=price,proto3"                                      json:"price,omitempty"`
	// timestamp is the timestamp at which the price was aggregated.
	Timestamp int64 `protobuf:"varint,4,opt,name=timestamp,proto3"                                  json:"timestamp,omitempty"`
}

func (m *Price) Reset()         { *m = Price{} }
func (m *Price) String() string { return proto.CompactTextString(m) }
func (*Price) ProtoMessage()    {}
