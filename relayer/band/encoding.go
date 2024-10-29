package band

import (
	"cosmossdk.io/x/tx/signing"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	tunneltypes "github.com/bandprotocol/chain/v3/x/tunnel/types"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Marshaler         codec.Codec
}

// MakeEncodingConfig creates an EncodingConfig for an amino based test configuration.
func MakeEncodingConfig() EncodingConfig {
	interfaceRegistry, err := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	interfaceRegistry.RegisterInterface(
		"tunnel.v1beta1.RouteI",
		(*tunneltypes.RouteI)(nil),
		&tunneltypes.TSSRoute{},
		&tunneltypes.AxelarRoute{},
	)

	interfaceRegistry.RegisterInterface(
		"tunnel.v1beta1.PacketContentI",
		(*tunneltypes.PacketContentI)(nil),
		&tunneltypes.TSSPacketContent{},
		&tunneltypes.AxelarPacketContent{},
	)

	cdc := codec.NewProtoCodec(interfaceRegistry)
	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         cdc,
	}
}
