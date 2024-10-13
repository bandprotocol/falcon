package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/falcon"
)

// configCmd returns a command that manages global configuration file
func configCmd(app *falcon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Manage global configuration file",
	}

	cmd.AddCommand(
		configShowCmd(app),
		configInitCmd(app),
	)
	return cmd
}

// configShowCmd returns the commands that prints current configuration
func configShowCmd(app *falcon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s", "list", "l"},
		Short:   "Display global configuration",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config show --home %s
$ %s cfg s`, appName, defaultHome, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = app
			return nil
		},
	}

	return cmd
}

// configInitCmd returns the commands that for initializing an empty config at the --home location
func configInitCmd(app *falcon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Create a default configuration at home directory path defined by --home",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config init --home %s
$ %s cfg i`, appName, defaultHome, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = app
			return nil
		},
	}

	return cmd
}
