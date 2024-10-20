package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/falcon"
)

// queryCmd represents the command for querying data from source and destination chains.
func queryCmd(app *falcon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Query commands on source and destination chains.",
	}

	cmd.AddCommand(
		queryTunnelCmd(app),
		queryBalanceCmd(app),
	)

	return cmd
}

// queryTunnelCmd returns a command that query tunnel information.
func queryTunnelCmd(app *falcon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tunnel [tunnel_id]",
		Aliases: []string{"t"},
		Short:   "Query commands on tunnel data",
		Args:    withUsage(cobra.ExactArgs(1)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s query tunnel 1`, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			tunnelID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			tunnel, err := app.QueryTunnelInfo(tunnelID)
			if err != nil {
				return err
			}

			out, err := json.Marshal(tunnel)
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
func queryBalanceCmd(app *falcon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance [chain_name] [key_name]",
		Aliases: []string{"b"},
		Short:   "Query commands on account balance",
		Args:    withUsage(cobra.ExactArgs(2)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s query balance eth test-key
$ %s q b eth test-key`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = app
			return nil
		},
	}

	return cmd
}
