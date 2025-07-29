package relayer

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/band/subscriber"
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
	PacketHandler  *PacketHandler

	triggerRelayerCh chan uint64
	signingIDCh      chan uint64
	tunnelIDCh       chan uint64
	tunnelRelayers   map[uint64]*TunnelRelayer
	bandLatestTunnel int
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
	triggerRelayerCh := make(chan uint64, 1000)
	signingIDCh := make(chan uint64, 1000)
	tunnelIDCh := make(chan uint64, 1000)

	packetHandler := NewPacketHandler(
		log,
		bandClient,
		triggerRelayerCh,
		signingIDCh,
		tunnelIDCh,
	)

	return &Scheduler{
		Log:                    log,
		CheckingPacketInterval: checkingPacketInterval,
		SyncTunnelsInterval:    syncTunnelsInterval,
		PenaltySkipRounds:      penaltySkipRounds,
		BandClient:             bandClient,
		ChainProviders:         chainProviders,
		PacketHandler:          packetHandler,
		triggerRelayerCh:       triggerRelayerCh,
		signingIDCh:            signingIDCh,
		tunnelIDCh:             tunnelIDCh,
		tunnelRelayers:         make(map[uint64]*TunnelRelayer),
		bandLatestTunnel:       0,
	}
}

// Start starts all tunnel relayers
func (s *Scheduler) Start(
	ctx context.Context,
	tunnelIDs []uint64,
	tunnelCreator string,
	subscriptionTimeout time.Duration,
) error {
	subscribers := []subscriber.Subscriber{
		subscriber.NewPacketSuccessSubscriber(s.Log, s.tunnelIDCh, subscriptionTimeout),
		subscriber.NewManualTriggerSubscriber(s.Log, s.tunnelIDCh, subscriptionTimeout),
		subscriber.NewSigningSuccessSubscriber(s.Log, s.signingIDCh, subscriptionTimeout),
		subscriber.NewSigningFailedSubscriber(s.Log, s.signingIDCh, subscriptionTimeout),
	}
	s.BandClient.SetSubscribers(subscribers)

	if err := s.BandClient.Subscribe(ctx); err != nil {
		s.Log.Error("Failed to subscribe to BandChain", zap.Error(err))
		return err
	}

	s.SyncTunnels(ctx, tunnelIDs, tunnelCreator)

	// listen events from BandChain
	for _, subscriber := range subscribers {
		go subscriber.HandleEvent(ctx)
	}

	// handle new packets and failed or successful signing IDs
	go s.PacketHandler.HandleNewPacket(ctx)
	go s.PacketHandler.HandleSigningResult()

	// handle trigger relayer event from packet handler
	go s.HandleTriggerTunnelRelayer(ctx)

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
	s.Log.Info("Executing tunnel relayers from the scheduler")

	for _, tr := range s.tunnelRelayers {
		if tr.penaltySkipRemaining > 0 {
			s.Log.Debug(
				"Skipping tunnel execution due to penalty from previous failure.",
				zap.Uint64("tunnel_id", tr.TunnelID),
				zap.Uint("penalty_skip_remaining", tr.penaltySkipRemaining),
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
			zap.Error(err),
			zap.Uint64("tunnel_id", tr.TunnelID),
		)

		return RelayStatusFailed
	}

	switch relayStatus {
	case RelayStatusExecuting:
		// Record metrics for the skipped task execution
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

		s.Log.Info("Relay packet successfully", zap.Uint64("tunnel_id", tr.TunnelID))
	}

	return relayStatus
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

	var newTunnelIDs []uint64
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

		newTunnelIDs = append(newTunnelIDs, tunnels[i].ID)
		s.tunnelRelayers[tunnels[i].ID] = &tr

		// update the metric for the number of tunnels per destination chain
		relayermetrics.IncTunnelsPerDestinationChain(tunnels[i].TargetChainID)

		s.Log.Info(
			"New tunnel synchronized successfully",
			zap.String("chain_name", tunnels[i].TargetChainID),
			zap.Uint64("tunnel_id", tunnels[i].ID),
		)
	}

	// update the valid tunnel IDs in the packet handler
	s.PacketHandler.UpdateValidTunnelIDs(newTunnelIDs)
	for _, tunnelID := range newTunnelIDs {
		s.tunnelIDCh <- tunnelID
	}

	s.bandLatestTunnel = len(tunnels)
}

// HandleTriggerTunnelRelayer triggers the tunnel relayer from the received tunnelID.
func (s *Scheduler) HandleTriggerTunnelRelayer(ctx context.Context) {
	for tunnelID := range s.triggerRelayerCh {
		tunnelRelayer, ok := s.tunnelRelayers[tunnelID]
		if !ok {
			return
		}

		s.Log.Info("Received trigger relayer event", zap.Uint64("tunnel_id", tunnelID))

		if tunnelRelayer.penaltySkipRemaining > 0 {
			s.Log.Info(
				"Skipping tunnel execution due to penalty from previous failure",
				zap.Uint64("tunnel_id", tunnelID),
			)
			return
		}

		go func() {
			status := s.TriggerTunnelRelayer(ctx, tunnelRelayer)
			if status != RelayStatusSuccess {
				s.Log.Info(
					"relay tunnel not success",
					zap.Uint64("tunnel_id", tunnelID),
					zap.String("status", string(status)),
				)
			}
		}()
	}
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
