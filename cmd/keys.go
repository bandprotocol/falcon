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
	PrivateKey   string
	Mnemonic     string
	CoinType     uint64
	Account      uint64
	Index        uint64
	RemoteSigner RemoteSignerInput
}

// RemoteSignerInput is the input that holds the parameters needed to configure a remote signer.
type RemoteSignerInput struct {
	Address string
	Url     string
	Key     *string
}

// KeysCmd represents the keys command
func KeysCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "Manage keys held by the relayer for each chain",
	}

	cmd.AddCommand(
		keysAddCmd(appCreator, defaultHome),
		keysDeleteCmd(appCreator, defaultHome),
		keysListCmd(appCreator, defaultHome),
		keysExportCmd(appCreator, defaultHome),
		keysShowCmd(appCreator, defaultHome),
	)

	return cmd
}

// keysAddCmd returns a command that adds a key to the keychain.
func keysAddCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [chain_name] [key_name]",
		Aliases: []string{"a"},
		Short:   "Add a key to the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(`
k add eth test-key
keys add eth test-key`),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			chainName := args[0]
			keyName := args[1]

			if err := app.InitTargetChain(chainName); err != nil {
				return err
			}

			input, err := parseKeysAddInputFromFlag(cmd)
			if err != nil {
				return err
			}

			// if no private key, mnemonic, or remote signer info is provided, prompt interactively
			if input.PrivateKey == "" && input.Mnemonic == "" && input.RemoteSigner.Address == "" &&
				input.RemoteSigner.Url == "" {
				input, err = showHuhPrompt()
				if err != nil {
					return err
				}
			}

			if err := validateAddKeyInput(input); err != nil {
				return err
			}

			key, err := addKey(app, chainName, keyName, input)
			if err != nil {
				return err
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

	cmd.Flags().String(flagRemoteAddress, "", "address of the remote signer key")
	cmd.Flags().String(flagRemoteUrl, "", "URL endpoint of the kms service")
	cmd.Flags().String(flagRemoteKey, "", "key for authenticating with the remote KMS signer service")

	return cmd
}

// keysDeleteCmd returns a command that delete the key from the keychain.
func keysDeleteCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [chain_name] [key_name]",
		Aliases: []string{"d"},
		Short:   "Delete a key from the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(`
k delete eth test-key
keys delete eth test-key`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			chainName := args[0]
			keyName := args[1]

			if err := app.InitTargetChain(chainName); err != nil {
				return err
			}

			return app.DeleteKey(chainName, keyName)
		},
	}

	return cmd
}

// keysListCmd returns a command that list keys associated with a particular chain from the keychain.
func keysListCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [chain_name]",
		Aliases: []string{"l"},
		Short:   "List keys from the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(`
k list eth
keys list eth`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			chainName := args[0]
			if err := app.InitTargetChain(chainName); err != nil {
				return err
			}

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
func keysExportCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export [chain_name] [key_name]",
		Aliases: []string{"e"},
		Short:   "Export a private key from the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(`
k export eth test-key
keys export eth test-key`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			chainName := args[0]
			keyName := args[1]

			if err := app.InitTargetChain(chainName); err != nil {
				return err
			}

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
func keysShowCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [chain_name] [key_name]",
		Aliases: []string{"s"},
		Short:   "Show a key from the keychain associated with a particular chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(`
k show eth test-key
keys show eth test-key`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			chainName := args[0]
			keyName := args[1]

			if err := app.InitTargetChain(chainName); err != nil {
				return err
			}

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

// validateAddKeyInput checks that the AddKeyInput is valid.
func validateAddKeyInput(input *AddKeyInput) error {
	hasPrivateKey := input.PrivateKey != ""
	hasMnemonic := input.Mnemonic != ""
	hasRemoteSigner := input.RemoteSigner.Address != "" || input.RemoteSigner.Url != ""

	// if a private key is provided, no other input should be present
	if hasPrivateKey && (hasMnemonic || hasRemoteSigner) {
		return fmt.Errorf("private key cannot be provided with mnemonic or remote signer")
	}

	// if a mnemonic is provided, no other input should be present
	if hasMnemonic && (hasPrivateKey || hasRemoteSigner) {
		return fmt.Errorf("mnemonic cannot be provided with private key or remote signer")
	}

	// if any remote-signer field is provided, it must be the only input
	if hasRemoteSigner {
		// both address and URL cannot be empty
		if input.RemoteSigner.Address == "" {
			return fmt.Errorf("remote signer address cannot be empty")
		}
		if input.RemoteSigner.Url == "" {
			return fmt.Errorf("remote signer URL cannot be empty")
		}
	}

	return nil
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

// parseKeysAddInputFromFlag parses the key addition input from the command line flags.
func parseKeysAddInputFromFlag(cmd *cobra.Command) (*AddKeyInput, error) {
	input := &AddKeyInput{}

	var err error
	input.Mnemonic, err = cmd.Flags().GetString(flagMnemonic)
	if err != nil {
		return nil, err
	}

	input.PrivateKey, err = cmd.Flags().GetString(flagPrivateKey)
	if err != nil {
		return nil, err
	}

	input.CoinType, err = cmd.Flags().GetUint64(flagCoinType)
	if err != nil {
		return nil, err
	}

	input.Account, err = cmd.Flags().GetUint64(flagWalletAccount)
	if err != nil {
		return nil, err
	}

	input.Index, err = cmd.Flags().GetUint64(flagWalletIndex)
	if err != nil {
		return nil, err
	}

	input.RemoteSigner.Address, err = cmd.Flags().GetString(flagRemoteAddress)
	if err != nil {
		return nil, err
	}

	input.RemoteSigner.Url, err = cmd.Flags().GetString(flagRemoteUrl)
	if err != nil {
		return nil, err
	}

	if cmd.Flags().Changed(flagRemoteKey) {
		remoteSignerKey, err := cmd.Flags().GetString(flagRemoteKey)
		if err != nil {
			return nil, err
		}
		input.RemoteSigner.Key = &remoteSignerKey
	}

	return input, nil
}

// addKey adds a new key to the keychain.
func addKey(
	app relayer.Application,
	chainName string,
	keyName string,
	input *AddKeyInput,
) (*chainstypes.Key, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid input: input is nil")
	}

	// Add key to the keychain
	if input.PrivateKey != "" {
		address, err := app.AddKeyByPrivateKey(chainName, keyName, input.PrivateKey)
		if err != nil {
			return nil, err
		}
		return chainstypes.NewKey("", address, ""), nil
	} else if input.RemoteSigner.Address != "" && input.RemoteSigner.Url != "" {
		if err := app.AddRemoteSignerKey(
			chainName,
			keyName,
			input.RemoteSigner.Address,
			input.RemoteSigner.Url,
			input.RemoteSigner.Key,
		); err != nil {
			return nil, err
		}
		return chainstypes.NewKey("", input.RemoteSigner.Address, ""), nil
	} else {
		mnemonic, address, err := app.AddKeyByMnemonic(
			chainName, keyName,
			input.Mnemonic,
			uint32(input.CoinType),
			uint(input.Account),
			uint(input.Index),
		)
		if err != nil {
			return nil, err
		}
		return chainstypes.NewKey(mnemonic, address, ""), nil
	}
}
