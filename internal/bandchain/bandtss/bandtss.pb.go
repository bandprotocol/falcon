package types

import (
	"bytes"
	"fmt"
	"io"
	"math"
	math_bits "math/bits"
	"time"

	github_com_cometbft_cometbft_libs_bytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"

	github_com_bandprotocol_chain_v3_pkg_tss "github.com/bandprotocol/falcon/internal/bandchain/tss"
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

// TransitionStatus is an enumeration of the possible statuses of a group transition process.
type TransitionStatus int32

const (
	// TRANSITION_STATUS_UNSPECIFIED is the status of a group transition that has not been specified.
	TRANSITION_STATUS_UNSPECIFIED TransitionStatus = 0
	// TRANSITION_STATUS_CREATING_GROUP is the status of a group transition that a new group
	// is being created.
	TRANSITION_STATUS_CREATING_GROUP TransitionStatus = 1
	// TRANSITION_STATUS_WAITING_SIGN is the status of a group transition that waits members in
	// a current group to sign the transition message.
	TRANSITION_STATUS_WAITING_SIGN TransitionStatus = 2
	// TRANSITION_STATUS_WAITING_EXECUTION is the status of a group transition that
	// a transition process is completed, either from a forceTransition or having a current-group
	// signature on a transition message, but waits for the execution time.
	TRANSITION_STATUS_WAITING_EXECUTION TransitionStatus = 3
)

var TransitionStatus_name = map[int32]string{
	0: "TRANSITION_STATUS_UNSPECIFIED",
	1: "TRANSITION_STATUS_CREATING_GROUP",
	2: "TRANSITION_STATUS_WAITING_SIGN",
	3: "TRANSITION_STATUS_WAITING_EXECUTION",
}

var TransitionStatus_value = map[string]int32{
	"TRANSITION_STATUS_UNSPECIFIED":       0,
	"TRANSITION_STATUS_CREATING_GROUP":    1,
	"TRANSITION_STATUS_WAITING_SIGN":      2,
	"TRANSITION_STATUS_WAITING_EXECUTION": 3,
}

func (x TransitionStatus) String() string {
	return proto.EnumName(TransitionStatus_name, int32(x))
}

func (TransitionStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2bc325518cc10c44, []int{0}
}

// Member maintains member information for monitoring their liveness activity.
type Member struct {
	// address is the address of the member.
	Address string `protobuf:"bytes,1,opt,name=address,proto3"                                                                                       json:"address,omitempty"`
	// group_id is the group ID that the member belongs to.
	GroupID github_com_bandprotocol_chain_v3_pkg_tss.GroupID `protobuf:"varint,2,opt,name=group_id,json=groupId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID" json:"group_id,omitempty"`
	// is_active is a flag to indicate whether a member is active or not.
	IsActive bool `protobuf:"varint,3,opt,name=is_active,json=isActive,proto3"                                                                      json:"is_active,omitempty"`
	// since is a block timestamp when a member status is changed (from active to inactive or vice versa).
	Since time.Time `protobuf:"bytes,4,opt,name=since,proto3,stdtime"                                                                                 json:"since"`
}

func (m *Member) Reset()         { *m = Member{} }
func (m *Member) String() string { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()    {}
func (*Member) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc325518cc10c44, []int{0}
}
func (m *Member) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Member) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Member.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Member) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Member.Merge(m, src)
}
func (m *Member) XXX_Size() int {
	return m.Size()
}
func (m *Member) XXX_DiscardUnknown() {
	xxx_messageInfo_Member.DiscardUnknown(m)
}

var xxx_messageInfo_Member proto.InternalMessageInfo

func (m *Member) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Member) GetGroupID() github_com_bandprotocol_chain_v3_pkg_tss.GroupID {
	if m != nil {
		return m.GroupID
	}
	return 0
}

func (m *Member) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

func (m *Member) GetSince() time.Time {
	if m != nil {
		return m.Since
	}
	return time.Time{}
}

// CuurentGroup is a bandtss current group information.
type CurrentGroup struct {
	// group_id is the ID of the current group.
	GroupID github_com_bandprotocol_chain_v3_pkg_tss.GroupID `protobuf:"varint,1,opt,name=group_id,json=groupId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID" json:"group_id,omitempty"`
	// active_time is the timestamp at which the group becomes the current group of the module.
	ActiveTime time.Time `protobuf:"bytes,2,opt,name=active_time,json=activeTime,proto3,stdtime"                                                           json:"active_time"`
}

func (m *CurrentGroup) Reset()         { *m = CurrentGroup{} }
func (m *CurrentGroup) String() string { return proto.CompactTextString(m) }
func (*CurrentGroup) ProtoMessage()    {}
func (*CurrentGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc325518cc10c44, []int{1}
}
func (m *CurrentGroup) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CurrentGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CurrentGroup.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CurrentGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CurrentGroup.Merge(m, src)
}
func (m *CurrentGroup) XXX_Size() int {
	return m.Size()
}
func (m *CurrentGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_CurrentGroup.DiscardUnknown(m)
}

