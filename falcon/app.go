package falcon

import (
	"context"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

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
