package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
)

// TransactionCmd returns a parent transaction command handler, where all child
// commands can submit transactions on destination chains.
func TransactionCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transact",
		Aliases: []string{"tx"},
		Short:   "Transaction commands to destination chain",
		Long: `Commands to create transactions on destination chains.

Make sure that chains are properly configured to relay over by using the 'chains list' command.`,
	}

	cmd.AddCommand(
		txRelayCmd(appCreator, defaultHome),
	)

	return cmd
}

// txRelayCmd returns a command that relays a specific message to the destination chain.
func txRelayCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "relay [tunnel_id]",
		Aliases: []string{"rly"},
		Short:   "Relay the next sequence message to the destination tunnel contract address",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(`
tx relay 1, 		# relay tunnelID 1's pending packets
tx relay 1 --force	# relay tunnelID 1's pending packets regardless of its active status on BandChain`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			if err := app.Init(cmd.Context()); err != nil {
				return err
			}

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

	registerCommonFlags(cmd, defaultHome)
	cmd.Flags().Bool(flagForce, false, "force relaying data from specific tunnelIDs")

	return cmd
}
