package types

import (
	"fmt"
	"io"
	"math"
	math_bits "math/bits"

	"github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types1 "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	types2 "github.com/bandprotocol/falcon/internal/bandchain/feeds"
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

// Tunnel contains the information of the tunnel that is created by the user
type Tunnel struct {
	// id is the tunnel ID
	ID uint64 `protobuf:"varint,1,opt,name=id,proto3"                                                                                   json:"id,omitempty"`
	// sequence is representing the sequence of the tunnel packet.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence,proto3"                                                                             json:"sequence,omitempty"`
	// route is the route for delivering the signal prices
	Route *types.Any `protobuf:"bytes,3,opt,name=route,proto3"                                                                                 json:"route,omitempty"`
	// fee_payer is the address of the fee payer
	FeePayer string `protobuf:"bytes,4,opt,name=fee_payer,json=feePayer,proto3"                                                               json:"fee_payer,omitempty"`
	// signal_deviations is the list of signal deviations
	SignalDeviations []SignalDeviation `protobuf:"bytes,5,rep,name=signal_deviations,json=signalDeviations,proto3"                                               json:"signal_deviations"`
	// interval is the interval for delivering the signal prices
	Interval uint64 `protobuf:"varint,6,opt,name=interval,proto3"                                                                             json:"interval,omitempty"`
	// total_deposit is the total deposit on the tunnel.
	TotalDeposit github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,7,rep,name=total_deposit,json=totalDeposit,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_deposit"`
	// is_active is the flag to indicate if the tunnel is active
	IsActive bool `protobuf:"varint,8,opt,name=is_active,json=isActive,proto3"                                                              json:"is_active,omitempty"`
	// created_at is the timestamp when the tunnel is created
	CreatedAt int64 `protobuf:"varint,9,opt,name=created_at,json=createdAt,proto3"                                                            json:"created_at,omitempty"`
	// creator is the address of the creator
	Creator string `protobuf:"bytes,10,opt,name=creator,proto3"                                                                              json:"creator,omitempty"`
}

func (m *Tunnel) Reset()         { *m = Tunnel{} }
func (m *Tunnel) String() string { return proto.CompactTextString(m) }
func (*Tunnel) ProtoMessage()    {}
func (*Tunnel) Descriptor() ([]byte, []int) {
	return fileDescriptor_6bb6151451ba2f25, []int{0}
}
func (m *Tunnel) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Tunnel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Tunnel.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Tunnel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tunnel.Merge(m, src)
}
func (m *Tunnel) XXX_Size() int {
	return m.Size()
}
func (m *Tunnel) XXX_DiscardUnknown() {
	xxx_messageInfo_Tunnel.DiscardUnknown(m)
}

var xxx_messageInfo_Tunnel proto.InternalMessageInfo

func (m *Tunnel) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Tunnel) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *Tunnel) GetRoute() *types.Any {
	if m != nil {
		return m.Route
	}
	return nil
}

func (m *Tunnel) GetFeePayer() string {
	if m != nil {
		return m.FeePayer
	}
	return ""
}

func (m *Tunnel) GetSignalDeviations() []SignalDeviation {
	if m != nil {
		return m.SignalDeviations
	}
	return nil
}

func (m *Tunnel) GetInterval() uint64 {
	if m != nil {
		return m.Interval
	}
	return 0
}

func (m *Tunnel) GetTotalDeposit() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.TotalDeposit
	}
	return nil
}

func (m *Tunnel) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

func (m *Tunnel) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func (m *Tunnel) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

// LatestPrices is the type for prices that tunnel produces
type LatestPrices struct {
	// tunnel_id is the tunnel ID
	TunnelID uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3"         json:"tunnel_id,omitempty"`
	// prices is the list of prices information from feeds module.
	Prices []types2.Price `protobuf:"bytes,2,rep,name=prices,proto3"                           json:"prices"`
	// last_interval is the last interval when the signal prices are produced by interval trigger
	LastInterval int64 `protobuf:"varint,3,opt,name=last_interval,json=lastInterval,proto3" json:"last_interval,omitempty"`
}

func (m *LatestPrices) Reset()         { *m = LatestPrices{} }
func (m *LatestPrices) String() string { return proto.CompactTextString(m) }
func (*LatestPrices) ProtoMessage()    {}
func (*LatestPrices) Descriptor() ([]byte, []int) {
	return fileDescriptor_6bb6151451ba2f25, []int{1}
}
func (m *LatestPrices) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LatestPrices) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LatestPrices.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LatestPrices) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LatestPrices.Merge(m, src)
}
func (m *LatestPrices) XXX_Size() int {
	return m.Size()
}
func (m *LatestPrices) XXX_DiscardUnknown() {
	xxx_messageInfo_LatestPrices.DiscardUnknown(m)
}

var xxx_messageInfo_LatestPrices proto.InternalMessageInfo

func (m *LatestPrices) GetTunnelID() uint64 {
	if m != nil {
		return m.TunnelID
	}
	return 0
}

func (m *LatestPrices) GetPrices() []types2.Price {
	if m != nil {
		return m.Prices
	}
	return nil
}

