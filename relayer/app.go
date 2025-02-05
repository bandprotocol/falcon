package relayer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal"
	"github.com/bandprotocol/falcon/relayer/band"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/types"
)

const (
	ConfigFolderName   = "config"
	ConfigFileName     = "config.toml"
	PassphraseFileName = "passphrase.hash"
	PassphraseEnvKey   = "PASSPHRASE"
)

// App is the main application struct.
type App struct {
	Log      *zap.Logger
	HomePath string
	Debug    bool
	Config   *Config

	TargetChains  chains.ChainProviders
	BandClient    band.Client
	EnvPassphrase string
}

// NewApp creates a new App instance.
func NewApp(
	log *zap.Logger,
	homePath string,
	debug bool,
	config *Config,
) *App {
	app := App{
		Log:      log,
		HomePath: homePath,
		Debug:    debug,
		Config:   config,
	}
	return &app
}

// Init initialize the application.
func (a *App) Init(ctx context.Context, logLevel, logFormat string) error {
	if a.Config == nil {
		if err := a.LoadConfigFile(); err != nil {
			return err
		}
	}

	// initialize logger, if not already initialized
	if a.Log == nil {
		if err := a.initLogger(logLevel, logFormat); err != nil {
			return err
		}
	}

	// load passphrase from .env file or system environment variables
	a.EnvPassphrase = a.loadEnvPassphrase()

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

// initLogger initializes the logger with the given log level.
func (a *App) initLogger(logLevel, logFormat string) error {
	if logLevel == "" && a.Config != nil {
		logLevel = a.Config.Global.LogLevel
	}

	// initialize logger only if user run command "start" or log level is "debug"
	if os.Args[1] == "start" || logLevel == "debug" {
		log, err := newRootLogger(logFormat, logLevel)
		if err != nil {
			return err
		}
		a.Log = log
	} else {
		a.Log = zap.NewNop()
	}

	return nil
}

// initTargetChains initializes the target chains.
func (a *App) initTargetChains() error {
	a.TargetChains = make(chains.ChainProviders)

	for chainName, chainConfig := range a.Config.TargetChains {
		cp, err := chainConfig.NewChainProvider(chainName, a.Log, a.HomePath, a.Debug)
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

// LoadConfigFile reads config file into a.Config if file is present.
func (a *App) LoadConfigFile() error {
	cfgPath := path.Join(a.HomePath, ConfigFolderName, ConfigFileName)

	// check if file doesn't exist, exit the function as the config may not be initialized.
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	// read the config from config path
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	// save configuration
	a.Config = cfg

	return nil
}

// InitConfigFile initializes the configuration to the given path.
func (a *App) InitConfigFile(homePath string, customFilePath string) error {
	cfgDir := path.Join(homePath, ConfigFolderName)
	cfgPath := path.Join(cfgDir, ConfigFileName)

	// check if the config file already exists
	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	if _, err := os.Stat(cfgPath); err == nil {
		return ErrConfigExist(cfgPath)
	} else if !os.IsNotExist(err) {
		return err
	}

	// Load config from given custom file path if exists
	var cfg *Config
	var err error
	switch {
	case customFilePath != "":
		cfg, err = LoadConfig(customFilePath) // Initialize with CustomConfig if file is provided
		if err != nil {
			return fmt.Errorf("LoadConfig file %v error %v", customFilePath, err)
		}
	default:
		cfg = DefaultConfig() // Initialize with DefaultConfig if no file is provided
	}

	// Marshal config object into bytes
	b, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	// Create the home folder if doesn't exist
	if err := internal.CheckAndCreateFolder(homePath); err != nil {
		return err
	}

	// Create the config folder if doesn't exist
	if err := internal.CheckAndCreateFolder(cfgDir); err != nil {
		return err
	}

	// Create the file and write the default config to the given location.
	f, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(b); err != nil {
		return err
	}

	return nil
}

// InitPassphrase hashes the provided passphrase and saves it to the given path.
func (a *App) InitPassphrase() error {
	// Load and hash the passphrase
	h := sha256.New()
	h.Write([]byte(a.EnvPassphrase))
	b := h.Sum(nil)

	cfgDir := path.Join(a.HomePath, ConfigFolderName)
	passphrasePath := path.Join(cfgDir, PassphraseFileName)

	// Create the file and write the hashed passphrase to the given location.
	f, err := os.Create(passphrasePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(b); err != nil {
		return err
	}

	return nil
}

// QueryTunnelInfo queries tunnel information by given tunnel ID
func (a *App) QueryTunnelInfo(ctx context.Context, tunnelID uint64) (*types.Tunnel, error) {
	if a.Config == nil {
		return nil, ErrConfigNotExist(a.HomePath)
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
		return nil, ErrConfigNotExist(a.HomePath)
	}

	return a.BandClient.GetTunnelPacket(ctx, tunnelID, sequence)
}

// AddChainConfig adds a new chain configuration to the config file.
func (a *App) AddChainConfig(chainName string, filePath string) error {
	if a.Config == nil {
		return ErrConfigNotExist(a.HomePath)
	}

	if _, ok := a.Config.TargetChains[chainName]; ok {
		return ErrChainNameExist(chainName)
	}

	chainProviderConfig, err := LoadChainConfig(filePath)
	if err != nil {
		return err
	}

	a.Config.TargetChains[chainName] = chainProviderConfig

	cfgDir := path.Join(a.HomePath, ConfigFolderName)
	cfgPath := path.Join(cfgDir, ConfigFileName)

	// Marshal config object into bytes
	b, err := toml.Marshal(a.Config)
	if err != nil {
		return err
	}

	return os.WriteFile(cfgPath, b, 0o600)
}

// DeleteChainConfig deletes the chain configuration from the config file.
func (a *App) DeleteChainConfig(chainName string) error {
	if a.Config == nil {
		return ErrConfigNotExist(a.HomePath)
	}

	if _, ok := a.Config.TargetChains[chainName]; !ok {
		return ErrChainNameNotExist(chainName)
	}

	delete(a.Config.TargetChains, chainName)

	cfgDir := path.Join(a.HomePath, ConfigFolderName)
	cfgPath := path.Join(cfgDir, ConfigFileName)

	// Marshal config object into bytes
	b, err := toml.Marshal(a.Config)
	if err != nil {
		return err
	}

	return os.WriteFile(cfgPath, b, 0o600)
}

// GetChainConfig retrieves the chain configuration by given chain name.
func (a *App) GetChainConfig(chainName string) (chains.ChainProviderConfig, error) {
	if a.Config == nil {
		return nil, ErrConfigNotExist(a.HomePath)
	}

	chainProviders := a.Config.TargetChains

	if _, ok := chainProviders[chainName]; !ok {
		return nil, ErrChainNameNotExist(chainName)
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
		return nil, ErrConfigNotExist(a.HomePath)
	}

	if err := a.ValidatePassphrase(a.EnvPassphrase); err != nil {
		return nil, err
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return nil, ErrChainNameNotExist(chainName)
	}

	keyOutput, err := cp.AddKey(keyName, mnemonic, privateKey, a.HomePath, coinType, account, index, a.EnvPassphrase)
	if err != nil {
		return nil, err
	}

	return keyOutput, nil
}

// DeleteKey deletes the key from the chain provider.
func (a *App) DeleteKey(chainName string, keyName string) error {
	if a.Config == nil {
		return ErrConfigNotExist(a.HomePath)
	}

	if err := a.ValidatePassphrase(a.EnvPassphrase); err != nil {
		return err
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return ErrChainNameNotExist(chainName)
	}

	return cp.DeleteKey(a.HomePath, keyName, a.EnvPassphrase)
}

// ExportKey exports the private key from the chain provider.
func (a *App) ExportKey(chainName string, keyName string) (string, error) {
	if a.Config == nil {
		return "", ErrConfigNotExist(a.HomePath)
	}

	if err := a.ValidatePassphrase(a.EnvPassphrase); err != nil {
		return "", err
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return "", ErrChainNameNotExist(chainName)
	}

	privateKey, err := cp.ExportPrivateKey(keyName, a.EnvPassphrase)
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
		return "", ErrConfigNotExist(a.HomePath)
	}

	cp, exist := a.TargetChains[chainName]
	if !exist {
		return "", ErrChainNameNotExist(chainName)
	}

	return cp.ShowKey(keyName)
}

// QueryBalance retrieves the balance of the key from the chain provider.
func (a *App) QueryBalance(ctx context.Context, chainName string, keyName string) (*big.Int, error) {
	if a.Config == nil {
		return nil, ErrConfigNotExist(a.HomePath)
	}

	cp, exist := a.TargetChains[chainName]

	if !exist {
		return nil, ErrChainNameNotExist(chainName)
	}

	return cp.QueryBalance(ctx, keyName)
}

// loadEnvPassphrase retrieves the passphrase string from the .env file or system environment variables.
// It first attempts to load the .env file. If the file is not found or cannot be loaded,
// it falls back to retrieving the "PASSPHRASE" variable from the system environment variables.
func (a *App) loadEnvPassphrase() string {
	// load passphrase from .env first. if not present, use env variable from command
	if err := godotenv.Load(); err != nil {
		a.Log.Debug(
			".env file not found, attempting to use system environment variables",
			zap.Error(err),
		)
	} else {
		a.Log.Debug("Loaded .env file successfully, attempting to use variable from .env file")
	}
	return os.Getenv(PassphraseEnvKey)
}

// ValidatePassphrase checks if the provided passphrase (from the environment)
// matches the hashed passphrase stored on disk.
func (a *App) ValidatePassphrase(envPassphrase string) error {
	// prepare bytes slices of hashed env passphrase
	h := sha256.New()
	h.Write([]byte(envPassphrase))
	envb := h.Sum(nil)

	// load passphrase from local disk
	cfgDir := path.Join(a.HomePath, ConfigFolderName)
	passphrasePath := path.Join(cfgDir, PassphraseFileName)

	b, err := os.ReadFile(passphrasePath)
	if err != nil {
		return err
	}

	if !bytes.Equal(envb, b) {
		return fmt.Errorf("invalid passphrase: the provided passphrase does not match the stored passphrase")
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
	if err := a.ValidatePassphrase(a.EnvPassphrase); err != nil {
		return err
	}

	// initialize target chain providers
	for chainName, chainProvider := range a.TargetChains {
		if err := chainProvider.LoadFreeSenders(a.HomePath, a.EnvPassphrase); err != nil {
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
			return ErrChainNameNotExist(tunnel.TargetChainID)
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

	if err := a.ValidatePassphrase(a.EnvPassphrase); err != nil {
		return err
	}

	chainProvider, ok := a.TargetChains[tunnel.TargetChainID]
	if !ok {
		return ErrChainNameNotExist(tunnel.TargetChainID)
	}

	if err := chainProvider.LoadFreeSenders(a.HomePath, a.EnvPassphrase); err != nil {
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
