package relayer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/band"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/config"
	"github.com/bandprotocol/falcon/relayer/store"
	"github.com/bandprotocol/falcon/relayer/types"
)

const (
	ConfigFolderName   = "config"
	ConfigFileName     = "config.toml"
	PassphraseFileName = "passphrase.hash"
)

// App is the main application struct.
type App struct {
	Log      *zap.Logger
	HomePath string
	Debug    bool
	Config   *config.Config
	Store    store.Store

	TargetChains chains.ChainProviders
	BandClient   band.Client
	Passphrase   string
}

// NewApp creates a new App instance.
func NewApp(
	log *zap.Logger,
	homePath string,
	debug bool,
	config *config.Config,
	passphrase string,
	store store.Store,
) *App {
	app := App{
		Log:        log,
		HomePath:   homePath,
		Debug:      debug,
		Config:     config,
		Store:      store,
		Passphrase: passphrase,
	}
	return &app
}

// Init initialize the application.
func (a *App) Init(ctx context.Context) error {
	// if config is not initialized, return
	if a.Config == nil {
		return nil
	}

	// initialize target chain clients
	if err := a.initTargetChains(); err != nil {
		return err
	}

	// initialize BandChain client
	if err := a.initBandClient(ctx); err != nil {
		return err
	}

	return nil
}

// initBandClient establishes connection to rpc endpoints.
func (a *App) initBandClient(ctx context.Context) error {
	a.BandClient = band.NewClient(nil, a.Log, &a.Config.BandChain)

	// connect to BandChain, if error occurs, log the error as debug and continue
	if err := a.BandClient.Init(ctx); err != nil {
		a.Log.Error("Cannot connect to BandChain", zap.Error(err))
		return err
	}

	return nil
}

// initTargetChains initializes the target chains.
func (a *App) initTargetChains() error {
	a.TargetChains = make(chains.ChainProviders)

	for chainName, chainConfig := range a.Config.TargetChains {
		wallet, err := a.Store.NewWallet(chainConfig.GetChainType(), chainName, a.Passphrase)
		if err != nil {
			a.Log.Error("Wallet registry not found",
				zap.Error(err),
				zap.String("chain_name", chainName),
			)
			return err
		}

		cp, err := chainConfig.NewChainProvider(chainName, a.Log, a.HomePath, a.Debug, wallet)
		if err != nil {
			a.Log.Error("Cannot create chain provider",
				zap.Error(err),
				zap.String("chain_name", chainName),
			)
			return err
		}

		a.TargetChains[chainName] = cp
	}

	return nil
}

// InitConfigFile initializes the configuration to the given path.
func (a *App) InitConfigFile(homePath string, customFilePath string) error {
	// Check if config already exists
	if ok, err := a.Store.HasConfig(); err != nil {
		return err
	} else if ok {
		return fmt.Errorf("config already exists")
	}

	// Load config from given custom file path if exists
	var cfg *config.Config
	switch {
	case customFilePath != "":
		b, err := os.ReadFile(customFilePath)
		if err != nil {
			return fmt.Errorf("cannot read a config file %s: %w", customFilePath, err)
		}

		cfg, err = config.ParseConfig(b)
		if err != nil {
			return fmt.Errorf("parsing config error %w", err)
		}
	default:
		cfg = config.DefaultConfig() // Initialize with DefaultConfig if no file is provided
	}

	return a.Store.SaveConfig(cfg)
}

// InitPassphrase hashes the provided passphrase and saves it to the given path.
func (a *App) InitPassphrase() error {
	// Load and hash the passphrase
	h := sha256.New()
	h.Write([]byte(a.Passphrase))
	hashedPassphrase := h.Sum(nil)

	return a.Store.SaveHashedPassphrase(hashedPassphrase)
}

// QueryTunnelInfo queries tunnel information by given tunnel ID
func (a *App) QueryTunnelInfo(ctx context.Context, tunnelID uint64) (*types.Tunnel, error) {
	if a.Config == nil {
		return nil, fmt.Errorf("config is not initialized")
	}

	tunnel, err := a.BandClient.GetTunnel(ctx, tunnelID)
	if err != nil {
		return nil, err
	}

	bandChainInfo := bandtypes.NewTunnel(
		tunnel.ID,
		tunnel.LatestSequence,
		tunnel.TargetAddress,
		tunnel.TargetChainID,
		tunnel.IsActive,
	)

	cp, ok := a.TargetChains[bandChainInfo.TargetChainID]
	if !ok {
		a.Log.Debug("Target chain provider not found", zap.String("chain_id", bandChainInfo.TargetChainID))
		return types.NewTunnel(bandChainInfo, nil), nil
	}

	tunnelChainInfo, err := cp.QueryTunnelInfo(ctx, tunnelID, bandChainInfo.TargetAddress)
	if err != nil {
		return nil, err
	}

	return types.NewTunnel(
		bandChainInfo,
		tunnelChainInfo,
	), nil
}

