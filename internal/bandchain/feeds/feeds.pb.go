package types

import (
	"fmt"
	"io"
	"math"
	math_bits "math/bits"

	"github.com/cosmos/gogoproto/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

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

func (PriceStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{0}
}

// SignalPriceStatus is a structure that defines the price status of a signal id.
type SignalPriceStatus int32

const (
	// SIGNAL_PRICE_STATUS_UNSPECIFIED is an unspecified signal price status.
	SIGNAL_PRICE_STATUS_UNSPECIFIED SignalPriceStatus = 0
	// SIGNAL_PRICE_STATUS_UNSUPPORTED is an unsupported signal price status.
	SIGNAL_PRICE_STATUS_UNSUPPORTED SignalPriceStatus = 1
	// SIGNAL_PRICE_STATUS_UNAVAILABLE is an unavailable signal price status.
	SIGNAL_PRICE_STATUS_UNAVAILABLE SignalPriceStatus = 2
	// SIGNAL_PRICE_STATUS_AVAILABLE is an available signal price status.
	SIGNAL_PRICE_STATUS_AVAILABLE SignalPriceStatus = 3
)

var SignalPriceStatus_name = map[int32]string{
	0: "SIGNAL_PRICE_STATUS_UNSPECIFIED",
	1: "SIGNAL_PRICE_STATUS_UNSUPPORTED",
	2: "SIGNAL_PRICE_STATUS_UNAVAILABLE",
	3: "SIGNAL_PRICE_STATUS_AVAILABLE",
}

var SignalPriceStatus_value = map[string]int32{
	"SIGNAL_PRICE_STATUS_UNSPECIFIED": 0,
	"SIGNAL_PRICE_STATUS_UNSUPPORTED": 1,
	"SIGNAL_PRICE_STATUS_UNAVAILABLE": 2,
	"SIGNAL_PRICE_STATUS_AVAILABLE":   3,
}

func (x SignalPriceStatus) String() string {
	return proto.EnumName(SignalPriceStatus_name, int32(x))
}

func (SignalPriceStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{1}
}

// Signal is the data structure that contains signal id and power of that signal.
type Signal struct {
	// id is the id of the signal.
	ID string `protobuf:"bytes,1,opt,name=id,proto3"     json:"id,omitempty"`
	// power is the power of the corresponding signal id.
	Power int64 `protobuf:"varint,2,opt,name=power,proto3" json:"power,omitempty"`
}

func (m *Signal) Reset()         { *m = Signal{} }
func (m *Signal) String() string { return proto.CompactTextString(m) }
func (*Signal) ProtoMessage()    {}
func (*Signal) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{0}
}
func (m *Signal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Signal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Signal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Signal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Signal.Merge(m, src)
}
func (m *Signal) XXX_Size() int {
	return m.Size()
}
func (m *Signal) XXX_DiscardUnknown() {
	xxx_messageInfo_Signal.DiscardUnknown(m)
}

var xxx_messageInfo_Signal proto.InternalMessageInfo

func (m *Signal) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Signal) GetPower() int64 {
	if m != nil {
		return m.Power
	}
	return 0
}

// Vote is the data structure that contains array of signals of a voter.
type Vote struct {
	// voter is the address of the voter of this signals.
	Voter string `protobuf:"bytes,1,opt,name=voter,proto3"   json:"voter,omitempty"`
	// signals is a list of signals submit by the voter.
	Signals []Signal `protobuf:"bytes,2,rep,name=signals,proto3" json:"signals"`
}

func (m *Vote) Reset()         { *m = Vote{} }
func (m *Vote) String() string { return proto.CompactTextString(m) }
func (*Vote) ProtoMessage()    {}
func (*Vote) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{1}
}
func (m *Vote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Vote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Vote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Vote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vote.Merge(m, src)
}
func (m *Vote) XXX_Size() int {
	return m.Size()
}
func (m *Vote) XXX_DiscardUnknown() {
	xxx_messageInfo_Vote.DiscardUnknown(m)
}

var xxx_messageInfo_Vote proto.InternalMessageInfo

func (m *Vote) GetVoter() string {
	if m != nil {
		return m.Voter
	}
	return ""
}

func (m *Vote) GetSignals() []Signal {
	if m != nil {
		return m.Signals
	}
	return nil
}

// Feed is a structure that holds a signal id, its total power, and its calculated interval.
type Feed struct {
	// signal_id is the unique string that identifies the unit of feed.
	SignalID string `protobuf:"bytes,1,opt,name=signal_id,json=signalId,proto3" json:"signal_id,omitempty"`
	// power is the power of the corresponding signal id.
	Power int64 `protobuf:"varint,2,opt,name=power,proto3"                  json:"power,omitempty"`
	// interval is the interval of the price feed.
	Interval int64 `protobuf:"varint,3,opt,name=interval,proto3"               json:"interval,omitempty"`
}

func (m *Feed) Reset()         { *m = Feed{} }
func (m *Feed) String() string { return proto.CompactTextString(m) }
func (*Feed) ProtoMessage()    {}
func (*Feed) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{2}
}
func (m *Feed) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Feed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Feed.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Feed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Feed.Merge(m, src)
}
func (m *Feed) XXX_Size() int {
	return m.Size()
}
func (m *Feed) XXX_DiscardUnknown() {
	xxx_messageInfo_Feed.DiscardUnknown(m)
}

var xxx_messageInfo_Feed proto.InternalMessageInfo

func (m *Feed) GetSignalID() string {
	if m != nil {
		return m.SignalID
	}
	return ""
}

func (m *Feed) GetPower() int64 {
	if m != nil {
		return m.Power
	}
	return 0
}

func (m *Feed) GetInterval() int64 {
	if m != nil {
		return m.Interval
	}
	return 0
}

// FeedWithDeviation is a structure that holds a signal id, its total power, and its calculated interval and deviation.
type FeedWithDeviation struct {
	// signal_id is the unique string that identifies the unit of feed.
	SignalID string `protobuf:"bytes,1,opt,name=signal_id,json=signalId,proto3"                         json:"signal_id,omitempty"`
	// power is the power of the corresponding signal id.
	Power int64 `protobuf:"varint,2,opt,name=power,proto3"                                          json:"power,omitempty"`
	// interval is the interval of the price feed.
	Interval int64 `protobuf:"varint,3,opt,name=interval,proto3"                                       json:"interval,omitempty"`
	// deviation_basis_point is the maximum deviation value the feed can tolerate, expressed in basis points.
	DeviationBasisPoint int64 `protobuf:"varint,4,opt,name=deviation_basis_point,json=deviationBasisPoint,proto3" json:"deviation_basis_point,omitempty"`
}

