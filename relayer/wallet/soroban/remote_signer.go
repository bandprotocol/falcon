package soroban

import (
	"encoding/json"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner is signer that uses KMS service to sign Soroban data.
type RemoteSigner struct {
	wallet.BaseRemoteSigner
}

// NewRemoteSigner creates a new RemoteSigner instance.
func NewRemoteSigner(name, address, url string, key string) (wallet.Signer, error) {
	base, err := wallet.NewBaseRemoteSigner(name, address, url, key)
	if err != nil {
		return nil, err
	}

	return &RemoteSigner{BaseRemoteSigner: *base}, nil
}

// Sign requests the remote KMS to sign the data and returns the tx blob.
func (r *RemoteSigner) Sign(payload []byte, tssPayload wallet.TssPayload) ([]byte, error) {
	var signerPayload SignerPayload
	if err := json.Unmarshal(payload, &signerPayload); err != nil {
		return nil, err
	}
	res, err := r.FkmsClient.SignSoroban(
		r.ContextWithKey(),
		&fkmsv1.SignSorobanRequest{
			SignerPayload: &fkmsv1.SorobanSignerPayload{
				SourceAccount:     signerPayload.SourceAccount,
				ContractAddress:   signerPayload.ContractAddress,
				Fee:               signerPayload.Fee,
				Sequence:          signerPayload.Sequence,
				NetworkPassphrase: signerPayload.NetworkPassphrase,
				RpcUrl:            signerPayload.RpcUrl,
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
