package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"context"
	"github.com/bandprotocol/falcon/relayer"
	"github.com/spf13/cobra"
	"time"
)

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
			// Create a parent context and ticker
			parentCtx := cmd.Context()
			ticker := time.NewTicker(app.Config.Global.SyncTunnelsInterval)
			defer ticker.Stop()

			tunnelIDs := []uint64{}
			for _, arg := range args {
				tunnelID, err := strconv.ParseUint(arg, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid tunnel ID '%s': %w", arg, err)
				}
				tunnelIDs = append(tunnelIDs, tunnelID)
			}

			ctx, cancel := context.WithCancel(parentCtx)
			defer cancel()

			// Create a channel to capture errors from app.Start
			errCh := make(chan error, 1)

			// Start relayer asynchronously
			go func() {
				errCh <- app.Start(ctx, tunnelIDs)
			}()

			// handle sync intervals and context cancellation
			for {
				select {
				case <-parentCtx.Done():
					// Stop the relayer gracefully
					cancel()
					app.Log.Info("Stopping tunnel relayer")
					return nil

				case <-ticker.C:
					// Restart the relayer with updated tunnels
					app.Log.Info("Sync interval reached, restarting tunnel relayer")
					cancel() // Cancel the previous context

					// Create a new context for the restarted relayer
					ctx, cancel = context.WithCancel(parentCtx)

					// Restart relayer asynchronously
					go func() {
						errCh <- app.Start(ctx, tunnelIDs)
					}()

				case err := <-errCh:
					// Handle errors from app.Start
					if err != nil {
						return err
					}
				}
			}
		},
	}

	return cmd
}
