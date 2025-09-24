package cmd

import (
	"github.com/spf13/cast"

	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/logger"
	"github.com/bandprotocol/falcon/relayer/store"
)

// AppCreator is an object that provide a function that allows to
// lazily initialize an application that being used in the falcon command.
type AppCreator struct{}

// NewApp creates a new App instance that being used in the falcon command.
func (a *AppCreator) NewApp(
	store store.Store,
	opts relayer.AppOptions,
) (relayer.Application, error) {
	passphrase := cast.ToString(opts.Get(EnvPassphrase))
	dbPath := cast.ToString(opts.Get(DbPath))

	cfg, err := store.GetConfig()
	if err != nil {
		return nil, err
	}

	logLevel := cast.ToString(opts.Get(FlagLogLevel))
	logFormat := cast.ToString(opts.Get(FlagLogFormat))

	log, err := initLogger(logLevel, logFormat)
	if err != nil {
		return nil, err
	}

	logWrapper := logger.NewZapLogWrapper(log)
	app := relayer.NewApp(logWrapper, cfg, passphrase, dbPath, store)
	return app, nil
}
