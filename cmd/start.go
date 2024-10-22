package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/relayer"
)

// startCmd represents the start command
func startCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start [tunnel_id...]",
		Aliases: []string{"st"},
		Short:   "Start the falcon tunnel relayer system.",
		Args:    withUsage(cobra.MinimumNArgs(0)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s start           # start relaying data from every tunnel being registered on source chain.
$ %s start 1 12      # start relaying data from specific tunnelIDs.`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			tunnelIDs := []uint64{}
			for _, arg := range args {
				tunnelID, err := strconv.ParseUint(arg, 10, 64)
				if err != nil {
					return err
				}

				tunnelIDs = append(tunnelIDs, tunnelID)
			}

			// TODO: add context to the function so that it
			return app.Start(cmd.Context(), tunnelIDs)
		},
	}

	return cmd
}
