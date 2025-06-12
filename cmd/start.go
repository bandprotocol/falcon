package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer"
)

// startCmd represents the start command
func startCmd(app *relayer.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"st"},
		Short:   "Start the falcon tunnel relayer program",
		Args:    withUsage(cobra.MinimumNArgs(0)),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s start                             # start relaying data from every tunnel being registered on source chain.
$ %s start --tunnel-ids 1,12           # start relaying data from specific tunnelIDs.
$ %s start --tunnel-creator 0xABC123   # start relaying data from tunnels created by a specific address`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsTunnelIDs, err := cmd.Flags().GetUintSlice(flagTunnelIds)
			if err != nil {
				return err
			}

			tunnelIDs := make([]uint64, len(argsTunnelIDs))
			for i, v := range argsTunnelIDs {
				tunnelIDs[i] = uint64(v)
			}

			tunnelCreator, err := cmd.Flags().GetString(flagTunnelCreator)
			if err != nil {
				return err
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

			if tunnelCreator != "" && len(tunnelIDs) != 0 {
				return fmt.Errorf(
					"the --tunnel-creator and --tunnel-ids flags cannot be used together, please specify only one of these options at a time",
				)
			}

			return app.Start(cmd.Context(), tunnelIDs, tunnelCreator)
		},
	}

	cmd.Flags().String(
		flagMetricsListenAddr,
		"",
		"address to use for metrics server. By default, "+
			"will be the metrics-listen-addr parameter in the global config. ",
	)
	cmd.Flags().UintSlice(
		flagTunnelIds,
		[]uint{},
		"specific tunnel IDs to relay",
	)
	cmd.Flags().String(
		flagTunnelCreator,
		"",
		"relay tunnels created by this creator address",
	)
	return cmd
}
