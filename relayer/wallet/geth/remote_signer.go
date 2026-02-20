package geth

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner is signer that uses KMS service to sign data.
type RemoteSigner struct {
	Name       string
	Address    common.Address
	FkmsClient fkmsv1.FkmsServiceClient
	Key        *string
}

// NewRemoteSigner creates a new RemoteSigner instance.
func NewRemoteSigner(name string, address common.Address, url string, key *string) (*RemoteSigner, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote signer at %s: %w", url, err)
	}

	fkmsClient := fkmsv1.NewFkmsServiceClient(conn)

	return &RemoteSigner{
		Name:       name,
		Address:    address,
		FkmsClient: fkmsClient,
		Key:        key,
	}, nil
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
	return r.Address.String()
}

// Sign requests the remote KMS to sign the data and returns the signature.
func (r *RemoteSigner) Sign(data []byte, preSignPayload wallet.PreSignPayload) ([]byte, error) {
	ctx := context.Background()
	if r.Key != nil {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("api-key", *r.Key))
	}

	tss := &fkmsv1.Tss{
		Message:    preSignPayload.TssMessage,
		RandomAddr: preSignPayload.RandomAddr,
		SignatureS: preSignPayload.Signature,
	}

	res, err := r.FkmsClient.SignEvm(
		ctx,
		&fkmsv1.SignEvmRequest{Address: strings.ToLower(r.Address.String()), TxMessage: data, Tss: tss},
	)
	if err != nil {
		return []byte{}, err
	}

	return res.Signature, nil
}
