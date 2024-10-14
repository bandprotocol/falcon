package falcon

import (
	"context"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/falcon/band"
	bandtypes "github.com/bandprotocol/falcon/falcon/band/types"
	"github.com/bandprotocol/falcon/falcon/chains"
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

// loadConfigFile reads config file into a.Config if file is present.
func (a *App) LoadConfigFile(ctx context.Context) error {
	return nil
}

// InitConfigFile initializes the configuration to the given path.
func (a *App) InitConfigFile(homePath string) error {
	return nil
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
