package types

import (
	"bytes"
	"fmt"
	"io"
	"math"
	math_bits "math/bits"
	"time"

	github_com_cometbft_cometbft_libs_bytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
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

// SigningStatus is an enumeration of the possible statuses of a signing.
type SigningStatus int32

const (
	// SIGNING_STATUS_UNSPECIFIED is the status of a signing that has not been specified.
	SIGNING_STATUS_UNSPECIFIED SigningStatus = 0
	// SIGNING_STATUS_WAITING is the status of a signing that is waiting to be signed in the protocol.
	SIGNING_STATUS_WAITING SigningStatus = 1
	// SIGNING_STATUS_SUCCESS is the status of a signing that has success in the protocol.
	SIGNING_STATUS_SUCCESS SigningStatus = 2
	// SIGNING_STATUS_FALLEN is the status of a signing that has fallen out of the protocol.
	SIGNING_STATUS_FALLEN SigningStatus = 3
)

var SigningStatus_name = map[int32]string{
	0: "SIGNING_STATUS_UNSPECIFIED",
	1: "SIGNING_STATUS_WAITING",
	2: "SIGNING_STATUS_SUCCESS",
	3: "SIGNING_STATUS_FALLEN",
}

var SigningStatus_value = map[string]int32{
	"SIGNING_STATUS_UNSPECIFIED": 0,
	"SIGNING_STATUS_WAITING":     1,
	"SIGNING_STATUS_SUCCESS":     2,
	"SIGNING_STATUS_FALLEN":      3,
}

func (x SigningStatus) String() string {
	return proto.EnumName(SigningStatus_name, int32(x))
}

func (SigningStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{0}
}

// GroupStatus is an enumeration of the possible statuses of a group.
type GroupStatus int32

const (
	// GROUP_STATUS_UNSPECIFIED is the status of a group that has not been specified.
	GROUP_STATUS_UNSPECIFIED GroupStatus = 0
	// GROUP_STATUS_ROUND_1 is the status of a group that is in the first round of the protocol.
	GROUP_STATUS_ROUND_1 GroupStatus = 1
	// GROUP_STATUS_ROUND_2 is the status of a group that is in the second round of the protocol.
	GROUP_STATUS_ROUND_2 GroupStatus = 2
	// GROUP_STATUS_ROUND_3 is the status of a group that is in the third round of the protocol.
	GROUP_STATUS_ROUND_3 GroupStatus = 3
	// GROUP_STATUS_ACTIVE is the status of a group that is actively participating in the protocol.
	GROUP_STATUS_ACTIVE GroupStatus = 4
	// GROUP_STATUS_EXPIRED is the status of a group that has expired in the protocol.
	GROUP_STATUS_EXPIRED GroupStatus = 5
	// GROUP_STATUS_FALLEN is the status of a group that has fallen out of the protocol.
	GROUP_STATUS_FALLEN GroupStatus = 6
)

var GroupStatus_name = map[int32]string{
	0: "GROUP_STATUS_UNSPECIFIED",
	1: "GROUP_STATUS_ROUND_1",
	2: "GROUP_STATUS_ROUND_2",
	3: "GROUP_STATUS_ROUND_3",
	4: "GROUP_STATUS_ACTIVE",
	5: "GROUP_STATUS_EXPIRED",
	6: "GROUP_STATUS_FALLEN",
}

var GroupStatus_value = map[string]int32{
	"GROUP_STATUS_UNSPECIFIED": 0,
	"GROUP_STATUS_ROUND_1":     1,
	"GROUP_STATUS_ROUND_2":     2,
	"GROUP_STATUS_ROUND_3":     3,
	"GROUP_STATUS_ACTIVE":      4,
	"GROUP_STATUS_EXPIRED":     5,
	"GROUP_STATUS_FALLEN":      6,
}

func (x GroupStatus) String() string {
	return proto.EnumName(GroupStatus_name, int32(x))
}

func (GroupStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{1}
}

// ComplaintStatus represents the status of a complaint.
type ComplaintStatus int32

const (
	// COMPLAINT_STATUS_UNSPECIFIED represents an undefined status of the complaint.
	COMPLAINT_STATUS_UNSPECIFIED ComplaintStatus = 0
	// COMPLAINT_STATUS_SUCCESS represents a successful complaint.
	COMPLAINT_STATUS_SUCCESS ComplaintStatus = 1
	// COMPLAINT_STATUS_FAILED represents a failed complaint.
	COMPLAINT_STATUS_FAILED ComplaintStatus = 2
)

var ComplaintStatus_name = map[int32]string{
	0: "COMPLAINT_STATUS_UNSPECIFIED",
	1: "COMPLAINT_STATUS_SUCCESS",
	2: "COMPLAINT_STATUS_FAILED",
}

var ComplaintStatus_value = map[string]int32{
	"COMPLAINT_STATUS_UNSPECIFIED": 0,
	"COMPLAINT_STATUS_SUCCESS":     1,
	"COMPLAINT_STATUS_FAILED":      2,
}

func (x ComplaintStatus) String() string {
	return proto.EnumName(ComplaintStatus_name, int32(x))
}

func (ComplaintStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{2}
}

// Group is a type representing a participant group in a Distributed Key Generation or signing process.
type Group struct {
	// id is the unique identifier of the group.
	ID GroupID `protobuf:"varint,1,opt,name=id,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID"               json:"id,omitempty"`
	// size is the number of members in the group.
	Size_ uint64 `protobuf:"varint,2,opt,name=size,proto3"                                                                                    json:"size,omitempty"`
	// threshold is the minimum number of members needed to generate a valid signature.
	Threshold uint64 `protobuf:"varint,3,opt,name=threshold,proto3"                                                                               json:"threshold,omitempty"`
	// pub_key is the public key generated by the group after successful completion of the DKG process.
	PubKey Point `protobuf:"bytes,4,opt,name=pub_key,json=pubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point" json:"pub_key,omitempty"`
	// status represents the current stage of the group in the DKG or signing process.
	Status GroupStatus `protobuf:"varint,5,opt,name=status,proto3,enum=band.tss.v1beta1.GroupStatus"                                                json:"status,omitempty"`
	// created_height is the block height when the group was created.
	CreatedHeight uint64 `protobuf:"varint,6,opt,name=created_height,json=createdHeight,proto3"                                                       json:"created_height,omitempty"`
	// module_owner is the module that creates this group.
	ModuleOwner string `protobuf:"bytes,7,opt,name=module_owner,json=moduleOwner,proto3"                                                            json:"module_owner,omitempty"`
}

func (m *Group) Reset()         { *m = Group{} }
func (m *Group) String() string { return proto.CompactTextString(m) }
func (*Group) ProtoMessage()    {}
func (*Group) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{0}
}
func (m *Group) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Group) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Group.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Group) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Group.Merge(m, src)
}
func (m *Group) XXX_Size() int {
	return m.Size()
}
func (m *Group) XXX_DiscardUnknown() {
	xxx_messageInfo_Group.DiscardUnknown(m)
}

var xxx_messageInfo_Group proto.InternalMessageInfo

func (m *Group) GetID() GroupID {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Group) GetSize_() uint64 {
	if m != nil {
		return m.Size_
	}
	return 0
}

func (m *Group) GetThreshold() uint64 {
	if m != nil {
		return m.Threshold
	}
	return 0
}

func (m *Group) GetPubKey() Point {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *Group) GetStatus() GroupStatus {
	if m != nil {
		return m.Status
	}
	return GROUP_STATUS_UNSPECIFIED
}

func (m *Group) GetCreatedHeight() uint64 {
	if m != nil {
		return m.CreatedHeight
	}
	return 0
}

func (m *Group) GetModuleOwner() string {
	if m != nil {
		return m.ModuleOwner
	}
	return ""
}

// GroupResult is a tss group result from querying tss group information.
type GroupResult struct {
	// group defines the group object containing group information.
	Group Group `protobuf:"bytes,1,opt,name=group,proto3"                                                                                 json:"group"`
	// dkg_context defines the DKG context data.
	DKGContext github_com_cometbft_cometbft_libs_bytes.HexBytes `protobuf:"bytes,2,opt,name=dkg_context,json=dkgContext,proto3,casttype=github.com/cometbft/cometbft/libs/bytes.HexBytes" json:"dkg_context,omitempty"`
	// members is the list of members in the group.
	Members []Member `protobuf:"bytes,3,rep,name=members,proto3"                                                                               json:"members"`
	// round1_infos is the list of Round 1 information.
	Round1Infos []Round1Info `protobuf:"bytes,4,rep,name=round1_infos,json=round1Infos,proto3"                                                         json:"round1_infos"`
	// round2_infos is the list of Round 2 information.
	Round2Infos []Round2Info `protobuf:"bytes,5,rep,name=round2_infos,json=round2Infos,proto3"                                                         json:"round2_infos"`
	// complaints_with_status is the list of complaints with status.
	ComplaintsWithStatus []ComplaintsWithStatus `protobuf:"bytes,6,rep,name=complaints_with_status,json=complaintsWithStatus,proto3"                                      json:"complaints_with_status"`
	// confirms is the list of confirms.
	Confirms []Confirm `protobuf:"bytes,7,rep,name=confirms,proto3"                                                                              json:"confirms"`
}

func (m *GroupResult) Reset()         { *m = GroupResult{} }
func (m *GroupResult) String() string { return proto.CompactTextString(m) }
func (*GroupResult) ProtoMessage()    {}
func (*GroupResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{1}
}
func (m *GroupResult) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GroupResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GroupResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GroupResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupResult.Merge(m, src)
}
func (m *GroupResult) XXX_Size() int {
	return m.Size()
}
func (m *GroupResult) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupResult.DiscardUnknown(m)
}

var xxx_messageInfo_GroupResult proto.InternalMessageInfo

func (m *GroupResult) GetGroup() Group {
	if m != nil {
		return m.Group
	}
	return Group{}
}

func (m *GroupResult) GetDKGContext() github_com_cometbft_cometbft_libs_bytes.HexBytes {
	if m != nil {
		return m.DKGContext
	}
	return nil
}

func (m *GroupResult) GetMembers() []Member {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *GroupResult) GetRound1Infos() []Round1Info {
	if m != nil {
		return m.Round1Infos
	}
	return nil
}

func (m *GroupResult) GetRound2Infos() []Round2Info {
	if m != nil {
		return m.Round2Infos
	}
	return nil
}

func (m *GroupResult) GetComplaintsWithStatus() []ComplaintsWithStatus {
	if m != nil {
		return m.ComplaintsWithStatus
	}
	return nil
}

func (m *GroupResult) GetConfirms() []Confirm {
	if m != nil {
		return m.Confirms
	}
	return nil
}

// Round1Info contains all necessary information for handling round 1 of the DKG process.
type Round1Info struct {
	// member_id is the unique identifier of a group member.
	MemberID MemberID `protobuf:"varint,1,opt,name=member_id,json=memberId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID"                      json:"member_id,omitempty"`
	// coefficients_commits is a list of commitments to the coefficients of the member's secret polynomial.
	CoefficientCommits Points `protobuf:"bytes,2,rep,name=coefficient_commits,json=coefficientCommits,proto3,castrepeated=github.com/bandprotocol/falcon/internal/bandchain/tss.Points" json:"coefficient_commits,omitempty"`
	// one_time_pub_key is the one-time public key used by the member to encrypt secret shares.
	OneTimePubKey Point `protobuf:"bytes,3,opt,name=one_time_pub_key,json=oneTimePubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"              json:"one_time_pub_key,omitempty"`
	// a0_signature is the member's signature on the first coefficient of its secret polynomial.
	A0Signature Signature `protobuf:"bytes,4,opt,name=a0_signature,json=a0Signature,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Signature"                json:"a0_signature,omitempty"`
	// one_time_signature is the member's signature on its one-time public key.
	OneTimeSignature Signature `protobuf:"bytes,5,opt,name=one_time_signature,json=oneTimeSignature,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Signature"     json:"one_time_signature,omitempty"`
}

func (m *Round1Info) Reset()         { *m = Round1Info{} }
func (m *Round1Info) String() string { return proto.CompactTextString(m) }
func (*Round1Info) ProtoMessage()    {}
func (*Round1Info) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{2}
}
func (m *Round1Info) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Round1Info) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Round1Info.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Round1Info) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Round1Info.Merge(m, src)
}
func (m *Round1Info) XXX_Size() int {
	return m.Size()
}
func (m *Round1Info) XXX_DiscardUnknown() {
	xxx_messageInfo_Round1Info.DiscardUnknown(m)
}

var xxx_messageInfo_Round1Info proto.InternalMessageInfo

func (m *Round1Info) GetMemberID() MemberID {
	if m != nil {
		return m.MemberID
	}
	return 0
}

func (m *Round1Info) GetCoefficientCommits() Points {
	if m != nil {
		return m.CoefficientCommits
	}
	return nil
}

func (m *Round1Info) GetOneTimePubKey() Point {
	if m != nil {
		return m.OneTimePubKey
	}
	return nil
}

func (m *Round1Info) GetA0Signature() Signature {
	if m != nil {
		return m.A0Signature
	}
	return nil
}

func (m *Round1Info) GetOneTimeSignature() Signature {
	if m != nil {
		return m.OneTimeSignature
	}
	return nil
}

// Round2Info contains all necessary information for handling round 2 of the DKG process.
type Round2Info struct {
	// member_id is the unique identifier of a group member.
	MemberID MemberID `protobuf:"varint,1,opt,name=member_id,json=memberId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID"                                      json:"member_id,omitempty"`
	// encrypted_secret_shares is a list of secret shares encrypted under the public keys of other members.
	EncryptedSecretShares EncSecretShares `protobuf:"bytes,2,rep,name=encrypted_secret_shares,json=encryptedSecretShares,proto3,castrepeated=github.com/bandprotocol/falcon/internal/bandchain/tss.EncSecretShares" json:"encrypted_secret_shares,omitempty"`
}

func (m *Round2Info) Reset()         { *m = Round2Info{} }
func (m *Round2Info) String() string { return proto.CompactTextString(m) }
func (*Round2Info) ProtoMessage()    {}
func (*Round2Info) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{3}
}
func (m *Round2Info) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Round2Info) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Round2Info.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Round2Info) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Round2Info.Merge(m, src)
}
func (m *Round2Info) XXX_Size() int {
	return m.Size()
}
func (m *Round2Info) XXX_DiscardUnknown() {
	xxx_messageInfo_Round2Info.DiscardUnknown(m)
}

