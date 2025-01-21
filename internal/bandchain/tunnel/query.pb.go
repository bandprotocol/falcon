package types

import (
	"context"
	"fmt"
	"io"
	"math"
	math_bits "math/bits"

	"github.com/cosmos/cosmos-sdk/types/query"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// TunnelStatusFilter defines a filter for tunnel status.
type TunnelStatusFilter int32

const (
	// TUNNEL_STATUS_FILTER_UNSPECIFIED defines an unspecified status.
	TUNNEL_STATUS_FILTER_UNSPECIFIED TunnelStatusFilter = 0
	// TUNNEL_STATUS_FILTER_ACTIVE defines an active tunnel.
	TUNNEL_STATUS_FILTER_ACTIVE TunnelStatusFilter = 1
	// TUNNEL_STATUS_FILTER_INACTIVE defines an inactive tunnel.
	TUNNEL_STATUS_FILTER_INACTIVE TunnelStatusFilter = 2
)

var TunnelStatusFilter_name = map[int32]string{
	0: "TUNNEL_STATUS_FILTER_UNSPECIFIED",
	1: "TUNNEL_STATUS_FILTER_ACTIVE",
	2: "TUNNEL_STATUS_FILTER_INACTIVE",
}

var TunnelStatusFilter_value = map[string]int32{
	"TUNNEL_STATUS_FILTER_UNSPECIFIED": 0,
	"TUNNEL_STATUS_FILTER_ACTIVE":      1,
	"TUNNEL_STATUS_FILTER_INACTIVE":    2,
}

func (x TunnelStatusFilter) String() string {
	return proto.EnumName(TunnelStatusFilter_name, int32(x))
}

func (TunnelStatusFilter) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{0}
}

// QueryTunnelsRequest is the request type for the Query/Tunnels RPC method.
type QueryTunnelsRequest struct {
	// status_filter is a flag to filter tunnels by status.
	StatusFilter TunnelStatusFilter `protobuf:"varint,1,opt,name=status_filter,json=statusFilter,proto3,enum=band.tunnel.v1beta1.TunnelStatusFilter" json:"status_filter,omitempty"`
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3"                                                                   json:"pagination,omitempty"`
}

