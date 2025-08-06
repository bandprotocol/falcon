package cmd

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/store"
)

// AppCreator is an object that provide a function that allows to
// lazily initialize an application that being used in the falcon command.
type AppCreator struct{}

// NewApp creates a new App instance that being used in the falcon command.
func (a *AppCreator) NewApp(
	log *zap.Logger,
	store store.Store,
	opts relayer.AppOptions,
) (relayer.Application, error) {
	passphrase := ""
	passphraseOpt := opts.Get(EnvPassphrase)
	if passphraseOpt != nil {
		var ok bool
		passphrase, ok = passphraseOpt.(string)
		if !ok {
			return nil, fmt.Errorf("passphrase require a string or empty string")
		}
	}

	cfg, err := store.GetConfig()
	if err != nil {
		return nil, err
	}

	app := relayer.NewApp(log, cfg, passphrase, store)
	return app, nil
}