func (m *FeedWithDeviation) Reset()         { *m = FeedWithDeviation{} }
func (m *FeedWithDeviation) String() string { return proto.CompactTextString(m) }
func (*FeedWithDeviation) ProtoMessage()    {}
func (*FeedWithDeviation) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{3}
}
func (m *FeedWithDeviation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FeedWithDeviation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FeedWithDeviation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FeedWithDeviation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeedWithDeviation.Merge(m, src)
}
func (m *FeedWithDeviation) XXX_Size() int {
	return m.Size()
}
func (m *FeedWithDeviation) XXX_DiscardUnknown() {
	xxx_messageInfo_FeedWithDeviation.DiscardUnknown(m)
}

var xxx_messageInfo_FeedWithDeviation proto.InternalMessageInfo

func (m *FeedWithDeviation) GetSignalID() string {
	if m != nil {
		return m.SignalID
	}
	return ""
}

func (m *FeedWithDeviation) GetPower() int64 {
	if m != nil {
		return m.Power
	}
	return 0
}

func (m *FeedWithDeviation) GetInterval() int64 {
	if m != nil {
		return m.Interval
	}
	return 0
}

func (m *FeedWithDeviation) GetDeviationBasisPoint() int64 {
	if m != nil {
		return m.DeviationBasisPoint
	}
	return 0
}

// CurrentFeeds is a structure that holds a list of currently supported feeds, and its last update time and block.
type CurrentFeeds struct {
	// feeds is a list of currently supported feeds.
	Feeds []Feed `protobuf:"bytes,1,rep,name=feeds,proto3"                                           json:"feeds"`
	// last_update_timestamp is the timestamp of the last time supported feeds list is updated.
	LastUpdateTimestamp int64 `protobuf:"varint,2,opt,name=last_update_timestamp,json=lastUpdateTimestamp,proto3" json:"last_update_timestamp,omitempty"`
	// last_update_block is the number of blocks of the last time supported feeds list is updated.
	LastUpdateBlock int64 `protobuf:"varint,3,opt,name=last_update_block,json=lastUpdateBlock,proto3"         json:"last_update_block,omitempty"`
}

func (m *CurrentFeeds) Reset()         { *m = CurrentFeeds{} }
func (m *CurrentFeeds) String() string { return proto.CompactTextString(m) }
func (*CurrentFeeds) ProtoMessage()    {}
func (*CurrentFeeds) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{4}
}
func (m *CurrentFeeds) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CurrentFeeds) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CurrentFeeds.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CurrentFeeds) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CurrentFeeds.Merge(m, src)
}
func (m *CurrentFeeds) XXX_Size() int {
	return m.Size()
}
func (m *CurrentFeeds) XXX_DiscardUnknown() {
	xxx_messageInfo_CurrentFeeds.DiscardUnknown(m)
}

var xxx_messageInfo_CurrentFeeds proto.InternalMessageInfo

func (m *CurrentFeeds) GetFeeds() []Feed {
	if m != nil {
		return m.Feeds
	}
	return nil
}

func (m *CurrentFeeds) GetLastUpdateTimestamp() int64 {
	if m != nil {
		return m.LastUpdateTimestamp
	}
	return 0
}

func (m *CurrentFeeds) GetLastUpdateBlock() int64 {
	if m != nil {
		return m.LastUpdateBlock
	}
	return 0
}

// CurrentFeedWithDeviations is a structure that holds a list of currently supported feed-with-deviations, and its
// last update time and block.
type CurrentFeedWithDeviations struct {
	// feeds is a list of currently supported feed-with-deviations.
	Feeds []FeedWithDeviation `protobuf:"bytes,1,rep,name=feeds,proto3"                                           json:"feeds"`
	// last_update_timestamp is the timestamp of the last time supported feeds list is updated.
	LastUpdateTimestamp int64 `protobuf:"varint,2,opt,name=last_update_timestamp,json=lastUpdateTimestamp,proto3" json:"last_update_timestamp,omitempty"`
	// last_update_block is the number of blocks of the last time supported feeds list is updated.
	LastUpdateBlock int64 `protobuf:"varint,3,opt,name=last_update_block,json=lastUpdateBlock,proto3"         json:"last_update_block,omitempty"`
}

func (m *CurrentFeedWithDeviations) Reset()         { *m = CurrentFeedWithDeviations{} }
func (m *CurrentFeedWithDeviations) String() string { return proto.CompactTextString(m) }
func (*CurrentFeedWithDeviations) ProtoMessage()    {}
func (*CurrentFeedWithDeviations) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{5}
}
func (m *CurrentFeedWithDeviations) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CurrentFeedWithDeviations) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CurrentFeedWithDeviations.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CurrentFeedWithDeviations) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CurrentFeedWithDeviations.Merge(m, src)
}
func (m *CurrentFeedWithDeviations) XXX_Size() int {
	return m.Size()
}
func (m *CurrentFeedWithDeviations) XXX_DiscardUnknown() {
	xxx_messageInfo_CurrentFeedWithDeviations.DiscardUnknown(m)
}

var xxx_messageInfo_CurrentFeedWithDeviations proto.InternalMessageInfo

func (m *CurrentFeedWithDeviations) GetFeeds() []FeedWithDeviation {
	if m != nil {
		return m.Feeds
	}
	return nil
}

func (m *CurrentFeedWithDeviations) GetLastUpdateTimestamp() int64 {
	if m != nil {
		return m.LastUpdateTimestamp
	}
	return 0
}

func (m *CurrentFeedWithDeviations) GetLastUpdateBlock() int64 {
	if m != nil {
		return m.LastUpdateBlock
	}
	return 0
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
func (*Price) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{6}
}
func (m *Price) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Price) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Price.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Price) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Price.Merge(m, src)
}
func (m *Price) XXX_Size() int {
	return m.Size()
}
func (m *Price) XXX_DiscardUnknown() {
	xxx_messageInfo_Price.DiscardUnknown(m)
}

var xxx_messageInfo_Price proto.InternalMessageInfo

func (m *Price) GetStatus() PriceStatus {
	if m != nil {
		return m.Status
	}
	return PRICE_STATUS_UNSPECIFIED
}

func (m *Price) GetSignalID() string {
	if m != nil {
		return m.SignalID
	}
	return ""
}

func (m *Price) GetPrice() uint64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *Price) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

