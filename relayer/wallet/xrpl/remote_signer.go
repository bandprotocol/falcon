package xrpl

import (
	"context"
	"encoding/hex"
	"fmt"

	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	// binarycodec "github.com/Peersyst/xrpl-go/binary-codec"
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

// remoteSign requests the remote KMS to sign the data and returns the tx blob.
func (r *RemoteSigner) remoteSign(signerPayload SignerPayload, tssPayload wallet.TssPayload) (string, error) {
	ctx := context.Background()
	if r.Key != nil {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("api-key", *r.Key))
	}

	res, err := r.FkmsClient.SignXrpl(
		ctx,
		&fkmsv1.SignXrplRequest{
			SignerPayload: &fkmsv1.XrplSignerPayload{
				Account:         signerPayload.Account,
				OracleId:        signerPayload.OracleId,
				Fee:             signerPayload.Fee,
				Sequence:        signerPayload.Sequence,
				LastUpdatedTime: signerPayload.LastUpdatedTime,
			},
			Tss: &fkmsv1.Tss{
				Message:    tssPayload.TssMessage,
				RandomAddr: tssPayload.RandomAddr,
				SignatureS: tssPayload.Signature,
			},
		},
	)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(res.TxBlob), nil
}
