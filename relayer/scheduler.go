package relayer

import (
	"context"
	"time"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/alert"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/band/subscriber"
	bandtypes "github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	"github.com/bandprotocol/falcon/relayer/config"
	"github.com/bandprotocol/falcon/relayer/logger"
)

// Scheduler is a struct to manage all tunnel relayers
type Scheduler struct {
	Log                    logger.Logger
	CheckingPacketInterval time.Duration
	SyncTunnelsInterval    time.Duration
	PenaltySkipRounds      uint
	SubscriptionTimeout    time.Duration

	BandClient     band.Client
	ChainProviders chains.ChainProviders

	Alert alert.Alert

	relayTunnelIDCh  chan uint64
	tunnelRelayers   map[uint64]*TunnelRelayer
	bandLatestTunnel int
	tunnelCreator    string
}

// NewScheduler creates a new Scheduler
func NewScheduler(
	log logger.Logger,
	config *config.Config,
	bandClient band.Client,
	chainProviders chains.ChainProviders,
	tunnelCreator string,
	alert alert.Alert,
) *Scheduler {
	relayTunnelIDCh := make(chan uint64, 1000)

	return &Scheduler{
		Log:                    log,
		CheckingPacketInterval: config.Global.CheckingPacketInterval,
		SyncTunnelsInterval:    config.Global.SyncTunnelsInterval,
		PenaltySkipRounds:      config.Global.PenaltySkipRounds,
		SubscriptionTimeout:    config.BandChain.Timeout,
		BandClient:             bandClient,
		ChainProviders:         chainProviders,
		Alert:                  alert,
		relayTunnelIDCh:        relayTunnelIDCh,
		tunnelRelayers:         make(map[uint64]*TunnelRelayer),
		bandLatestTunnel:       0,
		tunnelCreator:          tunnelCreator,
	}
}

// WithTunnels sets the tunnel relayer to the scheduler from the given tunnels.
func (s *Scheduler) WithTunnels(tunnels []bandtypes.Tunnel) *Scheduler {
	validTunnels := s.filterTunnels(tunnels)
	s.setTunnelRelayer(validTunnels)

	return s
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(ctx context.Context, isSyncTunnels bool) error {
	if err := s.initialize(ctx); err != nil {
		return err
	}

	// execute first time
	if isSyncTunnels {
		s.SyncTunnels(ctx)
	} else {
		s.Execute(ctx)
	}

	ticker := time.NewTicker(s.CheckingPacketInterval)
	defer ticker.Stop()

	syncTunnelTicker := time.NewTicker(s.SyncTunnelsInterval)
	defer syncTunnelTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.Log.Info("Stopping the scheduler")

			return nil
		case <-syncTunnelTicker.C:
			if isSyncTunnels {
				s.SyncTunnels(ctx)
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
				"tunnel_id", tr.TunnelID,
				"penalty_skip_remaining", tr.penaltySkipRemaining,
			)
			tr.penaltySkipRemaining -= 1
			continue
		}
		go func() { _ = s.TriggerTunnelRelayer(ctx, tr) }()
	}
}

// TriggerTunnelRelayer triggers the tunnel relayer to check and relay the packet
func (s *Scheduler) TriggerTunnelRelayer(ctx context.Context, tr *TunnelRelayer) (status RelayStatus) {
	chainName := tr.TargetChainProvider.GetChainName()
	startExecutionTaskTime := time.Now()

	// checkAndRelay tunnel's packets and update penalty if it fails to do so.
	relayStatus, err := tr.CheckAndRelay(ctx, false)
	if err != nil {
		tr.penaltySkipRemaining = s.PenaltySkipRounds

		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.ErrorTaskStatus)
		s.Log.Error(
			"Failed to execute, Penalty for the tunnel relayer",
			"tunnel_id", tr.TunnelID,
			err,
		)

		return RelayStatusFailed
	}

	switch relayStatus {
	case RelayStatusExecuting:
		// Record metrics for the executing task execution
		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.ExecutingTaskStatus)
	case RelayStatusSuccess:
		// Record execution time of finished task (ms)
		relayermetrics.ObserveFinishedTaskExecutionTime(
			tr.TunnelID,
			chainName,
			time.Since(startExecutionTaskTime).Milliseconds(),
		)

		// Record metrics for the finished task execution
		relayermetrics.IncTasksCount(tr.TunnelID, chainName, relayermetrics.FinishedTaskStatus)

		s.Log.Info("Relay packet successfully", "tunnel_id", tr.TunnelID)
	}

	return relayStatus
}

