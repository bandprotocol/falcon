package relayer

import (
	"context"

	"go.uber.org/zap"
)

// Scheduler is a struct to manage all tunnel relayers
type Scheduler struct {
	Log            *zap.Logger
	TunnelRelayers []TunnelRelayer
}

// NewScheduler creates a new Scheduler
func NewScheduler(
	log *zap.Logger,
	tunnelRelayers []TunnelRelayer,
) *Scheduler {
	return &Scheduler{
		Log:            log,
		TunnelRelayers: tunnelRelayers,
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(ctx context.Context) error {
	// TODO: Do a go routine for checkAndRelay for an interval.
	// TODO: handle panic in go routine.
	for _, tr := range s.TunnelRelayers {
		if err := tr.CheckAndRelay(ctx); err != nil {
			s.Log.Error("failed to start tunnel relayer", zap.Error(err))
		}
	}

	// TODO: Do a go routine for updating client's selected endpoint.

	return nil
}