var xxx_messageInfo_CurrentGroup proto.InternalMessageInfo

func (m *CurrentGroup) GetGroupID() github_com_bandprotocol_chain_v3_pkg_tss.GroupID {
	if m != nil {
		return m.GroupID
	}
	return 0
}

func (m *CurrentGroup) GetActiveTime() time.Time {
	if m != nil {
		return m.ActiveTime
	}
	return time.Time{}
}

// Signing is a bandtss signing information.
type Signing struct {
	// id is the unique identifier of the bandtss signing.
	ID SigningID `protobuf:"varint,1,opt,name=id,proto3,casttype=SigningID"                                                                                                          json:"id,omitempty"`
	// fee_per_signer is the tokens that will be paid per signer for this bandtss signing.
	FeePerSigner github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=fee_per_signer,json=feePerSigner,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins"                                          json:"fee_per_signer"`
	// requester is the address who pays the Bandtss signing.
	Requester string `protobuf:"bytes,3,opt,name=requester,proto3"                                                                                                                       json:"requester,omitempty"`
	// current_group_signing_id is a tss signing ID of a current group.
	CurrentGroupSigningID github_com_bandprotocol_chain_v3_pkg_tss.SigningID `protobuf:"varint,4,opt,name=current_group_signing_id,json=currentGroupSigningId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID"   json:"current_group_signing_id,omitempty"`
	// incoming_group_signing_id is a tss signing ID of an incoming group, if any.
	IncomingGroupSigningID github_com_bandprotocol_chain_v3_pkg_tss.SigningID `protobuf:"varint,5,opt,name=incoming_group_signing_id,json=incomingGroupSigningId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID" json:"incoming_group_signing_id,omitempty"`
}

func (m *Signing) Reset()         { *m = Signing{} }
func (m *Signing) String() string { return proto.CompactTextString(m) }
func (*Signing) ProtoMessage()    {}
func (*Signing) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc325518cc10c44, []int{2}
}
func (m *Signing) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Signing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Signing.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Signing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Signing.Merge(m, src)
}
func (m *Signing) XXX_Size() int {
	return m.Size()
}
func (m *Signing) XXX_DiscardUnknown() {
	xxx_messageInfo_Signing.DiscardUnknown(m)
}

var xxx_messageInfo_Signing proto.InternalMessageInfo

func (m *Signing) GetID() SigningID {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Signing) GetFeePerSigner() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.FeePerSigner
	}
	return nil
}

func (m *Signing) GetRequester() string {
	if m != nil {
		return m.Requester
	}
	return ""
}

func (m *Signing) GetCurrentGroupSigningID() github_com_bandprotocol_chain_v3_pkg_tss.SigningID {
	if m != nil {
		return m.CurrentGroupSigningID
	}
	return 0
}

func (m *Signing) GetIncomingGroupSigningID() github_com_bandprotocol_chain_v3_pkg_tss.SigningID {
	if m != nil {
		return m.IncomingGroupSigningID
	}
	return 0
}

// GroupTransition defines the group transition information of the current group and incoming group.
type GroupTransition struct {
	// signing_id is a tss signing ID of group transition signing request.
	SigningID github_com_bandprotocol_chain_v3_pkg_tss.SigningID `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID"                  json:"signing_id,omitempty"`
	// current_group_id is the ID of the group that will be replaced.
	CurrentGroupID github_com_bandprotocol_chain_v3_pkg_tss.GroupID `protobuf:"varint,2,opt,name=current_group_id,json=currentGroupId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID"         json:"current_group_id,omitempty"`
	// current_group_pub_key is the public key pair that used for sign & verify transition group msg.
	CurrentGroupPubKey github_com_bandprotocol_chain_v3_pkg_tss.Point `protobuf:"bytes,3,opt,name=current_group_pub_key,json=currentGroupPubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"   json:"current_group_pub_key,omitempty"`
	// new_group_id is the ID of the new group that be a new key candidate.
	IncomingGroupID github_com_bandprotocol_chain_v3_pkg_tss.GroupID `protobuf:"varint,4,opt,name=incoming_group_id,json=incomingGroupId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID"       json:"incoming_group_id,omitempty"`
	// incoming_group_pub_key is the public key of the group that will be the next key of this group
	IncomingGroupPubKey github_com_bandprotocol_chain_v3_pkg_tss.Point `protobuf:"bytes,5,opt,name=incoming_group_pub_key,json=incomingGroupPubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point" json:"incoming_group_pub_key,omitempty"`
	// status is an enumeration of the possible statuses of a group transition process.
	Status TransitionStatus `protobuf:"varint,6,opt,name=status,proto3,enum=band.bandtss.v1beta1.TransitionStatus"                                                                   json:"status,omitempty"`
	// exec_time is the time when the transition will be executed.
	ExecTime time.Time `protobuf:"bytes,7,opt,name=exec_time,json=execTime,proto3,stdtime"                                                                                      json:"exec_time"`
	// is_force_transition is a flag to indicate whether the current group signs the transition message
	// before the transition is executed or not.
	IsForceTransition bool `protobuf:"varint,8,opt,name=is_force_transition,json=isForceTransition,proto3"                                                                          json:"is_force_transition,omitempty"`
}

