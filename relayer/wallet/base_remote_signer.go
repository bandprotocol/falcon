package wallet

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	fkmsv1 "github.com/bandprotocol/falcon/proto/fkms/v1"
)

// BaseRemoteSigner provides shared fields and methods for remote signers
// across all chain types.
type BaseRemoteSigner struct {
	Name       string
	Address    string
	FkmsClient fkmsv1.FkmsServiceClient
	Key        string
}

// NewBaseRemoteSigner creates a BaseRemoteSigner with a gRPC connection to the KMS.
func NewBaseRemoteSigner(name, address, url string, key string) (*BaseRemoteSigner, error) {
	conn, err := newGRPCConn(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote signer at %s: %w", url, err)
	}

	fkmsClient := fkmsv1.NewFkmsServiceClient(conn)

	return &BaseRemoteSigner{
		Name:       name,
		Address:    address,
		FkmsClient: fkmsClient,
		Key:        key,
	}, nil
}

func newGRPCConn(url string) (*grpc.ClientConn, error) {
	var opts grpc.DialOption
	var host string

	if strings.HasPrefix(url, "https://") {
		opts = grpc.WithTransportCredentials(credentials.NewTLS(nil))
		host = strings.TrimPrefix(url, "https://")

		if !strings.Contains(host, ":") {
			host += ":443"
		}
	} else {
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
		host = strings.TrimPrefix(url, "http://")

		if !strings.Contains(host, ":") {
			host += ":80"
		}
	}

	return grpc.NewClient(host, opts)
}

// ExportPrivateKey always returns an error for remote signers.
func (r *BaseRemoteSigner) ExportPrivateKey() (string, error) {
	return "", fmt.Errorf("cannot extract private key from remote signer")
}

// GetName returns the signer's key name.
func (r *BaseRemoteSigner) GetName() string {
	return r.Name
}

// GetAddress returns the signer's address.
func (r *BaseRemoteSigner) GetAddress() string {
	return r.Address
}

// ContextWithKey returns a context with the API key metadata attached, if present.
func (r *BaseRemoteSigner) ContextWithKey() context.Context {
	ctx := context.Background()
	if r.Key != "" {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("api-key", r.Key))
	}
	return ctx
}
