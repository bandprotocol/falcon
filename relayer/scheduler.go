package relayer

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
)

// Scheduler is a struct to manage all tunnel relayers
type Scheduler struct {
	Log                              *zap.Logger
	TunnelRelayers                   []*TunnelRelayer
	CheckingPacketInterval           time.Duration
	SyncTunnelsInterval              time.Duration
	MaxCheckingPacketPenaltyDuration time.Duration
	ExponentialFactor                float64

	isErrorOnHolds []bool
	penaltyTaskCh  chan Task

	BandClient     band.Client
	ChainProviders chains.ChainProviders
}

// NewScheduler creates a new Scheduler
func NewScheduler(
	log *zap.Logger,
	tunnelRelayers []*TunnelRelayer,
	checkingPacketInterval time.Duration,
	syncTunnelsInterval time.Duration,
	maxCheckingPacketPenaltyDuration time.Duration,
	exponentialFactor float64,
	bandClient band.Client,
	chainProviders chains.ChainProviders,
) *Scheduler {
	return &Scheduler{
		Log:                              log,
		TunnelRelayers:                   tunnelRelayers,
		CheckingPacketInterval:           checkingPacketInterval,
		SyncTunnelsInterval:              syncTunnelsInterval,
		MaxCheckingPacketPenaltyDuration: maxCheckingPacketPenaltyDuration,
		ExponentialFactor:                exponentialFactor,
		isErrorOnHolds:                   make([]bool, len(tunnelRelayers)),
		penaltyTaskCh:                    make(chan Task, len(tunnelRelayers)),
		BandClient:                       bandClient,
		ChainProviders:                   chainProviders,
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.CheckingPacketInterval)
	syncTunnelTicker := time.NewTicker(s.SyncTunnelsInterval)

	// execute once we start the scheduler.
	s.Execute(ctx)
	defer ticker.Stop()
	defer syncTunnelTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.Log.Info("Stopping the scheduler")

			return nil
		case <-syncTunnelTicker.C:
			s.SyncTunnels(ctx)
		case <-ticker.C:
			s.Execute(ctx)
		case task := <-s.penaltyTaskCh:
			// Execute the task with penalty waiting period
			go func(task Task) {
				executeFn := func(ctx context.Context, t Task) {
					s.isErrorOnHolds[task.RelayerID] = false
					s.TriggerTunnelRelayer(ctx, task)
				}

				task.Wait(ctx, executeFn)
			}(task)
		}
	}
}

// Execute executes the task for the tunnel relayer
func (s *Scheduler) Execute(ctx context.Context) {
	// Execute the task for each tunnel relayer
	for i, tr := range s.TunnelRelayers {
		if s.isErrorOnHolds[i] {
			s.Log.Info(
				"Skipping this tunnel: the operation is on hold due to error on last round.",
				zap.Uint64("tunnel_id", tr.TunnelID),
				zap.Int("relayer_id", i),
			)

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
	s.Log.Info("Executing task", zap.Uint64("tunnel_id", tr.TunnelID))

	// if the tunnel relayer is executing, skip the round
	if tr.IsExecuting() {
		s.Log.Info(
			"Skipping this tunnel: tunnel relayer is executing on another process",
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

		return
	}

	// Check and relay the packet, if error occurs, set the error flag.
	if err := tr.CheckAndRelay(ctx); err != nil {
		s.isErrorOnHolds[task.RelayerID] = true
		newInterval := s.calculatePenaltyInterval(task.WaitingInterval)

		s.Log.Error(
			"Failed to execute, Penalty for the tunnel relayer",
			zap.Error(err),
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

		newTask := NewTask(task.RelayerID, newInterval)

		s.penaltyTaskCh <- newTask
		return
	}

	// If the task is successful, reset the error flag.
	s.isErrorOnHolds[task.RelayerID] = false
	if tr.isWaitingSignature {
		s.Log.Info(
			"the current packet must wait for a signature",
			zap.Uint64("tunnel_id", tr.TunnelID),
		)
	} else {
		s.Log.Info(
			"tunnel relayer is successfully executed",
			zap.Uint64("tunnel_id", tr.TunnelID),
		)
	}
}

func (s *Scheduler) SyncTunnels(ctx context.Context) {
	s.Log.Info("starting sync tunnels")
	tunnels, err := s.BandClient.GetTunnels(ctx)
	if err != nil {
		s.Log.Error("Failed to fetch tunnels from Band Chain", zap.Error(err))
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
		s.isErrorOnHolds = append(s.isErrorOnHolds, false)
		s.Log.Info(
			"New tunnel synchronized successfully",
			zap.String("chain_name", tunnels[i].TargetChainID),
			zap.Uint64("tunnel_id", tunnels[i].ID),
		)
	}
	for len(s.penaltyTaskCh) > 0 {
		// Wait for penalty tasks to drain
	}
	s.penaltyTaskCh = make(chan Task, len(s.TunnelRelayers))
}

// calculatePenaltyInterval applies exponential backoff with a max limit
func (s *Scheduler) calculatePenaltyInterval(interval time.Duration) time.Duration {
	newInterval := time.Duration(float64(interval) * s.ExponentialFactor)
	if newInterval > s.MaxCheckingPacketPenaltyDuration {
		newInterval = s.MaxCheckingPacketPenaltyDuration
	}
	return newInterval
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

// Wait waits for the task to be executed
func (t Task) Wait(ctx context.Context, executeFn func(ctx context.Context, t Task)) {
	select {
	case <-ctx.Done():
		// Do nothing
	case <-time.After(t.WaitingInterval):
		executeFn(ctx, t)
	}
}
