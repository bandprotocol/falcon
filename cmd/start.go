package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer"
)

// startCmd represents the start command
func startCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start [tunnel_id...]",
		Aliases: []string{"st"},
		Short:   "Start the falcon tunnel relayer program",
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

			metricsListenAddr, err := cmd.Flags().GetString(flagMetricsListenAddr)
			if err != nil {
				return err
			}

			// setup metrics server
			if metricsListenAddr == "" {
				metricsListenAddr = app.Config.Global.MetricsListenAddr
			}
			if metricsListenAddr != "" {
				if err := relayermetrics.StartMetricsServer(cmd.Context(), app.Log, metricsListenAddr); err != nil {
					return err
				}
			}

			return app.Start(cmd.Context(), tunnelIDs)
		},
	}

	cmd.Flags().String(
		flagMetricsListenAddr,
		"",
		"address to use for metrics server. By default, "+
			"will be the metrics-listen-addr parameter in the global config. ",
	)
	return cmd
}
