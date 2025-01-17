package chains

import (
	"context"
	"math/big"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

// ChainProvider defines the interface for the chain interaction with the destination chain.
type ChainProvider interface {
	KeyProvider
	// Init initialize to the chain.
	Init(ctx context.Context) error

	// QueryTunnelInfo queries the tunnel information from the destination chain.
	QueryTunnelInfo(
		ctx context.Context,
		tunnelID uint64,
		tunnelDestinationAddr string,
	) (*chainstypes.Tunnel, error)

	// RelayPacket relays the packet from the source chain to the destination chain.
	RelayPacket(ctx context.Context, packet *bandtypes.Packet) error

	// QueryBalance queries balance by given key name from the destination chain.
	QueryBalance(ctx context.Context, keyName string) (*big.Int, error)

	// Set Prometheus metrics in chain provider
	SetMetrics(m *relayermetrics.PrometheusMetrics)
}

// KeyProvider defines the interface for the key interaction with destination chain
type KeyProvider interface {
	// AddKey stores the private key with a given mnemonic and key name on the user's local disk.
	AddKey(
		keyName string,
		mnemonic string,
		privateKeyHex string,
		homePath string,
		coinType uint32,
		account uint,
		index uint,
		passphrase string,
	) (*chainstypes.Key, error)

	// IsKeyNameExist checks whether a key with the specified keyName already exists in storage.
	IsKeyNameExist(keyName string) bool

	// ExportPrivateKey exports private key of specified key name.
	ExportPrivateKey(keyName string, passphrase string) (string, error)

	// DeleteKey deletes the key information and private key
	DeleteKey(homePath, keyName, passphrase string) error

	// ListKeys lists all keys
	ListKeys() []*chainstypes.Key

	// ShowKey shows the address of the given key
	ShowKey(keyName string) string

	// LoadFreeSenders loads key info to prepare to relay the packet
	LoadFreeSenders(homePath, passphrase string) error
}

// BaseChainProvider is a base object for connecting with the chain network.
type BaseChainProvider struct {
	log *zap.Logger

	Config    ChainProviderConfig
	ChainName string
	ChainID   string

	debug bool
}
