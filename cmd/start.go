package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/falcon"
)

// startCmd represents the start command
func startCmd(app *falcon.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start [tunnel_id...]",
		Aliases: []string{"st"},
		Short:   "Start the falcon tunnel relayer system.",
		Args:    withUsage(cobra.MinimumNArgs(0)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s start           # start relaying data from every tunnel being registered on source chain.
$ %s start 1 12      # start relaying data from specific tunnelIDs.`, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = app
			return nil
		},
	}

	return cmd
}