// QueryTunnelPacketInfo queries tunnel packet information by given tunnel ID
func (a *App) QueryTunnelPacketInfo(ctx context.Context, tunnelID uint64, sequence uint64) (*bandtypes.Packet, error) {
	if a.Config == nil {
		return nil, fmt.Errorf("config is not initialized")
	}

	return a.BandClient.GetTunnelPacket(ctx, tunnelID, sequence)
}

// AddChainConfig adds a new chain configuration to the config file.
func (a *App) AddChainConfig(chainName string, filePath string) error {
	if a.Config == nil {
		return fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	if _, ok := a.Config.TargetChains[chainName]; ok {
		return fmt.Errorf("existing chain name : %s", chainName)
	}

	chainProviderConfig, err := config.LoadChainConfig(filePath)
	if err != nil {
		return err
	}

	a.Config.TargetChains[chainName] = chainProviderConfig
	return a.Store.SaveConfig(a.Config)
}

// DeleteChainConfig deletes the chain configuration from the config file.
func (a *App) DeleteChainConfig(chainName string) error {
	if a.Config == nil {
		return fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	if _, ok := a.Config.TargetChains[chainName]; !ok {
		return fmt.Errorf("not existing chain name : %s", chainName)
	}

	delete(a.Config.TargetChains, chainName)
	return a.Store.SaveConfig(a.Config)
}

// GetChainConfig retrieves the chain configuration by given chain name.
func (a *App) GetChainConfig(chainName string) (chains.ChainProviderConfig, error) {
	if a.Config == nil {
		return nil, fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	chainProviders := a.Config.TargetChains

	if _, ok := chainProviders[chainName]; !ok {
		return nil, fmt.Errorf("not existing chain name : %s", chainName)
	}

	return chainProviders[chainName], nil
}

// AddKey adds a new key to the chain provider.
func (a *App) AddKey(
	chainName string,
	keyName string,
	mnemonic string,
	privateKey string,
	coinType uint32,
	account uint,
	index uint,
) (*chainstypes.Key, error) {
	if a.Config == nil {
		return nil, fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	if err := a.ValidatePassphrase(a.Passphrase); err != nil {
		return nil, err
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return nil, fmt.Errorf("chain name does not exist: %s", chainName)
	}

	keyOutput, err := cp.AddKey(keyName, mnemonic, privateKey, a.HomePath, coinType, account, index, a.Passphrase)
	if err != nil {
		return nil, err
	}

	return keyOutput, nil
}

// DeleteKey deletes the key from the chain provider.
func (a *App) DeleteKey(chainName string, keyName string) error {
	if a.Config == nil {
		return fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	if err := a.ValidatePassphrase(a.Passphrase); err != nil {
		return err
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return fmt.Errorf("chain name does not exist: %s", chainName)
	}

	return cp.DeleteKey(a.HomePath, keyName, a.Passphrase)
}

// ExportKey exports the private key from the chain provider.
func (a *App) ExportKey(chainName string, keyName string) (string, error) {
	if a.Config == nil {
		return "", fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	if err := a.ValidatePassphrase(a.Passphrase); err != nil {
		return "", err
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return "", fmt.Errorf("chain name does not exist: %s", chainName)
	}

	privateKey, err := cp.ExportPrivateKey(keyName, a.Passphrase)
	if err != nil {
		return "", err
	}

	return privateKey, nil
}

// ListKeys retrieves the list of keys from the chain provider.
func (a *App) ListKeys(chainName string) ([]*chainstypes.Key, error) {
	if a.Config == nil {
		return nil, fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return nil, fmt.Errorf("chain name does not exist: %s", chainName)
	}

	return cp.ListKeys(), nil
}

// ShowKey retrieves the key information from the chain provider.
func (a *App) ShowKey(chainName string, keyName string) (string, error) {
	if a.Config == nil {
		return "", fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return "", fmt.Errorf("chain name does not exist: %s", chainName)
	}

	return cp.ShowKey(keyName)
}

// QueryBalance retrieves the balance of the key from the chain provider.
func (a *App) QueryBalance(ctx context.Context, chainName string, keyName string) (*big.Int, error) {
	if a.Config == nil {
		return nil, fmt.Errorf("config does not exist: %s", a.HomePath)
	}

	cp, exist := a.TargetChains[chainName]

	if !exist {
		return nil, fmt.Errorf("chain name does not exist: %s", chainName)
	}

	if !cp.IsKeyNameExist(keyName) {
		return nil, fmt.Errorf("key name does not exist: %s", chainName)
	}

	return cp.QueryBalance(ctx, keyName)
}

// ValidatePassphrase checks if the provided passphrase (from the environment)
// matches the hashed passphrase stored on disk.
func (a *App) ValidatePassphrase(envPassphrase string) error {
	// prepare bytes slices of hashed env passphrase
	h := sha256.New()
	h.Write([]byte(envPassphrase))
	hashedPassphrase := h.Sum(nil)

	// load passphrase from local disk
	storedHashedPassphrase, err := a.Store.GetHashedPassphrase()
	if err != nil {
		return err
	}

	if !bytes.Equal(hashedPassphrase, storedHashedPassphrase) {
		return fmt.Errorf("invalid passphrase: the provided passphrase does not match the stored hashed passphrase")
	}

	return nil
}

// Start starts the tunnel relayer program.
func (a *App) Start(ctx context.Context, tunnelIDs []uint64) error {
	a.Log.Info("Starting tunnel relayer")

	// query tunnels
	tunnels, err := a.getTunnels(ctx, tunnelIDs)
	if err != nil {
		a.Log.Error("Cannot get tunnels", zap.Error(err))
	}

	// validate passphrase
	if err := a.ValidatePassphrase(a.Passphrase); err != nil {
		return err
	}

	// initialize target chain providers
	for chainName, chainProvider := range a.TargetChains {
		if err := chainProvider.LoadFreeSenders(a.HomePath, a.Passphrase); err != nil {
			a.Log.Error("Cannot load keys in target chain",
				zap.Error(err),
				zap.String("chain_name", chainName),
			)
			return err
		}

		if err := chainProvider.Init(ctx); err != nil {
			a.Log.Error("Cannot initialize chain provider",
				zap.Error(err),
				zap.String("chain_name", chainName),
			)
			return err
		}
	}

	// initialize the tunnel relayer
	tunnelRelayers := []*TunnelRelayer{}
	for _, tunnel := range tunnels {
		chainProvider, ok := a.TargetChains[tunnel.TargetChainID]
		if !ok {
			return fmt.Errorf("target chain provider not found: %s", tunnel.TargetChainID)
		}

		tr := NewTunnelRelayer(
			a.Log,
			tunnel.ID,
			tunnel.TargetAddress,
			a.Config.Global.CheckingPacketInterval,
			a.BandClient,
			chainProvider,
		)
		tunnelRelayers = append(tunnelRelayers, &tr)
	}

	// start the tunnel relayers
	isSyncTunnelsAllowed := (len(tunnelIDs) == 0)
	scheduler := NewScheduler(
		a.Log,
		tunnelRelayers,
		a.Config.Global.CheckingPacketInterval,
		a.Config.Global.SyncTunnelsInterval,
		a.Config.Global.MaxCheckingPacketPenaltyDuration,
		a.Config.Global.PenaltyExponentialFactor,
		isSyncTunnelsAllowed,
		a.BandClient,
		a.TargetChains,
	)

	return scheduler.Start(ctx)
}

// Relay relays the packet from the source chain to the destination chain.
func (a *App) Relay(ctx context.Context, tunnelID uint64) error {
	a.Log.Debug("Query tunnel info on BandChain", zap.Uint64("tunnel_id", tunnelID))
	tunnel, err := a.BandClient.GetTunnel(ctx, tunnelID)
	if err != nil {
		return err
	}

	if err := a.ValidatePassphrase(a.Passphrase); err != nil {
		return err
	}

	chainProvider, ok := a.TargetChains[tunnel.TargetChainID]
	if !ok {
		return fmt.Errorf("target chain provider not found: %s", tunnel.TargetChainID)
	}

	if err := chainProvider.LoadFreeSenders(a.HomePath, a.Passphrase); err != nil {
		a.Log.Error("Cannot load keys in target chain",
			zap.Error(err),
			zap.String("chain_name", tunnel.TargetChainID),
		)
		return err
	}

	tr := NewTunnelRelayer(
		a.Log,
		tunnel.ID,
		tunnel.TargetAddress,
		a.Config.Global.CheckingPacketInterval,
		a.BandClient,
		chainProvider,
	)

	return tr.CheckAndRelay(ctx)
}

// GetTunnels retrieves the list of tunnels by given tunnel IDs. If no tunnel ID is provided,
// get all tunnels
func (a *App) getTunnels(ctx context.Context, tunnelIDs []uint64) ([]bandtypes.Tunnel, error) {
	if len(tunnelIDs) == 0 {
		return a.BandClient.GetTunnels(ctx)
	}

	tunnels := make([]bandtypes.Tunnel, 0, len(tunnelIDs))
	for _, tunnelID := range tunnelIDs {
		tunnel, err := a.BandClient.GetTunnel(ctx, tunnelID)
		if err != nil {
			return nil, err
		}

		tunnels = append(tunnels, *tunnel)
	}

	return tunnels, nil
}
