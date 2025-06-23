package chains

import (
	"context"
	"math/big"

	"go.uber.org/zap"

	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
)

// ChainProviders is a collection of ChainProvider interfaces (mapped by chainName)
type ChainProviders map[string]ChainProvider

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

	// GetChainName retrieves the chain name from the chain provider.
	GetChainName() string
}

// KeyProvider defines the interface for the key interaction with destination chain
type KeyProvider interface {
	// AddKeyByMnemonic adds a key using a mnemonic phrase.
	AddKeyByMnemonic(
		keyName string,
		mnemonic string,
		coinType uint32,
		account uint,
		index uint,
	) (*chainstypes.Key, error)

	// AddKeyByPrivateKey adds a key using a private key.
	AddKeyByPrivateKey(keyName string, privateKeyHex string) (*chainstypes.Key, error)

	// AddRemoteSignerKey adds a key using a remote signer’s address and a Falcon KMS URL.
	AddRemoteSignerKey(keyName string, addr string, url string) (*chainstypes.Key, error)

	// DeleteKey deletes the key information and private key
	DeleteKey(keyName string) error

	// ExportPrivateKey exports private key of specified key name.
	ExportPrivateKey(keyName string) (string, error)

	// ListKeys lists all keys
	ListKeys() []*chainstypes.Key

	// ShowKey shows the address of the given key
	ShowKey(keyName string) (string, error)

	// LoadSigners loads signers to prepare to relay the packet
	LoadSigners() error
}

// BaseChainProvider is a base object for connecting with the chain network.
type BaseChainProvider struct {
	log *zap.Logger

	Config    ChainProviderConfig
	ChainName string
	ChainID   string

	debug bool
}