var xxx_messageInfo_Round2Info proto.InternalMessageInfo

func (m *Round2Info) GetMemberID() MemberID {
	if m != nil {
		return m.MemberID
	}
	return 0
}

func (m *Round2Info) GetEncryptedSecretShares() EncSecretShares {
	if m != nil {
		return m.EncryptedSecretShares
	}
	return nil
}

// DE contains the public parts of a member's decryption and encryption keys.
type DE struct {
	// pub_d is the public value of own commitment (D).
	PubD Point `protobuf:"bytes,1,opt,name=pub_d,json=pubD,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point" json:"pub_d,omitempty"`
	// pub_e is the public value of own commitment (E).
	PubE Point `protobuf:"bytes,2,opt,name=pub_e,json=pubE,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point" json:"pub_e,omitempty"`
}

func (m *DE) Reset()         { *m = DE{} }
func (m *DE) String() string { return proto.CompactTextString(m) }
func (*DE) ProtoMessage()    {}
func (*DE) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{4}
}
func (m *DE) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DE) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DE.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DE) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DE.Merge(m, src)
}
func (m *DE) XXX_Size() int {
	return m.Size()
}
func (m *DE) XXX_DiscardUnknown() {
	xxx_messageInfo_DE.DiscardUnknown(m)
}

var xxx_messageInfo_DE proto.InternalMessageInfo

func (m *DE) GetPubD() Point {
	if m != nil {
		return m.PubD
	}
	return nil
}

func (m *DE) GetPubE() Point {
	if m != nil {
		return m.PubE
	}
	return nil
}

// DEQueue is a simple queue data structure contains index of existing DE objects of each member.
type DEQueue struct {
	// head is the current index of the first element in the queue.
	Head uint64 `protobuf:"varint,1,opt,name=head,proto3" json:"head,omitempty"`
	// tail is the current index of the last element in the queue.
	Tail uint64 `protobuf:"varint,2,opt,name=tail,proto3" json:"tail,omitempty"`
}

func (m *DEQueue) Reset()         { *m = DEQueue{} }
func (m *DEQueue) String() string { return proto.CompactTextString(m) }
func (*DEQueue) ProtoMessage()    {}
func (*DEQueue) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{5}
}
func (m *DEQueue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DEQueue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DEQueue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DEQueue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DEQueue.Merge(m, src)
}
func (m *DEQueue) XXX_Size() int {
	return m.Size()
}
func (m *DEQueue) XXX_DiscardUnknown() {
	xxx_messageInfo_DEQueue.DiscardUnknown(m)
}

var xxx_messageInfo_DEQueue proto.InternalMessageInfo

func (m *DEQueue) GetHead() uint64 {
	if m != nil {
		return m.Head
	}
	return 0
}

func (m *DEQueue) GetTail() uint64 {
	if m != nil {
		return m.Tail
	}
	return 0
}

// Signing contains all necessary information for handling a signing request.
type Signing struct {
	// id is the unique identifier of the signing.
	ID SigningID `protobuf:"varint,1,opt,name=id,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID"                            json:"id,omitempty"`
	// current_attempt is the latest round number that signing has been attempted.
	CurrentAttempt uint64 `protobuf:"varint,2,opt,name=current_attempt,json=currentAttempt,proto3"                                                                    json:"current_attempt,omitempty"`
	// group_id is the unique identifier of the group.
	GroupID GroupID `protobuf:"varint,3,opt,name=group_id,json=groupId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID"           json:"group_id,omitempty"`
	// group_pub_key is the public key of the group that sign this message.
	GroupPubKey Point `protobuf:"bytes,4,opt,name=group_pub_key,json=groupPubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"     json:"group_pub_key,omitempty"`
	// message is the message to be signed.
	Message github_com_cometbft_cometbft_libs_bytes.HexBytes `protobuf:"bytes,5,opt,name=message,proto3,casttype=github.com/cometbft/cometbft/libs/bytes.HexBytes"                                       json:"message,omitempty"`
	// group_pub_nonce is the public nonce generated by the group for this signing process.
	GroupPubNonce Point `protobuf:"bytes,6,opt,name=group_pub_nonce,json=groupPubNonce,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point" json:"group_pub_nonce,omitempty"`
	// signature is the group's signature on the message.
	Signature Signature `protobuf:"bytes,7,opt,name=signature,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Signature"                      json:"signature,omitempty"`
	// status represents the current stage of the signing in the signing process.
	Status SigningStatus `protobuf:"varint,8,opt,name=status,proto3,enum=band.tss.v1beta1.SigningStatus"                                                             json:"status,omitempty"`
	// created_height is the block height when the signing was created.
	CreatedHeight uint64 `protobuf:"varint,9,opt,name=created_height,json=createdHeight,proto3"                                                                      json:"created_height,omitempty"`
	// created_timestamp is the block timestamp when the signing was created.
	CreatedTimestamp time.Time `protobuf:"bytes,10,opt,name=created_timestamp,json=createdTimestamp,proto3,stdtime"                                                        json:"created_timestamp"`
}

func (m *Signing) Reset()         { *m = Signing{} }
func (m *Signing) String() string { return proto.CompactTextString(m) }
func (*Signing) ProtoMessage()    {}
func (*Signing) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{6}
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

func (m *Signing) GetCurrentAttempt() uint64 {
	if m != nil {
		return m.CurrentAttempt
	}
	return 0
}

func (m *Signing) GetGroupID() GroupID {
	if m != nil {
		return m.GroupID
	}
	return 0
}

func (m *Signing) GetGroupPubKey() Point {
	if m != nil {
		return m.GroupPubKey
	}
	return nil
}

func (m *Signing) GetMessage() github_com_cometbft_cometbft_libs_bytes.HexBytes {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *Signing) GetGroupPubNonce() Point {
	if m != nil {
		return m.GroupPubNonce
	}
	return nil
}

func (m *Signing) GetSignature() Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Signing) GetStatus() SigningStatus {
	if m != nil {
		return m.Status
	}
	return SIGNING_STATUS_UNSPECIFIED
}

func (m *Signing) GetCreatedHeight() uint64 {
	if m != nil {
		return m.CreatedHeight
	}
	return 0
}

func (m *Signing) GetCreatedTimestamp() time.Time {
	if m != nil {
		return m.CreatedTimestamp
	}
	return time.Time{}
}

// SigningAttempt contains a member that has been assigned to and expiration block height of
// the specific attempt.
type SigningAttempt struct {
	// signing_id is the unique identifier of the signing.
	SigningID SigningID `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID" json:"signing_id,omitempty"`
	// attempt is the number of round that this signing has been attempted.
	Attempt uint64 `protobuf:"varint,2,opt,name=attempt,proto3"                                                                                            json:"attempt,omitempty"`
	// expired_height is the block height when this signing attempt was expired.
	ExpiredHeight uint64 `protobuf:"varint,3,opt,name=expired_height,json=expiredHeight,proto3"                                                                  json:"expired_height,omitempty"`
	// assigned_members is a list of members assigned to the signing process.
	AssignedMembers []AssignedMember `protobuf:"bytes,4,rep,name=assigned_members,json=assignedMembers,proto3"                                                               json:"assigned_members"`
}

func (m *SigningAttempt) Reset()         { *m = SigningAttempt{} }
func (m *SigningAttempt) String() string { return proto.CompactTextString(m) }
func (*SigningAttempt) ProtoMessage()    {}
func (*SigningAttempt) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{7}
}
func (m *SigningAttempt) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SigningAttempt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SigningAttempt.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SigningAttempt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SigningAttempt.Merge(m, src)
}
func (m *SigningAttempt) XXX_Size() int {
	return m.Size()
}
func (m *SigningAttempt) XXX_DiscardUnknown() {
	xxx_messageInfo_SigningAttempt.DiscardUnknown(m)
}

var xxx_messageInfo_SigningAttempt proto.InternalMessageInfo

func (m *SigningAttempt) GetSigningID() SigningID {
	if m != nil {
		return m.SigningID
	}
	return 0
}

func (m *SigningAttempt) GetAttempt() uint64 {
	if m != nil {
		return m.Attempt
	}
	return 0
}

func (m *SigningAttempt) GetExpiredHeight() uint64 {
	if m != nil {
		return m.ExpiredHeight
	}
	return 0
}

func (m *SigningAttempt) GetAssignedMembers() []AssignedMember {
	if m != nil {
		return m.AssignedMembers
	}
	return nil
}

// AssignedMember is a type representing a member that has been assigned to a signing process.
type AssignedMember struct {
	// member_id is the unique identifier of the member.
	MemberID MemberID `protobuf:"varint,1,opt,name=member_id,json=memberId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID"        json:"member_id,omitempty"`
	// member is the human-readable name of the member.
	Address string `protobuf:"bytes,2,opt,name=address,proto3"                                                                                                 json:"address,omitempty"`
	// pub_key is the public part of a member.
	PubKey Point `protobuf:"bytes,3,opt,name=pub_key,json=pubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"                json:"pub_key,omitempty"`
	// pub_d is the public part of a member's decryption key.
	PubD Point `protobuf:"bytes,4,opt,name=pub_d,json=pubD,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"                    json:"pub_d,omitempty"`
	// pub_e is the public part of a member's encryption key.
	PubE Point `protobuf:"bytes,5,opt,name=pub_e,json=pubE,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"                    json:"pub_e,omitempty"`
	// binding_factor is the binding factor of the member for the signing process.
	BindingFactor Scalar `protobuf:"bytes,6,opt,name=binding_factor,json=bindingFactor,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Scalar" json:"binding_factor,omitempty"`
	// pub_nonce is the public nonce of the member for the signing process.
	PubNonce Point `protobuf:"bytes,7,opt,name=pub_nonce,json=pubNonce,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"            json:"pub_nonce,omitempty"`
}

func (m *AssignedMember) Reset()         { *m = AssignedMember{} }
func (m *AssignedMember) String() string { return proto.CompactTextString(m) }
func (*AssignedMember) ProtoMessage()    {}
func (*AssignedMember) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{8}
}
func (m *AssignedMember) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AssignedMember) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AssignedMember.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AssignedMember) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AssignedMember.Merge(m, src)
}
func (m *AssignedMember) XXX_Size() int {
	return m.Size()
}
func (m *AssignedMember) XXX_DiscardUnknown() {
	xxx_messageInfo_AssignedMember.DiscardUnknown(m)
}

var xxx_messageInfo_AssignedMember proto.InternalMessageInfo

func (m *AssignedMember) GetMemberID() MemberID {
	if m != nil {
		return m.MemberID
	}
	return 0
}

func (m *AssignedMember) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *AssignedMember) GetPubKey() Point {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *AssignedMember) GetPubD() Point {
	if m != nil {
		return m.PubD
	}
	return nil
}

func (m *AssignedMember) GetPubE() Point {
	if m != nil {
		return m.PubE
	}
	return nil
}

func (m *AssignedMember) GetBindingFactor() Scalar {
	if m != nil {
		return m.BindingFactor
	}
	return nil
}

func (m *AssignedMember) GetPubNonce() Point {
	if m != nil {
		return m.PubNonce
	}
	return nil
}

// PendingSignings is a list of all signing processes that are currently pending.
type PendingSignings struct {
	// signing_ids is a list of identifiers for the signing processes.
	SigningIds []uint64 `protobuf:"varint,1,rep,packed,name=signing_ids,json=signingIds,proto3" json:"signing_ids,omitempty"`
}

func (m *PendingSignings) Reset()         { *m = PendingSignings{} }
func (m *PendingSignings) String() string { return proto.CompactTextString(m) }
func (*PendingSignings) ProtoMessage()    {}
func (*PendingSignings) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{9}
}
func (m *PendingSignings) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PendingSignings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PendingSignings.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PendingSignings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingSignings.Merge(m, src)
}
func (m *PendingSignings) XXX_Size() int {
	return m.Size()
}
func (m *PendingSignings) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingSignings.DiscardUnknown(m)
}

var xxx_messageInfo_PendingSignings proto.InternalMessageInfo

func (m *PendingSignings) GetSigningIds() []uint64 {
	if m != nil {
		return m.SigningIds
	}
	return nil
}

// Member is a type representing a member of the group.
type Member struct {
	// id is the unique identifier of a member.
	ID MemberID `protobuf:"varint,1,opt,name=id,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID"                   json:"id,omitempty"`
	// group_id is the group id of this member.
	GroupID GroupID `protobuf:"varint,2,opt,name=group_id,json=groupId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID" json:"group_id,omitempty"`
	// address is the address of the member.
	Address string `protobuf:"bytes,3,opt,name=address,proto3"                                                                                       json:"address,omitempty"`
	// pub_key is the public key of the member.
	PubKey Point `protobuf:"bytes,4,opt,name=pub_key,json=pubKey,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"      json:"pub_key,omitempty"`
	// is_malicious is a boolean flag indicating whether the member is considered malicious.
	IsMalicious bool `protobuf:"varint,5,opt,name=is_malicious,json=isMalicious,proto3"                                                                json:"is_malicious,omitempty"`
	// is_active is a boolean flag indicating whether the member is currently active in the protocol.
	IsActive bool `protobuf:"varint,6,opt,name=is_active,json=isActive,proto3"                                                                      json:"is_active,omitempty"`
}

func (m *Member) Reset()         { *m = Member{} }
func (m *Member) String() string { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()    {}
func (*Member) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{10}
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

func (m *Member) GetID() MemberID {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Member) GetGroupID() GroupID {
	if m != nil {
		return m.GroupID
	}
	return 0
}

func (m *Member) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Member) GetPubKey() Point {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *Member) GetIsMalicious() bool {
	if m != nil {
		return m.IsMalicious
	}
	return false
}

func (m *Member) GetIsActive() bool {
	if m != nil {
		return m.IsActive
	}
	return false
}

// Confirm is a message type used to confirm participation in the protocol.
type Confirm struct {
	// member_id is the unique identifier of a group member.
	MemberID MemberID `protobuf:"varint,1,opt,name=member_id,json=memberId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID"           json:"member_id,omitempty"`
	// own_pub_key_sig is a signature over the member's own public key.
	OwnPubKeySig Signature `protobuf:"bytes,2,opt,name=own_pub_key_sig,json=ownPubKeySig,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Signature" json:"own_pub_key_sig,omitempty"`
}