func (m *LatestPrices) GetLastInterval() int64 {
	if m != nil {
		return m.LastInterval
	}
	return 0
}

// TotalFees is the type for the total fees collected by the tunnel
type TotalFees struct {
	// total_base_packet_fee is the total base packet fee collected by the tunnel
	TotalBasePacketFee github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=total_base_packet_fee,json=totalBasePacketFee,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_base_packet_fee"`
}

func (m *TotalFees) Reset()         { *m = TotalFees{} }
func (m *TotalFees) String() string { return proto.CompactTextString(m) }
func (*TotalFees) ProtoMessage()    {}
func (*TotalFees) Descriptor() ([]byte, []int) {
	return fileDescriptor_6bb6151451ba2f25, []int{2}
}
func (m *TotalFees) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TotalFees) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TotalFees.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TotalFees) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TotalFees.Merge(m, src)
}
func (m *TotalFees) XXX_Size() int {
	return m.Size()
}
func (m *TotalFees) XXX_DiscardUnknown() {
	xxx_messageInfo_TotalFees.DiscardUnknown(m)
}

var xxx_messageInfo_TotalFees proto.InternalMessageInfo

func (m *TotalFees) GetTotalBasePacketFee() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.TotalBasePacketFee
	}
	return nil
}

// Packet is the packet that tunnel produces
type Packet struct {
	// tunnel_id is the tunnel ID
	TunnelID uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3"   json:"tunnel_id,omitempty"`
	// sequence is representing the sequence of the tunnel packet.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence,proto3"                  json:"sequence,omitempty"`
	// prices is the list of prices information from feeds module.
	Prices []types2.Price `protobuf:"bytes,3,rep,name=prices,proto3"                     json:"prices"`
	// receipt represents the confirmation of the packet delivery to the destination via the specified route.
	Receipt *types.Any `protobuf:"bytes,4,opt,name=receipt,proto3"                    json:"receipt,omitempty"`
	// created_at is the timestamp when the packet is created
	CreatedAt int64 `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (m *Packet) Reset()         { *m = Packet{} }
func (m *Packet) String() string { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()    {}
func (*Packet) Descriptor() ([]byte, []int) {
	return fileDescriptor_6bb6151451ba2f25, []int{3}
}
func (m *Packet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Packet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Packet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Packet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet.Merge(m, src)
}
func (m *Packet) XXX_Size() int {
	return m.Size()
}
func (m *Packet) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet.DiscardUnknown(m)
}

var xxx_messageInfo_Packet proto.InternalMessageInfo

func (m *Packet) GetTunnelID() uint64 {
	if m != nil {
		return m.TunnelID
	}
	return 0
}

func (m *Packet) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *Packet) GetPrices() []types2.Price {
	if m != nil {
		return m.Prices
	}
	return nil
}

func (m *Packet) GetReceipt() *types.Any {
	if m != nil {
		return m.Receipt
	}
	return nil
}

func (m *Packet) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

// Deposit defines an amount deposited by an account address to the tunnel.
type Deposit struct {
	// tunnel_id defines the unique id of the tunnel.
	TunnelID uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3"                                     json:"tunnel_id,omitempty"`
	// depositor defines the deposit addresses from the proposals.
	Depositor string `protobuf:"bytes,2,opt,name=depositor,proto3"                                                    json:"depositor,omitempty"`
	// amount to be deposited by depositor.
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,3,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

func (m *Deposit) Reset()         { *m = Deposit{} }
func (m *Deposit) String() string { return proto.CompactTextString(m) }
func (*Deposit) ProtoMessage()    {}
func (*Deposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_6bb6151451ba2f25, []int{4}
}
func (m *Deposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Deposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Deposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Deposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deposit.Merge(m, src)
}
func (m *Deposit) XXX_Size() int {
	return m.Size()
}
func (m *Deposit) XXX_DiscardUnknown() {
	xxx_messageInfo_Deposit.DiscardUnknown(m)
}

var xxx_messageInfo_Deposit proto.InternalMessageInfo

func (m *Deposit) GetTunnelID() uint64 {
	if m != nil {
		return m.TunnelID
	}
	return 0
}

func (m *Deposit) GetDepositor() string {
	if m != nil {
		return m.Depositor
	}
	return ""
}

func (m *Deposit) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

// SignalDeviation is the type for a signal with soft and hard deviation
type SignalDeviation struct {
	// signal_id is the signal ID
	SignalID string `protobuf:"bytes,1,opt,name=signal_id,json=signalId,proto3"                   json:"signal_id,omitempty"`
	// soft_deviation_bps is the soft deviation in basis points
	SoftDeviationBPS uint64 `protobuf:"varint,2,opt,name=soft_deviation_bps,json=softDeviationBps,proto3" json:"soft_deviation_bps,omitempty"`
	// hard_deviation_bps is the hard deviation in basis points
	HardDeviationBPS uint64 `protobuf:"varint,3,opt,name=hard_deviation_bps,json=hardDeviationBps,proto3" json:"hard_deviation_bps,omitempty"`
}

func (m *SignalDeviation) Reset()         { *m = SignalDeviation{} }
func (m *SignalDeviation) String() string { return proto.CompactTextString(m) }
func (*SignalDeviation) ProtoMessage()    {}
func (*SignalDeviation) Descriptor() ([]byte, []int) {
	return fileDescriptor_6bb6151451ba2f25, []int{5}
}
func (m *SignalDeviation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SignalDeviation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SignalDeviation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SignalDeviation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignalDeviation.Merge(m, src)
}
func (m *SignalDeviation) XXX_Size() int {
	return m.Size()
}
func (m *SignalDeviation) XXX_DiscardUnknown() {
	xxx_messageInfo_SignalDeviation.DiscardUnknown(m)
}

var xxx_messageInfo_SignalDeviation proto.InternalMessageInfo

func (m *SignalDeviation) GetSignalID() string {
	if m != nil {
		return m.SignalID
	}
	return ""
}

func (m *SignalDeviation) GetSoftDeviationBPS() uint64 {
	if m != nil {
		return m.SoftDeviationBPS
	}
	return 0
}

func (m *SignalDeviation) GetHardDeviationBPS() uint64 {
	if m != nil {
		return m.HardDeviationBPS
	}
	return 0
}

// TunnelSignatureOrder defines a general signature order for sending signature to tss group.
type TunnelSignatureOrder struct {
	// sequence is the sequence of the packet
	Sequence uint64 `protobuf:"varint,1,opt,name=sequence,proto3"                                json:"sequence,omitempty"`
	// prices is the list of prices information from feeds module.
	Prices []types2.Price `protobuf:"bytes,2,rep,name=prices,proto3"                                   json:"prices"`
	// created_at is the timestamp when the packet is created
	CreatedAt int64 `protobuf:"varint,3,opt,name=created_at,json=createdAt,proto3"               json:"created_at,omitempty"`
	// encoder is the mode of encoding data.
	Encoder types2.Encoder `protobuf:"varint,4,opt,name=encoder,proto3,enum=band.feeds.v1beta1.Encoder" json:"encoder,omitempty"`
}

func (m *TunnelSignatureOrder) Reset()         { *m = TunnelSignatureOrder{} }
func (m *TunnelSignatureOrder) String() string { return proto.CompactTextString(m) }
func (*TunnelSignatureOrder) ProtoMessage()    {}
func (*TunnelSignatureOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_6bb6151451ba2f25, []int{6}
}
func (m *TunnelSignatureOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TunnelSignatureOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TunnelSignatureOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TunnelSignatureOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TunnelSignatureOrder.Merge(m, src)
}
func (m *TunnelSignatureOrder) XXX_Size() int {
	return m.Size()
}
func (m *TunnelSignatureOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_TunnelSignatureOrder.DiscardUnknown(m)
}

var xxx_messageInfo_TunnelSignatureOrder proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Tunnel)(nil), "band.tunnel.v1beta1.Tunnel")
	proto.RegisterType((*LatestPrices)(nil), "band.tunnel.v1beta1.LatestPrices")
	proto.RegisterType((*TotalFees)(nil), "band.tunnel.v1beta1.TotalFees")
	proto.RegisterType((*Packet)(nil), "band.tunnel.v1beta1.Packet")
	proto.RegisterType((*Deposit)(nil), "band.tunnel.v1beta1.Deposit")
	proto.RegisterType((*SignalDeviation)(nil), "band.tunnel.v1beta1.SignalDeviation")
	proto.RegisterType((*TunnelSignatureOrder)(nil), "band.tunnel.v1beta1.TunnelSignatureOrder")
}

