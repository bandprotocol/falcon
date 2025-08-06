package cmd

import (
	"fmt"

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
	passphrase, err := getOptString(opts, EnvPassphrase)
	if err != nil {
		return nil, err
	}

	cfg, err := store.GetConfig()
	if err != nil {
		return nil, err
	}

	logLevel, err := getOptString(opts, flagLogLevel)
	if err != nil {
		return nil, err
	}

	logFormat, err := getOptString(opts, flagLogFormat)
	if err != nil {
		return nil, err
	}

	log, err := initLogger(logLevel, logFormat)
	if err != nil {
		return nil, err
	}

	logWrapper := logger.NewZapLogWrapper(log)
	app := relayer.NewApp(logWrapper, cfg, passphrase, store)
	return app, nil
}

func getOptString(opts relayer.AppOptions, key string) (string, error) {
	res := ""
	opt := opts.Get(key)
	if opt != nil {
		var ok bool
		res, ok = opt.(string)
		if !ok {
			return "", fmt.Errorf("%s require a string or empty string", key)
		}
	}

	return res, nil
}
