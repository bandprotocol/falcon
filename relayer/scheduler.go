package relayer

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/band"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
)

// Scheduler is a struct to manage all tunnel relayers
type Scheduler struct {
	Log                    *zap.Logger
	CheckingPacketInterval time.Duration
	SyncTunnelsInterval    time.Duration
	PenaltySkipRounds      uint

	BandClient     band.Client
	ChainProviders chains.ChainProviders

	tunnelRelayers   map[uint64]*TunnelRelayer
	bandLatestTunnel int
}

// NewScheduler creates a new Scheduler
func NewScheduler(
	log *zap.Logger,
	checkingPacketInterval time.Duration,
	syncTunnelsInterval time.Duration,
	penaltySkipRounds uint,
	bandClient band.Client,
	chainProviders chains.ChainProviders,
) *Scheduler {
	return &Scheduler{
		Log:                    log,
		CheckingPacketInterval: checkingPacketInterval,
		SyncTunnelsInterval:    syncTunnelsInterval,
		PenaltySkipRounds:      penaltySkipRounds,
		BandClient:             bandClient,
		ChainProviders:         chainProviders,
		tunnelRelayers:         make(map[uint64]*TunnelRelayer),
		bandLatestTunnel:       0,
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(ctx context.Context, tunnelIDs []uint64, tunnelCreator string) error {
	s.SyncTunnels(ctx, tunnelIDs, tunnelCreator)

	go s.BandClient.HandleProducePacketSuccess(func(tunnelID uint64) {
		tunnelRelayer, ok := s.tunnelRelayers[tunnelID]
		if !ok {
			return
		}

		s.Log.Info("Received produce_packet_success event", zap.Uint64("tunnel_id", tunnelID))

		if tunnelRelayer.penaltySkipRemaining > 0 {
			s.Log.Info(
				"Skipping tunnel execution due to penalty from previous failure",
				zap.Uint64("tunnel_id", tunnelID),
			)
			return
		}

		go s.TriggerTunnelRelayer(ctx, tunnelRelayer)
	})

	ticker := time.NewTicker(s.CheckingPacketInterval)
	defer ticker.Stop()

	syncTunnelTicker := time.NewTicker(s.SyncTunnelsInterval)
	defer syncTunnelTicker.Stop()

	// execute once we start the scheduler.
	s.Execute(ctx)

	for {
		select {
		case <-ctx.Done():
			s.Log.Info("Stopping the scheduler")

			return nil
		case <-syncTunnelTicker.C:
			// sync tunnels only when no specific tunnel IDs are provided
			if len(tunnelIDs) == 0 {
				s.SyncTunnels(ctx, tunnelIDs, tunnelCreator)
			}
		case <-ticker.C:
			s.Execute(ctx)
		}
	}
}

// Execute executes the task for the tunnel relayer
func (s *Scheduler) Execute(ctx context.Context) {
	s.Log.Info("Executing tunnel relayers from the scheduler")

	for _, tr := range s.tunnelRelayers {
		if tr.penaltySkipRemaining > 0 {
			s.Log.Debug(
				"Skipping tunnel execution due to penalty from previous failure.",
				zap.Uint64("tunnel_id", tr.TunnelID),
				zap.Uint("penalty_skip_remaining", tr.penaltySkipRemaining),
			)
			tr.penaltySkipRemaining -= 1
			continue
		}
		go s.TriggerTunnelRelayer(ctx, tr)
	}
}

// TriggerTunnelRelayer triggers the tunnel relayer to check and relay the packet
func (s *Scheduler) TriggerTunnelRelayer(ctx context.Context, tr *TunnelRelayer) {
	chainName := tr.TargetChainProvider.GetChainName()
	startExecutionTaskTime := time.Now()

	isExecuting, err := tr.CheckAndRelay(ctx, false)

	switch {
	case err != nil:
		tr.penaltySkipRemaining = s.PenaltySkipRounds

		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.ErrorTaskStatus)

		s.Log.Error(
			"Failed to execute, Penalty for the tunnel relayer",
			zap.Error(err),
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

	case isExecuting:
		// Record metrics for the skipped task execution
		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.SkippedTaskStatus)

		s.Log.Info(
			"This tunnel relayer is already executing",
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

	default:
		// Record execution time of finished task (ms)
		relayermetrics.ObserveFinishedTaskExecutionTime(
			tr.TunnelID,
			chainName,
			time.Since(startExecutionTaskTime).Milliseconds(),
		)

		// Record metrics for the finished task execution
		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.FinishedTaskStatus)
	}
}

// SyncTunnels synchronizes the Bandchain's tunnels with the latest tunnels.
// If tunnel creator is provided, only tunnels created by that address will be synchronized.
func (s *Scheduler) SyncTunnels(ctx context.Context, tunnelIDs []uint64, tunnelCreator string) {
	s.Log.Info("Start syncing tunnels from Bandchain")
	tunnels, err := s.getTunnels(ctx, tunnelIDs)
	if err != nil {
		s.Log.Error("Failed to fetch tunnels from BandChain", zap.Error(err))
		return
	}

	if s.bandLatestTunnel == len(tunnels) {
		s.Log.Info("No new tunnels to sync")
		return
	}

	for i := s.bandLatestTunnel; i < len(tunnels); i++ {
		chainProvider, ok := s.ChainProviders[tunnels[i].TargetChainID]
		if !ok {
			s.Log.Warn(
				"Chain name not found in config",
				zap.String("chain_name", tunnels[i].TargetChainID),
				zap.Uint64("tunnel_id", tunnels[i].ID),
			)
			continue
		}

		// if tunnel creator is provided, check if the tunnel matches the creator
		if tunnelCreator != "" && tunnels[i].Creator != tunnelCreator {
			continue
		}

		tr := NewTunnelRelayer(
			s.Log,
			tunnels[i].ID,
			s.CheckingPacketInterval,
			s.BandClient,
			chainProvider,
		)

		s.tunnelRelayers[tunnels[i].ID] = &tr

		// update the metric for the number of tunnels per destination chain
		relayermetrics.IncTunnelsPerDestinationChain(tunnels[i].TargetChainID)

		s.Log.Info(
			"New tunnel synchronized successfully",
			zap.String("chain_name", tunnels[i].TargetChainID),
			zap.Uint64("tunnel_id", tunnels[i].ID),
		)
	}

	s.bandLatestTunnel = len(tunnels)
}

// getTunnels retrieves the list of tunnels by given tunnel IDs. If no tunnel ID is provided,
// get all tunnels
func (s *Scheduler) getTunnels(ctx context.Context, tunnelIDs []uint64) ([]bandtypes.Tunnel, error) {
	if len(tunnelIDs) == 0 {
		return s.BandClient.GetTunnels(ctx)
	}

	tunnels := make([]bandtypes.Tunnel, 0, len(tunnelIDs))
	for _, tunnelID := range tunnelIDs {
		tunnel, err := s.BandClient.GetTunnel(ctx, tunnelID)
		if err != nil {
			return nil, err
		}

		tunnels = append(tunnels, *tunnel)
	}

	return tunnels, nil
}