func (m *QueryTunnelsRequest) Reset()         { *m = QueryTunnelsRequest{} }
func (m *QueryTunnelsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTunnelsRequest) ProtoMessage()    {}
func (*QueryTunnelsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{0}
}
func (m *QueryTunnelsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTunnelsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTunnelsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTunnelsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTunnelsRequest.Merge(m, src)
}
func (m *QueryTunnelsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTunnelsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTunnelsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTunnelsRequest proto.InternalMessageInfo

func (m *QueryTunnelsRequest) GetStatusFilter() TunnelStatusFilter {
	if m != nil {
		return m.StatusFilter
	}
	return TUNNEL_STATUS_FILTER_UNSPECIFIED
}

func (m *QueryTunnelsRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryTunnelsResponse is the response type for the Query/Tunnels RPC method.
type QueryTunnelsResponse struct {
	// Tunnels is a list of tunnels.
	Tunnels []*Tunnel `protobuf:"bytes,1,rep,name=tunnels,proto3"    json:"tunnels,omitempty"`
	// pagination defines an optional pagination for the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryTunnelsResponse) Reset()         { *m = QueryTunnelsResponse{} }
func (m *QueryTunnelsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTunnelsResponse) ProtoMessage()    {}
func (*QueryTunnelsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{1}
}
func (m *QueryTunnelsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTunnelsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTunnelsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTunnelsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTunnelsResponse.Merge(m, src)
}
func (m *QueryTunnelsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTunnelsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTunnelsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTunnelsResponse proto.InternalMessageInfo

func (m *QueryTunnelsResponse) GetTunnels() []*Tunnel {
	if m != nil {
		return m.Tunnels
	}
	return nil
}

func (m *QueryTunnelsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryTunnelRequest is the request type for the Query/Tunnel RPC method.
type QueryTunnelRequest struct {
	// tunnel_id is the ID of the tunnel to query.
	TunnelId uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
}

func (m *QueryTunnelRequest) Reset()         { *m = QueryTunnelRequest{} }
func (m *QueryTunnelRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTunnelRequest) ProtoMessage()    {}
func (*QueryTunnelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{2}
}
func (m *QueryTunnelRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTunnelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTunnelRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTunnelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTunnelRequest.Merge(m, src)
}
func (m *QueryTunnelRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTunnelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTunnelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTunnelRequest proto.InternalMessageInfo

func (m *QueryTunnelRequest) GetTunnelId() uint64 {
	if m != nil {
		return m.TunnelId
	}
	return 0
}

// QueryTunnelResponse is the response type for the Query/Tunnel RPC method.
type QueryTunnelResponse struct {
	// tunnel is the tunnel with the given ID.
	Tunnel Tunnel `protobuf:"bytes,1,opt,name=tunnel,proto3" json:"tunnel"`
}

func (m *QueryTunnelResponse) Reset()         { *m = QueryTunnelResponse{} }
func (m *QueryTunnelResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTunnelResponse) ProtoMessage()    {}
func (*QueryTunnelResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{3}
}
func (m *QueryTunnelResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTunnelResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTunnelResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTunnelResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTunnelResponse.Merge(m, src)
}
func (m *QueryTunnelResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTunnelResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTunnelResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTunnelResponse proto.InternalMessageInfo

func (m *QueryTunnelResponse) GetTunnel() Tunnel {
	if m != nil {
		return m.Tunnel
	}
	return Tunnel{}
}

// QueryDepositsRequest is the request type for the Query/Deposits RPC method.
type QueryDepositsRequest struct {
	// tunnel_id is the ID of the tunnel to query deposits.
	TunnelId uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3"               json:"pagination,omitempty"`
}

func (m *QueryDepositsRequest) Reset()         { *m = QueryDepositsRequest{} }
func (m *QueryDepositsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryDepositsRequest) ProtoMessage()    {}
func (*QueryDepositsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{4}
}
func (m *QueryDepositsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDepositsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDepositsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDepositsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDepositsRequest.Merge(m, src)
}
func (m *QueryDepositsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryDepositsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDepositsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDepositsRequest proto.InternalMessageInfo

func (m *QueryDepositsRequest) GetTunnelId() uint64 {
	if m != nil {
		return m.TunnelId
	}
	return 0
}

func (m *QueryDepositsRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryDepositsResponse is the response type for the Query/Deposits RPC method.
type QueryDepositsResponse struct {
	// deposits is a list of deposits.
	Deposits []*Deposit `protobuf:"bytes,1,rep,name=deposits,proto3"   json:"deposits,omitempty"`
	// pagination defines an optional pagination for the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryDepositsResponse) Reset()         { *m = QueryDepositsResponse{} }
func (m *QueryDepositsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryDepositsResponse) ProtoMessage()    {}
func (*QueryDepositsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{5}
}
func (m *QueryDepositsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDepositsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDepositsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDepositsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDepositsResponse.Merge(m, src)
}
func (m *QueryDepositsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryDepositsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDepositsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDepositsResponse proto.InternalMessageInfo

func (m *QueryDepositsResponse) GetDeposits() []*Deposit {
	if m != nil {
		return m.Deposits
	}
	return nil
}

func (m *QueryDepositsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryDepositRequest is the request type for the Query/Deposit RPC method.
type QueryDepositRequest struct {
	// tunnel_id is the ID of the tunnel to query.
	TunnelId uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
	// depositor is the address of the depositor to query.
	Depositor string `protobuf:"bytes,2,opt,name=depositor,proto3"                json:"depositor,omitempty"`
}

func (m *QueryDepositRequest) Reset()         { *m = QueryDepositRequest{} }
func (m *QueryDepositRequest) String() string { return proto.CompactTextString(m) }
func (*QueryDepositRequest) ProtoMessage()    {}
func (*QueryDepositRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{6}
}
func (m *QueryDepositRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDepositRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDepositRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDepositRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDepositRequest.Merge(m, src)
}
func (m *QueryDepositRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryDepositRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDepositRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDepositRequest proto.InternalMessageInfo

func (m *QueryDepositRequest) GetTunnelId() uint64 {
	if m != nil {
		return m.TunnelId
	}
	return 0
}

func (m *QueryDepositRequest) GetDepositor() string {
	if m != nil {
		return m.Depositor
	}
	return ""
}

// QueryDepositResponse is the response type for the Query/Deposit RPC method.
type QueryDepositResponse struct {
	// deposit is the deposit with the given tunnel ID and depositor address.
	Deposit Deposit `protobuf:"bytes,1,opt,name=deposit,proto3" json:"deposit"`
}

func (m *QueryDepositResponse) Reset()         { *m = QueryDepositResponse{} }
func (m *QueryDepositResponse) String() string { return proto.CompactTextString(m) }
func (*QueryDepositResponse) ProtoMessage()    {}
func (*QueryDepositResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{7}
}
func (m *QueryDepositResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDepositResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDepositResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDepositResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDepositResponse.Merge(m, src)
}
func (m *QueryDepositResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryDepositResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDepositResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDepositResponse proto.InternalMessageInfo

func (m *QueryDepositResponse) GetDeposit() Deposit {
	if m != nil {
		return m.Deposit
	}
	return Deposit{}
}

// QueryPacketsRequest is the request type for the Query/Packets RPC method.
type QueryPacketsRequest struct {
	// tunnel_id is the ID of the tunnel to query packets.
	TunnelId uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3"               json:"pagination,omitempty"`
}

func (m *QueryPacketsRequest) Reset()         { *m = QueryPacketsRequest{} }
func (m *QueryPacketsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryPacketsRequest) ProtoMessage()    {}
func (*QueryPacketsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{8}
}
func (m *QueryPacketsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPacketsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPacketsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPacketsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPacketsRequest.Merge(m, src)
}
func (m *QueryPacketsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryPacketsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPacketsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPacketsRequest proto.InternalMessageInfo

func (m *QueryPacketsRequest) GetTunnelId() uint64 {
	if m != nil {
		return m.TunnelId
	}
	return 0
}

func (m *QueryPacketsRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryPacketsResponse is the response type for the Query/Packets RPC method.
type QueryPacketsResponse struct {
	// packets is a list of packets.
	Packets []*Packet `protobuf:"bytes,1,rep,name=packets,proto3"    json:"packets,omitempty"`
	// pagination defines an optional pagination for the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryPacketsResponse) Reset()         { *m = QueryPacketsResponse{} }
func (m *QueryPacketsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryPacketsResponse) ProtoMessage()    {}
func (*QueryPacketsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{9}
}
func (m *QueryPacketsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPacketsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPacketsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPacketsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPacketsResponse.Merge(m, src)
}
func (m *QueryPacketsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryPacketsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPacketsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPacketsResponse proto.InternalMessageInfo

func (m *QueryPacketsResponse) GetPackets() []*Packet {
	if m != nil {
		return m.Packets
	}
	return nil
}

func (m *QueryPacketsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryPacketRequest is the request type for the Query/Packet RPC method.
type QueryPacketRequest struct {
	// tunnel_id is the ID of the tunnel to query packets.
	TunnelId uint64 `protobuf:"varint,1,opt,name=tunnel_id,json=tunnelId,proto3" json:"tunnel_id,omitempty"`
	// sequence is the sequence of the packet to query.
	Sequence uint64 `protobuf:"varint,2,opt,name=sequence,proto3"                json:"sequence,omitempty"`
}

func (m *QueryPacketRequest) Reset()         { *m = QueryPacketRequest{} }
func (m *QueryPacketRequest) String() string { return proto.CompactTextString(m) }
func (*QueryPacketRequest) ProtoMessage()    {}
func (*QueryPacketRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{10}
}
func (m *QueryPacketRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPacketRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPacketRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPacketRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPacketRequest.Merge(m, src)
}
func (m *QueryPacketRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryPacketRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPacketRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPacketRequest proto.InternalMessageInfo

func (m *QueryPacketRequest) GetTunnelId() uint64 {
	if m != nil {
		return m.TunnelId
	}
	return 0
}

func (m *QueryPacketRequest) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

// QueryPacketResponse is the response type for the Query/Packet RPC method.
type QueryPacketResponse struct {
	// packet is the packet with the given tunnel ID and sequence.
	Packet *Packet `protobuf:"bytes,1,opt,name=packet,proto3" json:"packet,omitempty"`
}

func (m *QueryPacketResponse) Reset()         { *m = QueryPacketResponse{} }
func (m *QueryPacketResponse) String() string { return proto.CompactTextString(m) }
func (*QueryPacketResponse) ProtoMessage()    {}
func (*QueryPacketResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{11}
}
func (m *QueryPacketResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPacketResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPacketResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPacketResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPacketResponse.Merge(m, src)
}
func (m *QueryPacketResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryPacketResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPacketResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPacketResponse proto.InternalMessageInfo

func (m *QueryPacketResponse) GetPacket() *Packet {
	if m != nil {
		return m.Packet
	}
	return nil
}

// QueryTotalFeesRequest is the request type for the Query/TotalFees RPC method.
type QueryTotalFeesRequest struct {
}

func (m *QueryTotalFeesRequest) Reset()         { *m = QueryTotalFeesRequest{} }
func (m *QueryTotalFeesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryTotalFeesRequest) ProtoMessage()    {}
func (*QueryTotalFeesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{12}
}
func (m *QueryTotalFeesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTotalFeesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTotalFeesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTotalFeesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTotalFeesRequest.Merge(m, src)
}
func (m *QueryTotalFeesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryTotalFeesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTotalFeesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTotalFeesRequest proto.InternalMessageInfo

// QueryTotalFeesResponse is the response type for the Query/TotalFees RPC method.
type QueryTotalFeesResponse struct {
	// total_fees is the total fees collected by the tunnel.
	TotalFees TotalFees `protobuf:"bytes,1,opt,name=total_fees,json=totalFees,proto3" json:"total_fees"`
}

func (m *QueryTotalFeesResponse) Reset()         { *m = QueryTotalFeesResponse{} }
func (m *QueryTotalFeesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryTotalFeesResponse) ProtoMessage()    {}
func (*QueryTotalFeesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{13}
}
func (m *QueryTotalFeesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryTotalFeesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryTotalFeesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryTotalFeesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryTotalFeesResponse.Merge(m, src)
}
func (m *QueryTotalFeesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryTotalFeesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryTotalFeesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryTotalFeesResponse proto.InternalMessageInfo

func (m *QueryTotalFeesResponse) GetTotalFees() TotalFees {
	if m != nil {
		return m.TotalFees
	}
	return TotalFees{}
}

// QueryParamsRequest is the request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{14}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse is the response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params is the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f80b85392d1440ac, []int{15}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func init() {
	proto.RegisterEnum("band.tunnel.v1beta1.TunnelStatusFilter", TunnelStatusFilter_name, TunnelStatusFilter_value)
	proto.RegisterType((*QueryTunnelsRequest)(nil), "band.tunnel.v1beta1.QueryTunnelsRequest")
	proto.RegisterType((*QueryTunnelsResponse)(nil), "band.tunnel.v1beta1.QueryTunnelsResponse")
	proto.RegisterType((*QueryTunnelRequest)(nil), "band.tunnel.v1beta1.QueryTunnelRequest")
	proto.RegisterType((*QueryTunnelResponse)(nil), "band.tunnel.v1beta1.QueryTunnelResponse")
	proto.RegisterType((*QueryDepositsRequest)(nil), "band.tunnel.v1beta1.QueryDepositsRequest")
	proto.RegisterType((*QueryDepositsResponse)(nil), "band.tunnel.v1beta1.QueryDepositsResponse")
	proto.RegisterType((*QueryDepositRequest)(nil), "band.tunnel.v1beta1.QueryDepositRequest")
	proto.RegisterType((*QueryDepositResponse)(nil), "band.tunnel.v1beta1.QueryDepositResponse")
	proto.RegisterType((*QueryPacketsRequest)(nil), "band.tunnel.v1beta1.QueryPacketsRequest")
	proto.RegisterType((*QueryPacketsResponse)(nil), "band.tunnel.v1beta1.QueryPacketsResponse")
	proto.RegisterType((*QueryPacketRequest)(nil), "band.tunnel.v1beta1.QueryPacketRequest")
	proto.RegisterType((*QueryPacketResponse)(nil), "band.tunnel.v1beta1.QueryPacketResponse")
	proto.RegisterType((*QueryTotalFeesRequest)(nil), "band.tunnel.v1beta1.QueryTotalFeesRequest")
	proto.RegisterType((*QueryTotalFeesResponse)(nil), "band.tunnel.v1beta1.QueryTotalFeesResponse")
	proto.RegisterType((*QueryParamsRequest)(nil), "band.tunnel.v1beta1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "band.tunnel.v1beta1.QueryParamsResponse")
}

func init() { proto.RegisterFile("band/tunnel/v1beta1/query.proto", fileDescriptor_f80b85392d1440ac) }

var fileDescriptor_f80b85392d1440ac = []byte{
	// 944 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x56, 0xcf, 0x6f, 0x1b, 0x45,
	0x18, 0xf5, 0x84, 0x60, 0xc7, 0x53, 0x40, 0xd5, 0x34, 0xb4, 0x61, 0x13, 0x36, 0x66, 0xf9, 0x11,
	0x37, 0x81, 0x1d, 0x25, 0x01, 0x14, 0x10, 0x42, 0xb4, 0xa9, 0x8d, 0x8c, 0x42, 0x64, 0x36, 0x0e,
	0x07, 0x24, 0x64, 0xad, 0x9d, 0xa9, 0xbb, 0xc2, 0xd9, 0xd9, 0x7a, 0xd6, 0x15, 0xc5, 0x8a, 0x90,
	0x10, 0x87, 0x1e, 0x38, 0x20, 0x55, 0x02, 0x09, 0x2e, 0x48, 0x1c, 0xf9, 0x47, 0x7a, 0xac, 0xc4,
	0x85, 0x13, 0x42, 0x09, 0x7f, 0x08, 0xda, 0x99, 0x6f, 0xd6, 0x5e, 0x77, 0x59, 0x6f, 0xa4, 0xa8,
	0x37, 0x7b, 0xf6, 0x7d, 0xf3, 0xde, 0xf7, 0xbe, 0x99, 0xb7, 0x8b, 0x57, 0x3b, 0xae, 0x7f, 0x44,
	0xc3, 0xa1, 0xef, 0xb3, 0x3e, 0xbd, 0xb7, 0xd9, 0x61, 0xa1, 0xbb, 0x49, 0xef, 0x0e, 0xd9, 0xe0,
	0xbe, 0x1d, 0x0c, 0x78, 0xc8, 0xc9, 0x95, 0x08, 0x60, 0x2b, 0x80, 0x0d, 0x00, 0x63, 0xb1, 0xc7,
	0x7b, 0x5c, 0x3e, 0xa7, 0xd1, 0x2f, 0x05, 0x35, 0xd6, 0xbb, 0x5c, 0x1c, 0x73, 0x41, 0x3b, 0xae,
	0x60, 0x6a, 0x8f, 0x78, 0xc7, 0xc0, 0xed, 0x79, 0xbe, 0x1b, 0x7a, 0xdc, 0x07, 0xec, 0x4a, 0x8f,
	0xf3, 0x5e, 0x9f, 0x51, 0x37, 0xf0, 0xa8, 0xeb, 0xfb, 0x3c, 0x94, 0x0f, 0x05, 0x3c, 0xad, 0xa4,
	0xa9, 0x0a, 0xdc, 0x81, 0x7b, 0x9c, 0x89, 0x00, 0x95, 0x12, 0x61, 0xfd, 0x81, 0xf0, 0x95, 0xcf,
	0x22, 0x11, 0x2d, 0xb9, 0x2a, 0x1c, 0x76, 0x77, 0xc8, 0x44, 0x48, 0xf6, 0xf0, 0xf3, 0x22, 0x74,
	0xc3, 0xa1, 0x68, 0xdf, 0xf6, 0xfa, 0x21, 0x1b, 0x2c, 0xa1, 0x0a, 0xaa, 0xbe, 0xb0, 0xb5, 0x66,
	0xa7, 0x34, 0x6a, 0xab, 0xda, 0x03, 0x89, 0xaf, 0x4b, 0xb8, 0xf3, 0x9c, 0x98, 0xf8, 0x47, 0xea,
	0x18, 0x8f, 0x7b, 0x5b, 0x9a, 0xab, 0xa0, 0xea, 0xa5, 0xad, 0x37, 0x6c, 0x65, 0x84, 0x1d, 0x19,
	0x61, 0x2b, 0x33, 0xf5, 0x86, 0x4d, 0xb7, 0xc7, 0x40, 0x89, 0x33, 0x51, 0x69, 0xfd, 0x84, 0xf0,
	0x62, 0x52, 0xad, 0x08, 0xb8, 0x2f, 0x18, 0x79, 0x07, 0x97, 0x94, 0x26, 0xb1, 0x84, 0x2a, 0xcf,
	0x54, 0x2f, 0x6d, 0x2d, 0x67, 0x08, 0x75, 0x34, 0x96, 0x7c, 0x9c, 0xa2, 0x6b, 0x6d, 0xa6, 0x2e,
	0xc5, 0x99, 0x10, 0xb6, 0x89, 0xc9, 0x84, 0x2e, 0x6d, 0xe2, 0x32, 0x2e, 0x2b, 0xa6, 0xb6, 0x77,
	0x24, 0x0d, 0x9c, 0x77, 0x16, 0xd4, 0x42, 0xe3, 0xc8, 0x6a, 0x26, 0x8c, 0x8f, 0x3b, 0x79, 0x0f,
	0x17, 0x15, 0x44, 0x16, 0x64, 0x37, 0x72, 0x73, 0xfe, 0xd1, 0xdf, 0xab, 0x05, 0x07, 0x0a, 0xac,
	0x11, 0x98, 0x73, 0x8b, 0x05, 0x5c, 0x78, 0xa1, 0xc8, 0x23, 0xe3, 0xc2, 0x46, 0xf3, 0x0b, 0xc2,
	0x2f, 0x4e, 0xb1, 0x43, 0x47, 0x3b, 0x78, 0xe1, 0x08, 0xd6, 0x60, 0x38, 0x2b, 0xa9, 0x3d, 0x41,
	0xa1, 0x13, 0xa3, 0x2f, 0x6e, 0x3c, 0xda, 0x6b, 0x4d, 0x91, 0xc7, 0x98, 0x15, 0x5c, 0x06, 0x21,
	0x7c, 0x20, 0xb9, 0xcb, 0xce, 0x78, 0xc1, 0x6a, 0x25, 0xbd, 0x8e, 0x9b, 0xfd, 0x00, 0x97, 0x00,
	0x04, 0xf3, 0xcb, 0xec, 0x15, 0x06, 0xa8, 0x4b, 0xac, 0x6f, 0x40, 0x67, 0xd3, 0xed, 0x7e, 0xc5,
	0x9e, 0xf2, 0x00, 0xe3, 0xbb, 0x15, 0x93, 0x8f, 0xef, 0x56, 0xa0, 0x96, 0x32, 0xef, 0x96, 0x2a,
	0x73, 0x34, 0xf6, 0xe2, 0x86, 0xf7, 0x29, 0xdc, 0x2d, 0x20, 0xc8, 0xe3, 0x89, 0x81, 0x17, 0x44,
	0x84, 0xf3, 0xbb, 0x4c, 0x32, 0xcf, 0x3b, 0xf1, 0x7f, 0xeb, 0x93, 0x84, 0xc7, 0x71, 0x97, 0xdb,
	0xb8, 0xa8, 0x94, 0x67, 0xde, 0x3b, 0x28, 0x02, 0xa8, 0x75, 0x0d, 0xce, 0x7c, 0x8b, 0x87, 0x6e,
	0xbf, 0xce, 0x98, 0x9e, 0x98, 0xf5, 0x25, 0xbe, 0x3a, 0xfd, 0x00, 0x78, 0x76, 0x31, 0x0e, 0xa3,
	0xc5, 0xf6, 0x6d, 0xc6, 0x04, 0x70, 0x99, 0xe9, 0x77, 0x5c, 0xd7, 0xc2, 0x29, 0x29, 0x87, 0x7a,
	0xc1, 0x5a, 0x8c, 0x2d, 0x89, 0xc2, 0x5e, 0x93, 0x36, 0xe3, 0xce, 0xd4, 0xea, 0x38, 0x51, 0xd4,
	0x4b, 0x61, 0x46, 0x67, 0x11, 0x44, 0x27, 0x8a, 0x2a, 0x58, 0xff, 0x1e, 0x61, 0xf2, 0x64, 0xb8,
	0x93, 0xd7, 0x70, 0xa5, 0x75, 0xb8, 0xbf, 0x5f, 0xdb, 0x6b, 0x1f, 0xb4, 0x6e, 0xb4, 0x0e, 0x0f,
	0xda, 0xf5, 0xc6, 0x5e, 0xab, 0xe6, 0xb4, 0x0f, 0xf7, 0x0f, 0x9a, 0xb5, 0xdd, 0x46, 0xbd, 0x51,
	0xbb, 0x75, 0xb9, 0x40, 0x56, 0xf1, 0x72, 0x2a, 0xea, 0xc6, 0x6e, 0xab, 0xf1, 0x79, 0xed, 0x32,
	0x22, 0xaf, 0xe0, 0x97, 0x53, 0x01, 0x8d, 0x7d, 0x80, 0xcc, 0x19, 0xf3, 0x0f, 0x7e, 0x37, 0x0b,
	0x5b, 0x3f, 0x94, 0xf1, 0xb3, 0xb2, 0x33, 0xf2, 0x2d, 0x2e, 0x41, 0xf4, 0x93, 0x6a, 0x6a, 0x1b,
	0x29, 0xef, 0x32, 0xe3, 0x7a, 0x0e, 0xa4, 0xf2, 0xca, 0x5a, 0xfd, 0xee, 0xcf, 0x7f, 0x1f, 0xce,
	0xbd, 0x44, 0xae, 0xa5, 0xbf, 0x34, 0x05, 0x79, 0x80, 0x70, 0x51, 0x15, 0x91, 0xb5, 0x59, 0xdb,
	0x6a, 0xfe, 0xea, 0x6c, 0x20, 0xd0, 0x6f, 0x48, 0xfa, 0xd7, 0xc9, 0xab, 0xff, 0x43, 0x4f, 0x47,
	0xf1, 0x99, 0x3f, 0x21, 0x3f, 0x23, 0xbc, 0xa0, 0xc3, 0x96, 0x64, 0xf4, 0x38, 0xf5, 0x3a, 0x30,
	0xd6, 0xf3, 0x40, 0x41, 0xd0, 0xdb, 0x52, 0x90, 0x4d, 0xde, 0xcc, 0x21, 0x88, 0xc6, 0xb9, 0xfd,
	0x1b, 0xc2, 0x25, 0xd8, 0x2a, 0x6b, 0x4c, 0xc9, 0x34, 0x36, 0xae, 0xe7, 0x40, 0x82, 0xac, 0x8f,
	0xa4, 0xac, 0xf7, 0xc9, 0xce, 0x79, 0x64, 0xd1, 0x51, 0x1c, 0xdf, 0x27, 0xe4, 0x21, 0xc2, 0x25,
	0x08, 0xba, 0x2c, 0x89, 0xc9, 0x20, 0xce, 0x92, 0x38, 0x95, 0x9a, 0xd6, 0xb6, 0x94, 0xf8, 0x16,
	0xd9, 0xc8, 0x23, 0x51, 0x67, 0xe6, 0xaf, 0x08, 0x17, 0xd5, 0x46, 0x59, 0xa7, 0x2b, 0x11, 0x84,
	0x46, 0x75, 0x36, 0x10, 0x24, 0x7d, 0x28, 0x25, 0xed, 0x90, 0x77, 0xcf, 0x21, 0x89, 0x8e, 0x74,
	0x70, 0x9e, 0x44, 0x67, 0xbf, 0x1c, 0x87, 0x12, 0xc9, 0x38, 0x46, 0xd3, 0x71, 0x68, 0x6c, 0xe4,
	0xc2, 0x82, 0x4c, 0x4b, 0xca, 0x5c, 0x21, 0xc6, 0x13, 0x32, 0xe3, 0xdc, 0x24, 0xa3, 0xc8, 0xa7,
	0x28, 0xa2, 0xb2, 0x7d, 0x9a, 0x48, 0xc7, 0x6c, 0x9f, 0x26, 0x03, 0xd3, 0x32, 0xa5, 0x80, 0x25,
	0x72, 0x35, 0xfd, 0xdb, 0xfa, 0x66, 0xe3, 0xd1, 0xa9, 0x89, 0x1e, 0x9f, 0x9a, 0xe8, 0x9f, 0x53,
	0x13, 0xfd, 0x78, 0x66, 0x16, 0x1e, 0x9f, 0x99, 0x85, 0xbf, 0xce, 0xcc, 0xc2, 0x17, 0xb4, 0xe7,
	0x85, 0x77, 0x86, 0x1d, 0xbb, 0xcb, 0x8f, 0x69, 0xc4, 0x26, 0xbf, 0xb1, 0xbb, 0xbc, 0x4f, 0xbb,
	0x77, 0x5c, 0xcf, 0xa7, 0xf7, 0xb6, 0xe9, 0xd7, 0x7a, 0xcf, 0xf0, 0x7e, 0xc0, 0x44, 0xa7, 0x28,
	0x11, 0xdb, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0x75, 0xeb, 0x1c, 0x48, 0x61, 0x0c, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Tunnels is a RPC method that returns all tunnels.
	Tunnels(ctx context.Context, in *QueryTunnelsRequest, opts ...grpc.CallOption) (*QueryTunnelsResponse, error)
	// Tunnel is a RPC method that returns a tunnel by its ID.
	Tunnel(ctx context.Context, in *QueryTunnelRequest, opts ...grpc.CallOption) (*QueryTunnelResponse, error)
	// Deposits queries all deposits of a single tunnel.
	Deposits(ctx context.Context, in *QueryDepositsRequest, opts ...grpc.CallOption) (*QueryDepositsResponse, error)
	// Deposit queries single deposit information based tunnelID, depositAddr.
	Deposit(ctx context.Context, in *QueryDepositRequest, opts ...grpc.CallOption) (*QueryDepositResponse, error)
	// Packets is a RPC method that returns all packets of a tunnel.
	Packets(ctx context.Context, in *QueryPacketsRequest, opts ...grpc.CallOption) (*QueryPacketsResponse, error)
	// Packet is a RPC method that returns a packet by its tunnel ID and sequence.
	Packet(ctx context.Context, in *QueryPacketRequest, opts ...grpc.CallOption) (*QueryPacketResponse, error)
	// TotalFees is a RPC method that returns the total fees collected by the tunnel
	TotalFees(ctx context.Context, in *QueryTotalFeesRequest, opts ...grpc.CallOption) (*QueryTotalFeesResponse, error)
	// Params is a RPC method that returns all parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Tunnels(
	ctx context.Context,
	in *QueryTunnelsRequest,
	opts ...grpc.CallOption,
) (*QueryTunnelsResponse, error) {
	out := new(QueryTunnelsResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Tunnels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Tunnel(
	ctx context.Context,
	in *QueryTunnelRequest,
	opts ...grpc.CallOption,
) (*QueryTunnelResponse, error) {
	out := new(QueryTunnelResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Tunnel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Deposits(
	ctx context.Context,
	in *QueryDepositsRequest,
	opts ...grpc.CallOption,
) (*QueryDepositsResponse, error) {
	out := new(QueryDepositsResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Deposits", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Deposit(
	ctx context.Context,
	in *QueryDepositRequest,
	opts ...grpc.CallOption,
) (*QueryDepositResponse, error) {
	out := new(QueryDepositResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Deposit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Packets(
	ctx context.Context,
	in *QueryPacketsRequest,
	opts ...grpc.CallOption,
) (*QueryPacketsResponse, error) {
	out := new(QueryPacketsResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Packets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Packet(
	ctx context.Context,
	in *QueryPacketRequest,
	opts ...grpc.CallOption,
) (*QueryPacketResponse, error) {
	out := new(QueryPacketResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Packet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) TotalFees(
	ctx context.Context,
	in *QueryTotalFeesRequest,
	opts ...grpc.CallOption,
) (*QueryTotalFeesResponse, error) {
	out := new(QueryTotalFeesResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/TotalFees", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(
	ctx context.Context,
	in *QueryParamsRequest,
	opts ...grpc.CallOption,
) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Tunnels is a RPC method that returns all tunnels.
	Tunnels(context.Context, *QueryTunnelsRequest) (*QueryTunnelsResponse, error)
	// Tunnel is a RPC method that returns a tunnel by its ID.
	Tunnel(context.Context, *QueryTunnelRequest) (*QueryTunnelResponse, error)
	// Deposits queries all deposits of a single tunnel.
	Deposits(context.Context, *QueryDepositsRequest) (*QueryDepositsResponse, error)
	// Deposit queries single deposit information based tunnelID, depositAddr.
	Deposit(context.Context, *QueryDepositRequest) (*QueryDepositResponse, error)
	// Packets is a RPC method that returns all packets of a tunnel.
	Packets(context.Context, *QueryPacketsRequest) (*QueryPacketsResponse, error)
	// Packet is a RPC method that returns a packet by its tunnel ID and sequence.
	Packet(context.Context, *QueryPacketRequest) (*QueryPacketResponse, error)
	// TotalFees is a RPC method that returns the total fees collected by the tunnel
	TotalFees(context.Context, *QueryTotalFeesRequest) (*QueryTotalFeesResponse, error)
	// Params is a RPC method that returns all parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Tunnels(
	ctx context.Context,
	req *QueryTunnelsRequest,
) (*QueryTunnelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Tunnels not implemented")
}
func (*UnimplementedQueryServer) Tunnel(ctx context.Context, req *QueryTunnelRequest) (*QueryTunnelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Tunnel not implemented")
}

func (*UnimplementedQueryServer) Deposits(
	ctx context.Context,
	req *QueryDepositsRequest,
) (*QueryDepositsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposits not implemented")
}

func (*UnimplementedQueryServer) Deposit(
	ctx context.Context,
	req *QueryDepositRequest,
) (*QueryDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}

func (*UnimplementedQueryServer) Packets(
	ctx context.Context,
	req *QueryPacketsRequest,
) (*QueryPacketsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Packets not implemented")
}
func (*UnimplementedQueryServer) Packet(ctx context.Context, req *QueryPacketRequest) (*QueryPacketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Packet not implemented")
}

func (*UnimplementedQueryServer) TotalFees(
	ctx context.Context,
	req *QueryTotalFeesRequest,
) (*QueryTotalFeesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TotalFees not implemented")
}
func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Tunnels_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryTunnelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Tunnels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.tunnel.v1beta1.Query/Tunnels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Tunnels(ctx, req.(*QueryTunnelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Tunnel_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryTunnelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Tunnel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.tunnel.v1beta1.Query/Tunnel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Tunnel(ctx, req.(*QueryTunnelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Deposits_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryDepositsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Deposits(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.tunnel.v1beta1.Query/Deposits",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Deposits(ctx, req.(*QueryDepositsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Deposit_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.tunnel.v1beta1.Query/Deposit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Deposit(ctx, req.(*QueryDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Packets_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryPacketsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Packets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.tunnel.v1beta1.Query/Packets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Packets(ctx, req.(*QueryPacketsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Packet_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryPacketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Packet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.tunnel.v1beta1.Query/Packet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Packet(ctx, req.(*QueryPacketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TotalFees_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryTotalFeesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TotalFees(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.tunnel.v1beta1.Query/TotalFees",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TotalFees(ctx, req.(*QueryTotalFeesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.tunnel.v1beta1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "band.tunnel.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Tunnels",
			Handler:    _Query_Tunnels_Handler,
		},
		{
			MethodName: "Tunnel",
			Handler:    _Query_Tunnel_Handler,
		},
		{
			MethodName: "Deposits",
			Handler:    _Query_Deposits_Handler,
		},
		{
			MethodName: "Deposit",
			Handler:    _Query_Deposit_Handler,
		},
		{
			MethodName: "Packets",
			Handler:    _Query_Packets_Handler,
		},
		{
			MethodName: "Packet",
			Handler:    _Query_Packet_Handler,
		},
		{
			MethodName: "TotalFees",
			Handler:    _Query_TotalFees_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "band/tunnel/v1beta1/query.proto",
}

func (m *QueryTunnelsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTunnelsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTunnelsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.StatusFilter != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.StatusFilter))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryTunnelsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTunnelsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTunnelsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Tunnels) > 0 {
		for iNdEx := len(m.Tunnels) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Tunnels[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryTunnelRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTunnelRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTunnelRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TunnelId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.TunnelId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryTunnelResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTunnelResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTunnelResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Tunnel.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryDepositsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDepositsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDepositsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.TunnelId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.TunnelId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryDepositsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDepositsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDepositsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Deposits) > 0 {
		for iNdEx := len(m.Deposits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Deposits[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryDepositRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDepositRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDepositRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Depositor) > 0 {
		i -= len(m.Depositor)
		copy(dAtA[i:], m.Depositor)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Depositor)))
		i--
		dAtA[i] = 0x12
	}
	if m.TunnelId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.TunnelId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryDepositResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDepositResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDepositResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Deposit.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryPacketsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPacketsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPacketsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.TunnelId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.TunnelId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryPacketsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPacketsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPacketsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Packets) > 0 {
		for iNdEx := len(m.Packets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Packets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryPacketRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPacketRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPacketRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Sequence != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x10
	}
	if m.TunnelId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.TunnelId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryPacketResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPacketResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPacketResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Packet != nil {
		{
			size, err := m.Packet.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryTotalFeesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTotalFeesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTotalFeesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryTotalFeesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryTotalFeesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryTotalFeesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.TotalFees.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryTunnelsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.StatusFilter != 0 {
		n += 1 + sovQuery(uint64(m.StatusFilter))
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryTunnelsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Tunnels) > 0 {
		for _, e := range m.Tunnels {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryTunnelRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelId != 0 {
		n += 1 + sovQuery(uint64(m.TunnelId))
	}
	return n
}

func (m *QueryTunnelResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Tunnel.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryDepositsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelId != 0 {
		n += 1 + sovQuery(uint64(m.TunnelId))
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryDepositsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Deposits) > 0 {
		for _, e := range m.Deposits {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryDepositRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelId != 0 {
		n += 1 + sovQuery(uint64(m.TunnelId))
	}
	l = len(m.Depositor)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryDepositResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Deposit.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryPacketsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelId != 0 {
		n += 1 + sovQuery(uint64(m.TunnelId))
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryPacketsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Packets) > 0 {
		for _, e := range m.Packets {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryPacketRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TunnelId != 0 {
		n += 1 + sovQuery(uint64(m.TunnelId))
	}
	if m.Sequence != 0 {
		n += 1 + sovQuery(uint64(m.Sequence))
	}
	return n
}

func (m *QueryPacketResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Packet != nil {
		l = m.Packet.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryTotalFeesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryTotalFeesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.TotalFees.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryTunnelsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTunnelsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTunnelsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatusFilter", wireType)
			}
			m.StatusFilter = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StatusFilter |= TunnelStatusFilter(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTunnelsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTunnelsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTunnelsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tunnels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tunnels = append(m.Tunnels, &Tunnel{})
			if err := m.Tunnels[len(m.Tunnels)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTunnelRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTunnelRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTunnelRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelId", wireType)
			}
			m.TunnelId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TunnelId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTunnelResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTunnelResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTunnelResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tunnel", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Tunnel.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryDepositsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryDepositsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDepositsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelId", wireType)
			}
			m.TunnelId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TunnelId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryDepositsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryDepositsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDepositsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposits", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Deposits = append(m.Deposits, &Deposit{})
			if err := m.Deposits[len(m.Deposits)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryDepositRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryDepositRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDepositRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelId", wireType)
			}
			m.TunnelId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TunnelId |= uint64(b&0x7F) << shift
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
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Depositor = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryDepositResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryDepositResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDepositResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposit", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Deposit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryPacketsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryPacketsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPacketsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelId", wireType)
			}
			m.TunnelId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TunnelId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryPacketsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryPacketsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPacketsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Packets", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Packets = append(m.Packets, &Packet{})
			if err := m.Packets[len(m.Packets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryPacketRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryPacketRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPacketRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TunnelId", wireType)
			}
			m.TunnelId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TunnelId |= uint64(b&0x7F) << shift
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
					return ErrIntOverflowQuery
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
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryPacketResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryPacketResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPacketResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Packet", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Packet == nil {
				m.Packet = &Packet{}
			}
			if err := m.Packet.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTotalFeesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTotalFeesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTotalFeesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryTotalFeesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryTotalFeesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryTotalFeesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalFees", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalFees.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
