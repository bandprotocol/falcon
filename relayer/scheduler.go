package relayer

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
)

// Scheduler is a struct to manage all tunnel relayers
type Scheduler struct {
	Log                    *zap.Logger
	TunnelRelayers         []*TunnelRelayer
	CheckingPacketInterval time.Duration
	SyncTunnelsInterval    time.Duration
	PenaltySkipRounds      uint

	PenaltySkipRemaining []uint
	isSyncTunnelsAllowed bool

	BandClient     band.Client
	ChainProviders chains.ChainProviders
}

// NewScheduler creates a new Scheduler
func NewScheduler(
	log *zap.Logger,
	tunnelRelayers []*TunnelRelayer,
	checkingPacketInterval time.Duration,
	syncTunnelsInterval time.Duration,
	penaltyAttempts uint,
	isSyncTunnelsAllowed bool,
	bandClient band.Client,
	chainProviders chains.ChainProviders,
) *Scheduler {
	return &Scheduler{
		Log:                    log,
		TunnelRelayers:         tunnelRelayers,
		CheckingPacketInterval: checkingPacketInterval,
		SyncTunnelsInterval:    syncTunnelsInterval,
		PenaltySkipRounds:      penaltyAttempts,
		PenaltySkipRemaining:   make([]uint, len(tunnelRelayers)),
		isSyncTunnelsAllowed:   isSyncTunnelsAllowed,
		BandClient:             bandClient,
		ChainProviders:         chainProviders,
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.CheckingPacketInterval)
	defer ticker.Stop()

	syncTunnelTicker := time.NewTicker(s.SyncTunnelsInterval)
	defer syncTunnelTicker.Stop()

	// Mutex to prevent overlapping executions
	var executionMutex sync.Mutex

	// Execute once at the start
	executionMutex.Lock()
	// execute once we start the scheduler.
	s.Execute(ctx)
	executionMutex.Unlock()

	for {
		select {
		case <-ctx.Done():
			s.Log.Info("Stopping the scheduler")

			return nil
		case <-syncTunnelTicker.C:
			s.SyncTunnels(ctx)
		case <-ticker.C:
			if executionMutex.TryLock() {
				s.Execute(ctx)
				executionMutex.Unlock()
			}
		}
	}
}

// Execute executes the task for the tunnel relayer
func (s *Scheduler) Execute(ctx context.Context) {
	// Execute the task for each tunnel relayer
	for i, tr := range s.TunnelRelayers {
		if s.PenaltySkipRemaining[i] > 0 {
			s.Log.Debug(
				"Skipping tunnel execution due to penalty from previous failure.",
				zap.Uint64("tunnel_id", tr.TunnelID),
				zap.Int("relayer_id", i),
				zap.Uint("penalty_skip_remaining", s.PenaltySkipRemaining[i]),
			)
			s.PenaltySkipRemaining[i] -= 1

			continue
		}

		// Execute the task, if error occurs, wait for the next round.
		task := NewTask(i, s.CheckingPacketInterval)
		go s.TriggerTunnelRelayer(ctx, task)
	}
}

// TriggerTunnelRelayer triggers the tunnel relayer to check and relay the packet
func (s *Scheduler) TriggerTunnelRelayer(ctx context.Context, task Task) {
	tr := s.TunnelRelayers[task.RelayerID]

	// Check and relay the packet, if error occurs, set the error flag.
	if err, isExecuting := tr.CheckAndRelay(ctx); err != nil && !isExecuting {
		s.PenaltySkipRemaining[task.RelayerID] = s.PenaltySkipRounds

		s.Log.Error(
			"Failed to execute, Penalty for the tunnel relayer",
			zap.Error(err),
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

		return
	}

	s.Log.Info(
		"Tunnel relayer finished execution",
		zap.Uint64("tunnel_id", tr.TunnelID),
	)
}

// SyncTunnels synchronizes the Bandchain's tunnels with the latest tunnels.
func (s *Scheduler) SyncTunnels(ctx context.Context) {
	if !s.isSyncTunnelsAllowed {
		return
	}

	s.Log.Info("Start syncing new tunnels")
	tunnels, err := s.BandClient.GetTunnels(ctx)
	if err != nil {
		s.Log.Error("Failed to fetch tunnels from BandChain", zap.Error(err))
		return
	}
	oldTunnelCount := len(s.TunnelRelayers)

	if oldTunnelCount == len(tunnels) {
		s.Log.Info("No new tunnels to sync")
		return
	}

	for i := oldTunnelCount; i < len(tunnels); i++ {
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

		s.TunnelRelayers = append(s.TunnelRelayers, &tr)
		s.PenaltySkipRemaining = append(s.PenaltySkipRemaining, 0)
		s.Log.Info(
			"New tunnel synchronized successfully",
			zap.String("chain_name", tunnels[i].TargetChainID),
			zap.Uint64("tunnel_id", tunnels[i].ID),
		)
	}
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