func (m *GroupTransition) Reset()         { *m = GroupTransition{} }
func (m *GroupTransition) String() string { return proto.CompactTextString(m) }
func (*GroupTransition) ProtoMessage()    {}
func (*GroupTransition) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc325518cc10c44, []int{3}
}
func (m *GroupTransition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GroupTransition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GroupTransition.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GroupTransition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupTransition.Merge(m, src)
}
func (m *GroupTransition) XXX_Size() int {
	return m.Size()
}
func (m *GroupTransition) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupTransition.DiscardUnknown(m)
}

var xxx_messageInfo_GroupTransition proto.InternalMessageInfo

func (m *GroupTransition) GetSigningID() github_com_bandprotocol_chain_v3_pkg_tss.SigningID {
	if m != nil {
		return m.SigningID
	}
	return 0
}

func (m *GroupTransition) GetCurrentGroupID() github_com_bandprotocol_chain_v3_pkg_tss.GroupID {
	if m != nil {
		return m.CurrentGroupID
	}
	return 0
}

func (m *GroupTransition) GetCurrentGroupPubKey() github_com_bandprotocol_chain_v3_pkg_tss.Point {
	if m != nil {
		return m.CurrentGroupPubKey
	}
	return nil
}

func (m *GroupTransition) GetIncomingGroupID() github_com_bandprotocol_chain_v3_pkg_tss.GroupID {
	if m != nil {
		return m.IncomingGroupID
	}
	return 0
}

func (m *GroupTransition) GetIncomingGroupPubKey() github_com_bandprotocol_chain_v3_pkg_tss.Point {
	if m != nil {
		return m.IncomingGroupPubKey
	}
	return nil
}

func (m *GroupTransition) GetStatus() TransitionStatus {
	if m != nil {
		return m.Status
	}
	return TRANSITION_STATUS_UNSPECIFIED
}

func (m *GroupTransition) GetExecTime() time.Time {
	if m != nil {
		return m.ExecTime
	}
	return time.Time{}
}

func (m *GroupTransition) GetIsForceTransition() bool {
	if m != nil {
		return m.IsForceTransition
	}
	return false
}

// GroupTransitionSignatureOrder defines a general signature order for group transition.
type GroupTransitionSignatureOrder struct {
	// pub_key is the public key of new group that the current group needs to be signed.
	PubKey github_com_cometbft_cometbft_libs_bytes.HexBytes `protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3,casttype=github.com/cometbft/cometbft/libs/bytes.HexBytes" json:"pub_key,omitempty"`
	// transition_time is the timestamp at which the transition is executed and the public key is active.
	TransitionTime time.Time `protobuf:"bytes,2,opt,name=transition_time,json=transitionTime,proto3,stdtime"                                   json:"transition_time"`
}

func (m *GroupTransitionSignatureOrder) Reset()         { *m = GroupTransitionSignatureOrder{} }
func (m *GroupTransitionSignatureOrder) String() string { return proto.CompactTextString(m) }
func (*GroupTransitionSignatureOrder) ProtoMessage()    {}
func (*GroupTransitionSignatureOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bc325518cc10c44, []int{4}
}
func (m *GroupTransitionSignatureOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GroupTransitionSignatureOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GroupTransitionSignatureOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GroupTransitionSignatureOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupTransitionSignatureOrder.Merge(m, src)
}
func (m *GroupTransitionSignatureOrder) XXX_Size() int {
	return m.Size()
}
func (m *GroupTransitionSignatureOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupTransitionSignatureOrder.DiscardUnknown(m)
}

var xxx_messageInfo_GroupTransitionSignatureOrder proto.InternalMessageInfo

func (m *GroupTransitionSignatureOrder) GetPubKey() github_com_cometbft_cometbft_libs_bytes.HexBytes {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *GroupTransitionSignatureOrder) GetTransitionTime() time.Time {
	if m != nil {
		return m.TransitionTime
	}
	return time.Time{}
}