func init() { proto.RegisterFile("band/tunnel/v1beta1/tunnel.proto", fileDescriptor_6bb6151451ba2f25) }

var fileDescriptor_6bb6151451ba2f25 = []byte{
	// 898 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0x4f, 0x6f, 0x1b, 0x45,
	0x14, 0xf7, 0xd8, 0x8e, 0xed, 0x9d, 0xa6, 0x25, 0x1d, 0x0c, 0xda, 0xa6, 0xc2, 0xb6, 0x02, 0x07,
	0x83, 0x94, 0x5d, 0x9a, 0x2a, 0x45, 0xea, 0x2d, 0x4b, 0xa8, 0xb0, 0x04, 0x22, 0xda, 0x54, 0x42,
	0xe2, 0xb2, 0x1a, 0xef, 0x3e, 0xdb, 0xa3, 0x3a, 0x3b, 0xcb, 0xcc, 0xd8, 0x22, 0x17, 0xce, 0x5c,
	0x90, 0xe0, 0x03, 0x20, 0xf5, 0x58, 0x71, 0xe2, 0x90, 0xef, 0x40, 0x95, 0x53, 0xc5, 0x89, 0x03,
	0x0a, 0xc8, 0x39, 0x80, 0xc4, 0x97, 0x40, 0x3b, 0x33, 0x6b, 0x63, 0x2b, 0x6d, 0x64, 0xa9, 0x97,
	0xc4, 0xef, 0xbd, 0xdf, 0xbc, 0x7f, 0xbf, 0xdf, 0xcc, 0xe2, 0x4e, 0x9f, 0xa6, 0x89, 0xaf, 0x26,
	0x69, 0x0a, 0x63, 0x7f, 0x7a, 0xaf, 0x0f, 0x8a, 0xde, 0xb3, 0xa6, 0x97, 0x09, 0xae, 0x38, 0x79,
	0x33, 0x47, 0x78, 0xd6, 0x65, 0x11, 0xdb, 0xb7, 0xe9, 0x09, 0x4b, 0xb9, 0xaf, 0xff, 0x1a, 0xdc,
	0x76, 0x2b, 0xe6, 0xf2, 0x84, 0x4b, 0xbf, 0x4f, 0x25, 0xcc, 0x33, 0xc5, 0x9c, 0xa5, 0x36, 0x7e,
	0xc7, 0xc4, 0x23, 0x6d, 0xf9, 0xc6, 0xb0, 0xa1, 0xe6, 0x90, 0x0f, 0xb9, 0xf1, 0xe7, 0xbf, 0x8a,
	0x03, 0x43, 0xce, 0x87, 0x63, 0xf0, 0xb5, 0xd5, 0x9f, 0x0c, 0x7c, 0x9a, 0x9e, 0xda, 0x90, 0xe9,
	0x7a, 0x00, 0x90, 0xc8, 0x79, 0x29, 0x48, 0x63, 0x9e, 0x80, 0x28, 0xba, 0xb9, 0x02, 0xa1, 0x2d,
	0x13, 0xdf, 0xf9, 0xbe, 0x8a, 0x6b, 0x8f, 0xf5, 0x4c, 0xe4, 0x6d, 0x5c, 0x66, 0x89, 0x8b, 0x3a,
	0xa8, 0x5b, 0x0d, 0x6a, 0xb3, 0x8b, 0x76, 0xb9, 0x77, 0x18, 0x96, 0x59, 0x42, 0xb6, 0x71, 0x43,
	0xc2, 0xd7, 0x13, 0x48, 0x63, 0x70, 0xcb, 0x79, 0x34, 0x9c, 0xdb, 0xe4, 0x01, 0xde, 0x10, 0x7c,
	0xa2, 0xc0, 0xad, 0x74, 0x50, 0xf7, 0xc6, 0x5e, 0xd3, 0x33, 0xbd, 0x7a, 0x45, 0xaf, 0xde, 0x41,
	0x7a, 0x1a, 0xe0, 0xf3, 0xb3, 0xdd, 0x5a, 0x98, 0xc3, 0x7a, 0xa1, 0x81, 0x93, 0x7d, 0xec, 0x0c,
	0x00, 0xa2, 0x8c, 0x9e, 0x82, 0x70, 0xab, 0x1d, 0xd4, 0x75, 0x02, 0xf7, 0xb7, 0xb3, 0xdd, 0xa6,
	0x5d, 0xc7, 0x41, 0x92, 0x08, 0x90, 0xf2, 0x58, 0x09, 0x96, 0x0e, 0xc3, 0xc6, 0x00, 0xe0, 0x28,
	0x47, 0x92, 0x2f, 0xf1, 0x6d, 0xc9, 0x86, 0x29, 0x1d, 0x47, 0x09, 0x4c, 0x19, 0x55, 0x8c, 0xa7,
	0xd2, 0xdd, 0xe8, 0x54, 0xba, 0x37, 0xf6, 0xde, 0xf3, 0xae, 0xe0, 0xc7, 0x3b, 0xd6, 0xe8, 0xc3,
	0x02, 0x1c, 0x54, 0x9f, 0x5f, 0xb4, 0x4b, 0xe1, 0x96, 0x5c, 0x76, 0xcb, 0x7c, 0x46, 0x96, 0x2a,
	0x10, 0x53, 0x3a, 0x76, 0x6b, 0x66, 0xc6, 0xc2, 0x26, 0x13, 0x7c, 0x53, 0x71, 0xa5, 0x6b, 0x66,
	0x5c, 0x32, 0xe5, 0xd6, 0x75, 0xc1, 0x3b, 0x9e, 0x6d, 0x36, 0x27, 0x7a, 0x5e, 0xf0, 0x63, 0xce,
	0xd2, 0x60, 0x3f, 0xaf, 0xf2, 0xf3, 0x9f, 0xed, 0xee, 0x90, 0xa9, 0xd1, 0xa4, 0xef, 0xc5, 0xfc,
	0xc4, 0x12, 0x6d, 0xff, 0xed, 0xca, 0xe4, 0x89, 0xaf, 0x4e, 0x33, 0x90, 0xfa, 0x80, 0x7c, 0xf6,
	0xf7, 0x2f, 0x1f, 0xa0, 0x70, 0x53, 0x97, 0x39, 0x34, 0x55, 0xc8, 0x5d, 0xec, 0x30, 0x19, 0xd1,
	0x58, 0xb1, 0x29, 0xb8, 0x8d, 0x0e, 0xea, 0x36, 0xc2, 0x06, 0x93, 0x07, 0xda, 0x26, 0xef, 0x60,
	0x1c, 0x0b, 0xa0, 0x0a, 0x92, 0x88, 0x2a, 0xd7, 0xe9, 0xa0, 0x6e, 0x25, 0x74, 0xac, 0xe7, 0x40,
	0x91, 0x3d, 0x5c, 0xd7, 0x06, 0x17, 0x2e, 0xbe, 0x66, 0xb9, 0x05, 0xf0, 0x61, 0xf5, 0x9f, 0xa7,
	0x6d, 0xb4, 0xf3, 0x13, 0xc2, 0x9b, 0x9f, 0x51, 0x05, 0x52, 0x1d, 0x09, 0x16, 0x83, 0x24, 0xef,
	0x63, 0xc7, 0xec, 0x34, 0x9a, 0x8b, 0x63, 0x73, 0x76, 0xd1, 0x6e, 0x18, 0xd1, 0xf4, 0x0e, 0xc3,
	0x86, 0x09, 0xf7, 0x12, 0xf2, 0x11, 0xae, 0x65, 0xfa, 0x90, 0x5b, 0xb6, 0x1b, 0xd2, 0x94, 0x18,
	0xb9, 0x15, 0x0b, 0xd2, 0x69, 0x2d, 0x0f, 0x16, 0x4e, 0xde, 0xc5, 0x37, 0xc7, 0x54, 0xaa, 0x68,
	0x4e, 0x41, 0x45, 0x0f, 0xb4, 0x99, 0x3b, 0x7b, 0xd6, 0x67, 0xfb, 0xfb, 0x11, 0x61, 0xe7, 0x71,
	0xbe, 0xa6, 0x47, 0x00, 0x92, 0x7c, 0x8b, 0xdf, 0x32, 0xd4, 0xe4, 0x1c, 0x44, 0x19, 0x8d, 0x9f,
	0x80, 0x8a, 0x06, 0x00, 0x2e, 0xba, 0x8e, 0xa2, 0x0f, 0xd7, 0xa5, 0x28, 0x24, 0xba, 0x52, 0x40,
	0x25, 0x1c, 0xe9, 0x3a, 0x8f, 0x00, 0x6c, 0x4f, 0xff, 0x22, 0x5c, 0x33, 0xbe, 0x75, 0xb6, 0xf5,
	0xaa, 0x6b, 0xb5, 0xd8, 0x64, 0x65, 0xbd, 0x4d, 0x06, 0xb8, 0x2e, 0x20, 0x06, 0x96, 0x29, 0x7d,
	0xab, 0x5e, 0x76, 0x23, 0xc9, 0xf9, 0xd9, 0xee, 0x2d, 0xd3, 0x72, 0x68, 0xe0, 0xbd, 0xb0, 0x38,
	0xb8, 0xa2, 0xad, 0x8d, 0x15, 0x6d, 0xed, 0xfc, 0x81, 0x70, 0xbd, 0xd0, 0xe8, 0x1a, 0xe3, 0x3e,
	0xc0, 0x8e, 0xbd, 0x3f, 0x5c, 0xe8, 0x79, 0x5f, 0x25, 0xca, 0x05, 0x94, 0x8c, 0x70, 0x8d, 0x9e,
	0xf0, 0x49, 0xaa, 0xe6, 0xab, 0x78, 0xdd, 0xd7, 0xce, 0xe6, 0xb7, 0x64, 0x9e, 0x23, 0xfc, 0xc6,
	0xca, 0xab, 0x91, 0x8f, 0x69, 0x9f, 0x1d, 0x3b, 0xa6, 0x63, 0xc6, 0x34, 0xb8, 0x7c, 0x4c, 0x13,
	0xee, 0x25, 0x24, 0xc0, 0x44, 0xf2, 0x81, 0x5a, 0xbc, 0x4f, 0x51, 0x3f, 0x93, 0x86, 0xdf, 0xa0,
	0x39, 0xbb, 0x68, 0x6f, 0x1d, 0xf3, 0x81, 0x5a, 0xbc, 0x47, 0x47, 0xc7, 0xe1, 0x96, 0x5c, 0xf2,
	0x64, 0x39, 0x89, 0x64, 0x44, 0x45, 0xb2, 0x92, 0xa3, 0xb2, 0xc8, 0xf1, 0x29, 0x15, 0xc9, 0x72,
	0x8e, 0xd1, 0x92, 0x27, 0x93, 0x76, 0x98, 0x5f, 0x11, 0x6e, 0x1a, 0x2e, 0x74, 0xab, 0x6a, 0x22,
	0xe0, 0x0b, 0x91, 0x80, 0x58, 0x12, 0x1f, 0x7a, 0xa9, 0xf8, 0xd6, 0xbc, 0xc6, 0xcb, 0xc2, 0xa9,
	0xac, 0x3e, 0x4a, 0xfb, 0xb8, 0x6e, 0xbf, 0x4d, 0x5a, 0x9b, 0xb7, 0xf6, 0xee, 0x5e, 0x95, 0xf8,
	0x13, 0x03, 0x09, 0x0b, 0xec, 0xc3, 0xea, 0x77, 0x4f, 0xdb, 0xa5, 0xe0, 0xf3, 0x67, 0xb3, 0x16,
	0x7a, 0x3e, 0x6b, 0xa1, 0x17, 0xb3, 0x16, 0xfa, 0x6b, 0xd6, 0x42, 0x3f, 0x5c, 0xb6, 0x4a, 0x2f,
	0x2e, 0x5b, 0xa5, 0xdf, 0x2f, 0x5b, 0xa5, 0xaf, 0xfc, 0xff, 0x31, 0x9e, 0xe7, 0xd4, 0x62, 0x8f,
	0xf9, 0xd8, 0x8f, 0x47, 0x94, 0xa5, 0xfe, 0xf4, 0xbe, 0xff, 0x4d, 0xf1, 0x6d, 0xd7, 0xf4, 0xf7,
	0x6b, 0x1a, 0x71, 0xff, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x38, 0x1b, 0xc5, 0xfa, 0xf7, 0x07,
	0x00, 0x00,
}

