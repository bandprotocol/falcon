package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer"
)

// StartCmd represents the start command
func StartCmd(appCreator relayer.AppCreator, defaultHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"st"},
		Short:   "Start the tunnel relayer program",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(`
start                             # start relaying data from every tunnel being registered on source chain.
start --tunnel-ids 1,12           # start relaying data from specific tunnelIDs.
start --tunnel-creator 0xABC123   # start relaying data from tunnels created by a specific address`),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := createApp(cmd, appCreator, defaultHome)
			if err != nil {
				return err
			}
			defer syncLog(app.GetLog())

			if err := app.Init(cmd.Context()); err != nil {
				return err
			}

			cfg := app.GetConfig()
			if cfg == nil {
				return fmt.Errorf("config is not initialized")
			}

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
				metricsListenAddr = cfg.Global.MetricsListenAddr
			}
			if metricsListenAddr != "" {
				if err := relayermetrics.StartMetricsServer(cmd.Context(), app.GetLog(), metricsListenAddr); err != nil {
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