// SignalPrice is a structure that defines the signaled price of a signal id.
type SignalPrice struct {
	// status is the status of the signal price.
	Status SignalPriceStatus `protobuf:"varint,1,opt,name=status,proto3,enum=band.feeds.v1beta1.SignalPriceStatus" json:"status,omitempty"`
	// signal_id is the signal id of the price.
	SignalID string `protobuf:"bytes,2,opt,name=signal_id,json=signalId,proto3"                           json:"signal_id,omitempty"`
	// price is the price submitted by the validator.
	Price uint64 `protobuf:"varint,3,opt,name=price,proto3"                                            json:"price,omitempty"`
}

func (m *SignalPrice) Reset()         { *m = SignalPrice{} }
func (m *SignalPrice) String() string { return proto.CompactTextString(m) }
func (*SignalPrice) ProtoMessage()    {}
func (*SignalPrice) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{7}
}
func (m *SignalPrice) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SignalPrice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SignalPrice.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SignalPrice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignalPrice.Merge(m, src)
}
func (m *SignalPrice) XXX_Size() int {
	return m.Size()
}
func (m *SignalPrice) XXX_DiscardUnknown() {
	xxx_messageInfo_SignalPrice.DiscardUnknown(m)
}

var xxx_messageInfo_SignalPrice proto.InternalMessageInfo

func (m *SignalPrice) GetStatus() SignalPriceStatus {
	if m != nil {
		return m.Status
	}
	return SIGNAL_PRICE_STATUS_UNSPECIFIED
}

func (m *SignalPrice) GetSignalID() string {
	if m != nil {
		return m.SignalID
	}
	return ""
}

func (m *SignalPrice) GetPrice() uint64 {
	if m != nil {
		return m.Price
	}
	return 0
}

// ValidatorPrice is a structure that defines the price submitted by a validator for a signal id.
type ValidatorPrice struct {
	// signal_price_status is the status of a signal price submitted.
	SignalPriceStatus SignalPriceStatus `protobuf:"varint,1,opt,name=signal_price_status,json=signalPriceStatus,proto3,enum=band.feeds.v1beta1.SignalPriceStatus" json:"signal_price_status,omitempty"`
	// signal_id is the signal id of the price.
	SignalID string `protobuf:"bytes,2,opt,name=signal_id,json=signalId,proto3"                                                               json:"signal_id,omitempty"`
	// price is the price submitted by the validator.
	Price uint64 `protobuf:"varint,3,opt,name=price,proto3"                                                                                json:"price,omitempty"`
	// timestamp is the timestamp at which the price was submitted.
	Timestamp int64 `protobuf:"varint,4,opt,name=timestamp,proto3"                                                                            json:"timestamp,omitempty"`
	// block_height is the block height at which the price was submitted.
	BlockHeight int64 `protobuf:"varint,5,opt,name=block_height,json=blockHeight,proto3"                                                        json:"block_height,omitempty"`
}

func (m *ValidatorPrice) Reset()         { *m = ValidatorPrice{} }
func (m *ValidatorPrice) String() string { return proto.CompactTextString(m) }
func (*ValidatorPrice) ProtoMessage()    {}
func (*ValidatorPrice) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{8}
}
func (m *ValidatorPrice) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ValidatorPrice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ValidatorPrice.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ValidatorPrice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidatorPrice.Merge(m, src)
}
func (m *ValidatorPrice) XXX_Size() int {
	return m.Size()
}
func (m *ValidatorPrice) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidatorPrice.DiscardUnknown(m)
}

var xxx_messageInfo_ValidatorPrice proto.InternalMessageInfo

func (m *ValidatorPrice) GetSignalPriceStatus() SignalPriceStatus {
	if m != nil {
		return m.SignalPriceStatus
	}
	return SIGNAL_PRICE_STATUS_UNSPECIFIED
}

func (m *ValidatorPrice) GetSignalID() string {
	if m != nil {
		return m.SignalID
	}
	return ""
}

func (m *ValidatorPrice) GetPrice() uint64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *ValidatorPrice) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *ValidatorPrice) GetBlockHeight() int64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

// ValidatorPriceList is a structure that holds a list of validator prices of
// a validator and its address.
type ValidatorPriceList struct {
	// validator is the validator address.
	Validator string `protobuf:"bytes,1,opt,name=validator,proto3"                             json:"validator,omitempty"`
	// validators_prices is a list of validator prices.
	ValidatorPrices []ValidatorPrice `protobuf:"bytes,2,rep,name=validator_prices,json=validatorPrices,proto3" json:"validator_prices"`
}

func (m *ValidatorPriceList) Reset()         { *m = ValidatorPriceList{} }
func (m *ValidatorPriceList) String() string { return proto.CompactTextString(m) }
func (*ValidatorPriceList) ProtoMessage()    {}
func (*ValidatorPriceList) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{9}
}
func (m *ValidatorPriceList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ValidatorPriceList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ValidatorPriceList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ValidatorPriceList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidatorPriceList.Merge(m, src)
}
func (m *ValidatorPriceList) XXX_Size() int {
	return m.Size()
}
func (m *ValidatorPriceList) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidatorPriceList.DiscardUnknown(m)
}

var xxx_messageInfo_ValidatorPriceList proto.InternalMessageInfo

func (m *ValidatorPriceList) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func (m *ValidatorPriceList) GetValidatorPrices() []ValidatorPrice {
	if m != nil {
		return m.ValidatorPrices
	}
	return nil
}

// ReferenceSourceConfig is a structure that defines the information of reference price source.
type ReferenceSourceConfig struct {
	// registry_ipfs_hash is the hash of the reference registry.
	RegistryIPFSHash string `protobuf:"bytes,1,opt,name=registry_ipfs_hash,json=registryIpfsHash,proto3" json:"registry_ipfs_hash,omitempty"`
	// registry_version is the version of the reference registry.
	RegistryVersion string `protobuf:"bytes,2,opt,name=registry_version,json=registryVersion,proto3"    json:"registry_version,omitempty"`
}

