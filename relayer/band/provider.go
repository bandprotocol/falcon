package band

import (
	"context"
	"fmt"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"

	bandtsstypes "github.com/bandprotocol/chain/v3/x/bandtss/types"
	tunneltypes "github.com/bandprotocol/chain/v3/x/tunnel/types"
)

var _ BandQueryClient = &QueryClientProvider{}

type BandQueryClient interface {
	Tunnel(ctx context.Context, tunnelID uint64) (*tunneltypes.QueryTunnelResponse, error)
	Packet(ctx context.Context, tunnelID uint64, sequence uint64) (*tunneltypes.QueryPacketResponse, error)
	Signing(ctx context.Context, signingID uint64) (*bandtsstypes.QuerySigningResponse, error)
}

type QueryClientProvider struct {
	TunnelClient  tunneltypes.QueryClient
	BandtssClient bandtsstypes.QueryClient
}

func NewQueryClientProvider(c cosmosclient.Context) *QueryClientProvider {
	return &QueryClientProvider{
		TunnelClient:  tunneltypes.NewQueryClient(c),
		BandtssClient: bandtsstypes.NewQueryClient(c),
	}
}

func (t QueryClientProvider) Tunnel(ctx context.Context, tunnelID uint64) (*tunneltypes.QueryTunnelResponse, error) {
	res, err := t.TunnelClient.Tunnel(ctx, &tunneltypes.QueryTunnelRequest{
		TunnelId: tunnelID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tunnel: %w", err)
	}

	return res, nil
}

func (t QueryClientProvider) Packet(ctx context.Context, tunnelID uint64, sequence uint64) (*tunneltypes.QueryPacketResponse, error) {
	res, err := t.TunnelClient.Packet(ctx, &tunneltypes.QueryPacketRequest{
		TunnelId: tunnelID,
		Sequence: sequence,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch packet: %w", err)
	}

	return res, nil
}

func (t QueryClientProvider) Signing(ctx context.Context, signingID uint64) (*bandtsstypes.QuerySigningResponse, error) {
	res, err := t.BandtssClient.Signing(ctx, &bandtsstypes.QuerySigningRequest{
		SigningId: signingID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch signing: %w", err)
	}

	return res, nil
}
