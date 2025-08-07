package relayer

import (
	"context"
	"math/big"

	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/config"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/store"
	"github.com/bandprotocol/falcon/relayer/types"
)

// AppOptions defines an interface that is passed into an application constructor,
// typically used to set App options that are either supplied via config file
// or through CLI arguments/flags. Note, casting Get calls may not yield the expected
// types and could result in type assertion errors. It is recommend to either
// use the cast package or perform manual conversion for safety.
type AppOptions interface {
	Get(string) any
}

// AppCreator is a function that allows us to lazily initialize an
// application using various configurations.
type AppCreator func(
	store store.Store,
	appOpt AppOptions,
) (Application, error)

// Application is an interface that wraps the basic methods of the application.
type Application interface {
	Init(ctx context.Context) error

	GetConfig() *config.Config
	SaveConfig(cfg *config.Config) error

	GetLog() logger.ZapLogger
	GetPassphrase() string
	SavePassphrase(passphrase string) error

	AddChainConfig(chainName string, filePath string) error
	DeleteChainConfig(chainName string) error
	GetChainConfig(chainName string) (chains.ChainProviderConfig, error)

	AddKeyByPrivateKey(chainName string, keyName string, privateKey string) (*chainstypes.Key, error)
	AddKeyByMnemonic(
		chainName string,
		keyName string,
		mnemonic string,
		coinType uint32,
		account uint,
		index uint,
	) (*chainstypes.Key, error)
	AddRemoteSignerKey(chainName string, keyName string, address string, url string, key *string) (*chainstypes.Key, error)
	DeleteKey(chainName string, keyName string) error
	ListKeys(chainName string) ([]*chainstypes.Key, error)
	ExportKey(chainName string, keyName string) (string, error)
	ShowKey(chainName string, keyName string) (string, error)

	Relay(ctx context.Context, tunnelID uint64, isForce bool) error
	Start(ctx context.Context, tunnelIDs []uint64, tunnelCreator string) error
	QueryTunnelInfo(ctx context.Context, tunnelID uint64) (*types.Tunnel, error)
	QueryTunnelPacketInfo(ctx context.Context, tunnelID uint64, sequence uint64) (*bandtypes.Packet, error)
	QueryBalance(ctx context.Context, chainName string, keyName string) (*big.Int, error)
}
