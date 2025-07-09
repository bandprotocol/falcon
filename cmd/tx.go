package cmd

import (
	"fmt"
	"strconv"
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
		Short:   "Relay the next sequence message to the destination tunnel contract address",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s tx relay 1, 		# relay tunnelID 1's pending packets
$ %s tx relay 1 --force	# relay tunnelID 1's pending packets regardless of its active status on BandChain`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			tunnelID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			isForce, err := cmd.Flags().GetBool(flagForce)
			if err != nil {
				return err
			}

			return app.Relay(cmd.Context(), tunnelID, isForce)
		},
	}

	cmd.Flags().Bool(flagForce, false, "force relaying data from specific tunnelIDs")

	return cmd
}
