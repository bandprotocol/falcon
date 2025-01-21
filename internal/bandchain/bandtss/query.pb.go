package types

import (
	"context"
	"fmt"
	"io"
	"math"
	math_bits "math/bits"
	"time"

	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types1 "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	github_com_bandprotocol_chain_v3_pkg_tss "github.com/bandprotocol/falcon/internal/bandchain/tss"
	types "github.com/bandprotocol/falcon/internal/bandchain/tss"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MemberStatusFilter defines the query options for filtering members by their active status.
type MemberStatusFilter int32

const (
	// MEMBER_STATUS_FILTER_UNSPECIFIED defines a filter for unspecified active status.
	MEMBER_STATUS_FILTER_UNSPECIFIED MemberStatusFilter = 0
	// MEMBER_STATUS_FILTER_ACTIVE defines a filter for active status.
	MEMBER_STATUS_FILTER_ACTIVE MemberStatusFilter = 1
	// MEMBER_STATUS_FILTER_INACTIVE defines a filter for inactive status.
	MEMBER_STATUS_FILTER_INACTIVE MemberStatusFilter = 2
)

var MemberStatusFilter_name = map[int32]string{
	0: "MEMBER_STATUS_FILTER_UNSPECIFIED",
	1: "MEMBER_STATUS_FILTER_ACTIVE",
	2: "MEMBER_STATUS_FILTER_INACTIVE",
}

var MemberStatusFilter_value = map[string]int32{
	"MEMBER_STATUS_FILTER_UNSPECIFIED": 0,
	"MEMBER_STATUS_FILTER_ACTIVE":      1,
	"MEMBER_STATUS_FILTER_INACTIVE":    2,
}

func (x MemberStatusFilter) String() string {
	return proto.EnumName(MemberStatusFilter_name, int32(x))
}

func (MemberStatusFilter) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{0}
}

// QueryCountsRequest is request type for the Query/Count RPC method.
type QueryCountsRequest struct {
}

