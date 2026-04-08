package flow

import (
	"encoding/json"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner is a signer that uses the KMS service to sign Flow transactions.
type RemoteSigner struct {
	wallet.BaseRemoteSigner
}

// NewRemoteSigner creates a new Flow RemoteSigner.
// Its signature matches wallet.RemoteSignerFactory so it can be passed
// directly to wallet.NewRemoteOnlyAdapter.
func NewRemoteSigner(name, address, url string, key string) (wallet.Signer, error) {
	base, err := wallet.NewBaseRemoteSigner(name, address, url, key)
	if err != nil {
		return nil, err
	}

	return &RemoteSigner{BaseRemoteSigner: *base}, nil
}

// Sign requests the remote KMS to sign the Flow transaction and returns the signed tx blob.
func (r *RemoteSigner) Sign(payload []byte, tssPayload wallet.TssPayload) ([]byte, error) {
	var signerPayload SignerPayload
	if err := json.Unmarshal(payload, &signerPayload); err != nil {
		return nil, err
	}

	res, err := r.FkmsClient.SignFlow(
		r.ContextWithKey(),
		&fkmsv1.SignFlowRequest{
			SignerPayload: &fkmsv1.FlowSignerPayload{
				Address:         signerPayload.Address,
				ComputeLimit:    signerPayload.ComputeLimit,
				BlockId:         signerPayload.BlockID,
				KeyIndex:        signerPayload.KeyIndex,
				Sequence:        signerPayload.Sequence,
				ContractAddress: signerPayload.ContractAddress,
			},
			Tss: &fkmsv1.Tss{
				Message:    tssPayload.TssMessage,
				RandomAddr: tssPayload.RandomAddr,
				SignatureS: tssPayload.Signature,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return res.TxBlob, nil
}
