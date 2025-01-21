package bandtss

import (
	"github.com/cosmos/cosmos-sdk/types"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
)

// Signing is a bandtss signing information.
type Signing struct {
	// id is the unique identifier of the bandtss signing.
	ID SigningID `protobuf:"varint,1,opt,name=id,proto3,casttype=SigningID"                                                                                                          json:"id,omitempty"`
	// fee_per_signer is the tokens that will be paid per signer for this bandtss signing.
	FeePerSigner types.Coins `protobuf:"bytes,2,rep,name=fee_per_signer,json=feePerSigner,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins"                                          json:"fee_per_signer"`
	// requester is the address who pays the Bandtss signing.
	Requester string `protobuf:"bytes,3,opt,name=requester,proto3"                                                                                                                       json:"requester,omitempty"`
	// current_group_signing_id is a tss signing ID of a current group.
	CurrentGroupSigningID tsstypes.SigningID `protobuf:"varint,4,opt,name=current_group_signing_id,json=currentGroupSigningId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID"   json:"current_group_signing_id,omitempty"`
	// incoming_group_signing_id is a tss signing ID of an incoming group, if any.
	IncomingGroupSigningID tsstypes.SigningID `protobuf:"varint,5,opt,name=incoming_group_signing_id,json=incomingGroupSigningId,proto3,casttype=github.com/bandprotocol/falcon/internal/bandchain/tss.SigningID" json:"incoming_group_signing_id,omitempty"`
}
