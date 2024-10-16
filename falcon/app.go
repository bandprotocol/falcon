package falcon

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/falcon/band"
	bandtypes "github.com/bandprotocol/falcon/falcon/band/types"
	"github.com/bandprotocol/falcon/falcon/chains"
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
}

// NewApp creates a new App instance.
func NewApp(
	log *zap.Logger,
	viper *viper.Viper,
	homePath string,
	debug bool,
	config *Config,
) *App {
	return &App{
		Log:      log,
		Viper:    viper,
		HomePath: homePath,
		Debug:    debug,
		Config:   config,
	}
}

// InitLogger initializes the logger with the given log level.
func (a *App) InitLogger(configLogLevel string) error {
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

func (a *App) configPath() string {
	return path.Join(a.HomePath, "config", "config.toml")
}

// loadConfigFile reads config file into a.Config if file is present.
func (a *App) LoadConfigFile(ctx context.Context) error {
	cfgPath := a.configPath()
	if _, err := os.Stat(cfgPath); err != nil {
		// don't return error if file doesn't exist
		return nil
	}

	// read the config file bytes
	file, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// unmarshall them into the struct
	cfg := &Config{}
	err = toml.Unmarshal(file, cfg)
	if err != nil {
		return fmt.Errorf("error unmarshalling config: %w", err)
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

	var cfg Config
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

	// Create the file and write the default config to the given location.
	f, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	if _, err = f.Write(b); err != nil {
		return err
	}

	return nil
}

// Get config file from given home path.
func (a *App) GetConfigFile(homePath string) (string, error) {
	cfgPath := path.Join(homePath, "config", "config.toml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if _, err := os.Stat(homePath); os.IsNotExist(err) {
			return "", fmt.Errorf("home path does not exist: %s", homePath)
		}
		return "", fmt.Errorf("config does not exist: %s", cfgPath)
	}
	
	out, err := toml.Marshal(a.Config)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Start starts the tunnel relayer program.
func (a *App) Start(ctx context.Context, tunnelIDs []uint64) error {
	// initialize band client
	bandClient := band.NewClient(a.Log, a.Config.BandChainConfig.RpcEndpoints)

	// TODO: initialize target chain clients
	chainClients := make(map[string]chains.Client)

	// TODO: load the tunnel information from the bandchain.
	// If len(tunnelIDs == 0), load all tunnels info.
	tunnels := []*bandtypes.Tunnel{}

	// initialize the tunnel relayer
	tunnelRelayers := []TunnelRelayer{}
	for _, tunnel := range tunnels {
		chainClient := chainClients[tunnel.TargetChainID]

		tr := NewTunnelRelayer(
			a.Log,
			tunnel.ID,
			"",
			a.Config.CheckingPacketInterval,
			bandClient,
			chainClient,
		)
		tunnelRelayers = append(tunnelRelayers, tr)
	}

	// start the tunnel relayers
	scheduler := NewScheduler(a.Log, tunnelRelayers)
	return scheduler.Start(ctx)
}
