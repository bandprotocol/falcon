package relayer

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/band"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/types"
)

const (
	configFolderName = "config"
	configFileName   = "config.toml"
)

// App is the main application struct.
type App struct {
	Log      *zap.Logger
	Viper    *viper.Viper
	HomePath string
	Debug    bool
	Config   *Config

	targetChains chains.ChainProviders
}

// NewApp creates a new App instance.
func NewApp(
	log *zap.Logger,
	viper *viper.Viper,
	homePath string,
	debug bool,
	config *Config,
) *App {
	app := App{
		Log:      log,
		Viper:    viper,
		HomePath: homePath,
		Debug:    debug,
		Config:   config,
	}
	return &app
}

// Initialize the application.
func (a *App) Init(ctx context.Context) error {
	if a.Config == nil {
		if err := a.loadConfigFile(); err != nil {
			return err
		}
	}

	// initialize logger, if not already initialized
	if a.Log == nil {
		if err := a.initLogger(""); err != nil {
			return err
		}
	}

	// initialize target chains
	if err := a.initTargetChains(ctx); err != nil {
		return err
	}

	// TODO: initialize band client

	return nil
}

// InitLogger initializes the logger with the given log level.
func (a *App) initLogger(configLogLevel string) error {
	logLevel := a.Viper.GetString("log-level")
	if a.Viper.GetBool("debug") {
		logLevel = "debug"
	} else if logLevel == "" {
		logLevel = configLogLevel
	}

	log, err := newRootLogger(a.Viper.GetString("log-format"), logLevel)
	if err != nil {
		return err
	}

	a.Log = log
	return nil
}

// InitTargetChains initializes the target chains.
func (a *App) initTargetChains(ctx context.Context) error {
	a.targetChains = make(chains.ChainProviders)
	if a.Config == nil || a.Config.TargetChains == nil {
		a.Log.Error("target chains not found in config")
		return nil
	}

	for chainName, chainConfig := range a.Config.TargetChains {
		cp, err := chainConfig.NewChainProvider(chainName, a.Log, a.HomePath, a.Debug)
		if err != nil {
			a.Log.Error("cannot create chain provider",
				zap.Error(err),
				zap.String("chain_name", chainName),
			)
			return err
		}

		if err := cp.Init(ctx); err != nil {
			a.Log.Error("cannot initialize chain provider",
				zap.Error(err),
				zap.String("chain_name", chainName),
			)
			return err
		}
		a.targetChains[chainName] = cp
	}
	return nil
}

// loadConfigFile reads config file into a.Config if file is present.
func (a *App) loadConfigFile() error {
	cfgPath := path.Join(a.HomePath, configFolderName, configFileName)
	if _, err := os.Stat(cfgPath); err != nil {
		// don't return error if file doesn't exist
		return nil
	}

	// read the config from config path
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	if a.Log == nil {
		if err := a.initLogger(cfg.Global.LogLevel); err != nil {
			return err
		}
	}

	// save configuration
	a.Config = cfg

	return nil
}

// InitConfigFile initializes the configuration to the given path.
func (a *App) InitConfigFile(homePath string, customFilePath string) error {
	cfgDir := path.Join(homePath, configFolderName)
	cfgPath := path.Join(cfgDir, configFileName)

	// check if the config file already exists
	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	if _, err := os.Stat(cfgPath); err == nil {
		return fmt.Errorf("config already exists: %s", cfgPath)
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
	if _, err := os.Stat(homePath); os.IsNotExist(err) {
		if err = os.Mkdir(homePath, os.ModePerm); err != nil {
			return err
		}
	}

	// Create the config folder if doesn't exist
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		if err = os.Mkdir(cfgDir, os.ModePerm); err != nil {
			return err
		}
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

func (a *App) QueryTunnelInfo(ctx context.Context, tunnelID uint64) (*types.Tunnel, error) {
	if a.Config == nil {
		return nil, fmt.Errorf("config is not initialized")
	}

	// TODO: add band client part and change targetChain and targetAddr
	// bandClient := band.NewClient(a.Log, a.Config.BandChain.RpcEndpoints)

	targetChain := "testnet_evm"
	targetAddr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"

	var tunnelChainInfo *chainstypes.Tunnel
	cp, ok := a.targetChains[targetChain]
	if ok {
		var err error
		tunnelChainInfo, err = cp.QueryTunnelInfo(ctx, tunnelID, targetAddr)
		if err != nil {
			return nil, err
		}
	}

	return types.NewTunnel(
		tunnelID,
		targetChain,
		targetAddr,
		tunnelChainInfo,
	), nil
}

// Start starts the tunnel relayer program.
func (a *App) Start(ctx context.Context, tunnelIDs []uint64) error {
	a.Log.Info("starting tunnel relayer")

	// initialize band client
	bandClient := band.NewClient(a.Log, a.Config.BandChain.RpcEndpoints)

	// TODO: load the tunnel information from the bandchain.
	// If len(tunnelIDs == 0), load all tunnels info.
	tunnels := []*bandtypes.Tunnel{
		{
			ID:             1,
			LatestSequence: 0,
			TargetChainID:  "testnet_evm",
			TargetAddress:  "0x5FbDB2315678afecb367f032d93F642f64180aa3",
		},
	}

	// initialize the tunnel relayer
	tunnelRelayers := []TunnelRelayer{}
	for _, tunnel := range tunnels {
		chainProvider, ok := a.targetChains[tunnel.TargetChainID]
		if !ok {
			return fmt.Errorf("target chain provider not found: %s", tunnel.TargetChainID)
		}

		tr := NewTunnelRelayer(
			a.Log,
			tunnel.ID,
			tunnel.TargetAddress,
			a.Config.Global.CheckingPacketInterval,
			bandClient,
			chainProvider,
		)
		tunnelRelayers = append(tunnelRelayers, tr)
	}

	// start the tunnel relayers
	scheduler := NewScheduler(a.Log, tunnelRelayers)
	return scheduler.Start(ctx)
}
