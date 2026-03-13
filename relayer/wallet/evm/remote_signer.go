package evm

import (
	"fmt"
	"strings"

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

// Sign requests the remote KMS to sign the data and returns the signature.
func (r *RemoteSigner) Sign(payload []byte, _ wallet.TssPayload) ([]byte, error) {
	res, err := r.FkmsClient.SignEvm(
		r.ContextWithKey(),
		&fkmsv1.SignEvmRequest{Address: strings.ToLower(r.Address), Message: payload},
	)
	if err != nil {
		return nil, err
	}

	return res.Signature, nil
}
