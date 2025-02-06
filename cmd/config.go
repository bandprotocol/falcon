package cmd

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
)

// configCmd returns a command that manages global configuration file
func configCmd(app *relayer.App) *cobra.Command {
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
func configShowCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s", "list", "l"},
		Short:   "Prints current configuration",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config show --home %s
$ %s cfg list`, appName, defaultHome, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if app.Config == nil {
				return fmt.Errorf("config is not initialized")
			}

			b, err := toml.Marshal(app.Config)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	}
	return cmd
}

// configInitCmd returns the commands that initializes an empty config at the --home location
func configInitCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Create a default configuration at home directory path",
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

			if err := app.InitConfigFile(home, file); err != nil {
				return err
			}

			return app.InitPassphrase()
		},
	}
	cmd.Flags().StringP(flagFile, "f", "", "fetch toml data from specified file")
	return cmd
}