func (m *ReferenceSourceConfig) Reset()         { *m = ReferenceSourceConfig{} }
func (m *ReferenceSourceConfig) String() string { return proto.CompactTextString(m) }
func (*ReferenceSourceConfig) ProtoMessage()    {}
func (*ReferenceSourceConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{10}
}
func (m *ReferenceSourceConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReferenceSourceConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReferenceSourceConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReferenceSourceConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReferenceSourceConfig.Merge(m, src)
}
func (m *ReferenceSourceConfig) XXX_Size() int {
	return m.Size()
}
func (m *ReferenceSourceConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ReferenceSourceConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ReferenceSourceConfig proto.InternalMessageInfo

func (m *ReferenceSourceConfig) GetRegistryIPFSHash() string {
	if m != nil {
		return m.RegistryIPFSHash
	}
	return ""
}

func (m *ReferenceSourceConfig) GetRegistryVersion() string {
	if m != nil {
		return m.RegistryVersion
	}
	return ""
}

// FeedsSignatureOrder defines a general signature order for feed data.
type FeedsSignatureOrder struct {
	// signal_ids is the list of signal ids that require signatures.
	SignalIDs []string `protobuf:"bytes,1,rep,name=signal_ids,json=signalIds,proto3"                json:"signal_ids,omitempty"`
	// encoder is the mode of encoding feeds signature order.
	Encoder Encoder `protobuf:"varint,2,opt,name=encoder,proto3,enum=band.feeds.v1beta1.Encoder" json:"encoder,omitempty"`
}

func (m *FeedsSignatureOrder) Reset()         { *m = FeedsSignatureOrder{} }
func (m *FeedsSignatureOrder) String() string { return proto.CompactTextString(m) }
func (*FeedsSignatureOrder) ProtoMessage()    {}
func (*FeedsSignatureOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc3afe81d3b13674, []int{11}
}
func (m *FeedsSignatureOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FeedsSignatureOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FeedsSignatureOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FeedsSignatureOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeedsSignatureOrder.Merge(m, src)
}
func (m *FeedsSignatureOrder) XXX_Size() int {
	return m.Size()
}
func (m *FeedsSignatureOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_FeedsSignatureOrder.DiscardUnknown(m)
}

var xxx_messageInfo_FeedsSignatureOrder proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("band.feeds.v1beta1.PriceStatus", PriceStatus_name, PriceStatus_value)
	proto.RegisterEnum("band.feeds.v1beta1.SignalPriceStatus", SignalPriceStatus_name, SignalPriceStatus_value)
	proto.RegisterType((*Signal)(nil), "band.feeds.v1beta1.Signal")
	proto.RegisterType((*Vote)(nil), "band.feeds.v1beta1.Vote")
	proto.RegisterType((*Feed)(nil), "band.feeds.v1beta1.Feed")
	proto.RegisterType((*FeedWithDeviation)(nil), "band.feeds.v1beta1.FeedWithDeviation")
	proto.RegisterType((*CurrentFeeds)(nil), "band.feeds.v1beta1.CurrentFeeds")
	proto.RegisterType((*CurrentFeedWithDeviations)(nil), "band.feeds.v1beta1.CurrentFeedWithDeviations")
	proto.RegisterType((*Price)(nil), "band.feeds.v1beta1.Price")
	proto.RegisterType((*SignalPrice)(nil), "band.feeds.v1beta1.SignalPrice")
	proto.RegisterType((*ValidatorPrice)(nil), "band.feeds.v1beta1.ValidatorPrice")
	proto.RegisterType((*ValidatorPriceList)(nil), "band.feeds.v1beta1.ValidatorPriceList")
	proto.RegisterType((*ReferenceSourceConfig)(nil), "band.feeds.v1beta1.ReferenceSourceConfig")
	proto.RegisterType((*FeedsSignatureOrder)(nil), "band.feeds.v1beta1.FeedsSignatureOrder")
}

func init() { proto.RegisterFile("band/feeds/v1beta1/feeds.proto", fileDescriptor_fc3afe81d3b13674) }

