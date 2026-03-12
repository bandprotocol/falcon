package geth

import (
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
func NewRemoteSigner(name string, address common.Address, url string, key *string) (*RemoteSigner, error) {
	base, err := wallet.NewBaseRemoteSigner(name, address.String(), url, key)
	if err != nil {
		return nil, err
	}

	return &RemoteSigner{BaseRemoteSigner: *base}, nil
}

// remoteSign requests the remote KMS to sign the data and returns the signature.
func (r *RemoteSigner) remoteSign(data []byte) ([]byte, error) {
	res, err := r.FkmsClient.SignEvm(
		r.ContextWithKey(),
		&fkmsv1.SignEvmRequest{Address: strings.ToLower(r.Address), Message: data},
	)
	if err != nil {
		return []byte{}, err
	}

	return res.Signature, nil
}
