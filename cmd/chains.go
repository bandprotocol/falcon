package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
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
			_ = app
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
		Short:   "Return a chain's configuration data",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s ch s eth
$ %s chains show eth --yaml`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = app
			return nil
		},
	}

	return cmd
}

// chainsAddCmd returns a command that adds a new chain configuration.
func chainsAddCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [chain_type] [chain_name]",
		Aliases: []string{"a"},
		Short:   "Add a new chain to the configuration file by passing a configuration file (-f)",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s ch a evm eth --file chains/eth.toml 
$ %s chains add evm eth --file chains/eth.toml`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = app
			return nil
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
			_ = app
			return nil
		},
	}

	return cmd
}