func (m *Confirm) Reset()         { *m = Confirm{} }
func (m *Confirm) String() string { return proto.CompactTextString(m) }
func (*Confirm) ProtoMessage()    {}
func (*Confirm) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{11}
}
func (m *Confirm) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Confirm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Confirm.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Confirm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Confirm.Merge(m, src)
}
func (m *Confirm) XXX_Size() int {
	return m.Size()
}
func (m *Confirm) XXX_DiscardUnknown() {
	xxx_messageInfo_Confirm.DiscardUnknown(m)
}

var xxx_messageInfo_Confirm proto.InternalMessageInfo

func (m *Confirm) GetMemberID() MemberID {
	if m != nil {
		return m.MemberID
	}
	return 0
}

func (m *Confirm) GetOwnPubKeySig() Signature {
	if m != nil {
		return m.OwnPubKeySig
	}
	return nil
}

// Complaint is a message type used to issue a complaint against a member.
type Complaint struct {
	// complainant is the member issuing the complaint.
	Complainant MemberID `protobuf:"varint,1,opt,name=complainant,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID"        json:"complainant,omitempty"`
	// respondent is the member against whom the complaint is issued.
	Respondent MemberID `protobuf:"varint,2,opt,name=respondent,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID"         json:"respondent,omitempty"`
	// key_sym is a symmetric key between respondent's private key and respondent's public key.
	KeySym Point `protobuf:"bytes,3,opt,name=key_sym,json=keySym,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Point"    json:"key_sym,omitempty"`
	// signature is the complaint signature that can do a symmetric key validation and complaint verification.
	Signature ComplaintSignature `protobuf:"bytes,4,opt,name=signature,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.ComplaintSignature" json:"signature,omitempty"`
}

func (m *Complaint) Reset()         { *m = Complaint{} }
func (m *Complaint) String() string { return proto.CompactTextString(m) }
func (*Complaint) ProtoMessage()    {}
func (*Complaint) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{12}
}
func (m *Complaint) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Complaint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Complaint.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Complaint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Complaint.Merge(m, src)
}
func (m *Complaint) XXX_Size() int {
	return m.Size()
}
func (m *Complaint) XXX_DiscardUnknown() {
	xxx_messageInfo_Complaint.DiscardUnknown(m)
}

var xxx_messageInfo_Complaint proto.InternalMessageInfo

func (m *Complaint) GetComplainant() MemberID {
	if m != nil {
		return m.Complainant
	}
	return 0
}

func (m *Complaint) GetRespondent() MemberID {
	if m != nil {
		return m.Respondent
	}
	return 0
}

func (m *Complaint) GetKeySym() Point {
	if m != nil {
		return m.KeySym
	}
	return nil
}

func (m *Complaint) GetSignature() ComplaintSignature {
	if m != nil {
		return m.Signature
	}
	return nil
}

// ComplaintWithStatus contains information about a complaint with its status.
type ComplaintWithStatus struct {
	// complaint is the information about the complaint.
	Complaint Complaint `protobuf:"bytes,1,opt,name=complaint,proto3"                                                                    json:"complaint"`
	// complaint_status is the status of the complaint.
	ComplaintStatus ComplaintStatus `protobuf:"varint,2,opt,name=complaint_status,json=complaintStatus,proto3,enum=band.tss.v1beta1.ComplaintStatus" json:"complaint_status,omitempty"`
}

func (m *ComplaintWithStatus) Reset()         { *m = ComplaintWithStatus{} }
func (m *ComplaintWithStatus) String() string { return proto.CompactTextString(m) }
func (*ComplaintWithStatus) ProtoMessage()    {}
func (*ComplaintWithStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{13}
}
func (m *ComplaintWithStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ComplaintWithStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ComplaintWithStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ComplaintWithStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComplaintWithStatus.Merge(m, src)
}
func (m *ComplaintWithStatus) XXX_Size() int {
	return m.Size()
}
func (m *ComplaintWithStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ComplaintWithStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ComplaintWithStatus proto.InternalMessageInfo

func (m *ComplaintWithStatus) GetComplaint() Complaint {
	if m != nil {
		return m.Complaint
	}
	return Complaint{}
}

func (m *ComplaintWithStatus) GetComplaintStatus() ComplaintStatus {
	if m != nil {
		return m.ComplaintStatus
	}
	return COMPLAINT_STATUS_UNSPECIFIED
}

// ComplaintsWithStatus contains information about multiple complaints and their status from a single member.
type ComplaintsWithStatus struct {
	// member_id is the identifier of the member filing the complaints.
	MemberID MemberID `protobuf:"varint,1,opt,name=member_id,json=memberId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID" json:"member_id,omitempty"`
	// complaints_with_status is the list of complaints with their status from this member.
	ComplaintsWithStatus []ComplaintWithStatus `protobuf:"bytes,2,rep,name=complaints_with_status,json=complaintsWithStatus,proto3"                                                 json:"complaints_with_status"`
}

func (m *ComplaintsWithStatus) Reset()         { *m = ComplaintsWithStatus{} }
func (m *ComplaintsWithStatus) String() string { return proto.CompactTextString(m) }
func (*ComplaintsWithStatus) ProtoMessage()    {}
func (*ComplaintsWithStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{14}
}
func (m *ComplaintsWithStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ComplaintsWithStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ComplaintsWithStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ComplaintsWithStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComplaintsWithStatus.Merge(m, src)
}
func (m *ComplaintsWithStatus) XXX_Size() int {
	return m.Size()
}
func (m *ComplaintsWithStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ComplaintsWithStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ComplaintsWithStatus proto.InternalMessageInfo

func (m *ComplaintsWithStatus) GetMemberID() MemberID {
	if m != nil {
		return m.MemberID
	}
	return 0
}

func (m *ComplaintsWithStatus) GetComplaintsWithStatus() []ComplaintWithStatus {
	if m != nil {
		return m.ComplaintsWithStatus
	}
	return nil
}

// PendingProcessGroups is a list of groups that are waiting to be processed.
type PendingProcessGroups struct {
	// group_ids is a list of group IDs.
	GroupIDs []GroupID `protobuf:"varint,1,rep,packed,name=group_ids,json=groupIds,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.GroupID" json:"group_ids,omitempty"`
}

func (m *PendingProcessGroups) Reset()         { *m = PendingProcessGroups{} }
func (m *PendingProcessGroups) String() string { return proto.CompactTextString(m) }
func (*PendingProcessGroups) ProtoMessage()    {}
func (*PendingProcessGroups) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{15}
}
func (m *PendingProcessGroups) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PendingProcessGroups) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PendingProcessGroups.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PendingProcessGroups) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingProcessGroups.Merge(m, src)
}
func (m *PendingProcessGroups) XXX_Size() int {
	return m.Size()
}
func (m *PendingProcessGroups) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingProcessGroups.DiscardUnknown(m)
}

var xxx_messageInfo_PendingProcessGroups proto.InternalMessageInfo

func (m *PendingProcessGroups) GetGroupIDs() []GroupID {
	if m != nil {
		return m.GroupIDs
	}
	return nil
}

// PendingProcessSignigns is a list of signings that are waiting to be processed.
type PendingProcessSignings struct {
	// signing_ids is a list of signing IDs.
	SigningIDs []SigningID `protobuf:"varint,1,rep,packed,name=signing_ids,json=signingIds,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID" json:"signing_ids,omitempty"`
}

func (m *PendingProcessSignings) Reset()         { *m = PendingProcessSignings{} }
func (m *PendingProcessSignings) String() string { return proto.CompactTextString(m) }
func (*PendingProcessSignings) ProtoMessage()    {}
func (*PendingProcessSignings) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{16}
}
func (m *PendingProcessSignings) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PendingProcessSignings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PendingProcessSignings.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PendingProcessSignings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingProcessSignings.Merge(m, src)
}
func (m *PendingProcessSignings) XXX_Size() int {
	return m.Size()
}
func (m *PendingProcessSignings) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingProcessSignings.DiscardUnknown(m)
}

var xxx_messageInfo_PendingProcessSignings proto.InternalMessageInfo

func (m *PendingProcessSignings) GetSigningIDs() []SigningID {
	if m != nil {
		return m.SigningIDs
	}
	return nil
}

// PartialSignature contains information about a member's partial signature.
type PartialSignature struct {
	// signing_id is the unique identifier of the signing.
	SigningID SigningID `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID" json:"signing_id,omitempty"`
	// signing_attempt is the number of attempts for this signing.
	SigningAttempt uint64 `protobuf:"varint,2,opt,name=signing_attempt,json=signingAttempt,proto3"                                                                json:"signing_attempt,omitempty"`
	// member_id is the identifier of the member providing the partial signature.
	MemberID MemberID `protobuf:"varint,3,opt,name=member_id,json=memberId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.MemberID"    json:"member_id,omitempty"`
	// signature is the partial signature provided by this member.
	Signature Signature `protobuf:"bytes,4,opt,name=signature,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.Signature"                  json:"signature,omitempty"`
}

func (m *PartialSignature) Reset()         { *m = PartialSignature{} }
func (m *PartialSignature) String() string { return proto.CompactTextString(m) }
func (*PartialSignature) ProtoMessage()    {}
func (*PartialSignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{17}
}
func (m *PartialSignature) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PartialSignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PartialSignature.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PartialSignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PartialSignature.Merge(m, src)
}
func (m *PartialSignature) XXX_Size() int {
	return m.Size()
}
func (m *PartialSignature) XXX_DiscardUnknown() {
	xxx_messageInfo_PartialSignature.DiscardUnknown(m)
}

var xxx_messageInfo_PartialSignature proto.InternalMessageInfo

func (m *PartialSignature) GetSigningID() SigningID {
	if m != nil {
		return m.SigningID
	}
	return 0
}

func (m *PartialSignature) GetSigningAttempt() uint64 {
	if m != nil {
		return m.SigningAttempt
	}
	return 0
}

func (m *PartialSignature) GetMemberID() MemberID {
	if m != nil {
		return m.MemberID
	}
	return 0
}

func (m *PartialSignature) GetSignature() Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

// TextSignatureOrder defines a general text signature order.
type TextSignatureOrder struct {
	// message is the data that needs to be signed.
	Message github_com_cometbft_cometbft_libs_bytes.HexBytes `protobuf:"bytes,1,opt,name=message,proto3,casttype=github.com/cometbft/cometbft/libs/bytes.HexBytes" json:"message,omitempty"`
}

func (m *TextSignatureOrder) Reset()         { *m = TextSignatureOrder{} }
func (m *TextSignatureOrder) String() string { return proto.CompactTextString(m) }
func (*TextSignatureOrder) ProtoMessage()    {}
func (*TextSignatureOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{18}
}
func (m *TextSignatureOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TextSignatureOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TextSignatureOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TextSignatureOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TextSignatureOrder.Merge(m, src)
}
func (m *TextSignatureOrder) XXX_Size() int {
	return m.Size()
}
func (m *TextSignatureOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_TextSignatureOrder.DiscardUnknown(m)
}

var xxx_messageInfo_TextSignatureOrder proto.InternalMessageInfo

func (m *TextSignatureOrder) GetMessage() github_com_cometbft_cometbft_libs_bytes.HexBytes {
	if m != nil {
		return m.Message
	}
	return nil
}

// EVMSignature defines a signature in the EVM format.
type EVMSignature struct {
	// r_address is the address of the nonce for using in the contract.
	RAddress github_com_cometbft_cometbft_libs_bytes.HexBytes `protobuf:"bytes,1,opt,name=r_address,json=rAddress,proto3,casttype=github.com/cometbft/cometbft/libs/bytes.HexBytes" json:"r_address,omitempty"`
	// signature is the signature part for using in the contract.
	Signature github_com_cometbft_cometbft_libs_bytes.HexBytes `protobuf:"bytes,2,opt,name=signature,proto3,casttype=github.com/cometbft/cometbft/libs/bytes.HexBytes"               json:"signature,omitempty"`
}

func (m *EVMSignature) Reset()         { *m = EVMSignature{} }
func (m *EVMSignature) String() string { return proto.CompactTextString(m) }
func (*EVMSignature) ProtoMessage()    {}
func (*EVMSignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{19}
}
func (m *EVMSignature) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EVMSignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EVMSignature.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EVMSignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EVMSignature.Merge(m, src)
}
func (m *EVMSignature) XXX_Size() int {
	return m.Size()
}
func (m *EVMSignature) XXX_DiscardUnknown() {
	xxx_messageInfo_EVMSignature.DiscardUnknown(m)
}

var xxx_messageInfo_EVMSignature proto.InternalMessageInfo

func (m *EVMSignature) GetRAddress() github_com_cometbft_cometbft_libs_bytes.HexBytes {
	if m != nil {
		return m.RAddress
	}
	return nil
}

func (m *EVMSignature) GetSignature() github_com_cometbft_cometbft_libs_bytes.HexBytes {
	if m != nil {
		return m.Signature
	}
	return nil
}

// SigningResult is a tss signing result from querying tss signing information.
type SigningResult struct {
	// signing is the tss signing result.
	Signing Signing `protobuf:"bytes,1,opt,name=signing,proto3"                                                    json:"signing"`
	// current_signing_attempt is the current attempt information of the signing.
	CurrentSigningAttempt *SigningAttempt `protobuf:"bytes,2,opt,name=current_signing_attempt,json=currentSigningAttempt,proto3"         json:"current_signing_attempt,omitempty"`
	// evm_signature is the signature in the format that can use directly in EVM.
	EVMSignature *EVMSignature `protobuf:"bytes,3,opt,name=evm_signature,json=evmSignature,proto3"                            json:"evm_signature,omitempty"`
	// received_partial_signatures is a list of received partial signatures.
	ReceivedPartialSignatures []PartialSignature `protobuf:"bytes,4,rep,name=received_partial_signatures,json=receivedPartialSignatures,proto3" json:"received_partial_signatures"`
}