func (this *Tunnel) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Tunnel)
	if !ok {
		that2, ok := that.(Tunnel)
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
	if this.Sequence != that1.Sequence {
		return false
	}
	if !this.Route.Equal(that1.Route) {
		return false
	}
	if this.FeePayer != that1.FeePayer {
		return false
	}
	if len(this.SignalDeviations) != len(that1.SignalDeviations) {
		return false
	}
	for i := range this.SignalDeviations {
		if !this.SignalDeviations[i].Equal(&that1.SignalDeviations[i]) {
			return false
		}
	}
	if this.Interval != that1.Interval {
		return false
	}
	if len(this.TotalDeposit) != len(that1.TotalDeposit) {
		return false
	}
	for i := range this.TotalDeposit {
		if !this.TotalDeposit[i].Equal(&that1.TotalDeposit[i]) {
			return false
		}
	}
	if this.IsActive != that1.IsActive {
		return false
	}
	if this.CreatedAt != that1.CreatedAt {
		return false
	}
	if this.Creator != that1.Creator {
		return false
	}
	return true
}
func (this *LatestPrices) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*LatestPrices)
	if !ok {
		that2, ok := that.(LatestPrices)
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
	if this.TunnelID != that1.TunnelID {
		return false
	}
	if len(this.Prices) != len(that1.Prices) {
		return false
	}
	for i := range this.Prices {
		if !this.Prices[i].Equal(&that1.Prices[i]) {
			return false
		}
	}
	if this.LastInterval != that1.LastInterval {
		return false
	}
	return true
}
func (this *TotalFees) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TotalFees)
	if !ok {
		that2, ok := that.(TotalFees)
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
	if len(this.TotalBasePacketFee) != len(that1.TotalBasePacketFee) {
		return false
	}
	for i := range this.TotalBasePacketFee {
		if !this.TotalBasePacketFee[i].Equal(&that1.TotalBasePacketFee[i]) {
			return false
		}
	}
	return true
}
func (this *Packet) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Packet)
	if !ok {
		that2, ok := that.(Packet)
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
	if this.TunnelID != that1.TunnelID {
		return false
	}
	if this.Sequence != that1.Sequence {
		return false
	}
	if len(this.Prices) != len(that1.Prices) {
		return false
	}
	for i := range this.Prices {
		if !this.Prices[i].Equal(&that1.Prices[i]) {
			return false
		}
	}
	if !this.Receipt.Equal(that1.Receipt) {
		return false
	}
	if this.CreatedAt != that1.CreatedAt {
		return false
	}
	return true
}
func (this *Deposit) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Deposit)
	if !ok {
		that2, ok := that.(Deposit)
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
	if this.TunnelID != that1.TunnelID {
		return false
	}
	if this.Depositor != that1.Depositor {
		return false
	}
	if len(this.Amount) != len(that1.Amount) {
		return false
	}
	for i := range this.Amount {
		if !this.Amount[i].Equal(&that1.Amount[i]) {
			return false
		}
	}
	return true
}
func (this *SignalDeviation) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SignalDeviation)
	if !ok {
		that2, ok := that.(SignalDeviation)
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
	if this.SoftDeviationBPS != that1.SoftDeviationBPS {
		return false
	}
	if this.HardDeviationBPS != that1.HardDeviationBPS {
		return false
	}
	return true
}
func (this *TunnelSignatureOrder) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TunnelSignatureOrder)
	if !ok {
		that2, ok := that.(TunnelSignatureOrder)
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
	if this.Sequence != that1.Sequence {
		return false
	}
	if len(this.Prices) != len(that1.Prices) {
		return false
	}
	for i := range this.Prices {
		if !this.Prices[i].Equal(&that1.Prices[i]) {
			return false
		}
	}
	if this.CreatedAt != that1.CreatedAt {
		return false
	}
	if this.Encoder != that1.Encoder {
		return false
	}
	return true
}
func (m *Tunnel) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Tunnel) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Tunnel) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTunnel(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x52
	}
	if m.CreatedAt != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.CreatedAt))
		i--
		dAtA[i] = 0x48
	}
	if m.IsActive {
		i--
		if m.IsActive {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if len(m.TotalDeposit) > 0 {
		for iNdEx := len(m.TotalDeposit) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TotalDeposit[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTunnel(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.Interval != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.Interval))
		i--
		dAtA[i] = 0x30
	}
	if len(m.SignalDeviations) > 0 {
		for iNdEx := len(m.SignalDeviations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SignalDeviations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTunnel(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.FeePayer) > 0 {
		i -= len(m.FeePayer)
		copy(dAtA[i:], m.FeePayer)
		i = encodeVarintTunnel(dAtA, i, uint64(len(m.FeePayer)))
		i--
		dAtA[i] = 0x22
	}
	if m.Route != nil {
		{
			size, err := m.Route.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTunnel(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Sequence != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x10
	}
	if m.ID != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LatestPrices) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LatestPrices) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LatestPrices) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastInterval != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.LastInterval))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Prices) > 0 {
		for iNdEx := len(m.Prices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Prices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTunnel(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.TunnelID != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.TunnelID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TotalFees) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TotalFees) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TotalFees) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TotalBasePacketFee) > 0 {
		for iNdEx := len(m.TotalBasePacketFee) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TotalBasePacketFee[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTunnel(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Packet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Packet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Packet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CreatedAt != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.CreatedAt))
		i--
		dAtA[i] = 0x28
	}
	if m.Receipt != nil {
		{
			size, err := m.Receipt.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTunnel(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.Prices) > 0 {
		for iNdEx := len(m.Prices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Prices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTunnel(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.Sequence != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x10
	}
	if m.TunnelID != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.TunnelID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Deposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Deposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Deposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTunnel(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Depositor) > 0 {
		i -= len(m.Depositor)
		copy(dAtA[i:], m.Depositor)
		i = encodeVarintTunnel(dAtA, i, uint64(len(m.Depositor)))
		i--
		dAtA[i] = 0x12
	}
	if m.TunnelID != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.TunnelID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SignalDeviation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SignalDeviation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignalDeviation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.HardDeviationBPS != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.HardDeviationBPS))
		i--
		dAtA[i] = 0x18
	}
	if m.SoftDeviationBPS != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.SoftDeviationBPS))
		i--
		dAtA[i] = 0x10
	}
	if len(m.SignalID) > 0 {
		i -= len(m.SignalID)
		copy(dAtA[i:], m.SignalID)
		i = encodeVarintTunnel(dAtA, i, uint64(len(m.SignalID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TunnelSignatureOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TunnelSignatureOrder) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TunnelSignatureOrder) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Encoder != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.Encoder))
		i--
		dAtA[i] = 0x20
	}
	if m.CreatedAt != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.CreatedAt))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Prices) > 0 {
		for iNdEx := len(m.Prices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Prices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTunnel(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Sequence != 0 {
		i = encodeVarintTunnel(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTunnel(dAtA []byte, offset int, v uint64) int {
	offset -= sovTunnel(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Tunnel) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovTunnel(uint64(m.ID))
	}
	if m.Sequence != 0 {
		n += 1 + sovTunnel(uint64(m.Sequence))
	}
	if m.Route != nil {
		l = m.Route.Size()
		n += 1 + l + sovTunnel(uint64(l))
	}
	l = len(m.FeePayer)
	if l > 0 {
		n += 1 + l + sovTunnel(uint64(l))
	}
	if len(m.SignalDeviations) > 0 {
		for _, e := range m.SignalDeviations {
			l = e.Size()
			n += 1 + l + sovTunnel(uint64(l))
		}
	}
	if m.Interval != 0 {
		n += 1 + sovTunnel(uint64(m.Interval))
	}
	if len(m.TotalDeposit) > 0 {
		for _, e := range m.TotalDeposit {
			l = e.Size()
			n += 1 + l + sovTunnel(uint64(l))
		}
	}
	if m.IsActive {
		n += 2
	}
	if m.CreatedAt != 0 {
		n += 1 + sovTunnel(uint64(m.CreatedAt))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTunnel(uint64(l))
	}
	return n
}

func (m *LatestPrices) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelID != 0 {
		n += 1 + sovTunnel(uint64(m.TunnelID))
	}
	if len(m.Prices) > 0 {
		for _, e := range m.Prices {
			l = e.Size()
			n += 1 + l + sovTunnel(uint64(l))
		}
	}
	if m.LastInterval != 0 {
		n += 1 + sovTunnel(uint64(m.LastInterval))
	}
	return n
}

func (m *TotalFees) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.TotalBasePacketFee) > 0 {
		for _, e := range m.TotalBasePacketFee {
			l = e.Size()
			n += 1 + l + sovTunnel(uint64(l))
		}
	}
	return n
}

func (m *Packet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelID != 0 {
		n += 1 + sovTunnel(uint64(m.TunnelID))
	}
	if m.Sequence != 0 {
		n += 1 + sovTunnel(uint64(m.Sequence))
	}
	if len(m.Prices) > 0 {
		for _, e := range m.Prices {
			l = e.Size()
			n += 1 + l + sovTunnel(uint64(l))
		}
	}
	if m.Receipt != nil {
		l = m.Receipt.Size()
		n += 1 + l + sovTunnel(uint64(l))
	}
	if m.CreatedAt != 0 {
		n += 1 + sovTunnel(uint64(m.CreatedAt))
	}
	return n
}

func (m *Deposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelID != 0 {
		n += 1 + sovTunnel(uint64(m.TunnelID))
	}
	l = len(m.Depositor)
	if l > 0 {
		n += 1 + l + sovTunnel(uint64(l))
	}
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovTunnel(uint64(l))
		}
	}
	return n
}

