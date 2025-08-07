package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/config"
)

// ConfigCmd returns a command that manages global configuration file
func ConfigCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Manage global configuration file",
	}

	registerCommonFlags(cmd, defaultHome)

	cmd.AddCommand(
		configShowCmd(appCreator, defaultHome),
		configInitCmd(appCreator, defaultHome),
	)

	return cmd
}

// configShowCmd returns the commands that prints current configuration
func configShowCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s", "list", "l"},
		Short:   "Prints current configuration",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
config show --home %s
cfg list`, defaultHome)),
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

			b, err := toml.Marshal(cfg)
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
func configInitCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Create a default configuration at home directory path",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
config init --home %s
cfg i`, defaultHome)),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, err := cmd.Flags().GetString(flagFile)
			if err != nil {
				return err
			}

			// Parse the config from the file if file's path is given
			var cfg *config.Config
			if filePath != "" {
				b, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("cannot read a config file %s: %w", filePath, err)
				}

				cfg, err = config.ParseConfig(b)
				if err != nil {
					return fmt.Errorf("parsing config error %w", err)
				}
			}

			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			if err := app.Init(cmd.Context()); err != nil {
				return err
			}

			if err := app.SaveConfig(cfg); err != nil {
				return err
			}

			return app.SavePassphrase(app.GetPassphrase())
		},
	}
	cmd.Flags().StringP(flagFile, "f", "", "input config .toml file path")
	return cmd
}