func (m *SigningResult) Reset()         { *m = SigningResult{} }
func (m *SigningResult) String() string { return proto.CompactTextString(m) }
func (*SigningResult) ProtoMessage()    {}
func (*SigningResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{20}
}
func (m *SigningResult) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SigningResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SigningResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SigningResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SigningResult.Merge(m, src)
}
func (m *SigningResult) XXX_Size() int {
	return m.Size()
}
func (m *SigningResult) XXX_DiscardUnknown() {
	xxx_messageInfo_SigningResult.DiscardUnknown(m)
}

var xxx_messageInfo_SigningResult proto.InternalMessageInfo

func (m *SigningResult) GetSigning() Signing {
	if m != nil {
		return m.Signing
	}
	return Signing{}
}

func (m *SigningResult) GetCurrentSigningAttempt() *SigningAttempt {
	if m != nil {
		return m.CurrentSigningAttempt
	}
	return nil
}

func (m *SigningResult) GetEVMSignature() *EVMSignature {
	if m != nil {
		return m.EVMSignature
	}
	return nil
}

func (m *SigningResult) GetReceivedPartialSignatures() []PartialSignature {
	if m != nil {
		return m.ReceivedPartialSignatures
	}
	return nil
}

// SigningExpiration defines the expiration time of the signing.
type SigningExpiration struct {
	// signing_id is the id of the signing.
	SigningID SigningID `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID" json:"signing_id,omitempty"`
	// signing_attempt is the number of attempts of the signing.
	SigningAttempt uint64 `protobuf:"varint,2,opt,name=signing_attempt,json=signingAttempt,proto3"                                                                json:"signing_attempt,omitempty"`
}

func (m *SigningExpiration) Reset()         { *m = SigningExpiration{} }
func (m *SigningExpiration) String() string { return proto.CompactTextString(m) }
func (*SigningExpiration) ProtoMessage()    {}
func (*SigningExpiration) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{21}
}
func (m *SigningExpiration) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SigningExpiration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SigningExpiration.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SigningExpiration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SigningExpiration.Merge(m, src)
}
func (m *SigningExpiration) XXX_Size() int {
	return m.Size()
}
func (m *SigningExpiration) XXX_DiscardUnknown() {
	xxx_messageInfo_SigningExpiration.DiscardUnknown(m)
}

var xxx_messageInfo_SigningExpiration proto.InternalMessageInfo

func (m *SigningExpiration) GetSigningID() SigningID {
	if m != nil {
		return m.SigningID
	}
	return 0
}

func (m *SigningExpiration) GetSigningAttempt() uint64 {
	if m != nil {
		return m.SigningAttempt
	}
	return 0
}

// SigningExpirations is a list of signing expiration information that are waiting in the queue.
type SigningExpirations struct {
	// signing_expirations is a list of SigningExpiration object.
	SigningExpirations []SigningExpiration `protobuf:"bytes,1,rep,name=signing_expirations,json=signingExpirations,proto3" json:"signing_expirations"`
}

func (m *SigningExpirations) Reset()         { *m = SigningExpirations{} }
func (m *SigningExpirations) String() string { return proto.CompactTextString(m) }
func (*SigningExpirations) ProtoMessage()    {}
func (*SigningExpirations) Descriptor() ([]byte, []int) {
	return fileDescriptor_26231ff63bcc8f4b, []int{22}
}
func (m *SigningExpirations) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SigningExpirations) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SigningExpirations.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SigningExpirations) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SigningExpirations.Merge(m, src)
}
func (m *SigningExpirations) XXX_Size() int {
	return m.Size()
}
func (m *SigningExpirations) XXX_DiscardUnknown() {
	xxx_messageInfo_SigningExpirations.DiscardUnknown(m)
}

var xxx_messageInfo_SigningExpirations proto.InternalMessageInfo

func (m *SigningExpirations) GetSigningExpirations() []SigningExpiration {
	if m != nil {
		return m.SigningExpirations
	}
	return nil
}

func init() {
	proto.RegisterEnum("band.tss.v1beta1.SigningStatus", SigningStatus_name, SigningStatus_value)
	proto.RegisterEnum("band.tss.v1beta1.GroupStatus", GroupStatus_name, GroupStatus_value)
	proto.RegisterEnum("band.tss.v1beta1.ComplaintStatus", ComplaintStatus_name, ComplaintStatus_value)
	proto.RegisterType((*Group)(nil), "band.tss.v1beta1.Group")
	proto.RegisterType((*GroupResult)(nil), "band.tss.v1beta1.GroupResult")
	proto.RegisterType((*Round1Info)(nil), "band.tss.v1beta1.Round1Info")
	proto.RegisterType((*Round2Info)(nil), "band.tss.v1beta1.Round2Info")
	proto.RegisterType((*DE)(nil), "band.tss.v1beta1.DE")
	proto.RegisterType((*DEQueue)(nil), "band.tss.v1beta1.DEQueue")
	proto.RegisterType((*Signing)(nil), "band.tss.v1beta1.Signing")
	proto.RegisterType((*SigningAttempt)(nil), "band.tss.v1beta1.SigningAttempt")
	proto.RegisterType((*AssignedMember)(nil), "band.tss.v1beta1.AssignedMember")
	proto.RegisterType((*PendingSignings)(nil), "band.tss.v1beta1.PendingSignings")
	proto.RegisterType((*Member)(nil), "band.tss.v1beta1.Member")
	proto.RegisterType((*Confirm)(nil), "band.tss.v1beta1.Confirm")
	proto.RegisterType((*Complaint)(nil), "band.tss.v1beta1.Complaint")
	proto.RegisterType((*ComplaintWithStatus)(nil), "band.tss.v1beta1.ComplaintWithStatus")
	proto.RegisterType((*ComplaintsWithStatus)(nil), "band.tss.v1beta1.ComplaintsWithStatus")
	proto.RegisterType((*PendingProcessGroups)(nil), "band.tss.v1beta1.PendingProcessGroups")
	proto.RegisterType((*PendingProcessSignings)(nil), "band.tss.v1beta1.PendingProcessSignings")
	proto.RegisterType((*PartialSignature)(nil), "band.tss.v1beta1.PartialSignature")
	proto.RegisterType((*TextSignatureOrder)(nil), "band.tss.v1beta1.TextSignatureOrder")
	proto.RegisterType((*EVMSignature)(nil), "band.tss.v1beta1.EVMSignature")
	proto.RegisterType((*SigningResult)(nil), "band.tss.v1beta1.SigningResult")
	proto.RegisterType((*SigningExpiration)(nil), "band.tss.v1beta1.SigningExpiration")
	proto.RegisterType((*SigningExpirations)(nil), "band.tss.v1beta1.SigningExpirations")
}

func init() { proto.RegisterFile("band/tss/v1beta1/tss.proto", fileDescriptor_26231ff63bcc8f4b) }

