package band

import (
	"context"

	cosmosgrpc "github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc"

	bandtsstypes "github.com/bandprotocol/falcon/internal/bandchain/bandtss"
	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
)

var _ QueryClient = (*BandQueryClient)(nil)

type QueryClient interface {
	Tunnel(
		ctx context.Context,
		in *tunneltypes.QueryTunnelRequest,
		opts ...grpc.CallOption,
	) (*tunneltypes.QueryTunnelResponse, error)

	Tunnels(
		ctx context.Context,
		in *tunneltypes.QueryTunnelsRequest,
		opts ...grpc.CallOption,
	) (*tunneltypes.QueryTunnelsResponse, error)

	Packet(
		ctx context.Context,
		in *tunneltypes.QueryPacketRequest,
		opts ...grpc.CallOption,
	) (*tunneltypes.QueryPacketResponse, error)

	Signing(
		ctx context.Context,
		in *bandtsstypes.QuerySigningRequest,
		opts ...grpc.CallOption,
	) (*bandtsstypes.QuerySigningResponse, error)
}

type BandQueryClient struct {
	cc cosmosgrpc.ClientConn
}

func NewBandQueryClient(cc cosmosgrpc.ClientConn) *BandQueryClient {
	return &BandQueryClient{cc}
}

func (c *BandQueryClient) Tunnel(
	ctx context.Context,
	in *tunneltypes.QueryTunnelRequest,
	opts ...grpc.CallOption,
) (*tunneltypes.QueryTunnelResponse, error) {
	out := new(tunneltypes.QueryTunnelResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Tunnel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *BandQueryClient) Tunnels(
	ctx context.Context,
	in *tunneltypes.QueryTunnelsRequest,
	opts ...grpc.CallOption,
) (*tunneltypes.QueryTunnelsResponse, error) {
	out := new(tunneltypes.QueryTunnelsResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Tunnels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *BandQueryClient) Packet(
	ctx context.Context,
	in *tunneltypes.QueryPacketRequest,
	opts ...grpc.CallOption,
) (*tunneltypes.QueryPacketResponse, error) {
	out := new(tunneltypes.QueryPacketResponse)
	err := c.cc.Invoke(ctx, "/band.tunnel.v1beta1.Query/Packet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *BandQueryClient) Signing(
	ctx context.Context,
	in *bandtsstypes.QuerySigningRequest,
	opts ...grpc.CallOption,
) (*bandtsstypes.QuerySigningResponse, error) {
	out := new(bandtsstypes.QuerySigningResponse)
	err := c.cc.Invoke(ctx, "/band.bandtss.v1beta1.Query/Signing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
