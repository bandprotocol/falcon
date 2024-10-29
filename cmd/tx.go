package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
)

// transactionCmd returns a parent transaction command handler, where all child
// commands can submit transactions on destination chains.
func transactionCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transact",
		Aliases: []string{"tx"},
		Short:   "Transaction commands to destination chain",
		Long: strings.TrimSpace(`Commands to create transactions on destination chains.

Make sure that chains are properly configured to relay over by using the 'falcon chains list' command.`),
	}

	cmd.AddCommand(
		txRelayCmd(app),
	)

	return cmd
}

// txRelayCmd returns a command that relays a specific message to the destination chain.
func txRelayCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "relay [tunnel_id]",
		Aliases: []string{"rly"},
		Short:   "Relay a specific message to the destination chain",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s tx relay 1 1, # relay the message with tunnelID 1 and sequence 1`, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = app
			return nil
		},
	}

	return cmd
}
