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

	tunnelRelayers       []*TunnelRelayer
	bandLatestTunnel     int
	penaltySkipRemaining []uint
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
		tunnelRelayers:         []*TunnelRelayer{},
		bandLatestTunnel:       0,
		penaltySkipRemaining:   []uint{},
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(ctx context.Context, tunnelIDs []uint64, tunnelCreator string) error {
	s.SyncTunnels(ctx, tunnelIDs, tunnelCreator)

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
	// Execute the task for each tunnel relayer
	for i, tr := range s.tunnelRelayers {
		if s.penaltySkipRemaining[i] > 0 {
			s.Log.Debug(
				"Skipping tunnel execution due to penalty from previous failure.",
				zap.Uint64("tunnel_id", tr.TunnelID),
				zap.Int("relayer_id", i),
				zap.Uint("penalty_skip_remaining", s.penaltySkipRemaining[i]),
			)
			s.penaltySkipRemaining[i] -= 1

			continue
		}

		// Execute the task, if error occurs, wait for the next round.
		task := NewTask(i, s.CheckingPacketInterval)
		go s.TriggerTunnelRelayer(ctx, task)
	}
}

// TriggerTunnelRelayer triggers the tunnel relayer to check and relay the packet
func (s *Scheduler) TriggerTunnelRelayer(ctx context.Context, task Task) {
	tr := s.tunnelRelayers[task.RelayerID]
	chainName := tr.TargetChainProvider.GetChainName()
	startExecutionTaskTime := time.Now()

	isExecuting, err := tr.CheckAndRelay(ctx)

	switch {
	case err != nil:
		s.penaltySkipRemaining[task.RelayerID] = s.PenaltySkipRounds

		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.ErrorTaskStatus)

		s.Log.Error(
			"Failed to execute, Penalty for the tunnel relayer",
			zap.Error(err),
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

	case isExecuting:
		// Record metrics for the skipped task execution
		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.SkippedTaskStatus)

	default:
		// Record execution time of finished task (ms)
		relayermetrics.ObserveFinishedTaskExecutionTime(
			tr.TunnelID,
			chainName,
			time.Since(startExecutionTaskTime).Milliseconds(),
		)

		// Record metrics for the finished task execution
		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.FinishedTaskStatus)

		s.Log.Info(
			"Tunnel relayer finished execution",
			zap.Uint64("tunnel_id", tr.TunnelID),
		)
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

		s.tunnelRelayers = append(s.tunnelRelayers, &tr)
		s.penaltySkipRemaining = append(s.penaltySkipRemaining, 0)

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

// Task is a struct to manage the task for the tunnel relayer
type Task struct {
	RelayerID       int
	WaitingInterval time.Duration
}

// NewTask creates a new Task
func NewTask(relayerID int, waitingInterval time.Duration) Task {
	return Task{
		RelayerID:       relayerID,
		WaitingInterval: waitingInterval,
	}
}
