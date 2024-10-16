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
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}
			file, err := cmd.Flags().GetString(flagFile) 
			if err != nil {
				return err
			}
			return app.InitConfigFile(home, file)
		},
	}

	return configInitFlags(app.Viper, cmd)
}

// Command for printing current configuration
func configShowCmd(app *falcon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s", "list", "l"},
		Short:   "Prints current configuration",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config show --home %s
$ %s cfg list`, appName, defaultHome, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}
			out, err := app.GetConfigFile(home)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), out)
			return nil
		},
	}
	return cmd
}
