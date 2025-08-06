package relayer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/band/types"
	"github.com/bandprotocol/falcon/relayer/chains"
	chaintypes "github.com/bandprotocol/falcon/relayer/chains/types"
	"github.com/bandprotocol/falcon/relayer/logger"
)

type RelayStatus string

const (
	RelayStatusSuccess   RelayStatus = "success"
	RelayStatusExecuting RelayStatus = "executing"
	RelayStatusSkipped   RelayStatus = "skipped"
	RelayStatusFailed    RelayStatus = "failed"
)

// TunnelRelayer is a relayer that listens to the tunnel and relays the packet
type TunnelRelayer struct {
	Log                    logger.Logger
	TunnelID               uint64
	CheckingPacketInterval time.Duration
	BandClient             band.Client
	TargetChainProvider    chains.ChainProvider

	isTargetChainActive  bool
	penaltySkipRemaining uint
	mu                   *sync.Mutex
}

// NewTunnelRelayer creates a new TunnelRelayer
func NewTunnelRelayer(
	log logger.Logger,
	tunnelID uint64,
	checkingPacketInterval time.Duration,
	bandClient band.Client,
	targetChainProvider chains.ChainProvider,
) TunnelRelayer {
	return TunnelRelayer{
		Log:                    log.With(zap.Uint64("tunnel_id", tunnelID)),
		TunnelID:               tunnelID,
		CheckingPacketInterval: checkingPacketInterval,
		BandClient:             bandClient,
		TargetChainProvider:    targetChainProvider,
		isTargetChainActive:    false,
		penaltySkipRemaining:   0,
		mu:                     &sync.Mutex{},
	}
}

// CheckAndRelay checks the tunnel and relays the packet
func (t *TunnelRelayer) CheckAndRelay(
	ctx context.Context,
	isForce bool,
) (relayStatus RelayStatus, err error) {
	// if the tunnel relayer is executing, skip the round
	if !t.mu.TryLock() {
		t.Log.Debug("Skip this tunnel: tunnel relayer is executing on another process")
		return RelayStatusExecuting, nil
	}
	defer func() {
		t.mu.Unlock()

		// Recover from panic
		if r := recover(); r != nil {
			newErr, ok := r.(error)
			if !ok {
				newErr = fmt.Errorf("%v", r)
			}

			relayStatus = RelayStatusFailed
			err = newErr
		}
	}()

	t.Log.Debug("Executing task")
	isPacketRelayed := false
	for {
		// get next packet sequence to relay
		seq, err := t.getNextPacketSequence(ctx, isForce)
		if err != nil {
			return RelayStatusFailed, err
		}
		if seq == 0 {
			break
		}
		t.Log.Debug("Next packet sequence to relay", zap.Uint64("sequence", seq))
		// get packet of the sequence
		packet, err := t.getTunnelPacket(ctx, seq)
		if err != nil {
			return RelayStatusFailed, err
		}

		// relay the packet
		if err := t.relayPacket(ctx, packet); err != nil {
			return RelayStatusFailed, err
		}

		isPacketRelayed = true
	}

	if !isPacketRelayed {
		return RelayStatusSkipped, nil
	}
	return RelayStatusSuccess, nil
}

