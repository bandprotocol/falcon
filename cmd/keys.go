package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
)

const (
	privateKeyLabel = "Private key (provide an existing private key)"
	mnemonicLabel   = "Mnemonic (recover from an existing mnemonic phrase)"
	defaultLabel    = "Generate new address (no private key or mnemonic needed)"
)

// keysCmd represents the keys command
func keysCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "Manage keys held by the relayer for each chain",
	}

	cmd.AddCommand(
		keysAddCmd(app),
		keysDeleteCmd(app),
		keysListCmd(app),
		keysExportCmd(app),
		keysShowCmd(app),
	)

	return cmd
}

// keysAddCmd returns a command that adds a key to the keychain.
func keysAddCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [chain_name] [key_name]",
		Aliases: []string{"a"},
		Short:   "Add a key to the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys add eth test-key
$ %s k a eth test-key`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainName := args[0]
			keyName := args[1]
			mnemonic := ""
			privateKey := ""

			// Use huh to create a form for user input
			selection := ""
			selectionPrompt := huh.NewGroup(huh.NewSelect[string]().
				Title("Choose how to add a key").
				Options(
					huh.NewOption(privateKeyLabel, privateKeyLabel),
					huh.NewOption(mnemonicLabel, mnemonicLabel),
					huh.NewOption(defaultLabel, defaultLabel),
				).
				Value(&selection))

			form := huh.NewForm(selectionPrompt)
			if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
				return err
			}

			// Handle the selected option
			switch selection {
			case privateKeyLabel:
				privateKeyPrompt := huh.NewGroup(huh.NewInput().
					Title("Enter your private key").
					EchoMode(huh.EchoModePassword).
					Value(&privateKey))

				form := huh.NewForm(privateKeyPrompt)
				if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
					return err
				}

			case mnemonicLabel:
				mnemonicPrompt := huh.NewGroup(huh.NewInput().
					Title("Enter your mnemonic").
					EchoMode(huh.EchoModePassword).
					Value(&mnemonic))

				form := huh.NewForm(mnemonicPrompt)
				if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
					return err
				}
			}

			// Get additional flags
			coinType, err := cmd.Flags().GetInt32(flagCoinType)
			if err != nil {
				return err
			}
			if coinType < 0 {
				coinType = defaultCoinType
			}

			account, err := cmd.Flags().GetUint(flagAccount)
			if err != nil {
				return err
			}

			index, err := cmd.Flags().GetUint(flagAccountIndex)
			if err != nil {
				return err
			}

			// Add the key to the app
			keyOutput, err := app.AddKey(chainName, keyName, mnemonic, privateKey, uint32(coinType), account, index)
			if err != nil {
				return err
			}

			out, err := json.MarshalIndent(keyOutput, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}

	// Command flags
	cmd.Flags().Int32(flagCoinType, -1, "Coin type number for HD derivation")
	cmd.Flags().Uint(flagAccount, 0, "Account number within the HD derivation path")
	cmd.Flags().
		Uint(flagAccountIndex, 0, "Index number for the specific address within an account in the HD derivation path")
	return cmd
}

// keysDeleteCmd returns a command that delete the key from the keychain.
func keysDeleteCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [chain_name] [key_name]",
		Aliases: []string{"d"},
		Short:   "Delete a key from the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys delete eth test-key 
$ %s k d eth test-key`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainName := args[0]
			keyName := args[1]
			return app.DeleteKey(chainName, keyName)
		},
	}

	return cmd
}

// keysListCmd returns a command that list keys associated with a particular chain from the keychain.
func keysListCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [chain_name]",
		Aliases: []string{"l"},
		Short:   "List keys from the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys list eth
$ %s k l eth`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainName := args[0]
			keys, err := app.ListKeys(chainName)
			if err != nil {
				return err
			}

			out := "key(%s) -> %s"

			for _, key := range keys {
				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(out, key.KeyName, key.Address))
			}

			return nil
		},
	}

	return cmd
}

// keysExportCmd returns a command that export the private key from the keychain.
func keysExportCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export [chain_name] [key_name]",
		Aliases: []string{"e"},
		Short:   "Export a private key from the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys export eth test-key
$ %s k e eth test-key`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainName := args[0]
			keyName := args[1]

			privateKey, err := app.ExportKey(chainName, keyName)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), privateKey)
			return nil
		},
	}

	return cmd
}

// keysShowCmd a command that show the key information.
func keysShowCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [chain_name] [key_name]",
		Aliases: []string{"s"},
		Short:   "Show a key from the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys show eth test-key
$ %s k s eth test-key`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainName := args[0]
			keyName := args[1]

			address, err := app.ShowKey(chainName, keyName)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), address)

			return nil
		},
	}

	return cmd
}
