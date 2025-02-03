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

			enableMetricsServerFlag, err := cmd.Flags().GetBool(flagEnableMetricsServer)
			if err != nil {
				return err
			}

			metricsListenAddrFlag, err := cmd.Flags().GetString(flagMetricsListenAddr)
			if err != nil {
				return err
			}

			return app.Start(cmd.Context(), tunnelIDs, enableMetricsServerFlag, metricsListenAddrFlag)
		},
	}
	cmd.Flags().Bool(
		flagEnableMetricsServer,
		false,
		"Enables the metrics server. By default, this is determined by the enable-metrics-server parameter in the global config. "+
			"Use this flag to override the config setting.",
	)
	cmd.Flags().String(
		flagMetricsListenAddr,
		"",
		"address to use for metrics server. By default, "+
			"will be the metrics-listen-addr parameter in the global config. "+
			"Make sure to enable metrics server using --enable-metrics-server flag.",
	)
	return cmd
}
