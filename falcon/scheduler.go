package falcon

// Scheduler is a struct to manage all tunnel relayers
type Scheduler struct {
	TunnelRelayers []TunnelRelayer
}

// NewScheduler creates a new Scheduler
func NewScheduler(
	tunnelRelayers []TunnelRelayer,
) *Scheduler {
	return &Scheduler{
		TunnelRelayers: tunnelRelayers,
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start() error {
	return nil
}

// Stop stops all tunnel relayers
func (s *Scheduler) Stop() error {
	return nil
}