var fileDescriptor_26231ff63bcc8f4b = []byte{
	// 2007 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x59, 0x4f, 0x6f, 0xdb, 0xc8,
	0x15, 0x37, 0x25, 0xd9, 0x92, 0x9e, 0x64, 0x5b, 0x3b, 0x71, 0x62, 0xd9, 0x71, 0x25, 0x87, 0xc5,
	0x76, 0x8d, 0x45, 0x6b, 0xc5, 0x72, 0xb7, 0xdd, 0x6e, 0x16, 0x48, 0x25, 0x4b, 0xf1, 0x6a, 0xe3,
	0xd8, 0x0e, 0x25, 0x27, 0x5b, 0x17, 0x5b, 0x82, 0x22, 0xc7, 0x12, 0x61, 0x91, 0x54, 0x39, 0x23,
	0xc7, 0xee, 0xa5, 0xd7, 0x1c, 0x73, 0xda, 0x73, 0x81, 0xf6, 0x50, 0xec, 0xb9, 0xfd, 0x0e, 0x41,
	0x7b, 0x59, 0xa0, 0x28, 0xb0, 0xbd, 0x28, 0x85, 0x82, 0xa2, 0xe8, 0xad, 0x67, 0x9f, 0x0a, 0x0e,
	0x87, 0xa4, 0x64, 0xc9, 0x59, 0x47, 0x4a, 0x82, 0xde, 0x38, 0xef, 0xbd, 0xf9, 0xcd, 0x9b, 0x37,
	0xef, 0xaf, 0x04, 0xcb, 0x75, 0xc5, 0xd4, 0x72, 0x94, 0x90, 0xdc, 0xc9, 0x46, 0x1d, 0x53, 0x65,
	0xc3, 0xf9, 0x5e, 0x6f, 0xdb, 0x16, 0xb5, 0x50, 0xca, 0xe1, 0xad, 0x3b, 0x6b, 0xce, 0x5b, 0x5e,
	0x68, 0x58, 0x0d, 0x8b, 0x31, 0x73, 0xce, 0x97, 0x2b, 0xb7, 0xbc, 0xa4, 0x5a, 0xc4, 0xb0, 0x88,
	0xec, 0x32, 0xdc, 0x05, 0x67, 0x65, 0x1b, 0x96, 0xd5, 0x68, 0xe1, 0x1c, 0x5b, 0xd5, 0x3b, 0x47,
	0x39, 0xaa, 0x1b, 0x98, 0x50, 0xc5, 0x68, 0xbb, 0x02, 0xe2, 0xdf, 0x42, 0x30, 0xbd, 0x6d, 0x5b,
	0x9d, 0x36, 0xfa, 0x1c, 0x42, 0xba, 0x96, 0x16, 0x56, 0x85, 0xb5, 0x48, 0xf1, 0x93, 0x5e, 0x37,
	0x1b, 0xaa, 0x94, 0xce, 0xbb, 0xd9, 0xdb, 0x0d, 0x9d, 0x36, 0x3b, 0xf5, 0x75, 0xd5, 0x32, 0x72,
	0x8e, 0x3a, 0x6c, 0x97, 0x6a, 0xb5, 0x72, 0x6a, 0x53, 0xd1, 0xcd, 0xdc, 0xc9, 0x66, 0xae, 0x7d,
	0xdc, 0x60, 0x3a, 0x33, 0x94, 0x4a, 0x49, 0x0a, 0xe9, 0x1a, 0x42, 0x10, 0x21, 0xfa, 0x6f, 0x70,
	0x3a, 0xe4, 0xa0, 0x49, 0xec, 0x1b, 0xad, 0x40, 0x9c, 0x36, 0x6d, 0x4c, 0x9a, 0x56, 0x4b, 0x4b,
	0x87, 0x19, 0x23, 0x20, 0xa0, 0xfb, 0x10, 0x6d, 0x77, 0xea, 0xf2, 0x31, 0x3e, 0x4b, 0x47, 0x56,
	0x85, 0xb5, 0x64, 0x31, 0x7f, 0xde, 0xcd, 0xae, 0x5f, 0xf9, 0xf0, 0x7d, 0x4b, 0x37, 0xa9, 0x34,
	0xd3, 0xee, 0xd4, 0xef, 0xe3, 0x33, 0xf4, 0x11, 0xcc, 0x10, 0xaa, 0xd0, 0x0e, 0x49, 0x4f, 0xaf,
	0x0a, 0x6b, 0x73, 0xf9, 0xef, 0xad, 0x5f, 0xb4, 0xa4, 0xab, 0x6d, 0x95, 0x09, 0x49, 0x5c, 0x18,
	0xbd, 0x0f, 0x73, 0xaa, 0x8d, 0x15, 0x8a, 0x35, 0xb9, 0x89, 0xf5, 0x46, 0x93, 0xa6, 0x67, 0x98,
	0x9a, 0xb3, 0x9c, 0xfa, 0x19, 0x23, 0xa2, 0x5b, 0x90, 0x34, 0x2c, 0xad, 0xd3, 0xc2, 0xb2, 0xf5,
	0xc4, 0xc4, 0x76, 0x3a, 0xba, 0x2a, 0xac, 0xc5, 0xa5, 0x84, 0x4b, 0xdb, 0x73, 0x48, 0xe2, 0xb3,
	0x08, 0x24, 0xd8, 0x09, 0x12, 0x26, 0x9d, 0x16, 0x45, 0x9b, 0x30, 0xdd, 0x70, 0x96, 0xcc, 0xbc,
	0x89, 0xfc, 0xe2, 0x25, 0xfa, 0x14, 0x23, 0xcf, 0xbb, 0xd9, 0x29, 0xc9, 0x95, 0x45, 0x2a, 0x24,
	0xb4, 0xe3, 0x86, 0xac, 0x5a, 0x26, 0xc5, 0xa7, 0x94, 0xd9, 0x32, 0x59, 0x2c, 0xf6, 0xba, 0x59,
	0x28, 0xdd, 0xdf, 0xde, 0x72, 0xa9, 0x17, 0x5e, 0x48, 0xb5, 0x0c, 0x4c, 0xeb, 0x47, 0x34, 0xf8,
	0x68, 0xe9, 0x75, 0x92, 0xab, 0x9f, 0x51, 0x4c, 0xd6, 0x3f, 0xc3, 0xa7, 0x45, 0xe7, 0x43, 0x02,
	0xed, 0xb8, 0xc1, 0xf7, 0xa3, 0x8f, 0x21, 0x6a, 0x60, 0xa3, 0x8e, 0x6d, 0x92, 0x0e, 0xaf, 0x86,
	0xd7, 0x12, 0xf9, 0xf4, 0xb0, 0x6e, 0x0f, 0x98, 0x00, 0x57, 0xce, 0x13, 0x47, 0x65, 0x48, 0xda,
	0x56, 0xc7, 0xd4, 0x36, 0x64, 0xdd, 0x3c, 0xb2, 0x48, 0x3a, 0xc2, 0xb6, 0xaf, 0x0c, 0x6f, 0x97,
	0x98, 0x54, 0xc5, 0x3c, 0xb2, 0x38, 0x44, 0xc2, 0xf6, 0x29, 0x01, 0x4c, 0x9e, 0xc3, 0x4c, 0xbf,
	0x12, 0x26, 0x3f, 0x04, 0x93, 0x77, 0x61, 0xea, 0x70, 0x43, 0xb5, 0x8c, 0x76, 0x4b, 0xd1, 0x4d,
	0x4a, 0xe4, 0x27, 0x3a, 0x6d, 0xca, 0xdc, 0x05, 0x66, 0x18, 0xe0, 0x0f, 0x86, 0x01, 0xb7, 0x7c,
	0xf9, 0xc7, 0x3a, 0x6d, 0xba, 0xbe, 0xc0, 0xa1, 0x17, 0xd4, 0x11, 0x3c, 0x74, 0x07, 0x62, 0xaa,
	0x65, 0x1e, 0xe9, 0xb6, 0x41, 0xd2, 0x51, 0x86, 0xba, 0x34, 0x0a, 0x95, 0x49, 0x70, 0x20, 0x7f,
	0x83, 0xf8, 0xdf, 0x30, 0x40, 0x60, 0x09, 0xf4, 0x2b, 0x88, 0xbb, 0x86, 0x94, 0xfd, 0xa0, 0x2b,
	0xf4, 0xba, 0xd9, 0x98, 0x6b, 0x6b, 0x16, 0x7a, 0x1b, 0x57, 0xf6, 0x7e, 0x6f, 0x93, 0x14, 0x73,
	0x31, 0x2b, 0x1a, 0xd2, 0xe0, 0x9a, 0x6a, 0xe1, 0xa3, 0x23, 0x5d, 0xd5, 0xb1, 0x49, 0x65, 0xd5,
	0x32, 0x0c, 0x9d, 0x92, 0x74, 0x68, 0x35, 0xbc, 0x96, 0x2c, 0x6e, 0x7e, 0xfd, 0x22, 0x9b, 0x7b,
	0xbd, 0xd8, 0x22, 0x12, 0xea, 0xc3, 0xdb, 0x72, 0xe1, 0xd0, 0x2f, 0x21, 0x65, 0x99, 0x58, 0x76,
	0x92, 0x8a, 0xec, 0x85, 0x6f, 0x78, 0xec, 0xf0, 0x9d, 0xb5, 0x4c, 0x5c, 0xd3, 0x0d, 0xbc, 0xef,
	0x46, 0xf1, 0x2f, 0x20, 0xa9, 0xdc, 0x96, 0x89, 0xde, 0x30, 0x15, 0xda, 0xb1, 0x31, 0xcf, 0x0b,
	0x3f, 0x39, 0xef, 0x66, 0xf3, 0x57, 0x06, 0xae, 0x7a, 0xbb, 0xa5, 0x84, 0x72, 0xdb, 0x5f, 0x20,
	0x0d, 0x90, 0xaf, 0x77, 0x70, 0xc0, 0xf4, 0x44, 0x07, 0xa4, 0xb8, 0xf6, 0x3e, 0x45, 0xec, 0x09,
	0xfc, 0xc9, 0xf3, 0xef, 0xe4, 0xc9, 0x29, 0x2c, 0x62, 0x53, 0xb5, 0xcf, 0xda, 0x4e, 0x02, 0x23,
	0x58, 0xb5, 0x31, 0x95, 0x49, 0x53, 0xb1, 0xb1, 0xf7, 0xec, 0x9f, 0x7e, 0xfd, 0x22, 0xfb, 0xf1,
	0x95, 0x4f, 0x28, 0x9b, 0x6a, 0x95, 0x81, 0x54, 0x19, 0x86, 0x74, 0xdd, 0x07, 0xef, 0x27, 0x8b,
	0x5f, 0x09, 0x10, 0x2a, 0x95, 0xd1, 0x36, 0x4c, 0x3b, 0x0e, 0xe0, 0x5e, 0x6c, 0xbc, 0xe7, 0x8f,
	0xb4, 0x3b, 0xf5, 0x92, 0x07, 0x84, 0x79, 0xbe, 0x1b, 0x17, 0xa8, 0x2c, 0x6e, 0x40, 0xb4, 0x54,
	0x7e, 0xd8, 0xc1, 0x1d, 0xec, 0x94, 0xa3, 0x26, 0x56, 0xb8, 0xd1, 0x25, 0xf6, 0xed, 0xd0, 0xa8,
	0xa2, 0xb7, 0xbc, 0x12, 0xe5, 0x7c, 0x8b, 0xff, 0x99, 0x86, 0xa8, 0xf3, 0x7c, 0xba, 0xd9, 0x40,
	0x3b, 0x7d, 0xe5, 0xf0, 0x53, 0xbf, 0x1c, 0xbe, 0x9e, 0x63, 0xe8, 0x66, 0x83, 0x17, 0xc4, 0x0f,
	0x60, 0x5e, 0xed, 0xd8, 0xb6, 0x13, 0x8a, 0x0a, 0xa5, 0xd8, 0x68, 0x53, 0x7e, 0xf0, 0x1c, 0x27,
	0x17, 0x5c, 0x2a, 0x3a, 0x84, 0x18, 0xcb, 0xfe, 0x8e, 0x8f, 0xb0, 0x22, 0x59, 0xbc, 0xdb, 0xeb,
	0x66, 0xa3, 0xbc, 0xb8, 0x8e, 0x55, 0x90, 0xa3, 0x0c, 0xb0, 0xa2, 0xa1, 0x47, 0x30, 0xeb, 0x62,
	0x4f, 0x5e, 0x69, 0x13, 0x0c, 0x88, 0x07, 0xea, 0xae, 0x53, 0x43, 0x08, 0x51, 0x1a, 0x5e, 0x08,
	0xfd, 0x78, 0xac, 0xb2, 0xe4, 0x81, 0xa0, 0x43, 0x98, 0x0f, 0xf4, 0x34, 0x2d, 0x53, 0xc5, 0xac,
	0x10, 0x8f, 0x99, 0x54, 0x3c, 0x4d, 0x77, 0x1d, 0x20, 0x54, 0x83, 0x78, 0x10, 0xf0, 0xd1, 0x89,
	0x02, 0x3e, 0x00, 0x42, 0x3f, 0xf5, 0x1b, 0x8e, 0x18, 0x6b, 0x38, 0xb2, 0xc3, 0x75, 0x81, 0xfb,
	0xc3, 0x77, 0xb6, 0x1c, 0xf1, 0x51, 0x2d, 0xc7, 0x43, 0x78, 0xcf, 0x13, 0xf3, 0x1b, 0xb8, 0x34,
	0xb0, 0x5e, 0x62, 0x79, 0xdd, 0x6d, 0xf1, 0xd6, 0xbd, 0x16, 0x6f, 0xbd, 0xe6, 0x49, 0x14, 0x63,
	0x4e, 0x0d, 0x7a, 0xf6, 0x22, 0x2b, 0x48, 0x29, 0xbe, 0xdd, 0xe7, 0x89, 0x4f, 0x43, 0x30, 0xc7,
	0x75, 0xf2, 0x7c, 0xaf, 0x0e, 0x40, 0x5c, 0x4a, 0x90, 0xa1, 0xb6, 0x7a, 0xdd, 0x6c, 0xdc, 0xf7,
	0xe5, 0x31, 0x23, 0x20, 0xce, 0x61, 0x2b, 0x1a, 0x4a, 0x43, 0x74, 0x30, 0x00, 0xbc, 0xa5, 0x63,
	0x0a, 0x7c, 0xda, 0xd6, 0xed, 0xc0, 0x14, 0x6e, 0x93, 0x38, 0xcb, 0xa9, 0xbe, 0x29, 0x52, 0x0a,
	0x71, 0xf0, 0xb0, 0x26, 0x7b, 0x9d, 0x8b, 0xdb, 0x7a, 0xac, 0x0e, 0x1b, 0xbd, 0xc0, 0x25, 0x07,
	0x3a, 0x98, 0x79, 0x65, 0x80, 0x4a, 0xc4, 0xaf, 0x22, 0x30, 0x37, 0x28, 0xf9, 0xd6, 0x73, 0xb5,
	0x63, 0x06, 0x4d, 0xb3, 0x31, 0x21, 0xcc, 0x0c, 0x71, 0xc9, 0x5b, 0xf6, 0x37, 0xc2, 0xe1, 0x89,
	0x1b, 0x61, 0x3f, 0x2b, 0x47, 0xde, 0x54, 0x56, 0x9e, 0x9e, 0x2c, 0x2b, 0xa3, 0x43, 0x98, 0xab,
	0xeb, 0xa6, 0xe6, 0xf8, 0xd8, 0x91, 0xa2, 0x52, 0xcb, 0xe6, 0xa1, 0xbd, 0x79, 0xde, 0x7d, 0x8d,
	0x96, 0xa4, 0xaa, 0x2a, 0x2d, 0xc5, 0x96, 0x66, 0x39, 0xd4, 0x3d, 0x86, 0x84, 0xf6, 0x20, 0x1e,
	0x64, 0x8c, 0xe8, 0xd8, 0x8a, 0xc6, 0xda, 0x3c, 0x59, 0x88, 0x79, 0x98, 0xdf, 0xc7, 0xec, 0x04,
	0xee, 0xcb, 0x04, 0x65, 0x21, 0x11, 0xc4, 0x08, 0x49, 0x0b, 0xab, 0xe1, 0xb5, 0x88, 0x04, 0xbe,
	0x7f, 0x13, 0xf1, 0xdb, 0x10, 0xcc, 0x70, 0x27, 0xba, 0xdf, 0x57, 0x42, 0xee, 0xf8, 0x25, 0x64,
	0x0c, 0xbf, 0x71, 0x2a, 0x48, 0x7f, 0x61, 0x08, 0xbd, 0xe1, 0xc2, 0xd0, 0xe7, 0x8d, 0xe1, 0x4b,
	0xbd, 0x71, 0xf2, 0xb1, 0xec, 0x16, 0x24, 0x75, 0x22, 0x1b, 0x4a, 0x4b, 0x57, 0x75, 0x8b, 0x0f,
	0x67, 0x31, 0x29, 0xa1, 0x93, 0x07, 0x1e, 0x09, 0xdd, 0x84, 0xb8, 0x4e, 0x64, 0x45, 0xa5, 0xfa,
	0x89, 0x9b, 0xf4, 0x63, 0x52, 0x4c, 0x27, 0x05, 0xb6, 0x16, 0x9f, 0x0b, 0x10, 0xe5, 0xed, 0xf5,
	0x5b, 0x0f, 0xd0, 0x2f, 0x61, 0xde, 0x7a, 0x62, 0x7a, 0x95, 0xd2, 0x69, 0x12, 0x79, 0x43, 0x32,
	0x6e, 0xb5, 0x48, 0x5a, 0x4f, 0x4c, 0xb7, 0x5e, 0x56, 0xf5, 0x86, 0xf8, 0xef, 0x10, 0xc4, 0xfd,
	0xf9, 0x03, 0x3d, 0x86, 0x84, 0x37, 0x70, 0x28, 0x26, 0xe5, 0xd7, 0xf9, 0x68, 0xbc, 0x2b, 0xf4,
	0x23, 0xa1, 0x03, 0x00, 0x1b, 0x93, 0xb6, 0x65, 0x6a, 0xd8, 0xe4, 0x09, 0x77, 0x5c, 0xdc, 0x3e,
	0x20, 0xc7, 0x2b, 0x98, 0x51, 0xce, 0x8c, 0x49, 0x72, 0xd4, 0x31, 0x3e, 0xab, 0x9e, 0x19, 0xe8,
	0xcb, 0xfe, 0x8a, 0xec, 0x3a, 0xd9, 0xdd, 0xf3, 0x6e, 0xf6, 0xce, 0x95, 0xe1, 0x7c, 0x3b, 0x8e,
	0x2a, 0xcd, 0xe2, 0x1f, 0x04, 0xb8, 0xe6, 0x4b, 0xf4, 0x0d, 0x73, 0x77, 0x21, 0xee, 0x0f, 0x79,
	0x7c, 0x2c, 0xbf, 0xf9, 0x8a, 0x19, 0x91, 0xd7, 0x8e, 0x60, 0x0f, 0xda, 0x81, 0x94, 0xbf, 0xf0,
	0x66, 0xcd, 0x10, 0xab, 0xfe, 0xb7, 0x5e, 0x81, 0xc3, 0xeb, 0xff, 0xbc, 0x3a, 0x48, 0x10, 0xff,
	0x21, 0xc0, 0xc2, 0xa8, 0x81, 0xf4, 0xad, 0x3b, 0xba, 0x72, 0xe9, 0xe0, 0x1c, 0x62, 0x55, 0xf5,
	0xfd, 0x57, 0x5c, 0xe6, 0x6a, 0x73, 0xb3, 0xd8, 0x81, 0x05, 0x9e, 0x46, 0xf7, 0x6d, 0x4b, 0xc5,
	0x84, 0xb0, 0x04, 0x44, 0x9c, 0x97, 0xf7, 0x52, 0x1a, 0xcf, 0xa4, 0xc5, 0x9f, 0x3b, 0x57, 0xe3,
	0xf9, 0x89, 0x8c, 0x95, 0xd4, 0x62, 0x3c, 0xa9, 0x11, 0xf1, 0xb7, 0x70, 0x63, 0xf0, 0x58, 0x3f,
	0x89, 0xe3, 0x11, 0x49, 0xbc, 0x58, 0xea, 0x75, 0xb3, 0xe0, 0xf7, 0x2c, 0x64, 0xcc, 0x56, 0xa7,
	0xbf, 0x14, 0xfc, 0x3d, 0x04, 0xa9, 0x7d, 0xc5, 0xa6, 0xba, 0xd2, 0x0a, 0x46, 0xcf, 0x77, 0xd1,
	0x64, 0x7d, 0x00, 0xf3, 0xde, 0x19, 0x17, 0xa6, 0x0d, 0x32, 0xd8, 0xf1, 0x0d, 0x38, 0x57, 0xf8,
	0xcd, 0x3b, 0x57, 0x6d, 0x38, 0xb6, 0x27, 0xef, 0xb6, 0xc5, 0x5f, 0x03, 0xaa, 0xe1, 0xd3, 0x20,
	0xdc, 0xf7, 0x6c, 0x0d, 0xdb, 0xfd, 0x53, 0x88, 0xf0, 0x06, 0xa6, 0x90, 0x4f, 0x12, 0x7f, 0xf9,
	0xd3, 0x8f, 0xa2, 0xec, 0x67, 0x32, 0x93, 0x8a, 0x7f, 0x16, 0x20, 0x59, 0x7e, 0xf4, 0x20, 0x78,
	0xc6, 0x87, 0x10, 0xb7, 0x65, 0xaf, 0x68, 0x4e, 0x72, 0x5e, 0xcc, 0x2e, 0xf0, 0x5a, 0x2b, 0xf5,
	0x1b, 0x2b, 0x34, 0x01, 0x64, 0x9f, 0xa9, 0xfe, 0x15, 0x82, 0x59, 0xee, 0x22, 0xfc, 0xa7, 0xc8,
	0x9f, 0x41, 0x94, 0x3b, 0x01, 0xcf, 0x7a, 0x4b, 0x97, 0xce, 0x2a, 0xde, 0x2f, 0x7e, 0x5c, 0x1e,
	0x7d, 0x01, 0x8b, 0xde, 0x10, 0x3b, 0xca, 0xbd, 0x46, 0x76, 0xe0, 0x83, 0x23, 0x86, 0x74, 0x9d,
	0x03, 0x5c, 0x98, 0x3c, 0x0e, 0x60, 0x16, 0x9f, 0x18, 0x7d, 0x3f, 0xc5, 0x84, 0x19, 0x5e, 0x66,
	0x18, 0xaf, 0xff, 0x11, 0x8a, 0xa9, 0x5e, 0x37, 0x3b, 0xf0, 0x2c, 0x52, 0x12, 0x9f, 0x18, 0xc1,
	0x23, 0x35, 0xe1, 0xa6, 0x8d, 0x55, 0xac, 0x9f, 0x60, 0x4d, 0x6e, 0xbb, 0x81, 0x18, 0x9c, 0xe1,
	0x8d, 0x0d, 0xe2, 0xf0, 0x21, 0x17, 0x83, 0x96, 0x1b, 0x62, 0xc9, 0x03, 0xbb, 0xc8, 0x27, 0xe2,
	0xef, 0x04, 0x78, 0x8f, 0xdf, 0xa9, 0xec, 0x8c, 0x2b, 0x0a, 0xd5, 0x2d, 0xf3, 0xff, 0x2a, 0xd6,
	0xc5, 0x36, 0xa0, 0x21, 0x0d, 0x09, 0x3a, 0x84, 0x6b, 0xde, 0x76, 0x1c, 0x90, 0x59, 0x4a, 0x4c,
	0xe4, 0xbf, 0x7f, 0xe9, 0x7b, 0x06, 0x10, 0xdc, 0x36, 0x88, 0x0c, 0x61, 0x7f, 0xf8, 0x54, 0xf0,
	0x9d, 0x8f, 0x17, 0xb3, 0x0c, 0x2c, 0x57, 0x2b, 0xdb, 0xbb, 0x95, 0xdd, 0x6d, 0xb9, 0x5a, 0x2b,
	0xd4, 0x0e, 0xaa, 0xf2, 0xc1, 0x6e, 0x75, 0xbf, 0xbc, 0x55, 0xb9, 0x57, 0x29, 0x97, 0x52, 0x53,
	0x68, 0x19, 0x6e, 0x5c, 0xe0, 0x3f, 0x2e, 0x54, 0x6a, 0x95, 0xdd, 0xed, 0x94, 0x30, 0x82, 0x57,
	0x3d, 0xd8, 0xda, 0x2a, 0x57, 0xab, 0xa9, 0x10, 0x5a, 0x82, 0xeb, 0x17, 0x78, 0xf7, 0x0a, 0x3b,
	0x3b, 0xe5, 0xdd, 0x54, 0x78, 0x39, 0xf2, 0xf4, 0xf7, 0x99, 0xa9, 0x0f, 0xff, 0x2a, 0xf0, 0x1f,
	0xe4, 0xb9, 0x22, 0x2b, 0x90, 0xde, 0x96, 0xf6, 0x0e, 0xf6, 0x47, 0xab, 0x91, 0x86, 0x85, 0x01,
	0xae, 0xb4, 0x77, 0xb0, 0x5b, 0x92, 0x37, 0x52, 0xc2, 0x25, 0x9c, 0x7c, 0x2a, 0x74, 0x09, 0x67,
	0x33, 0x15, 0x46, 0x8b, 0x70, 0x6d, 0x80, 0x53, 0xd8, 0xaa, 0x55, 0x1e, 0x95, 0x53, 0x91, 0xa1,
	0x2d, 0xe5, 0x2f, 0xf6, 0x2b, 0x52, 0xb9, 0x94, 0x9a, 0x1e, 0xda, 0xc2, 0x6f, 0x33, 0xc3, 0x6f,
	0x43, 0x61, 0xfe, 0x42, 0x43, 0x81, 0x56, 0x61, 0x65, 0x6b, 0xef, 0xc1, 0xfe, 0x4e, 0xa1, 0xb2,
	0x5b, 0x1b, 0x7d, 0xa9, 0x15, 0x48, 0x0f, 0x49, 0x78, 0x16, 0x14, 0xd0, 0x4d, 0x58, 0x1c, 0xe2,
	0xde, 0x2b, 0x54, 0x76, 0xca, 0xa5, 0x54, 0xc8, 0x3d, 0xb5, 0xf8, 0xf9, 0x1f, 0x7b, 0x19, 0xe1,
	0x79, 0x2f, 0x23, 0x7c, 0xd3, 0xcb, 0x08, 0xff, 0xec, 0x65, 0x84, 0x67, 0x2f, 0x33, 0x53, 0xdf,
	0xbc, 0xcc, 0x4c, 0x7d, 0xfb, 0x32, 0x33, 0x75, 0xf8, 0xc3, 0xef, 0x74, 0xe3, 0x53, 0xf6, 0x37,
	0x17, 0x3d, 0x6b, 0x63, 0x52, 0x9f, 0x61, 0xec, 0xcd, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0x55,
	0x02, 0x06, 0x92, 0xff, 0x1a, 0x00, 0x00,
}