func init() {
	proto.RegisterEnum("band.bandtss.v1beta1.TransitionStatus", TransitionStatus_name, TransitionStatus_value)
	proto.RegisterType((*Member)(nil), "band.bandtss.v1beta1.Member")
	proto.RegisterType((*CurrentGroup)(nil), "band.bandtss.v1beta1.CurrentGroup")
	proto.RegisterType((*Signing)(nil), "band.bandtss.v1beta1.Signing")
	proto.RegisterType((*GroupTransition)(nil), "band.bandtss.v1beta1.GroupTransition")
	proto.RegisterType((*GroupTransitionSignatureOrder)(nil), "band.bandtss.v1beta1.GroupTransitionSignatureOrder")
}

func init() {
	proto.RegisterFile("band/bandtss/v1beta1/bandtss.proto", fileDescriptor_2bc325518cc10c44)
}

var fileDescriptor_2bc325518cc10c44 = []byte{
	// 957 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0x4d, 0x6f, 0xdb, 0x46,
	0x13, 0xd6, 0xca, 0x1f, 0x92, 0xd6, 0x86, 0xac, 0xac, 0x63, 0x43, 0xf6, 0xfb, 0x86, 0x54, 0x95,
	0xa2, 0x15, 0x0a, 0x84, 0x4c, 0x94, 0x9e, 0x7c, 0x68, 0x21, 0xc9, 0xb2, 0x4a, 0x14, 0x96, 0x05,
	0x4a, 0x46, 0x8b, 0x5c, 0x08, 0x7e, 0xac, 0x98, 0x85, 0x2d, 0xae, 0xc2, 0x5d, 0x19, 0xf6, 0x3f,
	0x08, 0xd0, 0x4b, 0x7e, 0x42, 0x81, 0x5e, 0x8a, 0x1e, 0x0a, 0x14, 0xe8, 0x8f, 0x08, 0x72, 0xca,
	0xb1, 0x97, 0xd2, 0x05, 0x7d, 0xe9, 0x6f, 0xf0, 0xa9, 0xe0, 0x92, 0xfa, 0xb4, 0x8b, 0x3a, 0x46,
	0x2e, 0xd2, 0xce, 0xce, 0xcc, 0xce, 0x3c, 0xcf, 0x33, 0x5c, 0x12, 0x96, 0x2d, 0xd3, 0x73, 0xd4,
	0xe8, 0x87, 0x33, 0xa6, 0x9e, 0x3d, 0xb3, 0x30, 0x37, 0x9f, 0x8d, 0x6d, 0x65, 0xe8, 0x53, 0x4e,
	0xd1, 0xc3, 0xc8, 0x54, 0xc6, 0x7b, 0x49, 0xcc, 0xee, 0x43, 0x97, 0xba, 0x54, 0x04, 0xa8, 0xd1,
	0x2a, 0x8e, 0xdd, 0x95, 0x5d, 0x4a, 0xdd, 0x53, 0xac, 0x0a, 0xcb, 0x1a, 0xf5, 0x55, 0x4e, 0x06,
	0x98, 0x71, 0x73, 0x30, 0x4c, 0x02, 0x24, 0x9b, 0xb2, 0x01, 0x65, 0xaa, 0x65, 0x32, 0x3c, 0xa9,
	0x67, 0x53, 0xe2, 0x25, 0xfe, 0x9d, 0xd8, 0x6f, 0xc4, 0x27, 0xc7, 0x46, 0xec, 0x2a, 0xff, 0x09,
	0xe0, 0xea, 0x21, 0x1e, 0x58, 0xd8, 0x47, 0x45, 0x98, 0x31, 0x1d, 0xc7, 0xc7, 0x8c, 0x15, 0x41,
	0x09, 0x54, 0x72, 0xfa, 0xd8, 0x44, 0x2f, 0x60, 0xd6, 0xf5, 0xe9, 0x68, 0x68, 0x10, 0xa7, 0x98,
	0x2e, 0x81, 0xca, 0x72, 0xfd, 0xeb, 0x30, 0x90, 0x33, 0xad, 0x68, 0x4f, 0xdb, 0xbf, 0x0e, 0xe4,
	0xa7, 0x2e, 0xe1, 0x2f, 0x47, 0x96, 0x62, 0xd3, 0x81, 0xc0, 0x29, 0xce, 0xb6, 0xe9, 0xa9, 0x6a,
	0xbf, 0x34, 0x89, 0xa7, 0x9e, 0x3d, 0x57, 0x87, 0x27, 0xae, 0x1a, 0x21, 0x4d, 0x72, 0xf4, 0x8c,
	0x38, 0x50, 0x73, 0xd0, 0xff, 0x60, 0x8e, 0x30, 0xc3, 0xb4, 0x39, 0x39, 0xc3, 0xc5, 0xa5, 0x12,
	0xa8, 0x64, 0xf5, 0x2c, 0x61, 0x35, 0x61, 0xa3, 0x3d, 0xb8, 0xc2, 0x88, 0x67, 0xe3, 0xe2, 0x72,
	0x09, 0x54, 0xd6, 0xaa, 0xbb, 0x4a, 0xcc, 0x84, 0x32, 0x66, 0x42, 0xe9, 0x8d, 0x99, 0xa8, 0x67,
	0xdf, 0x06, 0x72, 0xea, 0xcd, 0xa5, 0x0c, 0xf4, 0x38, 0x65, 0x6f, 0xf9, 0xef, 0x1f, 0x65, 0x50,
	0xfe, 0x0d, 0xc0, 0xf5, 0xc6, 0xc8, 0xf7, 0xb1, 0xc7, 0x45, 0xe9, 0x39, 0x2c, 0xe0, 0x23, 0x63,
	0x69, 0xc2, 0xb5, 0x18, 0x88, 0x11, 0x29, 0x24, 0xa8, 0xba, 0x6b, 0xd3, 0x30, 0x4e, 0x8c, 0x5c,
	0xe5, 0x60, 0x09, 0x66, 0xba, 0xc4, 0xf5, 0x88, 0xe7, 0xa2, 0xc7, 0x30, 0x3d, 0x69, 0x74, 0x33,
	0x0c, 0xe4, 0xb4, 0xe8, 0x31, 0x97, 0xb8, 0xb5, 0x7d, 0x3d, 0x4d, 0x1c, 0xf4, 0x0a, 0xe6, 0xfb,
	0x18, 0x1b, 0x43, 0xec, 0x1b, 0x8c, 0xb8, 0x1e, 0xf6, 0x8b, 0xe9, 0xd2, 0x52, 0x65, 0xad, 0xba,
	0xa3, 0x24, 0x5a, 0x47, 0x83, 0x31, 0x1e, 0x32, 0xa5, 0x41, 0x89, 0x57, 0x7f, 0x1a, 0x55, 0xfe,
	0xe5, 0x52, 0xae, 0xcc, 0xa0, 0x4d, 0xa6, 0x28, 0xfe, 0x7b, 0xc2, 0x9c, 0x13, 0x95, 0x5f, 0x0c,
	0x31, 0x13, 0x09, 0x4c, 0x5f, 0xef, 0x63, 0xdc, 0xc1, 0x7e, 0x57, 0x14, 0x40, 0xff, 0x87, 0x39,
	0x1f, 0xbf, 0x1a, 0x61, 0xc6, 0xb1, 0x2f, 0x64, 0xcb, 0xe9, 0xd3, 0x0d, 0xf4, 0x1a, 0xc0, 0xa2,
	0x1d, 0xb3, 0x6e, 0xc4, 0x6c, 0xb3, 0xb8, 0xe1, 0x88, 0xf5, 0x65, 0x01, 0xe6, 0x28, 0x0c, 0xe4,
	0xad, 0x59, 0x65, 0x26, 0x90, 0xae, 0x03, 0xb9, 0x7a, 0x67, 0x0d, 0xa6, 0x44, 0x6c, 0xd9, 0xb7,
	0x1c, 0xe6, 0xa0, 0x1f, 0x00, 0xdc, 0x21, 0x9e, 0x4d, 0x07, 0x51, 0xf5, 0x1b, 0xbd, 0xac, 0x88,
	0x5e, 0x3a, 0x61, 0x20, 0x6f, 0x6b, 0x49, 0xd0, 0x47, 0x69, 0x66, 0x9b, 0xdc, 0x76, 0x9a, 0x93,
	0x0c, 0xe5, 0xe5, 0x0a, 0xdc, 0x10, 0x8e, 0x9e, 0x6f, 0x7a, 0x8c, 0x70, 0x42, 0x3d, 0x64, 0x41,
	0x38, 0xd3, 0x57, 0x2c, 0x78, 0x23, 0x9c, 0x95, 0xfa, 0x9e, 0xad, 0xe4, 0xd8, 0x84, 0x8b, 0x21,
	0x2c, 0xcc, 0xab, 0x32, 0x79, 0x9e, 0x0f, 0xc2, 0x40, 0xce, 0xcf, 0xaa, 0x71, 0xcf, 0x47, 0x21,
	0x3f, 0x2b, 0x82, 0xe6, 0x20, 0x0c, 0xb7, 0xe6, 0x2b, 0x0e, 0x47, 0x96, 0x71, 0x82, 0x2f, 0xc4,
	0xc8, 0xac, 0xd7, 0xab, 0xd7, 0x81, 0xac, 0xdc, 0xb9, 0x48, 0x87, 0x12, 0x8f, 0xeb, 0x68, 0xb6,
	0x44, 0x67, 0x64, 0x7d, 0x8b, 0x2f, 0x10, 0x83, 0x0f, 0x16, 0x34, 0x9e, 0xcc, 0x59, 0x2b, 0x0c,
	0xe4, 0x8d, 0x39, 0x6d, 0xef, 0x09, 0x6d, 0x63, 0x4e, 0x52, 0xcd, 0x41, 0x2e, 0xdc, 0x5e, 0x28,
	0x3a, 0x06, 0xb7, 0x72, 0x6f, 0x70, 0x9b, 0x73, 0x45, 0x12, 0x74, 0x5f, 0xc1, 0x55, 0xc6, 0x4d,
	0x3e, 0x62, 0xc5, 0xd5, 0x12, 0xa8, 0xe4, 0xab, 0x9f, 0x29, 0xb7, 0xbd, 0x3c, 0x94, 0xe9, 0x30,
	0x75, 0x45, 0xb4, 0x9e, 0x64, 0xa1, 0x1a, 0xcc, 0xe1, 0x73, 0x6c, 0xc7, 0x97, 0x52, 0xe6, 0x03,
	0x2e, 0xa5, 0x6c, 0x94, 0x16, 0x39, 0x90, 0x02, 0x37, 0x09, 0x33, 0xfa, 0xd4, 0xb7, 0xb1, 0xc1,
	0x27, 0x75, 0x8a, 0x59, 0x71, 0x5f, 0x3f, 0x20, 0xec, 0x20, 0xf2, 0x4c, 0x1b, 0x28, 0xbf, 0x03,
	0xf0, 0xd1, 0xc2, 0x84, 0x47, 0x13, 0x69, 0xf2, 0x91, 0x8f, 0x8f, 0x7c, 0x07, 0xfb, 0xe8, 0x10,
	0x66, 0xc6, 0x74, 0x01, 0x41, 0xd7, 0x97, 0x0b, 0xaa, 0xd8, 0x74, 0x80, 0xb9, 0xd5, 0xe7, 0xd3,
	0xc5, 0x29, 0xb1, 0x98, 0x6a, 0x5d, 0x70, 0xcc, 0x94, 0x6f, 0xf0, 0x79, 0x3d, 0x5a, 0xe8, 0xab,
	0xc3, 0x98, 0xa3, 0x43, 0xb8, 0x31, 0xed, 0xeb, 0xc3, 0xaf, 0xdf, 0xfc, 0x34, 0x39, 0x72, 0xef,
	0xad, 0xbd, 0xfb, 0xfd, 0x49, 0xa6, 0x41, 0x3d, 0x8e, 0x3d, 0xfe, 0xc5, 0xaf, 0x00, 0x16, 0x16,
	0xc9, 0x45, 0x9f, 0xc0, 0x47, 0x3d, 0xbd, 0xd6, 0xee, 0x6a, 0x3d, 0xed, 0xa8, 0x6d, 0x74, 0x7b,
	0xb5, 0xde, 0x71, 0xd7, 0x38, 0x6e, 0x77, 0x3b, 0xcd, 0x86, 0x76, 0xa0, 0x35, 0xf7, 0x0b, 0x29,
	0xf4, 0x29, 0x2c, 0xdd, 0x0c, 0x69, 0xe8, 0xcd, 0x5a, 0x4f, 0x6b, 0xb7, 0x8c, 0x96, 0x7e, 0x74,
	0xdc, 0x29, 0x00, 0x54, 0x86, 0xd2, 0xcd, 0xa8, 0xef, 0x6a, 0x9a, 0x08, 0xea, 0x6a, 0xad, 0x76,
	0x21, 0x8d, 0x3e, 0x87, 0x8f, 0xff, 0x3d, 0xa6, 0xf9, 0x7d, 0xb3, 0x71, 0x1c, 0x39, 0x0a, 0x4b,
	0xbb, 0xcb, 0xaf, 0x7f, 0x92, 0x52, 0xf5, 0xf6, 0xcf, 0xa1, 0x04, 0xde, 0x86, 0x12, 0x78, 0x1f,
	0x4a, 0xe0, 0xaf, 0x50, 0x02, 0x6f, 0xae, 0xa4, 0xd4, 0xfb, 0x2b, 0x29, 0xf5, 0xc7, 0x95, 0x94,
	0x7a, 0xf1, 0xdf, 0xa3, 0x7f, 0x3e, 0xf9, 0x78, 0x11, 0x2f, 0x00, 0x6b, 0x55, 0x84, 0x3c, 0xff,
	0x27, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x1b, 0x02, 0x05, 0xd9, 0x08, 0x00, 0x00,
}

