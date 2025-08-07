package cmd

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
)

// ChainsCmd returns a command that manages chain configurations.
func ChainsCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chains",
		Aliases: []string{"ch"},
		Short:   "Manage chain configurations",
	}

	registerCommonFlags(cmd, defaultHome)

	cmd.AddCommand(
		chainsAddCmd(appCreator, defaultHome),
		chainsDeleteCmd(appCreator, defaultHome),
		chainsListCmd(appCreator, defaultHome),
		chainsShowCmd(appCreator, defaultHome),
	)

	return cmd
}

// chainsAddCmd returns a command that adds a new chain configuration.
func chainsAddCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [chain_name] [chain_config_path]",
		Aliases: []string{"a"},
		Short:   "Add a new chain to the configuration file by passing a configuration file",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(`
ch a evm chains/eth.toml
chains add evm chains/eth.toml`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			if err := app.Init(cmd.Context()); err != nil {
				return err
			}

			chainName := args[0]
			filePath := args[1]

			return app.AddChainConfig(chainName, filePath)
		},
	}
	return cmd
}

// chainsDeleteCmd returns a command that deletes a chain configuration.
func chainsDeleteCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [chain_name]",
		Aliases: []string{"d"},
		Short:   "Remove chain from config based on chain_name",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(`
ch delete eth
chains delete eth`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			if err := app.Init(cmd.Context()); err != nil {
				return err
			}

			chainName := args[0]

			return app.DeleteChainConfig(chainName)
		},
	}

	return cmd
}

// chainsListCmd returns a command that lists all chains that are currently configured.
// NOTE: this command only show the list of registered chainIDs, to see the configuration
// please see `chains show`.
func chainsListCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "Return list of chain names that are currently configured",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(`
ch list
chains list`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			if err := app.Init(cmd.Context()); err != nil {
				return err
			}

			cfg := app.GetConfig()
			if cfg == nil {
				return fmt.Errorf("config is not initialized")
			}

			i := 1
			for chainName, chainProviderConfig := range cfg.TargetChains {
				out := "%d: %s -> type(%s)"
				switch cp := chainProviderConfig.(type) {
				case *evm.EVMChainProviderConfig:
					fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(out, i, chainName, cp.ChainType.String()))
				default:
					return fmt.Errorf("unsupported chain provider type for chain: %s", chainName)
				}
				i++
			}
			return nil
		},
	}

	return cmd
}

// chainsShowCmd returns a command that shows a chain's configuration data.
func chainsShowCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [chain_name]",
		Aliases: []string{"s"},
		Short:   "Return chain's configuration data",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(`
ch s eth
chains show eth --yaml`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			if err := app.Init(cmd.Context()); err != nil {
				return err
			}

			chainName := args[0]

			chainConfig, err := app.GetChainConfig(chainName)
			if err != nil {
				return err
			}

			b, err := toml.Marshal(chainConfig)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(b))

			return nil
		},
	}

	return cmd
}
