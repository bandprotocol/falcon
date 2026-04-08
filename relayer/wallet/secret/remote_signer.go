package secret

import (
	"encoding/json"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner signs Cosmos SDK transactions via fkms.SignSecret.
type RemoteSigner struct {
	wallet.BaseRemoteSigner
}

// NewRemoteSigner creates a Secret RemoteSigner.
func NewRemoteSigner(name, address, url string, key string) (wallet.Signer, error) {
	base, err := wallet.NewBaseRemoteSigner(name, address, url, key)
	if err != nil {
		return nil, err
	}
	return &RemoteSigner{BaseRemoteSigner: *base}, nil
}

// Sign requests fkms to build & sign a Secret Network TxRaw, then returns it as tx_blob bytes.
func (r *RemoteSigner) Sign(payload []byte, tssPayload wallet.TssPayload) ([]byte, error) {
	var signerPayload SignerPayload
	if err := json.Unmarshal(payload, &signerPayload); err != nil {
		return nil, err
	}

	res, err := r.FkmsClient.SignSecret(
		r.ContextWithKey(),
		&fkmsv1.SignSecretRequest{
			SignerPayload: &fkmsv1.SecretSignerPayload{
				Sender:          signerPayload.Sender,
				ContractAddress: signerPayload.ContractAddress,
				ChainId:         signerPayload.ChainID,
				AccountNumber:   signerPayload.AccountNumber,
				Sequence:        signerPayload.Sequence,
				GasLimit:        signerPayload.GasLimit,
				GasPrices:       signerPayload.GasPrices,
				Memo:            signerPayload.Memo,
				CodeHash:        signerPayload.CodeHash,
				Pubkey:          signerPayload.PubKey,
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