var fileDescriptor_fc3afe81d3b13674 = []byte{
	// 989 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x56, 0x41, 0x6f, 0x1a, 0x47,
	0x14, 0x66, 0x31, 0x76, 0xcc, 0xb3, 0x6b, 0xe3, 0xb1, 0x13, 0x6d, 0x68, 0x02, 0xb6, 0x2b, 0x4b,
	0x8e, 0xd5, 0x80, 0xe2, 0xb4, 0xaa, 0x64, 0xb5, 0xaa, 0xc0, 0xac, 0xeb, 0x55, 0x2d, 0x8c, 0x76,
	0xc1, 0x51, 0x7b, 0x59, 0x2d, 0xec, 0x00, 0xa3, 0xe2, 0x5d, 0x34, 0x33, 0xd0, 0xe4, 0xd6, 0x63,
	0x0e, 0x3d, 0x54, 0xea, 0x1f, 0x88, 0x54, 0xf5, 0xd2, 0x5e, 0x7a, 0xf0, 0x4f, 0xe8, 0x21, 0xc7,
	0x28, 0xa7, 0x9e, 0x50, 0x85, 0x2f, 0xbd, 0xf5, 0x2f, 0x54, 0x3b, 0x33, 0x0b, 0xc6, 0x86, 0x56,
	0xaa, 0x64, 0xe5, 0xc6, 0x7c, 0xdf, 0xf7, 0xf6, 0x7d, 0xef, 0xcd, 0xdb, 0xc7, 0x42, 0xa6, 0xee,
	0xfa, 0x5e, 0xbe, 0x89, 0xb1, 0xc7, 0xf2, 0xfd, 0x27, 0x75, 0xcc, 0xdd, 0x27, 0xf2, 0x94, 0xeb,
	0xd2, 0x80, 0x07, 0x08, 0x85, 0x7c, 0x4e, 0x22, 0x8a, 0x4f, 0xdf, 0x6f, 0x04, 0xec, 0x3c, 0x60,
	0x8e, 0x50, 0xe4, 0xe5, 0x41, 0xca, 0xd3, 0x1b, 0xad, 0xa0, 0x15, 0x48, 0x3c, 0xfc, 0xa5, 0xd0,
	0xcd, 0x29, 0x49, 0xb0, 0xdf, 0x08, 0x3c, 0x4c, 0xa5, 0x62, 0xfb, 0x53, 0x58, 0xb0, 0x49, 0xcb,
	0x77, 0x3b, 0xe8, 0x1e, 0xc4, 0x89, 0xa7, 0x6b, 0x9b, 0xda, 0x6e, 0xb2, 0xb8, 0x30, 0x1c, 0x64,
	0xe3, 0x66, 0xc9, 0x8a, 0x13, 0x0f, 0x6d, 0xc0, 0x7c, 0x37, 0xf8, 0x16, 0x53, 0x3d, 0xbe, 0xa9,
	0xed, 0xce, 0x59, 0xf2, 0x70, 0x90, 0xf8, 0xeb, 0x55, 0x56, 0xdb, 0x7e, 0x0e, 0x89, 0xb3, 0x80,
	0x63, 0x94, 0x83, 0xf9, 0x7e, 0xc0, 0x31, 0x55, 0xe1, 0xfa, 0xdb, 0x8b, 0xc7, 0x1b, 0xca, 0x5e,
	0xc1, 0xf3, 0x28, 0x66, 0xcc, 0xe6, 0x94, 0xf8, 0x2d, 0x4b, 0xca, 0xd0, 0x01, 0xdc, 0x61, 0x22,
	0x2b, 0xd3, 0xe3, 0x9b, 0x73, 0xbb, 0x4b, 0xfb, 0xe9, 0xdc, 0xcd, 0x72, 0x73, 0xd2, 0x58, 0x31,
	0xf1, 0x7a, 0x90, 0x8d, 0x59, 0x51, 0x80, 0xca, 0x4c, 0x20, 0x71, 0x84, 0xb1, 0x87, 0x1e, 0x41,
	0x52, 0x12, 0xce, 0xc8, 0xfc, 0xf2, 0x70, 0x90, 0x5d, 0x94, 0xb1, 0x66, 0xc9, 0x5a, 0x94, 0xb4,
	0x39, 0xa3, 0x10, 0x94, 0x86, 0x45, 0xe2, 0x73, 0x4c, 0xfb, 0x6e, 0x47, 0x9f, 0x13, 0xc4, 0xe8,
	0xac, 0x52, 0xfd, 0xa2, 0xc1, 0x5a, 0x98, 0xeb, 0x19, 0xe1, 0xed, 0x12, 0xee, 0x13, 0x97, 0x93,
	0xc0, 0xbf, 0xd5, 0xc4, 0x68, 0x1f, 0xee, 0x7a, 0x51, 0x26, 0xa7, 0xee, 0x32, 0xc2, 0x9c, 0x6e,
	0x40, 0x7c, 0xae, 0x27, 0x84, 0x70, 0x7d, 0x44, 0x16, 0x43, 0xae, 0x12, 0x52, 0x63, 0xb3, 0xcb,
	0x87, 0x3d, 0x4a, 0xb1, 0xcf, 0x43, 0xcf, 0x0c, 0x7d, 0x04, 0xf3, 0xa2, 0xab, 0xba, 0x26, 0x1a,
	0xad, 0x4f, 0x6b, 0x74, 0xa8, 0x54, 0x6d, 0x96, 0xe2, 0xd0, 0x40, 0xc7, 0x65, 0xdc, 0xe9, 0x75,
	0x3d, 0x97, 0x63, 0x87, 0x93, 0x73, 0xcc, 0xb8, 0x7b, 0xde, 0x55, 0x25, 0xac, 0x87, 0x64, 0x4d,
	0x70, 0xd5, 0x88, 0x42, 0x7b, 0xb0, 0x76, 0x35, 0xa6, 0xde, 0x09, 0x1a, 0xdf, 0xa8, 0xca, 0x56,
	0xc7, 0xfa, 0x62, 0x08, 0x2b, 0xb3, 0xbf, 0x6b, 0x70, 0xff, 0x8a, 0xd9, 0x89, 0x06, 0x33, 0x54,
	0x98, 0x74, 0xbe, 0x33, 0xcb, 0xf9, 0x44, 0xd8, 0xbb, 0x28, 0xe3, 0x67, 0x0d, 0xe6, 0x2b, 0x94,
	0x34, 0x30, 0xfa, 0x04, 0x16, 0x18, 0x77, 0x79, 0x8f, 0x89, 0x89, 0x58, 0xd9, 0xcf, 0x4e, 0xf3,
	0x2c, 0xa4, 0xb6, 0x90, 0x59, 0x4a, 0x3e, 0x39, 0x4d, 0xf1, 0xff, 0x9c, 0xa6, 0xf0, 0x09, 0xc2,
	0x53, 0xc2, 0x92, 0x07, 0xf4, 0x00, 0x92, 0xe3, 0xea, 0xe4, 0x94, 0x8c, 0x01, 0xe5, 0xf3, 0x47,
	0x0d, 0x96, 0xe4, 0x03, 0xa5, 0xdb, 0xcf, 0xae, 0xb9, 0xdd, 0x99, 0xfd, 0x12, 0xde, 0x86, 0x67,
	0xe5, 0xea, 0x6f, 0x0d, 0x56, 0xce, 0xdc, 0x0e, 0xf1, 0x5c, 0x1e, 0x50, 0x69, 0xac, 0x06, 0xeb,
	0xea, 0xc9, 0x42, 0xe8, 0xfc, 0x1f, 0x97, 0x6b, 0xec, 0x3a, 0x74, 0xcb, 0x4d, 0x46, 0x5b, 0xb0,
	0x2c, 0x86, 0xc5, 0x69, 0x63, 0xd2, 0x6a, 0x73, 0x7d, 0x5e, 0x08, 0x96, 0x04, 0x76, 0x2c, 0x20,
	0x55, 0xf1, 0x6f, 0x1a, 0xa0, 0xc9, 0x8a, 0x4f, 0x08, 0xe3, 0xe8, 0x73, 0x48, 0xf6, 0x23, 0x54,
	0x6d, 0x94, 0xad, 0xb7, 0x17, 0x8f, 0x1f, 0xaa, 0x45, 0x3a, 0x8a, 0x98, 0xdc, 0xa8, 0xe3, 0x18,
	0x64, 0x43, 0x6a, 0x74, 0x90, 0x9d, 0x8b, 0xd6, 0xeb, 0xf6, 0xb4, 0x9e, 0x4d, 0x5a, 0x50, 0x2f,
	0xce, 0x6a, 0x7f, 0x02, 0x8d, 0xd6, 0xed, 0xf7, 0x1a, 0xdc, 0xb5, 0x70, 0x13, 0x53, 0xec, 0x37,
	0xb0, 0x1d, 0xf4, 0x68, 0x03, 0x1f, 0x06, 0x7e, 0x93, 0xb4, 0x50, 0x11, 0x10, 0xc5, 0x2d, 0xc2,
	0x38, 0x7d, 0xe1, 0x90, 0x6e, 0x93, 0x39, 0x6d, 0x97, 0xb5, 0x95, 0xfd, 0x8d, 0xe1, 0x20, 0x9b,
	0xb2, 0x14, 0x6b, 0x56, 0x8e, 0xec, 0x63, 0x97, 0xb5, 0xad, 0x54, 0xa4, 0x37, 0xbb, 0x4d, 0x16,
	0x22, 0xe8, 0x11, 0x8c, 0x30, 0xa7, 0x8f, 0x29, 0x23, 0x81, 0x2f, 0xef, 0xc7, 0x5a, 0x8d, 0xf0,
	0x33, 0x09, 0x2b, 0x3b, 0xdf, 0x69, 0xb0, 0x2e, 0xd6, 0x9b, 0xb8, 0x3a, 0xde, 0xa3, 0xf8, 0x94,
	0x7a, 0x98, 0xa2, 0x0f, 0x01, 0x46, 0x37, 0x2c, 0xf7, 0x46, 0xb2, 0xf8, 0xde, 0x70, 0x90, 0x4d,
	0x46, 0x57, 0xcc, 0xac, 0x64, 0x74, 0xc7, 0x0c, 0x7d, 0x0c, 0x77, 0xd4, 0x9f, 0xa1, 0xc8, 0xb6,
	0xb2, 0xff, 0xfe, 0xb4, 0x36, 0x19, 0x52, 0x62, 0x45, 0xda, 0x83, 0xc4, 0xcb, 0x57, 0xd9, 0xd8,
	0xde, 0x85, 0x06, 0x4b, 0x57, 0x87, 0xeb, 0x01, 0xe8, 0x15, 0xcb, 0x3c, 0x34, 0x1c, 0xbb, 0x5a,
	0xa8, 0xd6, 0x6c, 0xa7, 0x56, 0xb6, 0x2b, 0xc6, 0xa1, 0x79, 0x64, 0x1a, 0xa5, 0x54, 0x0c, 0x6d,
	0x43, 0xe6, 0x1a, 0xfb, 0x65, 0xf9, 0xf4, 0x59, 0xd9, 0xb1, 0xcd, 0x2f, 0xca, 0x85, 0x13, 0xc7,
	0x2c, 0xa5, 0x34, 0x94, 0x86, 0x7b, 0x13, 0x9a, 0xf2, 0x69, 0xd5, 0xb1, 0x8c, 0x42, 0xe9, 0xab,
	0x54, 0xfc, 0x06, 0x57, 0x38, 0x2b, 0x98, 0x27, 0x85, 0xe2, 0x89, 0x91, 0x9a, 0x43, 0x3b, 0xb0,
	0x75, 0x23, 0xce, 0x2c, 0x3b, 0x87, 0x35, 0xcb, 0x32, 0xca, 0x55, 0xe7, 0xc8, 0x30, 0x4a, 0x76,
	0x2a, 0x91, 0x4e, 0xbc, 0xfc, 0x29, 0x13, 0xdb, 0xfb, 0x55, 0x83, 0xb5, 0x1b, 0x2f, 0x0b, 0xfa,
	0x00, 0xb2, 0xca, 0xc9, 0xbf, 0xd4, 0x30, 0x5b, 0x54, 0xab, 0x54, 0x4e, 0xad, 0xaa, 0x11, 0x16,
	0x31, 0x53, 0x34, 0x76, 0x1c, 0x47, 0x5b, 0xf0, 0x70, 0x9a, 0xe8, 0x4a, 0x51, 0xd2, 0x6d, 0xf1,
	0xf8, 0xf5, 0x30, 0xa3, 0xbd, 0x19, 0x66, 0xb4, 0x3f, 0x87, 0x19, 0xed, 0x87, 0xcb, 0x4c, 0xec,
	0xcd, 0x65, 0x26, 0xf6, 0xc7, 0x65, 0x26, 0xf6, 0x75, 0xae, 0x45, 0x78, 0xbb, 0x57, 0xcf, 0x35,
	0x82, 0xf3, 0x7c, 0x78, 0x69, 0xe2, 0x6b, 0xa6, 0x11, 0x74, 0xf2, 0x8d, 0xb6, 0x4b, 0xfc, 0x7c,
	0xff, 0x69, 0xfe, 0xb9, 0xfa, 0xee, 0xe1, 0x2f, 0xba, 0x98, 0xd5, 0x17, 0x84, 0xe0, 0xe9, 0x3f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x4a, 0x51, 0x26, 0xa8, 0x77, 0x09, 0x00, 0x00,
}

