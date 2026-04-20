package evm

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner is signer that uses KMS service to sign EVM data.
type RemoteSigner struct {
	wallet.BaseRemoteSigner
}

// NewRemoteSigner creates a new RemoteSigner instance.
func NewRemoteSigner(name, address, url string, key string) (*RemoteSigner, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid EVM address: %s", address)
	}

	base, err := wallet.NewBaseRemoteSigner(name, common.HexToAddress(address).Hex(), url, key)
	if err != nil {
		return nil, err
	}

	return &RemoteSigner{BaseRemoteSigner: *base}, nil
}

// Sign JSON-unmarshals payload into SignerPayload and delegates transaction
// building, signing, and RLP encoding to the remote KMS. Returns the
// EIP-2718 encoded signed transaction bytes.
func (r *RemoteSigner) Sign(payload []byte, tssPayload wallet.TssPayload) ([]byte, error) {
	var signerPayload SignerPayload
	if err := json.Unmarshal(payload, &signerPayload); err != nil {
		return nil, err
	}

	res, err := r.FkmsClient.SignEvm(
		r.ContextWithKey(),
		&fkmsv1.SignEvmRequest{
			SignerPayload: &fkmsv1.EvmSignerPayload{
				Address:   signerPayload.Address,
				ChainId:   signerPayload.ChainID,
				Nonce:     signerPayload.Nonce,
				To:        signerPayload.To,
				GasLimit:  signerPayload.GasLimit,
				GasPrice:  signerPayload.GasPrice,
				GasFeeCap: signerPayload.GasFeeCap,
				GasTipCap: signerPayload.GasTipCap,
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
