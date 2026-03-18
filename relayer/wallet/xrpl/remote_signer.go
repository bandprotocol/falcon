package xrpl

import (
	"encoding/hex"
	"encoding/json"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner is signer that uses KMS service to sign XRPL data.
type RemoteSigner struct {
	wallet.BaseRemoteSigner
}

// NewRemoteSigner creates a new RemoteSigner instance.
func NewRemoteSigner(name, address, url string, key string) (*RemoteSigner, error) {
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
	res, err := r.FkmsClient.SignXrpl(
		r.ContextWithKey(),
		&fkmsv1.SignXrplRequest{
			SignerPayload: &fkmsv1.XrplSignerPayload{
				Account:  signerPayload.Account,
				OracleId: signerPayload.OracleID,
				Fee:      signerPayload.Fee,
				Sequence: uint64(signerPayload.Sequence),
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

	return []byte(hex.EncodeToString(res.TxBlob)), nil
}