// SyncTunnels synchronizes the Bandchain's tunnels with the latest tunnels.
// If tunnel creator is provided, only tunnels created by that address will be synchronized.
func (s *Scheduler) SyncTunnels(ctx context.Context) {
	s.Log.Info("Start syncing tunnels from Bandchain")
	tunnels, err := s.BandClient.GetTunnels(ctx)
	if err != nil {
		s.Log.Error("Failed to fetch tunnels from BandChain", err)
		return
	}

	newTunnels := tunnels[min(s.bandLatestTunnel, len(tunnels)):]

	validTunnels := s.filterTunnels(newTunnels)
	if len(validTunnels) == 0 {
		s.Log.Info("No new tunnels to sync")
		return
	}

	s.setTunnelRelayer(validTunnels)
	s.bandLatestTunnel = len(tunnels)

	// update the valid tunnel IDs in the packet handler
	for _, tunnel := range validTunnels {
		go func() {
			tr := s.tunnelRelayers[tunnel.ID]
			_ = s.TriggerTunnelRelayer(ctx, tr)
		}()
	}
}

// initialize initializes the scheduler and execute subroutines.
func (s *Scheduler) initialize(ctx context.Context) error {
	subscribers := []subscriber.Subscriber{
		subscriber.NewPacketSuccessSubscriber(s.Log, s.relayTunnelIDCh, s.SubscriptionTimeout),
		subscriber.NewManualTriggerSubscriber(s.Log, s.relayTunnelIDCh, s.SubscriptionTimeout),
	}
	s.BandClient.SetSubscribers(subscribers)

	if err := s.BandClient.Subscribe(ctx); err != nil {
		s.Log.Error("Failed to subscribe to BandChain", err)
		return err
	}

	// listen events from BandChain
	for _, subscriber := range subscribers {
		go subscriber.HandleEvent(ctx)
	}

	// handle trigger relayer event from packet handler
	go s.handleTriggerTunnelRelayer(ctx)

	return nil
}

// handleTriggerTunnelRelayer triggers the tunnel relayer from the received tunnelID.
func (s *Scheduler) handleTriggerTunnelRelayer(ctx context.Context) {
	for tunnelID := range s.relayTunnelIDCh {
		tunnelRelayer, ok := s.tunnelRelayers[tunnelID]
		if !ok {
			continue
		}

		s.Log.Info("Received trigger relayer event", "tunnel_id", tunnelID)

		if tunnelRelayer.penaltySkipRemaining > 0 {
			s.Log.Info(
				"Skipping tunnel execution due to penalty from previous failure",
				"tunnel_id", tunnelID,
			)
			continue
		}

		go func() {
			status := s.TriggerTunnelRelayer(ctx, tunnelRelayer)
			if status != RelayStatusSuccess {
				s.Log.Info(
					"Tunnel relay completed with non-success status",
					"tunnel_id", tunnelID,
					"status", string(status),
				)
			}
		}()
	}
}

// isSupportedTunnel checks if the tunnel is supported by the scheduler.
func (s *Scheduler) isSupportedTunnel(tunnel bandtypes.Tunnel) bool {
	if _, ok := s.ChainProviders[tunnel.TargetChainID]; !ok {
		return false
	}

	if s.tunnelCreator != "" && tunnel.Creator != s.tunnelCreator {
		return false
	}

	return true
}

// filterTunnels selects only the supported tunnel and returns the valid tunnels.
func (s *Scheduler) filterTunnels(tunnels []bandtypes.Tunnel) []bandtypes.Tunnel {
	var validTunnels []bandtypes.Tunnel
	for _, tunnel := range tunnels {
		if !s.isSupportedTunnel(tunnel) {
			s.Log.Warn(
				"The program does not support this tunnel",
				"chain_name", tunnel.TargetChainID,
				"tunnel_id", tunnel.ID,
			)
			continue
		}

		validTunnels = append(validTunnels, tunnel)
	}

	return validTunnels
}

// setTunnelRelayer sets the tunnel relayer from the given tunnels.
func (s *Scheduler) setTunnelRelayer(tunnels []bandtypes.Tunnel) {
	for _, tunnel := range tunnels {
		chainProvider := s.ChainProviders[tunnel.TargetChainID]
		tr := NewTunnelRelayer(
			s.Log,
			tunnel.ID,
			s.CheckingPacketInterval,
			s.BandClient,
			chainProvider,
			s.Alert,
		)

		s.tunnelRelayers[tunnel.ID] = &tr

		// update the metric for the number of tunnels per destination chain
		relayermetrics.IncTunnelsPerDestinationChain(tunnel.TargetChainID)
		s.Log.Info(
			"New tunnel is set into the scheduler",
			"chain_name", tunnel.TargetChainID,
			"tunnel_id", tunnel.ID,
		)
	}
}
