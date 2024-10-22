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
	return nil
}
