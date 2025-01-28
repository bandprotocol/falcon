package cmd

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/chains/evm"
)

// chainsCmd returns a command that manages chain configurations.
func chainsCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chains",
		Aliases: []string{"ch"},
		Short:   "Manage chain configurations",
	}

	cmd.AddCommand(
		chainsAddCmd(app),
		chainsDeleteCmd(app),
		chainsListCmd(app),
		chainsShowCmd(app),
	)

	return cmd
}

// chainsAddCmd returns a command that adds a new chain configuration.
func chainsAddCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [chain_name] [chain_config_path]",
		Aliases: []string{"a"},
		Short:   "Add a new chain to the configuration file by passing a configuration file",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s ch a evm chains/eth.toml 
$ %s chains add evm chains/eth.toml`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainName := args[0]
			filePath := args[1]

			return app.AddChainConfig(chainName, filePath)
		},
	}

	return cmd
}

// chainsDeleteCmd returns a command that deletes a chain configuration.
func chainsDeleteCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [chain_name]",
		Aliases: []string{"d"},
		Short:   "Remove chain from config based on chain_name",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s chains delete eth
$ %s ch d eth`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainName := args[0]

			return app.DeleteChainConfig(chainName)
		},
	}

	return cmd
}

// chainsListCmd returns a command that lists all chains that are currently configured.
// NOTE: this command only show the list of registered chainIDs, to see the configuration
// please see `chains show`.
func chainsListCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "Return list of chain names that are currently configured",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s chains list
$ %s ch l`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if app.Config == nil {
				return relayer.ErrConfigNotExist(app.HomePath)
			}

			i := 1
			for chainName, chainProviderConfig := range app.Config.TargetChains {
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
func chainsShowCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [chain_name]",
		Aliases: []string{"s"},
		Short:   "Return chain's configuration data",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s ch s eth
$ %s chains show eth --yaml`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
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
