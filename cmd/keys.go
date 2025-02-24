package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
	"github.com/bandprotocol/falcon/relayer/chains/types"
)

const (
	privateKeyLabel = "Private key (provide an existing private key)"
	mnemonicLabel   = "Mnemonic (recover from an existing mnemonic phrase)"
	defaultLabel    = "Generate new address (no private key or mnemonic needed)"
)

const (
	privateKeyResult = iota
	mnemonicResult
	defaultResult
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			chainName := args[0]
			keyName := args[1]
			mnemonic := ""
			privateKey := ""

			var (
				coinType, account, index          uint64
				coinTypeStr, accountStr, indexStr string
			)

			// Use huh to create a form for user input
			selection := 0
			selectionPrompt := huh.NewGroup(huh.NewSelect[int]().
				Title("Choose how to add a key").
				Options(
					huh.NewOption(privateKeyLabel, privateKeyResult),
					huh.NewOption(mnemonicLabel, mnemonicResult),
					huh.NewOption(defaultLabel, defaultResult),
				).
				Value(&selection))

			form := huh.NewForm(selectionPrompt)
			if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
				return err
			}

			// Coin type input
			coinTypeInput := huh.NewInput().
				Title("Enter a coin type").
				Description("Coin type number for HD derivation (default: 60; leave empty to use default)").
				Value(&coinTypeStr).Validate(
				func(s string) error {
					if s == "" {
						coinType = defaultCoinType
						return nil
					}
					var err error
					coinType, err = strconv.ParseUint(s, 10, 32)
					if err != nil {
						return fmt.Errorf("invalid coin type input (should be uint32)")
					}

					return nil
				},
			)

			// Account type input
			accountInput := huh.NewInput().
				Title("Enter an account").
				Description("Account number in the HD derivation path (default: 0; leave empty to use default)").
				Value(&accountStr).Validate(
				func(s string) error {
					if s == "" {
						account = 0
						return nil
					}
					var err error
					account, err = strconv.ParseUint(s, 10, 32)
					if err != nil {
						return fmt.Errorf("invalid account input (should be uint32)")
					}

					return nil
				},
			)

			// Index type input
			indexInput := huh.NewInput().
				Title("Enter an index").
				Description("Index number for the specific address within an account in the HD derivation path (default: 0; leave empty to use default)").
				Value(&indexStr).Validate(
				func(s string) error {
					if s == "" {
						index = 0
						return nil
					}
					var err error
					index, err = strconv.ParseUint(s, 10, 32)
					if err != nil {
						return fmt.Errorf("invalid index input (should be uint32)")
					}

					return nil
				},
			)

			// Handle the selected option
			var keyOutput *types.Key
			switch selection {
			case privateKeyResult:
				privateKeyPrompt := huh.NewGroup(huh.NewInput().
					Title("Enter your private key").
					Value(&privateKey))

				form := huh.NewForm(privateKeyPrompt)
				if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
					return err
				}

				keyOutput, err = app.AddKeyByPrivateKey(
					chainName,
					keyName,
					privateKey,
				)
				if err != nil {
					return err
				}

			case mnemonicResult:
				mnemonicPrompt := huh.NewGroup(huh.NewInput().
					Title("Enter your mnemonic").
					Value(&mnemonic),
					coinTypeInput,
					accountInput,
					indexInput,
				)

				form := huh.NewForm(mnemonicPrompt)
				if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
					return err
				}

				keyOutput, err = app.AddKeyByMnemonic(
					chainName,
					keyName,
					mnemonic,
					uint32(coinType),
					uint(account),
					uint(index),
				)
				if err != nil {
					return err
				}
			case defaultResult:
				defaultPrompt := huh.NewGroup(coinTypeInput, accountInput, indexInput)
				form := huh.NewForm(defaultPrompt)
				if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
					return err
				}

				keyOutput, err = app.AddKeyByMnemonic(
					chainName,
					keyName,
					mnemonic,
					uint32(coinType),
					uint(account),
					uint(index),
				)
				if err != nil {
					return err
				}
			}

			out, err := json.MarshalIndent(keyOutput, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}

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