func (this *Signal) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Signal)
	if !ok {
		that2, ok := that.(Signal)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.ID != that1.ID {
		return false
	}
	if this.Power != that1.Power {
		return false
	}
	return true
}
func (this *Vote) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Vote)
	if !ok {
		that2, ok := that.(Vote)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Voter != that1.Voter {
		return false
	}
	if len(this.Signals) != len(that1.Signals) {
		return false
	}
	for i := range this.Signals {
		if !this.Signals[i].Equal(&that1.Signals[i]) {
			return false
		}
	}
	return true
}
func (this *Feed) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Feed)
	if !ok {
		that2, ok := that.(Feed)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.SignalID != that1.SignalID {
		return false
	}
	if this.Power != that1.Power {
		return false
	}
	if this.Interval != that1.Interval {
		return false
	}
	return true
}
func (this *FeedWithDeviation) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*FeedWithDeviation)
	if !ok {
		that2, ok := that.(FeedWithDeviation)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.SignalID != that1.SignalID {
		return false
	}
	if this.Power != that1.Power {
		return false
	}
	if this.Interval != that1.Interval {
		return false
	}
	if this.DeviationBasisPoint != that1.DeviationBasisPoint {
		return false
	}
	return true
}
func (this *CurrentFeeds) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*CurrentFeeds)
	if !ok {
		that2, ok := that.(CurrentFeeds)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Feeds) != len(that1.Feeds) {
		return false
	}
	for i := range this.Feeds {
		if !this.Feeds[i].Equal(&that1.Feeds[i]) {
			return false
		}
	}
	if this.LastUpdateTimestamp != that1.LastUpdateTimestamp {
		return false
	}
	if this.LastUpdateBlock != that1.LastUpdateBlock {
		return false
	}
	return true
}
func (this *CurrentFeedWithDeviations) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*CurrentFeedWithDeviations)
	if !ok {
		that2, ok := that.(CurrentFeedWithDeviations)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Feeds) != len(that1.Feeds) {
		return false
	}
	for i := range this.Feeds {
		if !this.Feeds[i].Equal(&that1.Feeds[i]) {
			return false
		}
	}
	if this.LastUpdateTimestamp != that1.LastUpdateTimestamp {
		return false
	}
	if this.LastUpdateBlock != that1.LastUpdateBlock {
		return false
	}
	return true
}
func (this *Price) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Price)
	if !ok {
		that2, ok := that.(Price)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Status != that1.Status {
		return false
	}
	if this.SignalID != that1.SignalID {
		return false
	}
	if this.Price != that1.Price {
		return false
	}
	if this.Timestamp != that1.Timestamp {
		return false
	}
	return true
}
func (this *SignalPrice) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SignalPrice)
	if !ok {
		that2, ok := that.(SignalPrice)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Status != that1.Status {
		return false
	}
	if this.SignalID != that1.SignalID {
		return false
	}
	if this.Price != that1.Price {
		return false
	}
	return true
}
func (this *ValidatorPrice) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ValidatorPrice)
	if !ok {
		that2, ok := that.(ValidatorPrice)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.SignalPriceStatus != that1.SignalPriceStatus {
		return false
	}
	if this.SignalID != that1.SignalID {
		return false
	}
	if this.Price != that1.Price {
		return false
	}
	if this.Timestamp != that1.Timestamp {
		return false
	}
	if this.BlockHeight != that1.BlockHeight {
		return false
	}
	return true
}
func (this *ValidatorPriceList) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ValidatorPriceList)
	if !ok {
		that2, ok := that.(ValidatorPriceList)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Validator != that1.Validator {
		return false
	}
	if len(this.ValidatorPrices) != len(that1.ValidatorPrices) {
		return false
	}
	for i := range this.ValidatorPrices {
		if !this.ValidatorPrices[i].Equal(&that1.ValidatorPrices[i]) {
			return false
		}
	}
	return true
}
func (this *ReferenceSourceConfig) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ReferenceSourceConfig)
	if !ok {
		that2, ok := that.(ReferenceSourceConfig)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.RegistryIPFSHash != that1.RegistryIPFSHash {
		return false
	}
	if this.RegistryVersion != that1.RegistryVersion {
		return false
	}
	return true
}
func (m *Signal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Signal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Signal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Power != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x10
	}
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Vote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Vote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Vote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signals) > 0 {
		for iNdEx := len(m.Signals) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Signals[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintFeeds(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Voter) > 0 {
		i -= len(m.Voter)
		copy(dAtA[i:], m.Voter)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.Voter)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Feed) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Feed) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Feed) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Interval != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Interval))
		i--
		dAtA[i] = 0x18
	}
	if m.Power != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x10
	}
	if len(m.SignalID) > 0 {
		i -= len(m.SignalID)
		copy(dAtA[i:], m.SignalID)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.SignalID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *FeedWithDeviation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FeedWithDeviation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FeedWithDeviation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.DeviationBasisPoint != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.DeviationBasisPoint))
		i--
		dAtA[i] = 0x20
	}
	if m.Interval != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Interval))
		i--
		dAtA[i] = 0x18
	}
	if m.Power != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x10
	}
	if len(m.SignalID) > 0 {
		i -= len(m.SignalID)
		copy(dAtA[i:], m.SignalID)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.SignalID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *CurrentFeeds) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CurrentFeeds) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CurrentFeeds) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastUpdateBlock != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.LastUpdateBlock))
		i--
		dAtA[i] = 0x18
	}
	if m.LastUpdateTimestamp != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.LastUpdateTimestamp))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Feeds) > 0 {
		for iNdEx := len(m.Feeds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Feeds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintFeeds(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *CurrentFeedWithDeviations) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CurrentFeedWithDeviations) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CurrentFeedWithDeviations) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastUpdateBlock != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.LastUpdateBlock))
		i--
		dAtA[i] = 0x18
	}
	if m.LastUpdateTimestamp != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.LastUpdateTimestamp))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Feeds) > 0 {
		for iNdEx := len(m.Feeds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Feeds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintFeeds(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Price) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Price) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Price) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Timestamp != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x20
	}
	if m.Price != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Price))
		i--
		dAtA[i] = 0x18
	}
	if len(m.SignalID) > 0 {
		i -= len(m.SignalID)
		copy(dAtA[i:], m.SignalID)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.SignalID)))
		i--
		dAtA[i] = 0x12
	}
	if m.Status != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SignalPrice) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SignalPrice) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignalPrice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Price != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Price))
		i--
		dAtA[i] = 0x18
	}
	if len(m.SignalID) > 0 {
		i -= len(m.SignalID)
		copy(dAtA[i:], m.SignalID)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.SignalID)))
		i--
		dAtA[i] = 0x12
	}
	if m.Status != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ValidatorPrice) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ValidatorPrice) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ValidatorPrice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BlockHeight != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x28
	}
	if m.Timestamp != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x20
	}
	if m.Price != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Price))
		i--
		dAtA[i] = 0x18
	}
	if len(m.SignalID) > 0 {
		i -= len(m.SignalID)
		copy(dAtA[i:], m.SignalID)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.SignalID)))
		i--
		dAtA[i] = 0x12
	}
	if m.SignalPriceStatus != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.SignalPriceStatus))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ValidatorPriceList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ValidatorPriceList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ValidatorPriceList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ValidatorPrices) > 0 {
		for iNdEx := len(m.ValidatorPrices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ValidatorPrices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintFeeds(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ReferenceSourceConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReferenceSourceConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReferenceSourceConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RegistryVersion) > 0 {
		i -= len(m.RegistryVersion)
		copy(dAtA[i:], m.RegistryVersion)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.RegistryVersion)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.RegistryIPFSHash) > 0 {
		i -= len(m.RegistryIPFSHash)
		copy(dAtA[i:], m.RegistryIPFSHash)
		i = encodeVarintFeeds(dAtA, i, uint64(len(m.RegistryIPFSHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *FeedsSignatureOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FeedsSignatureOrder) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FeedsSignatureOrder) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Encoder != 0 {
		i = encodeVarintFeeds(dAtA, i, uint64(m.Encoder))
		i--
		dAtA[i] = 0x10
	}
	if len(m.SignalIDs) > 0 {
		for iNdEx := len(m.SignalIDs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.SignalIDs[iNdEx])
			copy(dAtA[i:], m.SignalIDs[iNdEx])
			i = encodeVarintFeeds(dAtA, i, uint64(len(m.SignalIDs[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintFeeds(dAtA []byte, offset int, v uint64) int {
	offset -= sovFeeds(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Signal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	if m.Power != 0 {
		n += 1 + sovFeeds(uint64(m.Power))
	}
	return n
}

func (m *Vote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Voter)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	if len(m.Signals) > 0 {
		for _, e := range m.Signals {
			l = e.Size()
			n += 1 + l + sovFeeds(uint64(l))
		}
	}
	return n
}

func (m *Feed) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SignalID)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	if m.Power != 0 {
		n += 1 + sovFeeds(uint64(m.Power))
	}
	if m.Interval != 0 {
		n += 1 + sovFeeds(uint64(m.Interval))
	}
	return n
}

func (m *FeedWithDeviation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SignalID)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	if m.Power != 0 {
		n += 1 + sovFeeds(uint64(m.Power))
	}
	if m.Interval != 0 {
		n += 1 + sovFeeds(uint64(m.Interval))
	}
	if m.DeviationBasisPoint != 0 {
		n += 1 + sovFeeds(uint64(m.DeviationBasisPoint))
	}
	return n
}

func (m *CurrentFeeds) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Feeds) > 0 {
		for _, e := range m.Feeds {
			l = e.Size()
			n += 1 + l + sovFeeds(uint64(l))
		}
	}
	if m.LastUpdateTimestamp != 0 {
		n += 1 + sovFeeds(uint64(m.LastUpdateTimestamp))
	}
	if m.LastUpdateBlock != 0 {
		n += 1 + sovFeeds(uint64(m.LastUpdateBlock))
	}
	return n
}

func (m *CurrentFeedWithDeviations) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Feeds) > 0 {
		for _, e := range m.Feeds {
			l = e.Size()
			n += 1 + l + sovFeeds(uint64(l))
		}
	}
	if m.LastUpdateTimestamp != 0 {
		n += 1 + sovFeeds(uint64(m.LastUpdateTimestamp))
	}
	if m.LastUpdateBlock != 0 {
		n += 1 + sovFeeds(uint64(m.LastUpdateBlock))
	}
	return n
}

func (m *Price) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != 0 {
		n += 1 + sovFeeds(uint64(m.Status))
	}
	l = len(m.SignalID)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	if m.Price != 0 {
		n += 1 + sovFeeds(uint64(m.Price))
	}
	if m.Timestamp != 0 {
		n += 1 + sovFeeds(uint64(m.Timestamp))
	}
	return n
}

