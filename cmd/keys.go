package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
	chainstypes "github.com/bandprotocol/falcon/relayer/chains/types"
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

// AddKeyInput is the input for adding a key to the keychain.
type AddKeyInput struct {
	PrivateKey string
	Mnemonic   string
	CoinType   uint64
	Account    uint64
	Index      uint64
}

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

			input := &AddKeyInput{}
			input.Mnemonic, err = cmd.Flags().GetString(flagMnemonic)
			if err != nil {
				return err
			}

			input.PrivateKey, err = cmd.Flags().GetString(flagPrivateKey)
			if err != nil {
				return err
			}

			input.CoinType, err = cmd.Flags().GetUint64(flagCoinType)
			if err != nil {
				return err
			}

			input.Account, err = cmd.Flags().GetUint64(flagWalletAccount)
			if err != nil {
				return err
			}

			input.Index, err = cmd.Flags().GetUint64(flagWalletIndex)
			if err != nil {
				return err
			}

			if input.PrivateKey == "" && input.Mnemonic == "" {
				input, err = showHuhPrompt()
				if err != nil {
					return err
				}
			}

			var key *chainstypes.Key
			if input.PrivateKey != "" {
				key, err = app.AddKeyByPrivateKey(chainName, keyName, input.PrivateKey)
				if err != nil {
					return err
				}
			} else {
				key, err = app.AddKeyByMnemonic(
					chainName, keyName,
					input.Mnemonic,
					uint32(input.CoinType),
					uint(input.Account),
					uint(input.Index),
				)
				if err != nil {
					return err
				}
			}

			out, err := json.MarshalIndent(key, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}

	cmd.Flags().String(flagPrivateKey, "", "add key with the given private key")
	cmd.Flags().String(flagMnemonic, "", "add key with the given mnemonic")
	cmd.Flags().Uint64(flagCoinType, defaultCoinType, "coin type number for HD derivation")
	cmd.Flags().Uint64(flagWalletAccount, 0, "account number in the HD derivation path")
	cmd.Flags().
		Uint64(flagWalletIndex, 0, "index number for the specific address within an account in the HD derivation path")

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

// showHuhPrompt shows a prompt to the user to input a private key, mnemonic for generating or
// inserting a user's key.
func showHuhPrompt() (input *AddKeyInput, err error) {
	input = &AddKeyInput{}
	var coinTypeStr, accountStr, indexStr string

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
		return nil, err
	}

	// Coin type input
	coinTypeInput := huh.NewInput().
		Title("Enter a coin type").
		Description("Coin type number for HD derivation (default: 60; leave empty to use default)").
		Value(&coinTypeStr).Validate(
		func(s string) (err error) {
			if s == "" {
				input.CoinType = defaultCoinType
				return nil
			}

			input.CoinType, err = strconv.ParseUint(s, 10, 32)
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
		func(s string) (err error) {
			if s == "" {
				input.Account = 0
				return nil
			}

			input.Account, err = strconv.ParseUint(s, 10, 32)
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
		func(s string) (err error) {
			if s == "" {
				input.Index = 0
				return nil
			}

			input.Index, err = strconv.ParseUint(s, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid index input (should be uint32)")
			}

			return nil
		},
	)

	// Handle the selected option
	switch selection {
	case privateKeyResult:
		privateKeyPrompt := huh.NewGroup(huh.NewInput().
			Title("Enter your private key").
			Value(&input.PrivateKey))

		form := huh.NewForm(privateKeyPrompt)
		if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
			return nil, err
		}

	case mnemonicResult:
		mnemonicPrompt := huh.NewGroup(huh.NewInput().
			Title("Enter your mnemonic").
			Value(&input.Mnemonic),
			coinTypeInput,
			accountInput,
			indexInput,
		)

		form := huh.NewForm(mnemonicPrompt)
		if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
			return nil, err
		}
	case defaultResult:
		defaultPrompt := huh.NewGroup(coinTypeInput, accountInput, indexInput)
		form := huh.NewForm(defaultPrompt)
		if err := form.WithTheme(huh.ThemeBase()).Run(); err != nil {
			return nil, err
		}
	}

	return input, nil
}
