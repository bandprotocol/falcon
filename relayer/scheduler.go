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
func (s *Scheduler) Start(ctx context.Context, tunnelIds []uint64) error {
	s.SyncTunnels(ctx, tunnelIds)

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
			if len(tunnelIds) == 0 {
				s.SyncTunnels(ctx, tunnelIds)
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

		// record metrics for the task execution for the current tunnel relayer
		relayermetrics.IncTasksCount(tr.TunnelID)
	}
}

// TriggerTunnelRelayer triggers the tunnel relayer to check and relay the packet
func (s *Scheduler) TriggerTunnelRelayer(ctx context.Context, task Task) {
	tr := s.tunnelRelayers[task.RelayerID]

	startExecutionTaskTime := time.Now()

	// Check and relay the packet, if error occurs, set the error flag.
	if isExecuting, err := tr.CheckAndRelay(ctx); err != nil && !isExecuting {
		s.penaltySkipRemaining[task.RelayerID] = s.PenaltySkipRounds

		s.Log.Error(
			"Failed to execute, Penalty for the tunnel relayer",
			zap.Error(err),
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

		return
	}

	// record the execution time of successful task.
	relayermetrics.ObserveTaskExecutionTime(
		tr.TunnelID,
		float64(time.Since(startExecutionTaskTime).Milliseconds()),
	)

	s.Log.Info(
		"Tunnel relayer finished execution",
		zap.Uint64("tunnel_id", tr.TunnelID),
	)
}

// SyncTunnels synchronizes the Bandchain's tunnels with the latest tunnels.
func (s *Scheduler) SyncTunnels(ctx context.Context, tunnelIds []uint64) {
	s.Log.Info("Start syncing tunnels from Bandchain")
	tunnels, err := s.getTunnels(ctx, tunnelIds)
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

		tr := NewTunnelRelayer(
			s.Log,
			tunnels[i].ID,
			tunnels[i].TargetAddress,
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
