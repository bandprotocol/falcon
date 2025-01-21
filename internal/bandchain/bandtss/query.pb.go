package bandtss

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
)

// QuerySingingRequest is the request type for the Query/Signing RPC method.
type QuerySigningRequest struct {
	// signing_id is the ID of the signing request.
	SigningId uint64 `protobuf:"varint,1,opt,name=signing_id,json=signingId,proto3" json:"signing_id,omitempty"`
}

func (m *QuerySigningRequest) Reset()         { *m = QuerySigningRequest{} }
func (m *QuerySigningRequest) String() string { return proto.CompactTextString(m) }
func (*QuerySigningRequest) ProtoMessage()    {}

// QuerySigningResponse is the response type for the Query/Signing RPC method.
type QuerySigningResponse struct {
	// fee_per_signer is the tokens that will be paid per signer for this bandtss signing.
	FeePerSigner types.Coins `protobuf:"bytes,1,rep,name=fee_per_signer,json=feePerSigner,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"fee_per_signer"`
	// requester is the address of requester who paid for bandtss signing.
	Requester string `protobuf:"bytes,2,opt,name=requester,proto3"                                                                              json:"requester,omitempty"`
	// current_group_signing_result is the signing result from the current group.
	CurrentGroupSigningResult *tsstypes.SigningResult `protobuf:"bytes,3,opt,name=current_group_signing_result,json=currentGroupSigningResult,proto3"                            json:"current_group_signing_result,omitempty"`
	// incoming_group_signing_result is the signing result from the incoming group.
	IncomingGroupSigningResult *tsstypes.SigningResult `protobuf:"bytes,4,opt,name=incoming_group_signing_result,json=incomingGroupSigningResult,proto3"                          json:"incoming_group_signing_result,omitempty"`
}

func (m *QuerySigningResponse) Reset()         { *m = QuerySigningResponse{} }
func (m *QuerySigningResponse) String() string { return proto.CompactTextString(m) }
func (*QuerySigningResponse) ProtoMessage()    {}
