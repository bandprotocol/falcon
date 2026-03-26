package icon

import (
	"encoding/json"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner signs ICON transactions via the KMS service.
type RemoteSigner struct {
	wallet.BaseRemoteSigner
}

// NewRemoteSigner creates a new ICON RemoteSigner.
// Its signature matches wallet.RemoteSignerFactory so it can be passed
// directly to wallet.NewRemoteOnlyAdapter.
func NewRemoteSigner(name, address, url string, key string) (wallet.Signer, error) {
	base, err := wallet.NewBaseRemoteSigner(name, address, url, key)
	if err != nil {
		return nil, err
	}
	return &RemoteSigner{BaseRemoteSigner: *base}, nil
}

// Sign requests the remote KMS to sign the ICON transaction.
func (r *RemoteSigner) Sign(payload []byte, tssPayload wallet.TssPayload) ([]byte, error) {
	var signerPayload SignerPayload
	if err := json.Unmarshal(payload, &signerPayload); err != nil {
		return nil, err
	}
	res, err := r.FkmsClient.SignIcon(
		r.ContextWithKey(),
		&fkmsv1.SignIconRequest{
			SignerPayload: &fkmsv1.IconSignerPayload{
				Relayer:         signerPayload.Relayer,
				ContractAddress: signerPayload.ContractAddress,
				StepLimit:       signerPayload.StepLimit,
				NetworkId:       signerPayload.NetworkID,
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

	return res.TxParams, nil
}