func (m *SignalPrice) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != 0 {
		n += 1 + sovFeeds(uint64(m.Status))
	}
	l = len(m.SignalID)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	if m.Price != 0 {
		n += 1 + sovFeeds(uint64(m.Price))
	}
	return n
}

func (m *ValidatorPrice) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SignalPriceStatus != 0 {
		n += 1 + sovFeeds(uint64(m.SignalPriceStatus))
	}
	l = len(m.SignalID)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	if m.Price != 0 {
		n += 1 + sovFeeds(uint64(m.Price))
	}
	if m.Timestamp != 0 {
		n += 1 + sovFeeds(uint64(m.Timestamp))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovFeeds(uint64(m.BlockHeight))
	}
	return n
}

func (m *ValidatorPriceList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	if len(m.ValidatorPrices) > 0 {
		for _, e := range m.ValidatorPrices {
			l = e.Size()
			n += 1 + l + sovFeeds(uint64(l))
		}
	}
	return n
}

func (m *ReferenceSourceConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RegistryIPFSHash)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	l = len(m.RegistryVersion)
	if l > 0 {
		n += 1 + l + sovFeeds(uint64(l))
	}
	return n
}

func (m *FeedsSignatureOrder) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SignalIDs) > 0 {
		for _, s := range m.SignalIDs {
			l = len(s)
			n += 1 + l + sovFeeds(uint64(l))
		}
	}
	if m.Encoder != 0 {
		n += 1 + sovFeeds(uint64(m.Encoder))
	}
	return n
}

