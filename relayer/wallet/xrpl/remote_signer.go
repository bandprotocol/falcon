package xrpl

import (
	"fmt"

	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner is a placeholder for XRPL remote signers.
type RemoteSigner struct {
	Name    string
	Address string
	Url     string
	Key     *string
}

// NewRemoteSigner creates a new RemoteSigner instance.
func NewRemoteSigner(name, address, url string, key *string) *RemoteSigner {
	return &RemoteSigner{
		Name:    name,
		Address: address,
		Url:     url,
		Key:     key,
	}
}

// ExportPrivateKey always returns an error for remote signer.
func (r *RemoteSigner) ExportPrivateKey() (string, error) {
	return "", fmt.Errorf("cannot extract private key from remote signer")
}

// GetName returns the signer's key name.
func (r *RemoteSigner) GetName() string {
	return r.Name
}

// GetAddress returns the signer's address.
func (r *RemoteSigner) GetAddress() (addr string) {
	return r.Address
}

// Sign is unsupported for XRPL remote signers.
func (r *RemoteSigner) Sign(data []byte) ([]byte, error) {
	return nil, fmt.Errorf("xrpl remote signer is not supported")
}
