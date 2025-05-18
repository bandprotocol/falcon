package geth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kms "github.com/bandprotocol/falcon/proto/kms"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner is signer that uses KMS service to sign data.
type RemoteSigner struct {
	Name      string
	Address   common.Address
	KmsClient kms.KmsEvmServiceClient
}

// NewRemoteSigner creates a new RemoteSigner instance.
func NewRemoteSigner(name string, address common.Address, url string) (*RemoteSigner, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote signer at %s: %w", url, err)
	}

	kmsClient := kms.NewKmsEvmServiceClient(conn)

	return &RemoteSigner{
		Name:      name,
		Address:   address,
		KmsClient: kmsClient,
	}, nil
}

// ExportPrivateKey always returns an error for remote signerห.
func (r *RemoteSigner) ExportPrivateKey() (string, error) {
	return "", fmt.Errorf("cannot extract private key from remote signer")
}

// GetName returns the signer's key name.
func (r *RemoteSigner) GetName() string {
	return r.Name
}

// GetAddress returns the signer's address.
func (r *RemoteSigner) GetAddress() (addr string) {
	return r.Address.String()
}

// Sign requests the remote KMS to sign the data and returns the signature.
func (r *RemoteSigner) Sign(data []byte) ([]byte, error) {
	res, err := r.KmsClient.SignEvm(
		context.Background(),
		&kms.SignEvmRequest{Addr: r.Address.Hex(), Message: data},
	)
	if err != nil {
		return []byte{}, err
	}

	return res.Signature, nil
}