func sovFeeds(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFeeds(x uint64) (n int) {
	return sovFeeds(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Signal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Signal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Signal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Vote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Vote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Voter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Voter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signals", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signals = append(m.Signals, Signal{})
			if err := m.Signals[len(m.Signals)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Feed) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Feed: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Feed: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignalID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Interval", wireType)
			}
			m.Interval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Interval |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FeedWithDeviation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FeedWithDeviation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FeedWithDeviation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignalID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Interval", wireType)
			}
			m.Interval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Interval |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeviationBasisPoint", wireType)
			}
			m.DeviationBasisPoint = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DeviationBasisPoint |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CurrentFeeds) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CurrentFeeds: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CurrentFeeds: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feeds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feeds = append(m.Feeds, Feed{})
			if err := m.Feeds[len(m.Feeds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastUpdateTimestamp", wireType)
			}
			m.LastUpdateTimestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastUpdateTimestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastUpdateBlock", wireType)
			}
			m.LastUpdateBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastUpdateBlock |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *CurrentFeedWithDeviations) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CurrentFeedWithDeviations: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CurrentFeedWithDeviations: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feeds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feeds = append(m.Feeds, FeedWithDeviation{})
			if err := m.Feeds[len(m.Feeds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastUpdateTimestamp", wireType)
			}
			m.LastUpdateTimestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastUpdateTimestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastUpdateBlock", wireType)
			}
			m.LastUpdateBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastUpdateBlock |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Price) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Price: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Price: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= PriceStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignalID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			m.Price = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Price |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SignalPrice) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SignalPrice: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SignalPrice: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= SignalPriceStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignalID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			m.Price = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Price |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ValidatorPrice) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ValidatorPrice: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ValidatorPrice: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalPriceStatus", wireType)
			}
			m.SignalPriceStatus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignalPriceStatus |= SignalPriceStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignalID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			m.Price = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Price |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ValidatorPriceList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ValidatorPriceList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ValidatorPriceList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorPrices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorPrices = append(m.ValidatorPrices, ValidatorPrice{})
			if err := m.ValidatorPrices[len(m.ValidatorPrices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ReferenceSourceConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ReferenceSourceConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReferenceSourceConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegistryIPFSHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RegistryIPFSHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RegistryVersion", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RegistryVersion = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FeedsSignatureOrder) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FeedsSignatureOrder: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FeedsSignatureOrder: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalIDs", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthFeeds
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFeeds
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignalIDs = append(m.SignalIDs, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Encoder", wireType)
			}
			m.Encoder = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Encoder |= Encoder(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFeeds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFeeds
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipFeeds(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFeeds
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFeeds
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthFeeds
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFeeds
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFeeds
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFeeds        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFeeds          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFeeds = fmt.Errorf("proto: unexpected end of group")
)