func (this *Group) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Group)
	if !ok {
		that2, ok := that.(Group)
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
	if this.Size_ != that1.Size_ {
		return false
	}
	if this.Threshold != that1.Threshold {
		return false
	}
	if !bytes.Equal(this.PubKey, that1.PubKey) {
		return false
	}
	if this.Status != that1.Status {
		return false
	}
	if this.CreatedHeight != that1.CreatedHeight {
		return false
	}
	if this.ModuleOwner != that1.ModuleOwner {
		return false
	}
	return true
}
func (this *GroupResult) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*GroupResult)
	if !ok {
		that2, ok := that.(GroupResult)
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
	if !this.Group.Equal(&that1.Group) {
		return false
	}
	if !bytes.Equal(this.DKGContext, that1.DKGContext) {
		return false
	}
	if len(this.Members) != len(that1.Members) {
		return false
	}
	for i := range this.Members {
		if !this.Members[i].Equal(&that1.Members[i]) {
			return false
		}
	}
	if len(this.Round1Infos) != len(that1.Round1Infos) {
		return false
	}
	for i := range this.Round1Infos {
		if !this.Round1Infos[i].Equal(&that1.Round1Infos[i]) {
			return false
		}
	}
	if len(this.Round2Infos) != len(that1.Round2Infos) {
		return false
	}
	for i := range this.Round2Infos {
		if !this.Round2Infos[i].Equal(&that1.Round2Infos[i]) {
			return false
		}
	}
	if len(this.ComplaintsWithStatus) != len(that1.ComplaintsWithStatus) {
		return false
	}
	for i := range this.ComplaintsWithStatus {
		if !this.ComplaintsWithStatus[i].Equal(&that1.ComplaintsWithStatus[i]) {
			return false
		}
	}
	if len(this.Confirms) != len(that1.Confirms) {
		return false
	}
	for i := range this.Confirms {
		if !this.Confirms[i].Equal(&that1.Confirms[i]) {
			return false
		}
	}
	return true
}
func (this *Round1Info) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Round1Info)
	if !ok {
		that2, ok := that.(Round1Info)
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
	if this.MemberID != that1.MemberID {
		return false
	}
	if len(this.CoefficientCommits) != len(that1.CoefficientCommits) {
		return false
	}
	for i := range this.CoefficientCommits {
		if !bytes.Equal(this.CoefficientCommits[i], that1.CoefficientCommits[i]) {
			return false
		}
	}
	if !bytes.Equal(this.OneTimePubKey, that1.OneTimePubKey) {
		return false
	}
	if !bytes.Equal(this.A0Signature, that1.A0Signature) {
		return false
	}
	if !bytes.Equal(this.OneTimeSignature, that1.OneTimeSignature) {
		return false
	}
	return true
}
func (this *Round2Info) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Round2Info)
	if !ok {
		that2, ok := that.(Round2Info)
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
	if this.MemberID != that1.MemberID {
		return false
	}
	if len(this.EncryptedSecretShares) != len(that1.EncryptedSecretShares) {
		return false
	}
	for i := range this.EncryptedSecretShares {
		if !bytes.Equal(this.EncryptedSecretShares[i], that1.EncryptedSecretShares[i]) {
			return false
		}
	}
	return true
}
func (this *DE) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DE)
	if !ok {
		that2, ok := that.(DE)
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
	if !bytes.Equal(this.PubD, that1.PubD) {
		return false
	}
	if !bytes.Equal(this.PubE, that1.PubE) {
		return false
	}
	return true
}
func (this *DEQueue) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DEQueue)
	if !ok {
		that2, ok := that.(DEQueue)
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
	if this.Head != that1.Head {
		return false
	}
	if this.Tail != that1.Tail {
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
	if this.CurrentAttempt != that1.CurrentAttempt {
		return false
	}
	if this.GroupID != that1.GroupID {
		return false
	}
	if !bytes.Equal(this.GroupPubKey, that1.GroupPubKey) {
		return false
	}
	if !bytes.Equal(this.Message, that1.Message) {
		return false
	}
	if !bytes.Equal(this.GroupPubNonce, that1.GroupPubNonce) {
		return false
	}
	if !bytes.Equal(this.Signature, that1.Signature) {
		return false
	}
	if this.Status != that1.Status {
		return false
	}
	if this.CreatedHeight != that1.CreatedHeight {
		return false
	}
	if !this.CreatedTimestamp.Equal(that1.CreatedTimestamp) {
		return false
	}
	return true
}
func (this *SigningAttempt) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SigningAttempt)
	if !ok {
		that2, ok := that.(SigningAttempt)
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
	if this.Attempt != that1.Attempt {
		return false
	}
	if this.ExpiredHeight != that1.ExpiredHeight {
		return false
	}
	if len(this.AssignedMembers) != len(that1.AssignedMembers) {
		return false
	}
	for i := range this.AssignedMembers {
		if !this.AssignedMembers[i].Equal(&that1.AssignedMembers[i]) {
			return false
		}
	}
	return true
}
func (this *AssignedMember) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*AssignedMember)
	if !ok {
		that2, ok := that.(AssignedMember)
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
	if this.MemberID != that1.MemberID {
		return false
	}
	if this.Address != that1.Address {
		return false
	}
	if !bytes.Equal(this.PubKey, that1.PubKey) {
		return false
	}
	if !bytes.Equal(this.PubD, that1.PubD) {
		return false
	}
	if !bytes.Equal(this.PubE, that1.PubE) {
		return false
	}
	if !bytes.Equal(this.BindingFactor, that1.BindingFactor) {
		return false
	}
	if !bytes.Equal(this.PubNonce, that1.PubNonce) {
		return false
	}
	return true
}
func (this *PendingSignings) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PendingSignings)
	if !ok {
		that2, ok := that.(PendingSignings)
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
	if len(this.SigningIds) != len(that1.SigningIds) {
		return false
	}
	for i := range this.SigningIds {
		if this.SigningIds[i] != that1.SigningIds[i] {
			return false
		}
	}
	return true
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
	if this.ID != that1.ID {
		return false
	}
	if this.GroupID != that1.GroupID {
		return false
	}
	if this.Address != that1.Address {
		return false
	}
	if !bytes.Equal(this.PubKey, that1.PubKey) {
		return false
	}
	if this.IsMalicious != that1.IsMalicious {
		return false
	}
	if this.IsActive != that1.IsActive {
		return false
	}
	return true
}
func (this *Confirm) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Confirm)
	if !ok {
		that2, ok := that.(Confirm)
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
	if this.MemberID != that1.MemberID {
		return false
	}
	if !bytes.Equal(this.OwnPubKeySig, that1.OwnPubKeySig) {
		return false
	}
	return true
}
func (this *Complaint) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Complaint)
	if !ok {
		that2, ok := that.(Complaint)
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
	if this.Complainant != that1.Complainant {
		return false
	}
	if this.Respondent != that1.Respondent {
		return false
	}
	if !bytes.Equal(this.KeySym, that1.KeySym) {
		return false
	}
	if !bytes.Equal(this.Signature, that1.Signature) {
		return false
	}
	return true
}
func (this *ComplaintWithStatus) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ComplaintWithStatus)
	if !ok {
		that2, ok := that.(ComplaintWithStatus)
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
	if !this.Complaint.Equal(&that1.Complaint) {
		return false
	}
	if this.ComplaintStatus != that1.ComplaintStatus {
		return false
	}
	return true
}
func (this *ComplaintsWithStatus) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ComplaintsWithStatus)
	if !ok {
		that2, ok := that.(ComplaintsWithStatus)
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
	if this.MemberID != that1.MemberID {
		return false
	}
	if len(this.ComplaintsWithStatus) != len(that1.ComplaintsWithStatus) {
		return false
	}
	for i := range this.ComplaintsWithStatus {
		if !this.ComplaintsWithStatus[i].Equal(&that1.ComplaintsWithStatus[i]) {
			return false
		}
	}
	return true
}
func (this *PendingProcessGroups) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PendingProcessGroups)
	if !ok {
		that2, ok := that.(PendingProcessGroups)
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
	if len(this.GroupIDs) != len(that1.GroupIDs) {
		return false
	}
	for i := range this.GroupIDs {
		if this.GroupIDs[i] != that1.GroupIDs[i] {
			return false
		}
	}
	return true
}
func (this *PendingProcessSignings) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PendingProcessSignings)
	if !ok {
		that2, ok := that.(PendingProcessSignings)
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
	if len(this.SigningIDs) != len(that1.SigningIDs) {
		return false
	}
	for i := range this.SigningIDs {
		if this.SigningIDs[i] != that1.SigningIDs[i] {
			return false
		}
	}
	return true
}
func (this *PartialSignature) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*PartialSignature)
	if !ok {
		that2, ok := that.(PartialSignature)
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
	if this.SigningAttempt != that1.SigningAttempt {
		return false
	}
	if this.MemberID != that1.MemberID {
		return false
	}
	if !bytes.Equal(this.Signature, that1.Signature) {
		return false
	}
	return true
}
func (this *TextSignatureOrder) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*TextSignatureOrder)
	if !ok {
		that2, ok := that.(TextSignatureOrder)
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
	if !bytes.Equal(this.Message, that1.Message) {
		return false
	}
	return true
}
func (this *EVMSignature) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*EVMSignature)
	if !ok {
		that2, ok := that.(EVMSignature)
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
	if !bytes.Equal(this.RAddress, that1.RAddress) {
		return false
	}
	if !bytes.Equal(this.Signature, that1.Signature) {
		return false
	}
	return true
}
func (this *SigningResult) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SigningResult)
	if !ok {
		that2, ok := that.(SigningResult)
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
	if !this.Signing.Equal(&that1.Signing) {
		return false
	}
	if !this.CurrentSigningAttempt.Equal(that1.CurrentSigningAttempt) {
		return false
	}
	if !this.EVMSignature.Equal(that1.EVMSignature) {
		return false
	}
	if len(this.ReceivedPartialSignatures) != len(that1.ReceivedPartialSignatures) {
		return false
	}
	for i := range this.ReceivedPartialSignatures {
		if !this.ReceivedPartialSignatures[i].Equal(&that1.ReceivedPartialSignatures[i]) {
			return false
		}
	}
	return true
}
func (this *SigningExpiration) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SigningExpiration)
	if !ok {
		that2, ok := that.(SigningExpiration)
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
	if this.SigningAttempt != that1.SigningAttempt {
		return false
	}
	return true
}
func (this *SigningExpirations) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SigningExpirations)
	if !ok {
		that2, ok := that.(SigningExpirations)
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
	if len(this.SigningExpirations) != len(that1.SigningExpirations) {
		return false
	}
	for i := range this.SigningExpirations {
		if !this.SigningExpirations[i].Equal(&that1.SigningExpirations[i]) {
			return false
		}
	}
	return true
}
func (m *Group) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Group) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Group) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ModuleOwner) > 0 {
		i -= len(m.ModuleOwner)
		copy(dAtA[i:], m.ModuleOwner)
		i = encodeVarintTss(dAtA, i, uint64(len(m.ModuleOwner)))
		i--
		dAtA[i] = 0x3a
	}
	if m.CreatedHeight != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.CreatedHeight))
		i--
		dAtA[i] = 0x30
	}
	if m.Status != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x28
	}
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintTss(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0x22
	}
	if m.Threshold != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x18
	}
	if m.Size_ != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Size_))
		i--
		dAtA[i] = 0x10
	}
	if m.ID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GroupResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GroupResult) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GroupResult) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Confirms) > 0 {
		for iNdEx := len(m.Confirms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Confirms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.ComplaintsWithStatus) > 0 {
		for iNdEx := len(m.ComplaintsWithStatus) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ComplaintsWithStatus[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Round2Infos) > 0 {
		for iNdEx := len(m.Round2Infos) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Round2Infos[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Round1Infos) > 0 {
		for iNdEx := len(m.Round1Infos) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Round1Infos[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Members) > 0 {
		for iNdEx := len(m.Members) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Members[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.DKGContext) > 0 {
		i -= len(m.DKGContext)
		copy(dAtA[i:], m.DKGContext)
		i = encodeVarintTss(dAtA, i, uint64(len(m.DKGContext)))
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Group.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTss(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Round1Info) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Round1Info) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Round1Info) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OneTimeSignature) > 0 {
		i -= len(m.OneTimeSignature)
		copy(dAtA[i:], m.OneTimeSignature)
		i = encodeVarintTss(dAtA, i, uint64(len(m.OneTimeSignature)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.A0Signature) > 0 {
		i -= len(m.A0Signature)
		copy(dAtA[i:], m.A0Signature)
		i = encodeVarintTss(dAtA, i, uint64(len(m.A0Signature)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.OneTimePubKey) > 0 {
		i -= len(m.OneTimePubKey)
		copy(dAtA[i:], m.OneTimePubKey)
		i = encodeVarintTss(dAtA, i, uint64(len(m.OneTimePubKey)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.CoefficientCommits) > 0 {
		for iNdEx := len(m.CoefficientCommits) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.CoefficientCommits[iNdEx])
			copy(dAtA[i:], m.CoefficientCommits[iNdEx])
			i = encodeVarintTss(dAtA, i, uint64(len(m.CoefficientCommits[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.MemberID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.MemberID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Round2Info) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Round2Info) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Round2Info) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EncryptedSecretShares) > 0 {
		for iNdEx := len(m.EncryptedSecretShares) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.EncryptedSecretShares[iNdEx])
			copy(dAtA[i:], m.EncryptedSecretShares[iNdEx])
			i = encodeVarintTss(dAtA, i, uint64(len(m.EncryptedSecretShares[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.MemberID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.MemberID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *DE) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DE) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DE) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PubE) > 0 {
		i -= len(m.PubE)
		copy(dAtA[i:], m.PubE)
		i = encodeVarintTss(dAtA, i, uint64(len(m.PubE)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.PubD) > 0 {
		i -= len(m.PubD)
		copy(dAtA[i:], m.PubD)
		i = encodeVarintTss(dAtA, i, uint64(len(m.PubD)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DEQueue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DEQueue) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DEQueue) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Tail != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Tail))
		i--
		dAtA[i] = 0x10
	}
	if m.Head != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Head))
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
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(
		m.CreatedTimestamp,
		dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.CreatedTimestamp):],
	)
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintTss(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x52
	if m.CreatedHeight != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.CreatedHeight))
		i--
		dAtA[i] = 0x48
	}
	if m.Status != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x40
	}
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintTss(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.GroupPubNonce) > 0 {
		i -= len(m.GroupPubNonce)
		copy(dAtA[i:], m.GroupPubNonce)
		i = encodeVarintTss(dAtA, i, uint64(len(m.GroupPubNonce)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintTss(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.GroupPubKey) > 0 {
		i -= len(m.GroupPubKey)
		copy(dAtA[i:], m.GroupPubKey)
		i = encodeVarintTss(dAtA, i, uint64(len(m.GroupPubKey)))
		i--
		dAtA[i] = 0x22
	}
	if m.GroupID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.GroupID))
		i--
		dAtA[i] = 0x18
	}
	if m.CurrentAttempt != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.CurrentAttempt))
		i--
		dAtA[i] = 0x10
	}
	if m.ID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SigningAttempt) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SigningAttempt) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SigningAttempt) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AssignedMembers) > 0 {
		for iNdEx := len(m.AssignedMembers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AssignedMembers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.ExpiredHeight != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.ExpiredHeight))
		i--
		dAtA[i] = 0x18
	}
	if m.Attempt != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Attempt))
		i--
		dAtA[i] = 0x10
	}
	if m.SigningID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.SigningID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AssignedMember) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AssignedMember) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AssignedMember) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PubNonce) > 0 {
		i -= len(m.PubNonce)
		copy(dAtA[i:], m.PubNonce)
		i = encodeVarintTss(dAtA, i, uint64(len(m.PubNonce)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.BindingFactor) > 0 {
		i -= len(m.BindingFactor)
		copy(dAtA[i:], m.BindingFactor)
		i = encodeVarintTss(dAtA, i, uint64(len(m.BindingFactor)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.PubE) > 0 {
		i -= len(m.PubE)
		copy(dAtA[i:], m.PubE)
		i = encodeVarintTss(dAtA, i, uint64(len(m.PubE)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.PubD) > 0 {
		i -= len(m.PubD)
		copy(dAtA[i:], m.PubD)
		i = encodeVarintTss(dAtA, i, uint64(len(m.PubD)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintTss(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTss(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	if m.MemberID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.MemberID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PendingSignings) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PendingSignings) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PendingSignings) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SigningIds) > 0 {
		dAtA4 := make([]byte, len(m.SigningIds)*10)
		var j3 int
		for _, num := range m.SigningIds {
			for num >= 1<<7 {
				dAtA4[j3] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j3++
			}
			dAtA4[j3] = uint8(num)
			j3++
		}
		i -= j3
		copy(dAtA[i:], dAtA4[:j3])
		i = encodeVarintTss(dAtA, i, uint64(j3))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
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
	if m.IsActive {
		i--
		if m.IsActive {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if m.IsMalicious {
		i--
		if m.IsMalicious {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintTss(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTss(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x1a
	}
	if m.GroupID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.GroupID))
		i--
		dAtA[i] = 0x10
	}
	if m.ID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Confirm) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Confirm) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Confirm) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OwnPubKeySig) > 0 {
		i -= len(m.OwnPubKeySig)
		copy(dAtA[i:], m.OwnPubKeySig)
		i = encodeVarintTss(dAtA, i, uint64(len(m.OwnPubKeySig)))
		i--
		dAtA[i] = 0x12
	}
	if m.MemberID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.MemberID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Complaint) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Complaint) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Complaint) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintTss(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.KeySym) > 0 {
		i -= len(m.KeySym)
		copy(dAtA[i:], m.KeySym)
		i = encodeVarintTss(dAtA, i, uint64(len(m.KeySym)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Respondent != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Respondent))
		i--
		dAtA[i] = 0x10
	}
	if m.Complainant != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.Complainant))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ComplaintWithStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ComplaintWithStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ComplaintWithStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ComplaintStatus != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.ComplaintStatus))
		i--
		dAtA[i] = 0x10
	}
	{
		size, err := m.Complaint.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTss(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *ComplaintsWithStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ComplaintsWithStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ComplaintsWithStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ComplaintsWithStatus) > 0 {
		for iNdEx := len(m.ComplaintsWithStatus) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ComplaintsWithStatus[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.MemberID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.MemberID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PendingProcessGroups) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PendingProcessGroups) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PendingProcessGroups) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.GroupIDs) > 0 {
		dAtA7 := make([]byte, len(m.GroupIDs)*10)
		var j6 int
		for _, num := range m.GroupIDs {
			for num >= 1<<7 {
				dAtA7[j6] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j6++
			}
			dAtA7[j6] = uint8(num)
			j6++
		}
		i -= j6
		copy(dAtA[i:], dAtA7[:j6])
		i = encodeVarintTss(dAtA, i, uint64(j6))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PendingProcessSignings) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PendingProcessSignings) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PendingProcessSignings) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SigningIDs) > 0 {
		dAtA9 := make([]byte, len(m.SigningIDs)*10)
		var j8 int
		for _, num := range m.SigningIDs {
			for num >= 1<<7 {
				dAtA9[j8] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j8++
			}
			dAtA9[j8] = uint8(num)
			j8++
		}
		i -= j8
		copy(dAtA[i:], dAtA9[:j8])
		i = encodeVarintTss(dAtA, i, uint64(j8))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PartialSignature) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PartialSignature) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PartialSignature) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintTss(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x22
	}
	if m.MemberID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.MemberID))
		i--
		dAtA[i] = 0x18
	}
	if m.SigningAttempt != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.SigningAttempt))
		i--
		dAtA[i] = 0x10
	}
	if m.SigningID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.SigningID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TextSignatureOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TextSignatureOrder) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TextSignatureOrder) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintTss(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EVMSignature) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EVMSignature) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EVMSignature) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintTss(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.RAddress) > 0 {
		i -= len(m.RAddress)
		copy(dAtA[i:], m.RAddress)
		i = encodeVarintTss(dAtA, i, uint64(len(m.RAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SigningResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SigningResult) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SigningResult) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ReceivedPartialSignatures) > 0 {
		for iNdEx := len(m.ReceivedPartialSignatures) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ReceivedPartialSignatures[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.EVMSignature != nil {
		{
			size, err := m.EVMSignature.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTss(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.CurrentSigningAttempt != nil {
		{
			size, err := m.CurrentSigningAttempt.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTss(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Signing.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTss(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *SigningExpiration) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SigningExpiration) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SigningExpiration) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SigningAttempt != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.SigningAttempt))
		i--
		dAtA[i] = 0x10
	}
	if m.SigningID != 0 {
		i = encodeVarintTss(dAtA, i, uint64(m.SigningID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SigningExpirations) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SigningExpirations) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SigningExpirations) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SigningExpirations) > 0 {
		for iNdEx := len(m.SigningExpirations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SigningExpirations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTss(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTss(dAtA []byte, offset int, v uint64) int {
	offset -= sovTss(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Group) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovTss(uint64(m.ID))
	}
	if m.Size_ != 0 {
		n += 1 + sovTss(uint64(m.Size_))
	}
	if m.Threshold != 0 {
		n += 1 + sovTss(uint64(m.Threshold))
	}
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovTss(uint64(m.Status))
	}
	if m.CreatedHeight != 0 {
		n += 1 + sovTss(uint64(m.CreatedHeight))
	}
	l = len(m.ModuleOwner)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *GroupResult) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Group.Size()
	n += 1 + l + sovTss(uint64(l))
	l = len(m.DKGContext)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	if len(m.Members) > 0 {
		for _, e := range m.Members {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	if len(m.Round1Infos) > 0 {
		for _, e := range m.Round1Infos {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	if len(m.Round2Infos) > 0 {
		for _, e := range m.Round2Infos {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	if len(m.ComplaintsWithStatus) > 0 {
		for _, e := range m.ComplaintsWithStatus {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	if len(m.Confirms) > 0 {
		for _, e := range m.Confirms {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	return n
}

func (m *Round1Info) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MemberID != 0 {
		n += 1 + sovTss(uint64(m.MemberID))
	}
	if len(m.CoefficientCommits) > 0 {
		for _, b := range m.CoefficientCommits {
			l = len(b)
			n += 1 + l + sovTss(uint64(l))
		}
	}
	l = len(m.OneTimePubKey)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.A0Signature)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.OneTimeSignature)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *Round2Info) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MemberID != 0 {
		n += 1 + sovTss(uint64(m.MemberID))
	}
	if len(m.EncryptedSecretShares) > 0 {
		for _, b := range m.EncryptedSecretShares {
			l = len(b)
			n += 1 + l + sovTss(uint64(l))
		}
	}
	return n
}

func (m *DE) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PubD)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.PubE)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *DEQueue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Head != 0 {
		n += 1 + sovTss(uint64(m.Head))
	}
	if m.Tail != 0 {
		n += 1 + sovTss(uint64(m.Tail))
	}
	return n
}

func (m *Signing) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovTss(uint64(m.ID))
	}
	if m.CurrentAttempt != 0 {
		n += 1 + sovTss(uint64(m.CurrentAttempt))
	}
	if m.GroupID != 0 {
		n += 1 + sovTss(uint64(m.GroupID))
	}
	l = len(m.GroupPubKey)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.GroupPubNonce)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	if m.Status != 0 {
		n += 1 + sovTss(uint64(m.Status))
	}
	if m.CreatedHeight != 0 {
		n += 1 + sovTss(uint64(m.CreatedHeight))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.CreatedTimestamp)
	n += 1 + l + sovTss(uint64(l))
	return n
}

func (m *SigningAttempt) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigningID != 0 {
		n += 1 + sovTss(uint64(m.SigningID))
	}
	if m.Attempt != 0 {
		n += 1 + sovTss(uint64(m.Attempt))
	}
	if m.ExpiredHeight != 0 {
		n += 1 + sovTss(uint64(m.ExpiredHeight))
	}
	if len(m.AssignedMembers) > 0 {
		for _, e := range m.AssignedMembers {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	return n
}

func (m *AssignedMember) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MemberID != 0 {
		n += 1 + sovTss(uint64(m.MemberID))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.PubD)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.PubE)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.BindingFactor)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.PubNonce)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *PendingSignings) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SigningIds) > 0 {
		l = 0
		for _, e := range m.SigningIds {
			l += sovTss(uint64(e))
		}
		n += 1 + sovTss(uint64(l)) + l
	}
	return n
}