func (this *Member) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Member)
	if !ok {
		that2, ok := that.(Member)
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
	if this.Address != that1.Address {
		return false
	}
	if this.GroupID != that1.GroupID {
		return false
	}
	if this.IsActive != that1.IsActive {
		return false
	}
	if !this.Since.Equal(that1.Since) {
		return false
	}
	return true
}
func (this *CurrentGroup) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*CurrentGroup)
	if !ok {
		that2, ok := that.(CurrentGroup)
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
	if this.GroupID != that1.GroupID {
		return false
	}
	if !this.ActiveTime.Equal(that1.ActiveTime) {
		return false
	}
	return true
}
func (this *Signing) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Signing)
	if !ok {
		that2, ok := that.(Signing)
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
	if len(this.FeePerSigner) != len(that1.FeePerSigner) {
		return false
	}
	for i := range this.FeePerSigner {
		if !this.FeePerSigner[i].Equal(&that1.FeePerSigner[i]) {
			return false
		}
	}
	if this.Requester != that1.Requester {
		return false
	}
	if this.CurrentGroupSigningID != that1.CurrentGroupSigningID {
		return false
	}
	if this.IncomingGroupSigningID != that1.IncomingGroupSigningID {
		return false
	}
	return true
}
func (this *GroupTransition) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GroupTransition)
	if !ok {
		that2, ok := that.(GroupTransition)
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
	if this.CurrentGroupID != that1.CurrentGroupID {
		return false
	}
	if !bytes.Equal(this.CurrentGroupPubKey, that1.CurrentGroupPubKey) {
		return false
	}
	if this.IncomingGroupID != that1.IncomingGroupID {
		return false
	}
	if !bytes.Equal(this.IncomingGroupPubKey, that1.IncomingGroupPubKey) {
		return false
	}
	if this.Status != that1.Status {
		return false
	}
	if !this.ExecTime.Equal(that1.ExecTime) {
		return false
	}
	if this.IsForceTransition != that1.IsForceTransition {
		return false
	}
	return true
}
func (this *GroupTransitionSignatureOrder) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GroupTransitionSignatureOrder)
	if !ok {
		that2, ok := that.(GroupTransitionSignatureOrder)
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
	if !bytes.Equal(this.PubKey, that1.PubKey) {
		return false
	}
	if !this.TransitionTime.Equal(that1.TransitionTime) {
		return false
	}
	return true
}
func (m *Member) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Member) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Member) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(
		m.Since,
		dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Since):],
	)
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintBandtss(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	if m.IsActive {
		i--
		if m.IsActive {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.GroupID != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.GroupID))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintBandtss(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *CurrentGroup) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CurrentGroup) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CurrentGroup) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(
		m.ActiveTime,
		dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ActiveTime):],
	)
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintBandtss(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if m.GroupID != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.GroupID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Signing) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Signing) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Signing) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IncomingGroupSigningID != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.IncomingGroupSigningID))
		i--
		dAtA[i] = 0x28
	}
	if m.CurrentGroupSigningID != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.CurrentGroupSigningID))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Requester) > 0 {
		i -= len(m.Requester)
		copy(dAtA[i:], m.Requester)
		i = encodeVarintBandtss(dAtA, i, uint64(len(m.Requester)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.FeePerSigner) > 0 {
		for iNdEx := len(m.FeePerSigner) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.FeePerSigner[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBandtss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.ID != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GroupTransition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GroupTransition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GroupTransition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IsForceTransition {
		i--
		if m.IsForceTransition {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	n3, err3 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(
		m.ExecTime,
		dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ExecTime):],
	)
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintBandtss(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x3a
	if m.Status != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x30
	}
	if len(m.IncomingGroupPubKey) > 0 {
		i -= len(m.IncomingGroupPubKey)
		copy(dAtA[i:], m.IncomingGroupPubKey)
		i = encodeVarintBandtss(dAtA, i, uint64(len(m.IncomingGroupPubKey)))
		i--
		dAtA[i] = 0x2a
	}
	if m.IncomingGroupID != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.IncomingGroupID))
		i--
		dAtA[i] = 0x20
	}
	if len(m.CurrentGroupPubKey) > 0 {
		i -= len(m.CurrentGroupPubKey)
		copy(dAtA[i:], m.CurrentGroupPubKey)
		i = encodeVarintBandtss(dAtA, i, uint64(len(m.CurrentGroupPubKey)))
		i--
		dAtA[i] = 0x1a
	}
	if m.CurrentGroupID != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.CurrentGroupID))
		i--
		dAtA[i] = 0x10
	}
	if m.SigningID != 0 {
		i = encodeVarintBandtss(dAtA, i, uint64(m.SigningID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GroupTransitionSignatureOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GroupTransitionSignatureOrder) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GroupTransitionSignatureOrder) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n4, err4 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(
		m.TransitionTime,
		dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.TransitionTime):],
	)
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintBandtss(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x12
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintBandtss(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBandtss(dAtA []byte, offset int, v uint64) int {
	offset -= sovBandtss(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Member) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovBandtss(uint64(l))
	}
	if m.GroupID != 0 {
		n += 1 + sovBandtss(uint64(m.GroupID))
	}
	if m.IsActive {
		n += 2
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Since)
	n += 1 + l + sovBandtss(uint64(l))
	return n
}

func (m *CurrentGroup) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.GroupID != 0 {
		n += 1 + sovBandtss(uint64(m.GroupID))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ActiveTime)
	n += 1 + l + sovBandtss(uint64(l))
	return n
}