func (m *QueryCountsRequest) Reset()         { *m = QueryCountsRequest{} }
func (m *QueryCountsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryCountsRequest) ProtoMessage()    {}
func (*QueryCountsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{0}
}
func (m *QueryCountsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCountsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCountsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCountsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCountsRequest.Merge(m, src)
}
func (m *QueryCountsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryCountsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCountsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCountsRequest proto.InternalMessageInfo

// QueryCountsResponse is response type for the Query/Count RPC method.
type QueryCountsResponse struct {
	// signing_count is total number of signing request submitted to bandtss module
	SigningCount uint64 `protobuf:"varint,1,opt,name=signing_count,json=signingCount,proto3" json:"signing_count,omitempty"`
}

func (m *QueryCountsResponse) Reset()         { *m = QueryCountsResponse{} }
func (m *QueryCountsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryCountsResponse) ProtoMessage()    {}
func (*QueryCountsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{1}
}
func (m *QueryCountsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCountsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCountsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCountsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCountsResponse.Merge(m, src)
}
func (m *QueryCountsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryCountsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCountsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCountsResponse proto.InternalMessageInfo

func (m *QueryCountsResponse) GetSigningCount() uint64 {
	if m != nil {
		return m.SigningCount
	}
	return 0
}

// QueryMembersRequest is the request type for the Query/Members RPC method.
type QueryMembersRequest struct {
	// status define type of filter on member's status.
	Status MemberStatusFilter `protobuf:"varint,1,opt,name=status,proto3,enum=band.bandtss.v1beta1.MemberStatusFilter" json:"status,omitempty"`
	// is_incoming_group is a flag to indicate whether user query members in the incoming group
	// or the current group.
	IsIncomingGroup bool `protobuf:"varint,2,opt,name=is_incoming_group,json=isIncomingGroup,proto3"              json:"is_incoming_group,omitempty"`
	// pagination defines pagination settings for the request.
	Pagination *query.PageRequest `protobuf:"bytes,3,opt,name=pagination,proto3"                                           json:"pagination,omitempty"`
}

func (m *QueryMembersRequest) Reset()         { *m = QueryMembersRequest{} }
func (m *QueryMembersRequest) String() string { return proto.CompactTextString(m) }
func (*QueryMembersRequest) ProtoMessage()    {}
func (*QueryMembersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{2}
}
func (m *QueryMembersRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMembersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMembersRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMembersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMembersRequest.Merge(m, src)
}
func (m *QueryMembersRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryMembersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMembersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMembersRequest proto.InternalMessageInfo

func (m *QueryMembersRequest) GetStatus() MemberStatusFilter {
	if m != nil {
		return m.Status
	}
	return MEMBER_STATUS_FILTER_UNSPECIFIED
}

func (m *QueryMembersRequest) GetIsIncomingGroup() bool {
	if m != nil {
		return m.IsIncomingGroup
	}
	return false
}

func (m *QueryMembersRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryMembersResponse is the response type for the Query/Members RPC method.
type QueryMembersResponse struct {
	// members are those individuals who correspond to the provided is_active status.
	Members []*Member `protobuf:"bytes,1,rep,name=members,proto3"    json:"members,omitempty"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryMembersResponse) Reset()         { *m = QueryMembersResponse{} }
func (m *QueryMembersResponse) String() string { return proto.CompactTextString(m) }
func (*QueryMembersResponse) ProtoMessage()    {}
func (*QueryMembersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{3}
}
func (m *QueryMembersResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMembersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMembersResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMembersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMembersResponse.Merge(m, src)
}
func (m *QueryMembersResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryMembersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMembersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMembersResponse proto.InternalMessageInfo

func (m *QueryMembersResponse) GetMembers() []*Member {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *QueryMembersResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryMemberRequest is the request type for the Query/Member RPC method.
type QueryMemberRequest struct {
	// address is the member address.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryMemberRequest) Reset()         { *m = QueryMemberRequest{} }
func (m *QueryMemberRequest) String() string { return proto.CompactTextString(m) }
func (*QueryMemberRequest) ProtoMessage()    {}
func (*QueryMemberRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{4}
}
func (m *QueryMemberRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMemberRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMemberRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMemberRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMemberRequest.Merge(m, src)
}
func (m *QueryMemberRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryMemberRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMemberRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMemberRequest proto.InternalMessageInfo

func (m *QueryMemberRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

// QueryMemberResponse is the response type for the Query/Member RPC method.
type QueryMemberResponse struct {
	// current_group_member is the member detail.
	CurrentGroupMember Member `protobuf:"bytes,1,opt,name=current_group_member,json=currentGroupMember,proto3"   json:"current_group_member"`
	// incoming_group_member is the member detail.
	IncomingGroupMember Member `protobuf:"bytes,2,opt,name=incoming_group_member,json=incomingGroupMember,proto3" json:"incoming_group_member"`
}

func (m *QueryMemberResponse) Reset()         { *m = QueryMemberResponse{} }
func (m *QueryMemberResponse) String() string { return proto.CompactTextString(m) }
func (*QueryMemberResponse) ProtoMessage()    {}
func (*QueryMemberResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{5}
}
func (m *QueryMemberResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMemberResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMemberResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMemberResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMemberResponse.Merge(m, src)
}
func (m *QueryMemberResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryMemberResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMemberResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMemberResponse proto.InternalMessageInfo

func (m *QueryMemberResponse) GetCurrentGroupMember() Member {
	if m != nil {
		return m.CurrentGroupMember
	}
	return Member{}
}

func (m *QueryMemberResponse) GetIncomingGroupMember() Member {
	if m != nil {
		return m.IncomingGroupMember
	}
	return Member{}
}

// QueryCurrentGroupRequest is the request type for the Query/CurrentGroup RPC method.
type QueryCurrentGroupRequest struct {
}

func (m *QueryCurrentGroupRequest) Reset()         { *m = QueryCurrentGroupRequest{} }
func (m *QueryCurrentGroupRequest) String() string { return proto.CompactTextString(m) }
func (*QueryCurrentGroupRequest) ProtoMessage()    {}
func (*QueryCurrentGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{6}
}
func (m *QueryCurrentGroupRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCurrentGroupRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCurrentGroupRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCurrentGroupRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCurrentGroupRequest.Merge(m, src)
}
func (m *QueryCurrentGroupRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryCurrentGroupRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCurrentGroupRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCurrentGroupRequest proto.InternalMessageInfo

// QueryCurrentGroupResponse is the response type for the Query/CurrentGroup RPC method.
type QueryCurrentGroupResponse struct {
	// group_id is the ID of the current group.
	GroupID github_com_bandprotocol_chain_v3_pkg_tss.GroupID `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID" json:"group_id,omitempty"`
	// size is the number of members in the group.
	Size_ uint64 `protobuf:"varint,2,opt,name=size,proto3"                                                                                         json:"size,omitempty"`
	// threshold is the minimum number of members needed to generate a valid signature.
	Threshold uint64 `protobuf:"varint,3,opt,name=threshold,proto3"                                                                                    json:"threshold,omitempty"`
	// pub_key is the public key generated by the group.
	PubKey github_com_bandprotocol_chain_v3_pkg_tss.Point `protobuf:"bytes,4,opt,name=pub_key,json=pubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"      json:"pub_key,omitempty"`
	// status is the status of the current group.
	Status types.GroupStatus `protobuf:"varint,5,opt,name=status,proto3,enum=band.tss.v1beta1.GroupStatus"                                                     json:"status,omitempty"`
	// active_time is the timestamp at which the group becomes the current group of the module.
	ActiveTime time.Time `protobuf:"bytes,6,opt,name=active_time,json=activeTime,proto3,stdtime"                                                           json:"active_time"`
}

func (m *QueryCurrentGroupResponse) Reset()         { *m = QueryCurrentGroupResponse{} }
func (m *QueryCurrentGroupResponse) String() string { return proto.CompactTextString(m) }
func (*QueryCurrentGroupResponse) ProtoMessage()    {}
func (*QueryCurrentGroupResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{7}
}
func (m *QueryCurrentGroupResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryCurrentGroupResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryCurrentGroupResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryCurrentGroupResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryCurrentGroupResponse.Merge(m, src)
}
func (m *QueryCurrentGroupResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryCurrentGroupResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryCurrentGroupResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryCurrentGroupResponse proto.InternalMessageInfo

func (m *QueryCurrentGroupResponse) GetGroupID() github_com_bandprotocol_chain_v3_pkg_tss.GroupID {
	if m != nil {
		return m.GroupID
	}
	return 0
}

func (m *QueryCurrentGroupResponse) GetSize_() uint64 {
	if m != nil {
		return m.Size_
	}
	return 0
}

func (m *QueryCurrentGroupResponse) GetThreshold() uint64 {
	if m != nil {
		return m.Threshold
	}
	return 0
}

func (m *QueryCurrentGroupResponse) GetPubKey() github_com_bandprotocol_chain_v3_pkg_tss.Point {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *QueryCurrentGroupResponse) GetStatus() types.GroupStatus {
	if m != nil {
		return m.Status
	}
	return types.GROUP_STATUS_UNSPECIFIED
}

func (m *QueryCurrentGroupResponse) GetActiveTime() time.Time {
	if m != nil {
		return m.ActiveTime
	}
	return time.Time{}
}

// QueryIncomingGroupRequest is the request type for the Query/IncomingGroup RPC method.
type QueryIncomingGroupRequest struct {
}

func (m *QueryIncomingGroupRequest) Reset()         { *m = QueryIncomingGroupRequest{} }
func (m *QueryIncomingGroupRequest) String() string { return proto.CompactTextString(m) }
func (*QueryIncomingGroupRequest) ProtoMessage()    {}
func (*QueryIncomingGroupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{8}
}
func (m *QueryIncomingGroupRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryIncomingGroupRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryIncomingGroupRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryIncomingGroupRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryIncomingGroupRequest.Merge(m, src)
}
func (m *QueryIncomingGroupRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryIncomingGroupRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryIncomingGroupRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryIncomingGroupRequest proto.InternalMessageInfo

// QueryIncomingGroupResponse is the response type for the Query/IncomingGroup RPC method.
type QueryIncomingGroupResponse struct {
	// group_id is the ID of the incoming group.
	GroupID github_com_bandprotocol_chain_v3_pkg_tss.GroupID `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID" json:"group_id,omitempty"`
	// size is the number of members in the group.
	Size_ uint64 `protobuf:"varint,2,opt,name=size,proto3"                                                                                         json:"size,omitempty"`
	// threshold is the minimum number of members needed to generate a valid signature.
	Threshold uint64 `protobuf:"varint,3,opt,name=threshold,proto3"                                                                                    json:"threshold,omitempty"`
	// pub_key is the public key generated by the group.
	PubKey github_com_bandprotocol_chain_v3_pkg_tss.Point `protobuf:"bytes,4,opt,name=pub_key,json=pubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"      json:"pub_key,omitempty"`
	// status is the status of the incoming group.
	Status types.GroupStatus `protobuf:"varint,5,opt,name=status,proto3,enum=band.tss.v1beta1.GroupStatus"                                                     json:"status,omitempty"`
}

func (m *QueryIncomingGroupResponse) Reset()         { *m = QueryIncomingGroupResponse{} }
func (m *QueryIncomingGroupResponse) String() string { return proto.CompactTextString(m) }
func (*QueryIncomingGroupResponse) ProtoMessage()    {}
func (*QueryIncomingGroupResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{9}
}
func (m *QueryIncomingGroupResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryIncomingGroupResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryIncomingGroupResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryIncomingGroupResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryIncomingGroupResponse.Merge(m, src)
}
func (m *QueryIncomingGroupResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryIncomingGroupResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryIncomingGroupResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryIncomingGroupResponse proto.InternalMessageInfo

func (m *QueryIncomingGroupResponse) GetGroupID() github_com_bandprotocol_chain_v3_pkg_tss.GroupID {
	if m != nil {
		return m.GroupID
	}
	return 0
}

func (m *QueryIncomingGroupResponse) GetSize_() uint64 {
	if m != nil {
		return m.Size_
	}
	return 0
}

func (m *QueryIncomingGroupResponse) GetThreshold() uint64 {
	if m != nil {
		return m.Threshold
	}
	return 0
}

func (m *QueryIncomingGroupResponse) GetPubKey() github_com_bandprotocol_chain_v3_pkg_tss.Point {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *QueryIncomingGroupResponse) GetStatus() types.GroupStatus {
	if m != nil {
		return m.Status
	}
	return types.GROUP_STATUS_UNSPECIFIED
}

// QuerySingingRequest is the request type for the Query/Signing RPC method.
type QuerySigningRequest struct {
	// signing_id is the ID of the signing request.
	SigningId uint64 `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3" json:"signing_id,omitempty"`
}

func (m *QuerySigningRequest) Reset()         { *m = QuerySigningRequest{} }
func (m *QuerySigningRequest) String() string { return proto.CompactTextString(m) }
func (*QuerySigningRequest) ProtoMessage()    {}
func (*QuerySigningRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{10}
}
func (m *QuerySigningRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySigningRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySigningRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySigningRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySigningRequest.Merge(m, src)
}
func (m *QuerySigningRequest) XXX_Size() int {
	return m.Size()
}
func (m *QuerySigningRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySigningRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySigningRequest proto.InternalMessageInfo

func (m *QuerySigningRequest) GetSigningId() uint64 {
	if m != nil {
		return m.SigningId
	}
	return 0
}

// QuerySigningResponse is the response type for the Query/Signing RPC method.
type QuerySigningResponse struct {
	// fee_per_signer is the tokens that will be paid per signer for this bandtss signing.
	FeePerSigner github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=fee_per_signer,json=feePerSigner,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"fee_per_signer"`
	// requester is the address of requester who paid for bandtss signing.
	Requester string `protobuf:"bytes,2,opt,name=requester,proto3"                                                                              json:"requester,omitempty"`
	// current_group_signing_result is the signing result from the current group.
	CurrentGroupSigningResult *types.SigningResult `protobuf:"bytes,3,opt,name=current_group_signing_result,json=currentGroupSigningResult,proto3"                            json:"current_group_signing_result,omitempty"`
	// incoming_group_signing_result is the signing result from the incoming group.
	IncomingGroupSigningResult *types.SigningResult `protobuf:"bytes,4,opt,name=incoming_group_signing_result,json=incomingGroupSigningResult,proto3"                          json:"incoming_group_signing_result,omitempty"`
}

func (m *QuerySigningResponse) Reset()         { *m = QuerySigningResponse{} }
func (m *QuerySigningResponse) String() string { return proto.CompactTextString(m) }
func (*QuerySigningResponse) ProtoMessage()    {}
func (*QuerySigningResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{11}
}
func (m *QuerySigningResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuerySigningResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuerySigningResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuerySigningResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuerySigningResponse.Merge(m, src)
}
func (m *QuerySigningResponse) XXX_Size() int {
	return m.Size()
}
func (m *QuerySigningResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QuerySigningResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QuerySigningResponse proto.InternalMessageInfo

func (m *QuerySigningResponse) GetFeePerSigner() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.FeePerSigner
	}
	return nil
}

func (m *QuerySigningResponse) GetRequester() string {
	if m != nil {
		return m.Requester
	}
	return ""
}

func (m *QuerySigningResponse) GetCurrentGroupSigningResult() *types.SigningResult {
	if m != nil {
		return m.CurrentGroupSigningResult
	}
	return nil
}

func (m *QuerySigningResponse) GetIncomingGroupSigningResult() *types.SigningResult {
	if m != nil {
		return m.IncomingGroupSigningResult
	}
	return nil
}

// QueryGroupTransitionRequest is the request type for the Query/GroupTransition RPC method.
type QueryGroupTransitionRequest struct {
}

func (m *QueryGroupTransitionRequest) Reset()         { *m = QueryGroupTransitionRequest{} }
func (m *QueryGroupTransitionRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGroupTransitionRequest) ProtoMessage()    {}
func (*QueryGroupTransitionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{12}
}
func (m *QueryGroupTransitionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGroupTransitionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGroupTransitionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGroupTransitionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGroupTransitionRequest.Merge(m, src)
}
func (m *QueryGroupTransitionRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGroupTransitionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGroupTransitionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGroupTransitionRequest proto.InternalMessageInfo

// QueryGroupTransitionResponse is the response type for the Query/GroupTransition RPC method.
type QueryGroupTransitionResponse struct {
	// group_transition is the group transition information.
	GroupTransition *GroupTransition `protobuf:"bytes,1,opt,name=group_transition,json=groupTransition,proto3" json:"group_transition,omitempty"`
}

func (m *QueryGroupTransitionResponse) Reset()         { *m = QueryGroupTransitionResponse{} }
func (m *QueryGroupTransitionResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGroupTransitionResponse) ProtoMessage()    {}
func (*QueryGroupTransitionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{13}
}
func (m *QueryGroupTransitionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGroupTransitionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGroupTransitionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGroupTransitionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGroupTransitionResponse.Merge(m, src)
}
func (m *QueryGroupTransitionResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGroupTransitionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGroupTransitionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGroupTransitionResponse proto.InternalMessageInfo

func (m *QueryGroupTransitionResponse) GetGroupTransition() *GroupTransition {
	if m != nil {
		return m.GroupTransition
	}
	return nil
}

// QueryParamsRequest is request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{14}
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

// QueryParamsResponse is response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// Params is the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d619290a87c09054, []int{15}
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
	proto.RegisterEnum("band.bandtss.v1beta1.MemberStatusFilter", MemberStatusFilter_name, MemberStatusFilter_value)
	proto.RegisterType((*QueryCountsRequest)(nil), "band.bandtss.v1beta1.QueryCountsRequest")
	proto.RegisterType((*QueryCountsResponse)(nil), "band.bandtss.v1beta1.QueryCountsResponse")
	proto.RegisterType((*QueryMembersRequest)(nil), "band.bandtss.v1beta1.QueryMembersRequest")
	proto.RegisterType((*QueryMembersResponse)(nil), "band.bandtss.v1beta1.QueryMembersResponse")
	proto.RegisterType((*QueryMemberRequest)(nil), "band.bandtss.v1beta1.QueryMemberRequest")
	proto.RegisterType((*QueryMemberResponse)(nil), "band.bandtss.v1beta1.QueryMemberResponse")
	proto.RegisterType((*QueryCurrentGroupRequest)(nil), "band.bandtss.v1beta1.QueryCurrentGroupRequest")
	proto.RegisterType((*QueryCurrentGroupResponse)(nil), "band.bandtss.v1beta1.QueryCurrentGroupResponse")
	proto.RegisterType((*QueryIncomingGroupRequest)(nil), "band.bandtss.v1beta1.QueryIncomingGroupRequest")
	proto.RegisterType((*QueryIncomingGroupResponse)(nil), "band.bandtss.v1beta1.QueryIncomingGroupResponse")
	proto.RegisterType((*QuerySigningRequest)(nil), "band.bandtss.v1beta1.QuerySigningRequest")
	proto.RegisterType((*QuerySigningResponse)(nil), "band.bandtss.v1beta1.QuerySigningResponse")
	proto.RegisterType((*QueryGroupTransitionRequest)(nil), "band.bandtss.v1beta1.QueryGroupTransitionRequest")
	proto.RegisterType((*QueryGroupTransitionResponse)(nil), "band.bandtss.v1beta1.QueryGroupTransitionResponse")
	proto.RegisterType((*QueryParamsRequest)(nil), "band.bandtss.v1beta1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "band.bandtss.v1beta1.QueryParamsResponse")
}

func init() { proto.RegisterFile("band/bandtss/v1beta1/query.proto", fileDescriptor_d619290a87c09054) }

var fileDescriptor_d619290a87c09054 = []byte{
	// 1268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x57, 0xb1, 0x6f, 0xdb, 0x56,
	0x13, 0x17, 0x15, 0x45, 0x8a, 0x9f, 0x9d, 0xd8, 0xdf, 0x8b, 0x3f, 0x54, 0xa6, 0x6d, 0x51, 0x66,
	0x52, 0x47, 0x31, 0x50, 0xd2, 0x56, 0xda, 0x0e, 0x59, 0xda, 0xc8, 0x91, 0x03, 0x35, 0x75, 0xa0,
	0x50, 0x4a, 0x86, 0x2c, 0x2a, 0x25, 0x3d, 0xd3, 0x84, 0x25, 0x92, 0xe6, 0x23, 0x8d, 0xba, 0x41,
	0x3a, 0x04, 0x1d, 0x9a, 0xcd, 0x40, 0x87, 0x0c, 0x1d, 0xb3, 0x75, 0xe8, 0x3f, 0xd1, 0x25, 0x40,
	0x97, 0x00, 0x5d, 0x3a, 0x39, 0x85, 0xdd, 0xb5, 0x5b, 0xa7, 0x4c, 0x05, 0xdf, 0x3b, 0xca, 0xa2,
	0xc4, 0xd0, 0xea, 0xde, 0xc5, 0x96, 0xde, 0xdd, 0xbd, 0xfb, 0xdd, 0xdd, 0xbb, 0xdf, 0x9d, 0x50,
	0xb1, 0xad, 0x5b, 0x5d, 0x35, 0xf8, 0xe3, 0x51, 0xaa, 0x1e, 0x6c, 0xb4, 0x89, 0xa7, 0x6f, 0xa8,
	0xfb, 0x3e, 0x71, 0x0f, 0x15, 0xc7, 0xb5, 0x3d, 0x1b, 0xcf, 0x07, 0x42, 0x05, 0x34, 0x14, 0xd0,
	0x10, 0xe7, 0x0d, 0xdb, 0xb0, 0x99, 0x82, 0x1a, 0x7c, 0xe2, 0xba, 0xa2, 0x64, 0xd8, 0xb6, 0xd1,
	0x23, 0x2a, 0xfb, 0xd6, 0xf6, 0x77, 0x54, 0xcf, 0xec, 0x13, 0xea, 0xe9, 0x7d, 0x07, 0x14, 0xd6,
	0x3a, 0x36, 0xed, 0xdb, 0x54, 0x6d, 0xeb, 0x94, 0x70, 0x2f, 0x03, 0x9f, 0x8e, 0x6e, 0x98, 0x96,
	0xee, 0x99, 0xb6, 0x05, 0xba, 0x85, 0x61, 0xdd, 0x50, 0xab, 0x63, 0x9b, 0xa1, 0x7c, 0x09, 0x9c,
	0xe9, 0x8e, 0xa9, 0xea, 0x96, 0x65, 0x7b, 0xcc, 0x98, 0x82, 0x54, 0x64, 0x81, 0x0d, 0x07, 0x15,
	0xc0, 0xe7, 0x32, 0x39, 0x36, 0xe8, 0x30, 0xc4, 0x24, 0x1d, 0x83, 0x58, 0x84, 0x9a, 0xa0, 0x23,
	0xcf, 0x23, 0xfc, 0x30, 0x88, 0x61, 0xd3, 0xf6, 0x2d, 0x8f, 0x6a, 0x64, 0xdf, 0x27, 0xd4, 0x93,
	0x6f, 0xa3, 0xab, 0x91, 0x53, 0xea, 0xd8, 0x16, 0x25, 0xf8, 0x1a, 0xba, 0x4c, 0x4d, 0xc3, 0x32,
	0x2d, 0xa3, 0xd5, 0x09, 0x24, 0x79, 0xa1, 0x28, 0x94, 0x32, 0xda, 0x0c, 0x1c, 0x32, 0x6d, 0xf9,
	0x57, 0x01, 0x8c, 0xb7, 0x49, 0xbf, 0x4d, 0xdc, 0xf0, 0x4e, 0xfc, 0x39, 0xca, 0x52, 0x4f, 0xf7,
	0x7c, 0xca, 0xac, 0xae, 0x94, 0x4b, 0x4a, 0x5c, 0x55, 0x14, 0x6e, 0xd5, 0x60, 0x9a, 0x5b, 0x66,
	0xcf, 0x23, 0xae, 0x06, 0x76, 0x78, 0x0d, 0xfd, 0xcf, 0xa4, 0x2d, 0xd3, 0xea, 0xd8, 0xfd, 0x00,
	0x82, 0xe1, 0xda, 0xbe, 0x93, 0x4f, 0x17, 0x85, 0xd2, 0x25, 0x6d, 0xd6, 0xa4, 0x35, 0x38, 0xbf,
	0x17, 0x1c, 0xe3, 0x2d, 0x84, 0xce, 0xaa, 0x91, 0xbf, 0x50, 0x14, 0x4a, 0xd3, 0xe5, 0x55, 0x85,
	0x97, 0x43, 0x09, 0xca, 0xa1, 0xf0, 0x07, 0x12, 0xba, 0xad, 0xeb, 0x06, 0x01, 0xa4, 0xda, 0x90,
	0xa5, 0xfc, 0x52, 0x40, 0xf3, 0xd1, 0x68, 0x20, 0x17, 0x9f, 0xa2, 0x5c, 0x9f, 0x1f, 0xe5, 0x85,
	0xe2, 0x85, 0xd2, 0x74, 0x79, 0x29, 0x29, 0x1e, 0x2d, 0x54, 0xc6, 0xf7, 0x22, 0xc0, 0xd2, 0x0c,
	0xd8, 0x8d, 0x73, 0x81, 0x71, 0xa7, 0x11, 0x64, 0x0a, 0x54, 0x0e, 0x1c, 0x40, 0x96, 0xf3, 0x28,
	0xa7, 0x77, 0xbb, 0x2e, 0xa1, 0x3c, 0xcd, 0x53, 0x5a, 0xf8, 0x55, 0xfe, 0x25, 0x5a, 0x97, 0x41,
	0x20, 0x4d, 0x34, 0xdf, 0xf1, 0x5d, 0x97, 0x58, 0x1e, 0xcf, 0x68, 0x8b, 0x23, 0x65, 0xe6, 0xe7,
	0x44, 0x55, 0xc9, 0xbc, 0x3e, 0x96, 0x52, 0x1a, 0x06, 0x7b, 0x96, 0x79, 0x2e, 0xc1, 0x8f, 0xd1,
	0xff, 0xa3, 0x85, 0x0a, 0xaf, 0x4d, 0x4f, 0x7c, 0xed, 0x55, 0x73, 0xb8, 0xa2, 0x5c, 0x24, 0x8b,
	0x28, 0xcf, 0x5f, 0xe6, 0x90, 0xcb, 0xf0, 0xd5, 0xfe, 0x9d, 0x46, 0x0b, 0x31, 0x42, 0x88, 0xf3,
	0x09, 0xba, 0xc4, 0x81, 0x98, 0x5d, 0xfe, 0x6e, 0x2b, 0x9f, 0x9d, 0x1c, 0x4b, 0x39, 0xa6, 0x54,
	0xbb, 0xfb, 0xee, 0x58, 0x5a, 0x37, 0x4c, 0x6f, 0xd7, 0x6f, 0x2b, 0x1d, 0xbb, 0xcf, 0x9a, 0x86,
	0x35, 0x48, 0xc7, 0xee, 0xa9, 0x9d, 0x5d, 0xdd, 0xb4, 0xd4, 0x83, 0x5b, 0xaa, 0xb3, 0x67, 0xb0,
	0x16, 0x04, 0x1b, 0x2d, 0xc7, 0x2e, 0xac, 0x75, 0x31, 0x46, 0x19, 0x6a, 0x7e, 0x43, 0x58, 0x70,
	0x19, 0x8d, 0x7d, 0xc6, 0x4b, 0x68, 0xca, 0xdb, 0x75, 0x09, 0xdd, 0xb5, 0x7b, 0x5d, 0xf6, 0x00,
	0x33, 0xda, 0xd9, 0x01, 0xbe, 0x8f, 0x72, 0x8e, 0xdf, 0x6e, 0xed, 0x91, 0xc3, 0x7c, 0xa6, 0x28,
	0x94, 0x66, 0x2a, 0xe5, 0x77, 0xc7, 0x92, 0x32, 0x31, 0x82, 0xba, 0x6d, 0x5a, 0x9e, 0x96, 0x75,
	0xfc, 0xf6, 0x7d, 0x72, 0x88, 0x3f, 0x19, 0xb4, 0xd6, 0x45, 0xd6, 0x5a, 0xcb, 0x3c, 0xbb, 0xc3,
	0x99, 0x65, 0x90, 0x79, 0x57, 0x0d, 0xfa, 0xa9, 0x8a, 0xa6, 0xf5, 0x8e, 0x67, 0x1e, 0x90, 0x56,
	0xc0, 0x71, 0xf9, 0x2c, 0xab, 0x8c, 0xa8, 0x70, 0x4e, 0x52, 0x42, 0x02, 0x54, 0x9a, 0x21, 0x01,
	0x56, 0x2e, 0x05, 0x75, 0x39, 0x7a, 0x2b, 0x09, 0x1a, 0xe2, 0x86, 0x81, 0x48, 0x5e, 0x84, 0xac,
	0x47, 0x1a, 0x30, 0xac, 0xc9, 0xcf, 0x69, 0x24, 0xc6, 0x49, 0xff, 0x2b, 0x4a, 0x7c, 0x51, 0xe4,
	0x8f, 0xa1, 0x4b, 0x1b, 0x9c, 0x53, 0xc3, 0xbe, 0x5e, 0x46, 0x28, 0xa4, 0xde, 0x30, 0x55, 0xda,
	0x14, 0x9c, 0xd4, 0xba, 0xf2, 0x5f, 0x69, 0xa0, 0xa9, 0x81, 0x19, 0x24, 0x78, 0x1f, 0x5d, 0xd9,
	0x21, 0xa4, 0xe5, 0x10, 0xb7, 0x15, 0x68, 0xb3, 0xbe, 0x0e, 0xd8, 0x6a, 0x21, 0x42, 0x39, 0x21,
	0xa0, 0x4d, 0xdb, 0xb4, 0x2a, 0xeb, 0x41, 0x95, 0x7f, 0x7a, 0x2b, 0x95, 0x86, 0x02, 0x87, 0x39,
	0xc6, 0xff, 0x7d, 0x44, 0xbb, 0x7b, 0xaa, 0x77, 0xe8, 0x10, 0xca, 0x0c, 0xa8, 0x36, 0xb3, 0x43,
	0x48, 0x9d, 0xb8, 0x0d, 0xe6, 0x20, 0xc8, 0xb1, 0xcb, 0x51, 0x43, 0xbb, 0x4f, 0x69, 0x67, 0x07,
	0xf8, 0x2b, 0xb4, 0x14, 0xa5, 0x9b, 0x30, 0x2c, 0x97, 0x50, 0xbf, 0xe7, 0x01, 0x55, 0x4b, 0xe3,
	0xc9, 0x3a, 0x8b, 0xcc, 0xef, 0x79, 0xda, 0xc2, 0x30, 0xe7, 0x44, 0x44, 0xb8, 0x8d, 0x96, 0x47,
	0xa8, 0x67, 0xc4, 0x45, 0x66, 0x32, 0x17, 0x62, 0x84, 0x7f, 0x22, 0x32, 0x79, 0x19, 0x2d, 0xb2,
	0x74, 0x33, 0x51, 0xd3, 0xd5, 0x2d, 0x6a, 0x06, 0xa4, 0x1c, 0xbe, 0x7a, 0x07, 0x2d, 0xc5, 0x8b,
	0xa1, 0x2a, 0x75, 0x34, 0xc7, 0x91, 0x79, 0x03, 0x19, 0xf0, 0xed, 0x87, 0xf1, 0xc4, 0x38, 0x7a,
	0xd1, 0xac, 0x11, 0x3d, 0x18, 0xcc, 0xf1, 0xba, 0xee, 0xea, 0xfd, 0xc1, 0x1c, 0x7f, 0x08, 0x8f,
	0x29, 0x3c, 0x05, 0xf7, 0xb7, 0x51, 0xd6, 0x61, 0x27, 0xc9, 0x24, 0xcf, 0xad, 0x80, 0x8d, 0xc1,
	0x62, 0xed, 0x3b, 0x01, 0xe1, 0xf1, 0x19, 0x8d, 0xaf, 0xa3, 0xe2, 0x76, 0x75, 0xbb, 0x52, 0xd5,
	0x5a, 0x8d, 0xe6, 0x9d, 0xe6, 0xa3, 0x46, 0x6b, 0xab, 0xf6, 0x65, 0xb3, 0xaa, 0xb5, 0x1e, 0x3d,
	0x68, 0xd4, 0xab, 0x9b, 0xb5, 0xad, 0x5a, 0xf5, 0xee, 0x5c, 0x0a, 0x4b, 0x68, 0x31, 0x56, 0xeb,
	0xce, 0x66, 0xb3, 0xf6, 0xb8, 0x3a, 0x27, 0xe0, 0x15, 0xb4, 0x1c, 0xab, 0x50, 0x7b, 0x00, 0x2a,
	0x69, 0x31, 0xf3, 0xfd, 0xab, 0x42, 0xaa, 0xfc, 0x62, 0x0a, 0x5d, 0x64, 0xa1, 0xe1, 0x6f, 0x51,
	0x96, 0xaf, 0x29, 0xf8, 0x3d, 0x1b, 0xc5, 0xf8, 0x7e, 0x23, 0xde, 0x9c, 0x40, 0x93, 0xe7, 0x4a,
	0x96, 0x9e, 0xff, 0xf6, 0xe7, 0x0f, 0xe9, 0x05, 0xfc, 0xc1, 0xd8, 0x22, 0xd5, 0xe1, 0x5e, 0x9f,
	0x0b, 0x28, 0x07, 0xcb, 0x01, 0x4e, 0xba, 0x37, 0xba, 0x0e, 0x89, 0x6b, 0x93, 0xa8, 0x02, 0x86,
	0x22, 0xc3, 0x20, 0xe2, 0xfc, 0x18, 0x86, 0x70, 0xab, 0x78, 0x21, 0xa0, 0x2c, 0x4c, 0xde, 0xd2,
	0xb9, 0x17, 0x4f, 0x92, 0x85, 0xe8, 0x92, 0x20, 0xaf, 0x31, 0x04, 0xd7, 0xb1, 0xfc, 0x3e, 0x04,
	0xea, 0x53, 0xd8, 0x33, 0x9e, 0xe1, 0x97, 0x02, 0x9a, 0x19, 0x9e, 0xc0, 0x58, 0x49, 0xca, 0xf6,
	0xf8, 0x1c, 0x17, 0xd5, 0x89, 0xf5, 0x01, 0xdd, 0x2a, 0x43, 0x57, 0xc4, 0x85, 0xf1, 0x1a, 0x0d,
	0x53, 0x0d, 0xfe, 0x51, 0x40, 0x97, 0xa3, 0x6b, 0x62, 0x92, 0xab, 0xb8, 0x79, 0x26, 0xae, 0x4f,
	0x6e, 0x00, 0xe0, 0x6e, 0x30, 0x70, 0x2b, 0x58, 0x1a, 0x03, 0x17, 0x65, 0x29, 0x7c, 0x24, 0xa0,
	0x1c, 0xb0, 0x4c, 0xe2, 0x43, 0x8a, 0x4e, 0x86, 0xc4, 0x87, 0x34, 0x32, 0x0d, 0x64, 0x85, 0x61,
	0x29, 0xe1, 0xd5, 0x31, 0x2c, 0x40, 0x91, 0x54, 0x7d, 0x7a, 0x36, 0x66, 0x9e, 0xe1, 0x57, 0x02,
	0x9a, 0x1d, 0xa1, 0x1e, 0xbc, 0x91, 0xe0, 0x2f, 0x9e, 0x0e, 0xc5, 0xf2, 0xbf, 0x31, 0x01, 0xa8,
	0x37, 0x19, 0xd4, 0x6b, 0x78, 0x65, 0xfc, 0x07, 0xcc, 0x08, 0x73, 0x06, 0x0c, 0xc0, 0xa9, 0x2a,
	0xf1, 0xed, 0x47, 0x98, 0x31, 0xf1, 0xed, 0x47, 0xd9, 0x32, 0x81, 0x01, 0x38, 0x25, 0x56, 0xbe,
	0x78, 0x7d, 0x52, 0x10, 0xde, 0x9c, 0x14, 0x84, 0x3f, 0x4e, 0x0a, 0xc2, 0xd1, 0x69, 0x21, 0xf5,
	0xe6, 0xb4, 0x90, 0xfa, 0xfd, 0xb4, 0x90, 0x7a, 0x72, 0xfe, 0xf6, 0xf2, 0xf5, 0xe0, 0x52, 0x36,
	0x50, 0xdb, 0x59, 0xa6, 0x72, 0xeb, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x08, 0x2c, 0x11, 0xcf,
	0xd5, 0x0e, 0x00, 0x00,
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
	// Counts queries the number of existing signing.
	Counts(ctx context.Context, in *QueryCountsRequest, opts ...grpc.CallOption) (*QueryCountsResponse, error)
	// Members queries all members.
	Members(ctx context.Context, in *QueryMembersRequest, opts ...grpc.CallOption) (*QueryMembersResponse, error)
	// Member queries the member information of the given address.
	Member(ctx context.Context, in *QueryMemberRequest, opts ...grpc.CallOption) (*QueryMemberResponse, error)
	// CurrentGroup queries the current group information.
	CurrentGroup(
		ctx context.Context,
		in *QueryCurrentGroupRequest,
		opts ...grpc.CallOption,
	) (*QueryCurrentGroupResponse, error)
	// IncomingGroup queries the incoming group information.
	IncomingGroup(
		ctx context.Context,
		in *QueryIncomingGroupRequest,
		opts ...grpc.CallOption,
	) (*QueryIncomingGroupResponse, error)
	// Signing queries the signing result of the given signing request ID.
	Signing(ctx context.Context, in *QuerySigningRequest, opts ...grpc.CallOption) (*QuerySigningResponse, error)
	// GroupTransition queries the group transition information.
	GroupTransition(
		ctx context.Context,
		in *QueryGroupTransitionRequest,
		opts ...grpc.CallOption,
	) (*QueryGroupTransitionResponse, error)
	// Params queries parameters of bandtss module
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Counts(
	ctx context.Context,
	in *QueryCountsRequest,
	opts ...grpc.CallOption,
) (*QueryCountsResponse, error) {
	out := new(QueryCountsResponse)
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/Counts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Members(
	ctx context.Context,
	in *QueryMembersRequest,
	opts ...grpc.CallOption,
) (*QueryMembersResponse, error) {
	out := new(QueryMembersResponse)
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/Members", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Member(
	ctx context.Context,
	in *QueryMemberRequest,
	opts ...grpc.CallOption,
) (*QueryMemberResponse, error) {
	out := new(QueryMemberResponse)
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/Member", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CurrentGroup(
	ctx context.Context,
	in *QueryCurrentGroupRequest,
	opts ...grpc.CallOption,
) (*QueryCurrentGroupResponse, error) {
	out := new(QueryCurrentGroupResponse)
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/CurrentGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) IncomingGroup(
	ctx context.Context,
	in *QueryIncomingGroupRequest,
	opts ...grpc.CallOption,
) (*QueryIncomingGroupResponse, error) {
	out := new(QueryIncomingGroupResponse)
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/IncomingGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Signing(
	ctx context.Context,
	in *QuerySigningRequest,
	opts ...grpc.CallOption,
) (*QuerySigningResponse, error) {
	out := new(QuerySigningResponse)
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/Signing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GroupTransition(
	ctx context.Context,
	in *QueryGroupTransitionRequest,
	opts ...grpc.CallOption,
) (*QueryGroupTransitionResponse, error) {
	out := new(QueryGroupTransitionResponse)
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/GroupTransition", in, out, opts...)
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
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Counts queries the number of existing signing.
	Counts(context.Context, *QueryCountsRequest) (*QueryCountsResponse, error)
	// Members queries all members.
	Members(context.Context, *QueryMembersRequest) (*QueryMembersResponse, error)
	// Member queries the member information of the given address.
	Member(context.Context, *QueryMemberRequest) (*QueryMemberResponse, error)
	// CurrentGroup queries the current group information.
	CurrentGroup(context.Context, *QueryCurrentGroupRequest) (*QueryCurrentGroupResponse, error)
	// IncomingGroup queries the incoming group information.
	IncomingGroup(context.Context, *QueryIncomingGroupRequest) (*QueryIncomingGroupResponse, error)
	// Signing queries the signing result of the given signing request ID.
	Signing(context.Context, *QuerySigningRequest) (*QuerySigningResponse, error)
	// GroupTransition queries the group transition information.
	GroupTransition(context.Context, *QueryGroupTransitionRequest) (*QueryGroupTransitionResponse, error)
	// Params queries parameters of bandtss module
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Counts(ctx context.Context, req *QueryCountsRequest) (*QueryCountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Counts not implemented")
}

func (*UnimplementedQueryServer) Members(
	ctx context.Context,
	req *QueryMembersRequest,
) (*QueryMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Members not implemented")
}
func (*UnimplementedQueryServer) Member(ctx context.Context, req *QueryMemberRequest) (*QueryMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Member not implemented")
}

func (*UnimplementedQueryServer) CurrentGroup(
	ctx context.Context,
	req *QueryCurrentGroupRequest,
) (*QueryCurrentGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CurrentGroup not implemented")
}

func (*UnimplementedQueryServer) IncomingGroup(
	ctx context.Context,
	req *QueryIncomingGroupRequest,
) (*QueryIncomingGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IncomingGroup not implemented")
}

func (*UnimplementedQueryServer) Signing(
	ctx context.Context,
	req *QuerySigningRequest,
) (*QuerySigningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Signing not implemented")
}

func (*UnimplementedQueryServer) GroupTransition(
	ctx context.Context,
	req *QueryGroupTransitionRequest,
) (*QueryGroupTransitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GroupTransition not implemented")
}
func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Counts_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryCountsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Counts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.bandtss.v1beta1.Query/Counts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Counts(ctx, req.(*QueryCountsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Members_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Members(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.bandtss.v1beta1.Query/Members",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Members(ctx, req.(*QueryMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Member_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Member(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.bandtss.v1beta1.Query/Member",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Member(ctx, req.(*QueryMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CurrentGroup_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryCurrentGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CurrentGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.bandtss.v1beta1.Query/CurrentGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CurrentGroup(ctx, req.(*QueryCurrentGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_IncomingGroup_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryIncomingGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).IncomingGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.bandtss.v1beta1.Query/IncomingGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).IncomingGroup(ctx, req.(*QueryIncomingGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Signing_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QuerySigningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Signing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.bandtss.v1beta1.Query/Signing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Signing(ctx, req.(*QuerySigningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GroupTransition_Handler(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor,
) (interface{}, error) {
	in := new(QueryGroupTransitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GroupTransition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/band.bandtss.v1beta1.Query/GroupTransition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GroupTransition(ctx, req.(*QueryGroupTransitionRequest))
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
		FullMethod: "/band.bandtss.v1beta1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "band.bandtss.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Counts",
			Handler:    _Query_Counts_Handler,
		},
		{
			MethodName: "Members",
			Handler:    _Query_Members_Handler,
		},
		{
			MethodName: "Member",
			Handler:    _Query_Member_Handler,
		},
		{
			MethodName: "CurrentGroup",
			Handler:    _Query_CurrentGroup_Handler,
		},
		{
			MethodName: "IncomingGroup",
			Handler:    _Query_IncomingGroup_Handler,
		},
		{
			MethodName: "Signing",
			Handler:    _Query_Signing_Handler,
		},
		{
			MethodName: "GroupTransition",
			Handler:    _Query_GroupTransition_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "band/bandtss/v1beta1/query.proto",
}

func (m *QueryCountsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCountsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCountsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryCountsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCountsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCountsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SigningCount != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.SigningCount))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryMembersRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMembersRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMembersRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		dAtA[i] = 0x1a
	}
	if m.IsIncomingGroup {
		i--
		if m.IsIncomingGroup {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if m.Status != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryMembersResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMembersResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMembersResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Members) > 0 {
		for iNdEx := len(m.Members) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Members[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryMemberRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMemberRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMemberRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryMemberResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMemberResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMemberResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.IncomingGroupMember.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.CurrentGroupMember.MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryCurrentGroupRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCurrentGroupRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCurrentGroupRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryCurrentGroupResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryCurrentGroupResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryCurrentGroupResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n5, err5 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(
		m.ActiveTime,
		dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ActiveTime):],
	)
	if err5 != nil {
		return 0, err5
	}
	i -= n5
	i = encodeVarintQuery(dAtA, i, uint64(n5))
	i--
	dAtA[i] = 0x32
	if m.Status != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x28
	}
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0x22
	}
	if m.Threshold != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x18
	}
	if m.Size_ != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Size_))
		i--
		dAtA[i] = 0x10
	}
	if m.GroupID != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.GroupID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryIncomingGroupRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryIncomingGroupRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryIncomingGroupRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryIncomingGroupResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryIncomingGroupResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryIncomingGroupResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Status != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x28
	}
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0x22
	}
	if m.Threshold != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x18
	}
	if m.Size_ != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Size_))
		i--
		dAtA[i] = 0x10
	}
	if m.GroupID != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.GroupID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QuerySigningRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySigningRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySigningRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SigningId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.SigningId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QuerySigningResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuerySigningResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuerySigningResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IncomingGroupSigningResult != nil {
		{
			size, err := m.IncomingGroupSigningResult.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.CurrentGroupSigningResult != nil {
		{
			size, err := m.CurrentGroupSigningResult.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Requester) > 0 {
		i -= len(m.Requester)
		copy(dAtA[i:], m.Requester)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Requester)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.FeePerSigner) > 0 {
		for iNdEx := len(m.FeePerSigner) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.FeePerSigner[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryGroupTransitionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGroupTransitionRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGroupTransitionRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryGroupTransitionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGroupTransitionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGroupTransitionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.GroupTransition != nil {
		{
			size, err := m.GroupTransition.MarshalToSizedBuffer(dAtA[:i])
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
func (m *QueryCountsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryCountsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigningCount != 0 {
		n += 1 + sovQuery(uint64(m.SigningCount))
	}
	return n
}

func (m *QueryMembersRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != 0 {
		n += 1 + sovQuery(uint64(m.Status))
	}
	if m.IsIncomingGroup {
		n += 2
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryMembersResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Members) > 0 {
		for _, e := range m.Members {
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

func (m *QueryMemberRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryMemberResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.CurrentGroupMember.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = m.IncomingGroupMember.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryCurrentGroupRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryCurrentGroupResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GroupID != 0 {
		n += 1 + sovQuery(uint64(m.GroupID))
	}
	if m.Size_ != 0 {
		n += 1 + sovQuery(uint64(m.Size_))
	}
	if m.Threshold != 0 {
		n += 1 + sovQuery(uint64(m.Threshold))
	}
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovQuery(uint64(m.Status))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ActiveTime)
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryIncomingGroupRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryIncomingGroupResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GroupID != 0 {
		n += 1 + sovQuery(uint64(m.GroupID))
	}
	if m.Size_ != 0 {
		n += 1 + sovQuery(uint64(m.Size_))
	}
	if m.Threshold != 0 {
		n += 1 + sovQuery(uint64(m.Threshold))
	}
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovQuery(uint64(m.Status))
	}
	return n
}

func (m *QuerySigningRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigningId != 0 {
		n += 1 + sovQuery(uint64(m.SigningId))
	}
	return n
}

func (m *QuerySigningResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.FeePerSigner) > 0 {
		for _, e := range m.FeePerSigner {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	l = len(m.Requester)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.CurrentGroupSigningResult != nil {
		l = m.CurrentGroupSigningResult.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.IncomingGroupSigningResult != nil {
		l = m.IncomingGroupSigningResult.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryGroupTransitionRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryGroupTransitionResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GroupTransition != nil {
		l = m.GroupTransition.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
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
func (m *QueryCountsRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryCountsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCountsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryCountsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryCountsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCountsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningCount", wireType)
			}
			m.SigningCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningCount |= uint64(b&0x7F) << shift
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
func (m *QueryMembersRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMembersRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMembersRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= MemberStatusFilter(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsIncomingGroup", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
			m.IsIncomingGroup = bool(v != 0)
		case 3:
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
func (m *QueryMembersResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMembersResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMembersResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Members", wireType)
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
			m.Members = append(m.Members, &Member{})
			if err := m.Members[len(m.Members)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryMemberRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMemberRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMemberRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
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
func (m *QueryMemberResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMemberResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMemberResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentGroupMember", wireType)
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
			if err := m.CurrentGroupMember.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncomingGroupMember", wireType)
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
			if err := m.IncomingGroupMember.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryCurrentGroupRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryCurrentGroupRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCurrentGroupRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryCurrentGroupResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryCurrentGroupResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryCurrentGroupResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupID", wireType)
			}
			m.GroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GroupID |= github_com_bandprotocol_chain_v3_pkg_tss.GroupID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Size_", wireType)
			}
			m.Size_ = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Size_ |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Threshold", wireType)
			}
			m.Threshold = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Threshold |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKey = append(m.PubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKey == nil {
				m.PubKey = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= types.GroupStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActiveTime", wireType)
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
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.ActiveTime, dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryIncomingGroupRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryIncomingGroupRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryIncomingGroupRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryIncomingGroupResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryIncomingGroupResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryIncomingGroupResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupID", wireType)
			}
			m.GroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GroupID |= github_com_bandprotocol_chain_v3_pkg_tss.GroupID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Size_", wireType)
			}
			m.Size_ = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Size_ |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Threshold", wireType)
			}
			m.Threshold = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Threshold |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKey = append(m.PubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKey == nil {
				m.PubKey = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= types.GroupStatus(b&0x7F) << shift
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
func (m *QuerySigningRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QuerySigningRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySigningRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningId", wireType)
			}
			m.SigningId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningId |= uint64(b&0x7F) << shift
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
func (m *QuerySigningResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QuerySigningResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuerySigningResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeePerSigner", wireType)
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
			m.FeePerSigner = append(m.FeePerSigner, types1.Coin{})
			if err := m.FeePerSigner[len(m.FeePerSigner)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requester", wireType)
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
			m.Requester = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentGroupSigningResult", wireType)
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
			if m.CurrentGroupSigningResult == nil {
				m.CurrentGroupSigningResult = &types.SigningResult{}
			}
			if err := m.CurrentGroupSigningResult.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncomingGroupSigningResult", wireType)
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
			if m.IncomingGroupSigningResult == nil {
				m.IncomingGroupSigningResult = &types.SigningResult{}
			}
			if err := m.IncomingGroupSigningResult.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryGroupTransitionRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGroupTransitionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGroupTransitionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryGroupTransitionResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGroupTransitionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGroupTransitionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupTransition", wireType)
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
			if m.GroupTransition == nil {
				m.GroupTransition = &GroupTransition{}
			}
			if err := m.GroupTransition.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