func (m *SignalDeviation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SignalID)
	if l > 0 {
		n += 1 + l + sovTunnel(uint64(l))
	}
	if m.SoftDeviationBPS != 0 {
		n += 1 + sovTunnel(uint64(m.SoftDeviationBPS))
	}
	if m.HardDeviationBPS != 0 {
		n += 1 + sovTunnel(uint64(m.HardDeviationBPS))
	}
	return n
}

func (m *TunnelSignatureOrder) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Sequence != 0 {
		n += 1 + sovTunnel(uint64(m.Sequence))
	}
	if len(m.Prices) > 0 {
		for _, e := range m.Prices {
			l = e.Size()
			n += 1 + l + sovTunnel(uint64(l))
		}
	}
	if m.CreatedAt != 0 {
		n += 1 + sovTunnel(uint64(m.CreatedAt))
	}
	if m.Encoder != 0 {
		n += 1 + sovTunnel(uint64(m.Encoder))
	}
	return n
}

func sovTunnel(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTunnel(x uint64) (n int) {
	return sovTunnel(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Tunnel) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTunnel
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
			return fmt.Errorf("proto: Tunnel: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Tunnel: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Route", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Route == nil {
				m.Route = &types.Any{}
			}
			if err := m.Route.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeePayer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeePayer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalDeviations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignalDeviations = append(m.SignalDeviations, SignalDeviation{})
			if err := m.SignalDeviations[len(m.SignalDeviations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Interval", wireType)
			}
			m.Interval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Interval |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalDeposit", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TotalDeposit = append(m.TotalDeposit, types1.Coin{})
			if err := m.TotalDeposit[len(m.TotalDeposit)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsActive", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsActive = bool(v != 0)
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedAt |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTunnel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTunnel
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
func (m *LatestPrices) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTunnel
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
			return fmt.Errorf("proto: LatestPrices: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LatestPrices: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelID", wireType)
			}
			m.TunnelID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TunnelID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prices = append(m.Prices, types2.Price{})
			if err := m.Prices[len(m.Prices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastInterval", wireType)
			}
			m.LastInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastInterval |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTunnel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTunnel
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
func (m *TotalFees) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTunnel
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
			return fmt.Errorf("proto: TotalFees: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TotalFees: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalBasePacketFee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TotalBasePacketFee = append(m.TotalBasePacketFee, types1.Coin{})
			if err := m.TotalBasePacketFee[len(m.TotalBasePacketFee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTunnel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTunnel
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
func (m *Packet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTunnel
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
			return fmt.Errorf("proto: Packet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Packet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelID", wireType)
			}
			m.TunnelID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TunnelID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prices = append(m.Prices, types2.Price{})
			if err := m.Prices[len(m.Prices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receipt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Receipt == nil {
				m.Receipt = &types.Any{}
			}
			if err := m.Receipt.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedAt |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTunnel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTunnel
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
func (m *Deposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTunnel
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
			return fmt.Errorf("proto: Deposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Deposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelID", wireType)
			}
			m.TunnelID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TunnelID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Depositor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Depositor = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types1.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTunnel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTunnel
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
func (m *SignalDeviation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTunnel
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
			return fmt.Errorf("proto: SignalDeviation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SignalDeviation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignalID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignalID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SoftDeviationBPS", wireType)
			}
			m.SoftDeviationBPS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SoftDeviationBPS |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HardDeviationBPS", wireType)
			}
			m.HardDeviationBPS = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HardDeviationBPS |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTunnel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTunnel
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
func (m *TunnelSignatureOrder) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTunnel
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
			return fmt.Errorf("proto: TunnelSignatureOrder: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TunnelSignatureOrder: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
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
				return ErrInvalidLengthTunnel
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTunnel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prices = append(m.Prices, types2.Price{})
			if err := m.Prices[len(m.Prices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedAt |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Encoder", wireType)
			}
			m.Encoder = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTunnel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Encoder |= types2.Encoder(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTunnel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTunnel
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
func skipTunnel(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTunnel
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
					return 0, ErrIntOverflowTunnel
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
					return 0, ErrIntOverflowTunnel
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
				return 0, ErrInvalidLengthTunnel
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTunnel
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTunnel
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTunnel        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTunnel          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTunnel = fmt.Errorf("proto: unexpected end of group")
)