func (m *Signing) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovBandtss(uint64(m.ID))
	}
	if len(m.FeePerSigner) > 0 {
		for _, e := range m.FeePerSigner {
			l = e.Size()
			n += 1 + l + sovBandtss(uint64(l))
		}
	}
	l = len(m.Requester)
	if l > 0 {
		n += 1 + l + sovBandtss(uint64(l))
	}
	if m.CurrentGroupSigningID != 0 {
		n += 1 + sovBandtss(uint64(m.CurrentGroupSigningID))
	}
	if m.IncomingGroupSigningID != 0 {
		n += 1 + sovBandtss(uint64(m.IncomingGroupSigningID))
	}
	return n
}

func (m *GroupTransition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigningID != 0 {
		n += 1 + sovBandtss(uint64(m.SigningID))
	}
	if m.CurrentGroupID != 0 {
		n += 1 + sovBandtss(uint64(m.CurrentGroupID))
	}
	l = len(m.CurrentGroupPubKey)
	if l > 0 {
		n += 1 + l + sovBandtss(uint64(l))
	}
	if m.IncomingGroupID != 0 {
		n += 1 + sovBandtss(uint64(m.IncomingGroupID))
	}
	l = len(m.IncomingGroupPubKey)
	if l > 0 {
		n += 1 + l + sovBandtss(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovBandtss(uint64(m.Status))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.ExecTime)
	n += 1 + l + sovBandtss(uint64(l))
	if m.IsForceTransition {
		n += 2
	}
	return n
}

func (m *GroupTransitionSignatureOrder) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovBandtss(uint64(l))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.TransitionTime)
	n += 1 + l + sovBandtss(uint64(l))
	return n
}

