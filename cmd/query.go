package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
)

// QueryCmd represents the command for querying data from source and destination chains.
func QueryCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Query commands on source and destination chains",
	}

	cmd.AddCommand(
		queryTunnelCmd(appCreator, defaultHome),
		queryPacketCmd(appCreator, defaultHome),
		queryBalanceCmd(appCreator, defaultHome),
	)

	return cmd
}

// queryTunnelCmd returns a command that query tunnel information.
func queryTunnelCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tunnel [tunnel_id]",
		Aliases: []string{"t"},
		Short:   "Query commands on tunnel data",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: "query tunnel 1",
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

			tunnel, err := app.QueryTunnelInfo(cmd.Context(), tunnelID)
			if err != nil {
				return err
			}

			out, err := json.MarshalIndent(tunnel, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}

	return cmd
}

// queryPacketCmd returns a command that query packet information.
func queryPacketCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "packet [tunnel_id] [sequence]",
		Aliases: []string{"p"},
		Short:   "Query commands on packet data",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: "query packet 1 1",
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

			sequence, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			packet, err := app.QueryTunnelPacketInfo(cmd.Context(), tunnelID, sequence)
			if err != nil {
				return err
			}

			out, err := json.MarshalIndent(packet, "", "  ")
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}

	return cmd
}

// queryBalanceCmd returns a command that query balance of the given account.
func queryBalanceCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance [chain_name] [key_name]",
		Aliases: []string{"b"},
		Short:   "Query commands on account balance",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: "query balance eth test-key",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			if err := app.Init(cmd.Context()); err != nil {
				return err
			}

			chainName := args[0]
			keyName := args[1]
			bal, err := app.QueryBalance(cmd.Context(), chainName, keyName)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), bal.String())
			return nil
		},
	}

	return cmd
}
