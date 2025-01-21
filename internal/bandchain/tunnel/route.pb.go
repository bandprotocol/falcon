package types

import (
	"fmt"
	"io"
	"math"
	math_bits "math/bits"

	"github.com/cosmos/gogoproto/proto"

	github_com_bandprotocol_chain_v3_x_bandtss_types "github.com/bandprotocol/falcon/internal/bandchain/bandtss"
	types "github.com/bandprotocol/falcon/internal/bandchain/feeds"
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

// TSSRoute represents a route for TSS packets and implements the RouteI interface.
type TSSRoute struct {
	// destination_chain_id is the destination chain ID
	DestinationChainID string `protobuf:"bytes,1,opt,name=destination_chain_id,json=destinationChainId,proto3"                 json:"destination_chain_id,omitempty"`
	// destination_contract_address is the destination contract address
	DestinationContractAddress string `protobuf:"bytes,2,opt,name=destination_contract_address,json=destinationContractAddress,proto3" json:"destination_contract_address,omitempty"`
	// encoder is the mode of encoding packet data.
	Encoder types.Encoder `protobuf:"varint,3,opt,name=encoder,proto3,enum=band.feeds.v1beta1.Encoder"                     json:"encoder,omitempty"`
}

func (m *TSSRoute) Reset()         { *m = TSSRoute{} }
func (m *TSSRoute) String() string { return proto.CompactTextString(m) }
func (*TSSRoute) ProtoMessage()    {}
func (*TSSRoute) Descriptor() ([]byte, []int) {
	return fileDescriptor_543238289d94b7a6, []int{0}
}
func (m *TSSRoute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TSSRoute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TSSRoute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TSSRoute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TSSRoute.Merge(m, src)
}
func (m *TSSRoute) XXX_Size() int {
	return m.Size()
}
func (m *TSSRoute) XXX_DiscardUnknown() {
	xxx_messageInfo_TSSRoute.DiscardUnknown(m)
}

var xxx_messageInfo_TSSRoute proto.InternalMessageInfo

func (m *TSSRoute) GetDestinationChainID() string {
	if m != nil {
		return m.DestinationChainID
	}
	return ""
}

func (m *TSSRoute) GetDestinationContractAddress() string {
	if m != nil {
		return m.DestinationContractAddress
	}
	return ""
}

func (m *TSSRoute) GetEncoder() types.Encoder {
	if m != nil {
		return m.Encoder
	}
	return types.ENCODER_UNSPECIFIED
}

// TSSPacketReceipt represents a receipt for a TSS packet and implements the PacketReceiptI interface.
type TSSPacketReceipt struct {
	// signing_id is the signing ID
	SigningID github_com_bandprotocol_chain_v3_x_bandtss_types.SigningID `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3,casttype=github.com/bandprotocol/chain/v3/x/bandtss/types.SigningID" json:"signing_id,omitempty"`
}

func (m *TSSPacketReceipt) Reset()         { *m = TSSPacketReceipt{} }
func (m *TSSPacketReceipt) String() string { return proto.CompactTextString(m) }
func (*TSSPacketReceipt) ProtoMessage()    {}
func (*TSSPacketReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_543238289d94b7a6, []int{1}
}
func (m *TSSPacketReceipt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TSSPacketReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TSSPacketReceipt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TSSPacketReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TSSPacketReceipt.Merge(m, src)
}
func (m *TSSPacketReceipt) XXX_Size() int {
	return m.Size()
}
func (m *TSSPacketReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_TSSPacketReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_TSSPacketReceipt proto.InternalMessageInfo

func (m *TSSPacketReceipt) GetSigningID() github_com_bandprotocol_chain_v3_x_bandtss_types.SigningID {
	if m != nil {
		return m.SigningID
	}
	return 0
}

// IBCRoute represents a route for IBC packets and implements the RouteI interface.
type IBCRoute struct {
	// channel_id is the IBC channel ID
	ChannelID string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (m *IBCRoute) Reset()         { *m = IBCRoute{} }
func (m *IBCRoute) String() string { return proto.CompactTextString(m) }
func (*IBCRoute) ProtoMessage()    {}
func (*IBCRoute) Descriptor() ([]byte, []int) {
	return fileDescriptor_543238289d94b7a6, []int{2}
}
func (m *IBCRoute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IBCRoute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IBCRoute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IBCRoute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IBCRoute.Merge(m, src)
}
func (m *IBCRoute) XXX_Size() int {
	return m.Size()
}
func (m *IBCRoute) XXX_DiscardUnknown() {
	xxx_messageInfo_IBCRoute.DiscardUnknown(m)
}

var xxx_messageInfo_IBCRoute proto.InternalMessageInfo

func (m *IBCRoute) GetChannelID() string {
	if m != nil {
		return m.ChannelID
	}
	return ""
}

// IBCPacketReceipt represents a receipt for a IBC packet and implements the PacketReceiptI interface.
type IBCPacketReceipt struct {
	// sequence is representing the sequence of the IBC packet.
	Sequence uint64 `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

func (m *IBCPacketReceipt) Reset()         { *m = IBCPacketReceipt{} }
func (m *IBCPacketReceipt) String() string { return proto.CompactTextString(m) }
func (*IBCPacketReceipt) ProtoMessage()    {}
func (*IBCPacketReceipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_543238289d94b7a6, []int{3}
}
func (m *IBCPacketReceipt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IBCPacketReceipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IBCPacketReceipt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IBCPacketReceipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IBCPacketReceipt.Merge(m, src)
}
func (m *IBCPacketReceipt) XXX_Size() int {
	return m.Size()
}
func (m *IBCPacketReceipt) XXX_DiscardUnknown() {
	xxx_messageInfo_IBCPacketReceipt.DiscardUnknown(m)
}

var xxx_messageInfo_IBCPacketReceipt proto.InternalMessageInfo

func (m *IBCPacketReceipt) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

// TunnelPricesPacketData represents the IBC packet payload for the tunnel packet.
type TunnelPricesPacketData struct {
	// tunnel_id is the tunnel ID
	TunnelID uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3"   json:"tunnel_id,omitempty"`
	// sequence is representing the sequence of the tunnel packet.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence,proto3"                  json:"sequence,omitempty"`
	// prices is the list of prices information from feeds module.
	Prices []types.Price `protobuf:"bytes,3,rep,name=prices,proto3"                     json:"prices"`
	// created_at is the timestamp when the packet is created
	CreatedAt int64 `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (m *TunnelPricesPacketData) Reset()         { *m = TunnelPricesPacketData{} }
func (m *TunnelPricesPacketData) String() string { return proto.CompactTextString(m) }
func (*TunnelPricesPacketData) ProtoMessage()    {}
func (*TunnelPricesPacketData) Descriptor() ([]byte, []int) {
	return fileDescriptor_543238289d94b7a6, []int{4}
}
func (m *TunnelPricesPacketData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TunnelPricesPacketData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TunnelPricesPacketData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TunnelPricesPacketData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TunnelPricesPacketData.Merge(m, src)
}
func (m *TunnelPricesPacketData) XXX_Size() int {
	return m.Size()
}
func (m *TunnelPricesPacketData) XXX_DiscardUnknown() {
	xxx_messageInfo_TunnelPricesPacketData.DiscardUnknown(m)
}

var xxx_messageInfo_TunnelPricesPacketData proto.InternalMessageInfo

func (m *TunnelPricesPacketData) GetTunnelID() uint64 {
	if m != nil {
		return m.TunnelID
	}
	return 0
}

func (m *TunnelPricesPacketData) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *TunnelPricesPacketData) GetPrices() []types.Price {
	if m != nil {
		return m.Prices
	}
	return nil
}

func (m *TunnelPricesPacketData) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func init() {
	proto.RegisterType((*TSSRoute)(nil), "band.tunnel.v1beta1.TSSRoute")
	proto.RegisterType((*TSSPacketReceipt)(nil), "band.tunnel.v1beta1.TSSPacketReceipt")
	proto.RegisterType((*IBCRoute)(nil), "band.tunnel.v1beta1.IBCRoute")
	proto.RegisterType((*IBCPacketReceipt)(nil), "band.tunnel.v1beta1.IBCPacketReceipt")
	proto.RegisterType((*TunnelPricesPacketData)(nil), "band.tunnel.v1beta1.TunnelPricesPacketData")
}

func init() { proto.RegisterFile("band/tunnel/v1beta1/route.proto", fileDescriptor_543238289d94b7a6) }

var fileDescriptor_543238289d94b7a6 = []byte{
	// 542 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xc1, 0x6e, 0xd3, 0x30,
	0x18, 0xc7, 0xeb, 0xb5, 0x2a, 0xad, 0x81, 0x69, 0x32, 0xd3, 0xd4, 0x15, 0x48, 0xaa, 0x9e, 0x8a,
	0xc4, 0x12, 0x6d, 0x13, 0x42, 0xea, 0x89, 0xa5, 0x41, 0x22, 0x42, 0x48, 0x53, 0xda, 0x13, 0x97,
	0xca, 0xb5, 0x4d, 0x1b, 0xd8, 0xec, 0x12, 0xbb, 0x13, 0xbc, 0x05, 0xe2, 0x29, 0x78, 0x80, 0x49,
	0xbc, 0xc2, 0xb4, 0xd3, 0x8e, 0x70, 0x89, 0x50, 0xfa, 0x16, 0x9c, 0x50, 0x6c, 0xaf, 0x25, 0xa8,
	0x07, 0x6e, 0xf5, 0xff, 0xfb, 0xe5, 0xfb, 0x7f, 0xff, 0xcf, 0x2e, 0x74, 0x27, 0x98, 0x53, 0x5f,
	0x2d, 0x38, 0x67, 0x67, 0xfe, 0xc5, 0xe1, 0x84, 0x29, 0x7c, 0xe8, 0xa7, 0x62, 0xa1, 0x98, 0x37,
	0x4f, 0x85, 0x12, 0xe8, 0x41, 0x01, 0x78, 0x06, 0xf0, 0x2c, 0xd0, 0xde, 0x27, 0x42, 0x9e, 0x0b,
	0x39, 0xd6, 0x88, 0x6f, 0x0e, 0x86, 0x6f, 0xef, 0x4e, 0xc5, 0x54, 0x18, 0xbd, 0xf8, 0x65, 0xd5,
	0x8e, 0xb6, 0x79, 0xc7, 0x18, 0x95, 0x2b, 0x17, 0xc6, 0x89, 0xa0, 0x2c, 0xb5, 0x84, 0xb3, 0x81,
	0xd0, 0x27, 0x53, 0xef, 0xfe, 0x04, 0xb0, 0x31, 0x1a, 0x0e, 0xe3, 0x62, 0x34, 0xf4, 0x0a, 0xee,
	0x52, 0x26, 0x55, 0xc2, 0xb1, 0x4a, 0x04, 0x1f, 0x93, 0x19, 0x4e, 0xf8, 0x38, 0xa1, 0x2d, 0xd0,
	0x01, 0xbd, 0x66, 0xb0, 0x97, 0x67, 0x2e, 0x0a, 0xd7, 0xf5, 0x41, 0x51, 0x8e, 0xc2, 0x18, 0xd1,
	0x7f, 0x35, 0x8a, 0x5e, 0xc0, 0x47, 0xa5, 0x4e, 0x82, 0xab, 0x14, 0x13, 0x35, 0xc6, 0x94, 0xa6,
	0x4c, 0xca, 0xd6, 0x56, 0xd1, 0x31, 0x6e, 0xff, 0xfd, 0xa5, 0x45, 0x4e, 0x0c, 0x81, 0x9e, 0xc1,
	0x3b, 0x36, 0x49, 0xab, 0xda, 0x01, 0xbd, 0xed, 0xa3, 0x87, 0x9e, 0x5e, 0x99, 0x19, 0xde, 0x46,
	0xf1, 0x5e, 0x1a, 0x24, 0xbe, 0x65, 0xfb, 0xf0, 0xfa, 0xf2, 0xa0, 0xae, 0xd3, 0x44, 0xdd, 0xaf,
	0x00, 0xee, 0x8c, 0x86, 0xc3, 0x53, 0x4c, 0x3e, 0x30, 0x15, 0x33, 0xc2, 0x92, 0xb9, 0x42, 0xef,
	0x21, 0x94, 0xc9, 0x94, 0x27, 0x7c, 0x7a, 0x9b, 0xac, 0x16, 0xbc, 0xce, 0x33, 0xb7, 0x39, 0x34,
	0x6a, 0x14, 0xfe, 0xce, 0xdc, 0xfe, 0x34, 0x51, 0xb3, 0xc5, 0xc4, 0x23, 0xe2, 0xdc, 0x2f, 0x5c,
	0xf5, 0xae, 0x88, 0x38, 0xf3, 0xf5, 0x4a, 0xfc, 0x8b, 0x63, 0xff, 0x93, 0xd6, 0x95, 0x94, 0xbe,
	0xfa, 0x3c, 0x67, 0xd2, 0x5b, 0x7d, 0x1d, 0x37, 0x6d, 0xfb, 0x88, 0xf6, 0xd1, 0xf5, 0xe5, 0xc1,
	0x76, 0xc9, 0x3e, 0xea, 0x86, 0xb0, 0x11, 0x05, 0x03, 0xb3, 0xef, 0xa7, 0x10, 0x92, 0x19, 0x2e,
	0x9e, 0xc0, 0x7a, 0xcb, 0xf7, 0x8b, 0x59, 0x06, 0x46, 0x2d, 0xba, 0x59, 0x20, 0xa2, 0xa5, 0x68,
	0x01, 0xdc, 0x89, 0x82, 0x41, 0x39, 0x59, 0x1b, 0x36, 0x24, 0xfb, 0xb8, 0x60, 0x9c, 0x30, 0x93,
	0x2b, 0x5e, 0x9d, 0x37, 0x4e, 0xf2, 0x1d, 0xc0, 0xbd, 0x91, 0x7e, 0x80, 0xa7, 0x69, 0x42, 0x98,
	0x34, 0xe5, 0x10, 0x2b, 0x8c, 0x9e, 0xc0, 0xa6, 0x79, 0x9a, 0xeb, 0x1d, 0xdd, 0xcb, 0x33, 0xb7,
	0x61, 0xf0, 0x28, 0x8c, 0x1b, 0xa6, 0x1c, 0xd1, 0x92, 0xeb, 0x56, 0xd9, 0x15, 0x3d, 0x87, 0xf5,
	0xb9, 0x6e, 0xdd, 0xaa, 0x76, 0xaa, 0xbd, 0xbb, 0x47, 0xfb, 0x9b, 0xae, 0x50, 0x9b, 0x07, 0xb5,
	0xab, 0xcc, 0xad, 0xc4, 0x16, 0x47, 0x8f, 0x21, 0x24, 0x29, 0xc3, 0x8a, 0xd1, 0x31, 0x56, 0xad,
	0x5a, 0x07, 0xf4, 0xaa, 0x71, 0xd3, 0x2a, 0x27, 0x2a, 0x78, 0xf3, 0x2d, 0x77, 0xc0, 0x55, 0xee,
	0x80, 0x9b, 0xdc, 0x01, 0xbf, 0x72, 0x07, 0x7c, 0x59, 0x3a, 0x95, 0x9b, 0xa5, 0x53, 0xf9, 0xb1,
	0x74, 0x2a, 0x6f, 0xfd, 0xff, 0xb8, 0x3c, 0xfb, 0xcf, 0xd4, 0x77, 0x37, 0xa9, 0x6b, 0xe2, 0xf8,
	0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xce, 0xc1, 0x97, 0x14, 0xb5, 0x03, 0x00, 0x00,
}

func (this *TSSRoute) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TSSRoute)
	if !ok {
		that2, ok := that.(TSSRoute)
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
	if this.DestinationChainID != that1.DestinationChainID {
		return false
	}
	if this.DestinationContractAddress != that1.DestinationContractAddress {
		return false
	}
	if this.Encoder != that1.Encoder {
		return false
	}
	return true
}
func (this *TSSPacketReceipt) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TSSPacketReceipt)
	if !ok {
		that2, ok := that.(TSSPacketReceipt)
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
	if this.SigningID != that1.SigningID {
		return false
	}
	return true
}
func (this *IBCRoute) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*IBCRoute)
	if !ok {
		that2, ok := that.(IBCRoute)
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
	if this.ChannelID != that1.ChannelID {
		return false
	}
	return true
}
func (this *IBCPacketReceipt) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*IBCPacketReceipt)
	if !ok {
		that2, ok := that.(IBCPacketReceipt)
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
	return true
}
func (this *TunnelPricesPacketData) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TunnelPricesPacketData)
	if !ok {
		that2, ok := that.(TunnelPricesPacketData)
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
	if this.CreatedAt != that1.CreatedAt {
		return false
	}
	return true
}
func (m *TSSRoute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TSSRoute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TSSRoute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Encoder != 0 {
		i = encodeVarintRoute(dAtA, i, uint64(m.Encoder))
		i--
		dAtA[i] = 0x18
	}
	if len(m.DestinationContractAddress) > 0 {
		i -= len(m.DestinationContractAddress)
		copy(dAtA[i:], m.DestinationContractAddress)
		i = encodeVarintRoute(dAtA, i, uint64(len(m.DestinationContractAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DestinationChainID) > 0 {
		i -= len(m.DestinationChainID)
		copy(dAtA[i:], m.DestinationChainID)
		i = encodeVarintRoute(dAtA, i, uint64(len(m.DestinationChainID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TSSPacketReceipt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TSSPacketReceipt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TSSPacketReceipt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SigningID != 0 {
		i = encodeVarintRoute(dAtA, i, uint64(m.SigningID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *IBCRoute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IBCRoute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IBCRoute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChannelID) > 0 {
		i -= len(m.ChannelID)
		copy(dAtA[i:], m.ChannelID)
		i = encodeVarintRoute(dAtA, i, uint64(len(m.ChannelID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *IBCPacketReceipt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IBCPacketReceipt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IBCPacketReceipt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Sequence != 0 {
		i = encodeVarintRoute(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TunnelPricesPacketData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TunnelPricesPacketData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TunnelPricesPacketData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CreatedAt != 0 {
		i = encodeVarintRoute(dAtA, i, uint64(m.CreatedAt))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Prices) > 0 {
		for iNdEx := len(m.Prices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Prices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRoute(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.Sequence != 0 {
		i = encodeVarintRoute(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x10
	}
	if m.TunnelID != 0 {
		i = encodeVarintRoute(dAtA, i, uint64(m.TunnelID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintRoute(dAtA []byte, offset int, v uint64) int {
	offset -= sovRoute(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TSSRoute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DestinationChainID)
	if l > 0 {
		n += 1 + l + sovRoute(uint64(l))
	}
	l = len(m.DestinationContractAddress)
	if l > 0 {
		n += 1 + l + sovRoute(uint64(l))
	}
	if m.Encoder != 0 {
		n += 1 + sovRoute(uint64(m.Encoder))
	}
	return n
}

func (m *TSSPacketReceipt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigningID != 0 {
		n += 1 + sovRoute(uint64(m.SigningID))
	}
	return n
}

func (m *IBCRoute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChannelID)
	if l > 0 {
		n += 1 + l + sovRoute(uint64(l))
	}
	return n
}

func (m *IBCPacketReceipt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Sequence != 0 {
		n += 1 + sovRoute(uint64(m.Sequence))
	}
	return n
}

func (m *TunnelPricesPacketData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelID != 0 {
		n += 1 + sovRoute(uint64(m.TunnelID))
	}
	if m.Sequence != 0 {
		n += 1 + sovRoute(uint64(m.Sequence))
	}
	if len(m.Prices) > 0 {
		for _, e := range m.Prices {
			l = e.Size()
			n += 1 + l + sovRoute(uint64(l))
		}
	}
	if m.CreatedAt != 0 {
		n += 1 + sovRoute(uint64(m.CreatedAt))
	}
	return n
}

func sovRoute(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRoute(x uint64) (n int) {
	return sovRoute(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TSSRoute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRoute
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
			return fmt.Errorf("proto: TSSRoute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TSSRoute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationChainID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoute
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
				return ErrInvalidLengthRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationChainID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoute
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
				return ErrInvalidLengthRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Encoder", wireType)
			}
			m.Encoder = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Encoder |= types.Encoder(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRoute
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
func (m *TSSPacketReceipt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRoute
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
			return fmt.Errorf("proto: TSSPacketReceipt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TSSPacketReceipt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningID", wireType)
			}
			m.SigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningID |= github_com_bandprotocol_chain_v3_x_bandtss_types.SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRoute
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
func (m *IBCRoute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRoute
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
			return fmt.Errorf("proto: IBCRoute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IBCRoute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChannelID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoute
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
				return ErrInvalidLengthRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChannelID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRoute
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
func (m *IBCPacketReceipt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRoute
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
			return fmt.Errorf("proto: IBCPacketReceipt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IBCPacketReceipt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoute
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
		default:
			iNdEx = preIndex
			skippy, err := skipRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRoute
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
func (m *TunnelPricesPacketData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRoute
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
			return fmt.Errorf("proto: TunnelPricesPacketData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TunnelPricesPacketData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelID", wireType)
			}
			m.TunnelID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoute
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
					return ErrIntOverflowRoute
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
					return ErrIntOverflowRoute
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
				return ErrInvalidLengthRoute
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prices = append(m.Prices, types.Price{})
			if err := m.Prices[len(m.Prices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			m.CreatedAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoute
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
			skippy, err := skipRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRoute
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
func skipRoute(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRoute
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
					return 0, ErrIntOverflowRoute
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
					return 0, ErrIntOverflowRoute
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
				return 0, ErrInvalidLengthRoute
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRoute
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRoute
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRoute        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRoute          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRoute = fmt.Errorf("proto: unexpected end of group")
)
