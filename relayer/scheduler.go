package relayer

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Scheduler is a struct to manage all tunnel relayers
type Scheduler struct {
	Log                              *zap.Logger
	TunnelRelayers                   []TunnelRelayer
	CheckingPacketInterval           time.Duration
	MaxCheckingPacketPenaltyDuration time.Duration
	ExponentialFactor                float64

	relayerTaskCh  chan int
	isErrorOnHolds []bool
	penaltyTaskCh  chan Task
}

// NewScheduler creates a new Scheduler
func NewScheduler(
	log *zap.Logger,
	tunnelRelayers []TunnelRelayer,
	checkingPacketInterval time.Duration,
	maxCheckingPacketPenaltyDuration time.Duration,
	exponentialFactor float64,
) *Scheduler {
	return &Scheduler{
		Log:                              log,
		TunnelRelayers:                   tunnelRelayers,
		CheckingPacketInterval:           checkingPacketInterval,
		MaxCheckingPacketPenaltyDuration: maxCheckingPacketPenaltyDuration,
		ExponentialFactor:                exponentialFactor,
		relayerTaskCh:                    make(chan int),
		isErrorOnHolds:                   make([]bool, len(tunnelRelayers)),
		penaltyTaskCh:                    make(chan Task),
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.CheckingPacketInterval)
	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			// Execute the task for each tunnel relayer
			for i := range s.TunnelRelayers {
				if s.isErrorOnHolds[i] {
					tr := s.TunnelRelayers[i]
					s.Log.Info(
						"Skipping this tunnel: the operation is on hold due to error on last round.",
						zap.Uint64("tunnel_id", tr.TunnelID),
						zap.Int("relayer_id", i),
					)

					continue
				}

				// Execute the task, if error occurs, wait for the next round.
				task := NewTask(i, s.CheckingPacketInterval)
				s.Execute(ctx, task)
			}
		case task := <-s.penaltyTaskCh:
			// Execute the task with penalty waiting period
			go func(task Task) {
				executeFn := func(ctx context.Context, t Task) {
					s.isErrorOnHolds[task.RelayerID] = false
					s.Execute(ctx, task)
				}

				task.Wait(ctx, executeFn)
			}(task)
		}
	}
}

// Execute executes the task
func (s *Scheduler) Execute(ctx context.Context, task Task) {
	tr := s.TunnelRelayers[task.RelayerID]
	s.Log.Info("Executing task", zap.Uint64("tunnel_id", tr.TunnelID))

	// if the tunnel relayer is executing, skip the round
	if tr.IsExecuting() {
		s.Log.Debug(
			"Skipping this tunnel: tunnel relayer is executing on another process",
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

		return
	}

	// Check and relay the packet, if error occurs, set the error flag.
	if err := tr.CheckAndRelay(ctx); err != nil {
		s.isErrorOnHolds[task.RelayerID] = true
		newInterval := time.Duration(float64(task.WaitingInterval) * s.ExponentialFactor)
		if newInterval > s.MaxCheckingPacketPenaltyDuration {
			newInterval = s.MaxCheckingPacketPenaltyDuration
		}

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
	s.Log.Info(
		"tunnel relayer is successfully executed",
		zap.Uint64("tunnel_id", tr.TunnelID),
	)
}

type Task struct {
	RelayerID       int
	WaitingInterval time.Duration
}

func NewTask(relayerID int, waitingInterval time.Duration) Task {
	return Task{
		RelayerID:       relayerID,
		WaitingInterval: waitingInterval,
	}
}

func (t Task) Wait(ctx context.Context, executeFn func(ctx context.Context, t Task)) {
	ticker := time.NewTicker(t.WaitingInterval)
	select {
	case <-ctx.Done():
		// Do nothing
	case <-ticker.C:
		executeFn(ctx, t)
	}
}
