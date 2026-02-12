package xrpl

import (
	"context"
	"fmt"

	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
	"github.com/bandprotocol/falcon/relayer/wallet"
)

var _ wallet.Signer = (*RemoteSigner)(nil)

// RemoteSigner is a placeholder for XRPL remote signers.
type RemoteSigner struct {
	Name       string
	Address    string
	Key        *string
	FkmsClient fkmsv1.FkmsServiceClient
}

// NewRemoteSigner creates a new RemoteSigner instance.
func NewRemoteSigner(name, address, url string, key *string) (*RemoteSigner, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote signer at %s: %w", url, err)
	}

	fkmsClient := fkmsv1.NewFkmsServiceClient(conn)
	return &RemoteSigner{
		Name:       name,
		Address:    address,
		Key:        key,
		FkmsClient: fkmsClient,
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
	return r.Address
}

// Sign is unsupported for XRPL remote signers.
func (r *RemoteSigner) Sign(data []byte, preSignPayload *wallet.PreSignPayload) ([]byte, error) {
	ctx := context.Background()
	if r.Key != nil {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("api-key", *r.Key))
	}

	res, err := r.FkmsClient.SignXrpl(
		ctx,
		&fkmsv1.SignXrplRequest{Address: r.Address, TxMessage: data, Tss: &fkmsv1.Tss{
			Message:    preSignPayload.TssMessage,
			RandomAddr: preSignPayload.RandomAddr,
			SignatureS: preSignPayload.Signature,
		}},
	)
	if err != nil {
		return []byte{}, err
	}

	return res.TxBlob, nil
}