func (m *Member) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovTss(uint64(m.ID))
	}
	if m.GroupID != 0 {
		n += 1 + sovTss(uint64(m.GroupID))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	if m.IsMalicious {
		n += 2
	}
	if m.IsActive {
		n += 2
	}
	return n
}

func (m *Confirm) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MemberID != 0 {
		n += 1 + sovTss(uint64(m.MemberID))
	}
	l = len(m.OwnPubKeySig)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *Complaint) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Complainant != 0 {
		n += 1 + sovTss(uint64(m.Complainant))
	}
	if m.Respondent != 0 {
		n += 1 + sovTss(uint64(m.Respondent))
	}
	l = len(m.KeySym)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *ComplaintWithStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Complaint.Size()
	n += 1 + l + sovTss(uint64(l))
	if m.ComplaintStatus != 0 {
		n += 1 + sovTss(uint64(m.ComplaintStatus))
	}
	return n
}

func (m *ComplaintsWithStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MemberID != 0 {
		n += 1 + sovTss(uint64(m.MemberID))
	}
	if len(m.ComplaintsWithStatus) > 0 {
		for _, e := range m.ComplaintsWithStatus {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	return n
}

func (m *PendingProcessGroups) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.GroupIDs) > 0 {
		l = 0
		for _, e := range m.GroupIDs {
			l += sovTss(uint64(e))
		}
		n += 1 + sovTss(uint64(l)) + l
	}
	return n
}