func sovBandtss(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBandtss(x uint64) (n int) {
	return sovBandtss(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Member) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBandtss
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
			return fmt.Errorf("proto: Member: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Member: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupID", wireType)
			}
			m.GroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsActive", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Since", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Since, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBandtss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBandtss
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
func (m *CurrentGroup) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBandtss
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
			return fmt.Errorf("proto: CurrentGroup: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CurrentGroup: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupID", wireType)
			}
			m.GroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActiveTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
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
			skippy, err := skipBandtss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBandtss
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
func (m *Signing) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBandtss
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
			return fmt.Errorf("proto: Signing: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Signing: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeePerSigner", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeePerSigner = append(m.FeePerSigner, types.Coin{})
			if err := m.FeePerSigner[len(m.FeePerSigner)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requester", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Requester = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentGroupSigningID", wireType)
			}
			m.CurrentGroupSigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentGroupSigningID |= github_com_bandprotocol_chain_v3_pkg_tss.SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncomingGroupSigningID", wireType)
			}
			m.IncomingGroupSigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IncomingGroupSigningID |= github_com_bandprotocol_chain_v3_pkg_tss.SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBandtss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBandtss
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
func (m *GroupTransition) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBandtss
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
			return fmt.Errorf("proto: GroupTransition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GroupTransition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningID", wireType)
			}
			m.SigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningID |= github_com_bandprotocol_chain_v3_pkg_tss.SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentGroupID", wireType)
			}
			m.CurrentGroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentGroupID |= github_com_bandprotocol_chain_v3_pkg_tss.GroupID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentGroupPubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CurrentGroupPubKey = append(m.CurrentGroupPubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.CurrentGroupPubKey == nil {
				m.CurrentGroupPubKey = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncomingGroupID", wireType)
			}
			m.IncomingGroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IncomingGroupID |= github_com_bandprotocol_chain_v3_pkg_tss.GroupID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncomingGroupPubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IncomingGroupPubKey = append(m.IncomingGroupPubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.IncomingGroupPubKey == nil {
				m.IncomingGroupPubKey = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= TransitionStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExecTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.ExecTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsForceTransition", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
			m.IsForceTransition = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipBandtss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBandtss
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
func (m *GroupTransitionSignatureOrder) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBandtss
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
			return fmt.Errorf("proto: GroupTransitionSignatureOrder: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GroupTransitionSignatureOrder: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKey = append(m.PubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKey == nil {
				m.PubKey = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransitionTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBandtss
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
				return ErrInvalidLengthBandtss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBandtss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.TransitionTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBandtss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBandtss
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
func skipBandtss(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBandtss
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
					return 0, ErrIntOverflowBandtss
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
					return 0, ErrIntOverflowBandtss
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
				return 0, ErrInvalidLengthBandtss
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBandtss
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBandtss
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBandtss        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBandtss          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBandtss = fmt.Errorf("proto: unexpected end of group")
)