// getNextPacketSequence returns the next packet sequence to relay. Sequence 0 is returned
// if the tunnel status on BandChain is inactive (and not being forced) or the target contract
// is inactive or the current packet is already relayed.
func (t *TunnelRelayer) getNextPacketSequence(ctx context.Context, isForce bool) (uint64, error) {
	// Query tunnel info from BandChain
	tunnelInfo, err := t.BandClient.GetTunnel(ctx, t.TunnelID)
	if err != nil {
		t.Log.Error("Failed to get tunnel", zap.Error(err))
		return 0, err
	}

	// exit if the tunnel is not active and isForce is false
	if !isForce && !tunnelInfo.IsActive {
		t.Log.Debug("Tunnel is not active on BandChain")
		return 0, nil
	}

	// Query tunnel info from TargetChain
	targetContractInfo, err := t.TargetChainProvider.QueryTunnelInfo(
		ctx,
		t.TunnelID,
		tunnelInfo.TargetAddress,
	)
	if err != nil {
		t.Log.Error("Failed to get target contract info", zap.Error(err))
		return 0, err
	}

	t.updateRelayerMetrics(tunnelInfo, targetContractInfo)

	// check if the target contract is active
	t.isTargetChainActive = targetContractInfo.IsActive
	if !t.isTargetChainActive {
		t.Log.Debug("Tunnel is not active on target chain")
		return 0, nil
	}

	// end process if current packet is already relayed
	latestSeq := targetContractInfo.LatestSequence
	nextSeq := latestSeq + 1
	if tunnelInfo.LatestSequence < nextSeq {
		t.Log.Debug("No new packet to relay", zap.Uint64("sequence", latestSeq))
		return 0, nil
	}

	return nextSeq, nil
}

// updateRelayerMetrics updates the metrics for the relayer.
func (t *TunnelRelayer) updateRelayerMetrics(
	tunnelInfo *types.Tunnel,
	targetContractInfo *chaintypes.Tunnel,
) {
	// update the metric for unrelayed packets based on the difference
	// between the latest sequences on BandChain and the target chain
	unrelayedPackets := tunnelInfo.LatestSequence - targetContractInfo.LatestSequence
	relayermetrics.SetUnrelayedPackets(t.TunnelID, unrelayedPackets)

	// update the metric for the number of active target contracts
	if targetContractInfo.IsActive && !t.isTargetChainActive {
		relayermetrics.IncActiveTargetContractsCount(tunnelInfo.TargetChainID)
	} else if !targetContractInfo.IsActive && t.isTargetChainActive {
		relayermetrics.DecActiveTargetContractsCount(tunnelInfo.TargetChainID)
	}
}

// relayPacket relays the packet to the target chain.
func (t *TunnelRelayer) relayPacket(ctx context.Context, packet *types.Packet) error {
	t.Log.Info("Relaying packet", zap.Uint64("sequence", packet.Sequence))

	// Relay the packet to the target chain
	if err := t.TargetChainProvider.RelayPacket(ctx, packet); err != nil {
		t.Log.Error("Failed to relay packet", zap.Error(err), zap.Uint64("sequence", packet.Sequence))
		return err
	}

	// Increment the metric for successfully relayed packets
	relayermetrics.IncPacketsRelayedSuccess(t.TunnelID)
	t.Log.Info("Successfully relayed packet", zap.Uint64("sequence", packet.Sequence))

	return nil
}

// getTunnelPacket polls BandChain for the packet with the given sequence
// until its TSS signing status becomes SUCCESS, then returns it.
func (t *TunnelRelayer) getTunnelPacket(ctx context.Context, seq uint64) (*types.Packet, error) {
	for {
		// get packet of the sequence
		packet, err := t.BandClient.GetTunnelPacket(ctx, t.TunnelID, seq)
		if err != nil {
			t.Log.Error("Failed to get packet", zap.Error(err), zap.Uint64("sequence", seq))
			return nil, err
		}
		// Check signing status; if it is waiting, wait for the completion of the EVM signature.
		// If it is not success (Failed or Undefined), return error.
		signing := packet.CurrentGroupSigning
		if signing == nil ||
			signing.SigningStatus == tsstypes.SIGNING_STATUS_FALLEN {
			signing = packet.IncomingGroupSigning
		}

		if signing.SigningStatus == tsstypes.SIGNING_STATUS_WAITING {
			t.Log.Debug(
				"The current packet must wait for the completion of the EVM signature",
				zap.Uint64("sequence", seq),
			)
			// wait 1 secs for each block
			time.Sleep(time.Second)
			continue
		} else if signing.SigningStatus != tsstypes.SIGNING_STATUS_SUCCESS {
			err := fmt.Errorf("signing status is not success")
			t.Log.Error("Failed to relay packet", zap.Error(err), zap.Uint64("sequence", seq))
			return nil, err
		}

		return packet, nil
	}
}