func (m *PendingProcessSignings) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SigningIDs) > 0 {
		l = 0
		for _, e := range m.SigningIDs {
			l += sovTss(uint64(e))
		}
		n += 1 + sovTss(uint64(l)) + l
	}
	return n
}

func (m *PartialSignature) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigningID != 0 {
		n += 1 + sovTss(uint64(m.SigningID))
	}
	if m.SigningAttempt != 0 {
		n += 1 + sovTss(uint64(m.SigningAttempt))
	}
	if m.MemberID != 0 {
		n += 1 + sovTss(uint64(m.MemberID))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *TextSignatureOrder) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *EVMSignature) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RAddress)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovTss(uint64(l))
	}
	return n
}

func (m *SigningResult) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Signing.Size()
	n += 1 + l + sovTss(uint64(l))
	if m.CurrentSigningAttempt != nil {
		l = m.CurrentSigningAttempt.Size()
		n += 1 + l + sovTss(uint64(l))
	}
	if m.EVMSignature != nil {
		l = m.EVMSignature.Size()
		n += 1 + l + sovTss(uint64(l))
	}
	if len(m.ReceivedPartialSignatures) > 0 {
		for _, e := range m.ReceivedPartialSignatures {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	return n
}

func (m *SigningExpiration) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SigningID != 0 {
		n += 1 + sovTss(uint64(m.SigningID))
	}
	if m.SigningAttempt != 0 {
		n += 1 + sovTss(uint64(m.SigningAttempt))
	}
	return n
}

func (m *SigningExpirations) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SigningExpirations) > 0 {
		for _, e := range m.SigningExpirations {
			l = e.Size()
			n += 1 + l + sovTss(uint64(l))
		}
	}
	return n
}

func sovTss(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTss(x uint64) (n int) {
	return sovTss(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Group) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: Group: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Group: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= GroupID(b&0x7F) << shift
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
					return ErrIntOverflowTss
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
					return ErrIntOverflowTss
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
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
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
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= GroupStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedHeight", wireType)
			}
			m.CreatedHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ModuleOwner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ModuleOwner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *GroupResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: GroupResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GroupResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Group", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Group.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DKGContext", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DKGContext = append(m.DKGContext[:0], dAtA[iNdEx:postIndex]...)
			if m.DKGContext == nil {
				m.DKGContext = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Members", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Members = append(m.Members, Member{})
			if err := m.Members[len(m.Members)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Round1Infos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Round1Infos = append(m.Round1Infos, Round1Info{})
			if err := m.Round1Infos[len(m.Round1Infos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Round2Infos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Round2Infos = append(m.Round2Infos, Round2Info{})
			if err := m.Round2Infos[len(m.Round2Infos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ComplaintsWithStatus", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ComplaintsWithStatus = append(m.ComplaintsWithStatus, ComplaintsWithStatus{})
			if err := m.ComplaintsWithStatus[len(m.ComplaintsWithStatus)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Confirms", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Confirms = append(m.Confirms, Confirm{})
			if err := m.Confirms[len(m.Confirms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *Round1Info) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: Round1Info: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Round1Info: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberID", wireType)
			}
			m.MemberID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MemberID |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CoefficientCommits", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CoefficientCommits = append(m.CoefficientCommits, make([]byte, postIndex-iNdEx))
			copy(m.CoefficientCommits[len(m.CoefficientCommits)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OneTimePubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OneTimePubKey = append(m.OneTimePubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.OneTimePubKey == nil {
				m.OneTimePubKey = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field A0Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.A0Signature = append(m.A0Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.A0Signature == nil {
				m.A0Signature = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OneTimeSignature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OneTimeSignature = append(m.OneTimeSignature[:0], dAtA[iNdEx:postIndex]...)
			if m.OneTimeSignature == nil {
				m.OneTimeSignature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *Round2Info) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: Round2Info: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Round2Info: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberID", wireType)
			}
			m.MemberID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MemberID |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EncryptedSecretShares", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EncryptedSecretShares = append(m.EncryptedSecretShares, make([]byte, postIndex-iNdEx))
			copy(m.EncryptedSecretShares[len(m.EncryptedSecretShares)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *DE) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: DE: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DE: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubD", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubD = append(m.PubD[:0], dAtA[iNdEx:postIndex]...)
			if m.PubD == nil {
				m.PubD = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubE", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubE = append(m.PubE[:0], dAtA[iNdEx:postIndex]...)
			if m.PubE == nil {
				m.PubE = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *DEQueue) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: DEQueue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DEQueue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Head", wireType)
			}
			m.Head = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Head |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tail", wireType)
			}
			m.Tail = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Tail |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
				return ErrIntOverflowTss
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
					return ErrIntOverflowTss
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentAttempt", wireType)
			}
			m.CurrentAttempt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CurrentAttempt |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupID", wireType)
			}
			m.GroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GroupID |= GroupID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupPubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GroupPubKey = append(m.GroupPubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.GroupPubKey == nil {
				m.GroupPubKey = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = append(m.Message[:0], dAtA[iNdEx:postIndex]...)
			if m.Message == nil {
				m.Message = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupPubNonce", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GroupPubNonce = append(m.GroupPubNonce[:0], dAtA[iNdEx:postIndex]...)
			if m.GroupPubNonce == nil {
				m.GroupPubNonce = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= SigningStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedHeight", wireType)
			}
			m.CreatedHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.CreatedTimestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *SigningAttempt) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: SigningAttempt: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SigningAttempt: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningID", wireType)
			}
			m.SigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningID |= SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attempt", wireType)
			}
			m.Attempt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Attempt |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpiredHeight", wireType)
			}
			m.ExpiredHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpiredHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssignedMembers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AssignedMembers = append(m.AssignedMembers, AssignedMember{})
			if err := m.AssignedMembers[len(m.AssignedMembers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *AssignedMember) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: AssignedMember: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AssignedMember: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberID", wireType)
			}
			m.MemberID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MemberID |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKey = append(m.PubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKey == nil {
				m.PubKey = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubD", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubD = append(m.PubD[:0], dAtA[iNdEx:postIndex]...)
			if m.PubD == nil {
				m.PubD = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubE", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubE = append(m.PubE[:0], dAtA[iNdEx:postIndex]...)
			if m.PubE == nil {
				m.PubE = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BindingFactor", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BindingFactor = append(m.BindingFactor[:0], dAtA[iNdEx:postIndex]...)
			if m.BindingFactor == nil {
				m.BindingFactor = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubNonce", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubNonce = append(m.PubNonce[:0], dAtA[iNdEx:postIndex]...)
			if m.PubNonce == nil {
				m.PubNonce = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *PendingSignings) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: PendingSignings: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PendingSignings: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTss
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.SigningIds = append(m.SigningIds, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTss
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTss
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTss
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.SigningIds) == 0 {
					m.SigningIds = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTss
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.SigningIds = append(m.SigningIds, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningIds", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *Member) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupID", wireType)
			}
			m.GroupID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GroupID |= GroupID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
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
				return fmt.Errorf("proto: wrong wireType = %d for field IsMalicious", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
			m.IsMalicious = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsActive", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *Confirm) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: Confirm: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Confirm: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberID", wireType)
			}
			m.MemberID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MemberID |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OwnPubKeySig", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OwnPubKeySig = append(m.OwnPubKeySig[:0], dAtA[iNdEx:postIndex]...)
			if m.OwnPubKeySig == nil {
				m.OwnPubKeySig = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *Complaint) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: Complaint: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Complaint: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Complainant", wireType)
			}
			m.Complainant = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Complainant |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Respondent", wireType)
			}
			m.Respondent = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Respondent |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeySym", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KeySym = append(m.KeySym[:0], dAtA[iNdEx:postIndex]...)
			if m.KeySym == nil {
				m.KeySym = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *ComplaintWithStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: ComplaintWithStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ComplaintWithStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Complaint", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Complaint.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ComplaintStatus", wireType)
			}
			m.ComplaintStatus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ComplaintStatus |= ComplaintStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *ComplaintsWithStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: ComplaintsWithStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ComplaintsWithStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberID", wireType)
			}
			m.MemberID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MemberID |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ComplaintsWithStatus", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ComplaintsWithStatus = append(m.ComplaintsWithStatus, ComplaintWithStatus{})
			if err := m.ComplaintsWithStatus[len(m.ComplaintsWithStatus)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *PendingProcessGroups) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: PendingProcessGroups: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PendingProcessGroups: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v GroupID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTss
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= GroupID(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.GroupIDs = append(m.GroupIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTss
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTss
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTss
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.GroupIDs) == 0 {
					m.GroupIDs = make([]GroupID, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v GroupID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTss
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= GroupID(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.GroupIDs = append(m.GroupIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field GroupIDs", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *PendingProcessSignings) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: PendingProcessSignings: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PendingProcessSignings: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v SigningID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTss
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= SigningID(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.SigningIDs = append(m.SigningIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTss
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTss
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTss
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.SigningIDs) == 0 {
					m.SigningIDs = make([]SigningID, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v SigningID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTss
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= SigningID(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.SigningIDs = append(m.SigningIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningIDs", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *PartialSignature) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: PartialSignature: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PartialSignature: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningID", wireType)
			}
			m.SigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningID |= SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningAttempt", wireType)
			}
			m.SigningAttempt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningAttempt |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MemberID", wireType)
			}
			m.MemberID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MemberID |= MemberID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *TextSignatureOrder) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: TextSignatureOrder: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TextSignatureOrder: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = append(m.Message[:0], dAtA[iNdEx:postIndex]...)
			if m.Message == nil {
				m.Message = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *EVMSignature) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: EVMSignature: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EVMSignature: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RAddress = append(m.RAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.RAddress == nil {
				m.RAddress = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *SigningResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: SigningResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SigningResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signing", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Signing.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CurrentSigningAttempt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CurrentSigningAttempt == nil {
				m.CurrentSigningAttempt = &SigningAttempt{}
			}
			if err := m.CurrentSigningAttempt.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EVMSignature", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.EVMSignature == nil {
				m.EVMSignature = &EVMSignature{}
			}
			if err := m.EVMSignature.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceivedPartialSignatures", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceivedPartialSignatures = append(m.ReceivedPartialSignatures, PartialSignature{})
			if err := m.ReceivedPartialSignatures[len(m.ReceivedPartialSignatures)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *SigningExpiration) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: SigningExpiration: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SigningExpiration: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningID", wireType)
			}
			m.SigningID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningID |= SigningID(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningAttempt", wireType)
			}
			m.SigningAttempt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SigningAttempt |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func (m *SigningExpirations) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTss
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
			return fmt.Errorf("proto: SigningExpirations: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SigningExpirations: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningExpirations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTss
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
				return ErrInvalidLengthTss
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTss
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SigningExpirations = append(m.SigningExpirations, SigningExpiration{})
			if err := m.SigningExpirations[len(m.SigningExpirations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTss
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
func skipTss(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTss
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
					return 0, ErrIntOverflowTss
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
					return 0, ErrIntOverflowTss
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
				return 0, ErrInvalidLengthTss
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTss
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTss
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTss        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTss          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTss = fmt.Errorf("proto: unexpected end of group")
)
