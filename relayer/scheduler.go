package relayer

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
)

const penaltyTaskChSize = 1000

// Scheduler is a struct to manage all tunnel relayers
type Scheduler struct {
	Log                              *zap.Logger
	TunnelRelayers                   []*TunnelRelayer
	BandLatestTunnel                 int
	CheckingPacketInterval           time.Duration
	SyncTunnelsInterval              time.Duration
	MaxCheckingPacketPenaltyDuration time.Duration
	ExponentialFactor                float64

	isErrorOnHolds       []bool
	isSyncTunnelsAllowed bool
	penaltyTaskCh        chan Task

	BandClient     band.Client
	ChainProviders chains.ChainProviders
	ChainNames     map[string]bool
}

// NewScheduler creates a new Scheduler
func NewScheduler(
	log *zap.Logger,
	tunnelRelayers []*TunnelRelayer,
	bandLatestTunnel int,
	checkingPacketInterval time.Duration,
	syncTunnelsInterval time.Duration,
	maxCheckingPacketPenaltyDuration time.Duration,
	exponentialFactor float64,
	isSyncTunnelsAllowed bool,
	bandClient band.Client,
	chainProviders chains.ChainProviders,
	chainNames map[string]bool,
) *Scheduler {
	return &Scheduler{
		Log:                              log,
		TunnelRelayers:                   tunnelRelayers,
		BandLatestTunnel:                 bandLatestTunnel,
		CheckingPacketInterval:           checkingPacketInterval,
		SyncTunnelsInterval:              syncTunnelsInterval,
		MaxCheckingPacketPenaltyDuration: maxCheckingPacketPenaltyDuration,
		ExponentialFactor:                exponentialFactor,
		isErrorOnHolds:                   make([]bool, len(tunnelRelayers)),
		isSyncTunnelsAllowed:             isSyncTunnelsAllowed,
		penaltyTaskCh:                    make(chan Task, penaltyTaskChSize),
		BandClient:                       bandClient,
		ChainProviders:                   chainProviders,
		ChainNames:                       chainNames,
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.CheckingPacketInterval)
	defer ticker.Stop()

	syncTunnelTicker := time.NewTicker(s.SyncTunnelsInterval)
	defer syncTunnelTicker.Stop()

	relayermetrics.AddTunnellCount(uint64(len(s.TunnelRelayers)))
	relayermetrics.AddDestinationChainCount(uint64(len(s.ChainNames)))

	// execute once we start the scheduler.
	s.Execute(ctx)

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
			s.Log.Debug(
				"Skipping this tunnel: the operation is on hold due to error on last round.",
				zap.Uint64("tunnel_id", tr.TunnelID),
				zap.Int("relayer_id", i),
			)

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
	tr := s.TunnelRelayers[task.RelayerID]

	// if the tunnel relayer is executing, skip the round
	if tr.IsExecuting() {
		s.Log.Debug(
			"Skipping this tunnel: tunnel relayer is executing on another process",
			zap.Uint64("tunnel_id", tr.TunnelID),
		)
		return
	}

	s.Log.Info("Executing task", zap.Uint64("tunnel_id", tr.TunnelID))
	startExecutionTaskTime := time.Now()

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

	// record the execution time of successful task.
	relayermetrics.ObserveTaskExecutionTime(
		tr.TunnelID,
		float64(time.Since(startExecutionTaskTime).Milliseconds()),
	)

	// If the task is successful, reset the error flag.
	s.isErrorOnHolds[task.RelayerID] = false

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

	if s.BandLatestTunnel == len(tunnels) {
		s.Log.Info("No new tunnels to sync")
		return
	}

	oldTunnelRelayerCount := len(s.TunnelRelayers)
	oldDestinationChainCount := len(s.ChainNames)

	for i := s.BandLatestTunnel; i < len(tunnels); i++ {
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

		s.ChainNames[tunnels[i].TargetChainID] = true

		s.TunnelRelayers = append(s.TunnelRelayers, &tr)
		s.isErrorOnHolds = append(s.isErrorOnHolds, false)

		s.Log.Info(
			"New tunnel synchronized successfully",
			zap.String("chain_name", tunnels[i].TargetChainID),
			zap.Uint64("tunnel_id", tunnels[i].ID),
		)
	}

	s.BandLatestTunnel = len(tunnels)

	// update metrics for the number of destination chains and tunnels after synchronization
	relayermetrics.AddDestinationChainCount(uint64(len(s.ChainNames) - oldDestinationChainCount))
	relayermetrics.AddTunnellCount(uint64(len(s.TunnelRelayers) - oldTunnelRelayerCount))
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
