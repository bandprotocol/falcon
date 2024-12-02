package relayer

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/chains"
)

// TunnelRelayer is a relayer that listens to the tunnel and relays the packet
type TunnelRelayer struct {
	Log                    *zap.Logger
	TunnelID               uint64
	ContractAddress        string
	CheckingPacketInterval time.Duration
	BandClient             band.Client
	TargetChainProvider    chains.ChainProvider

	isExecuting bool
}

// NewTunnelRelayer creates a new TunnelRelayer
func NewTunnelRelayer(
	log *zap.Logger,
	tunnelID uint64,
	contractAddress string,
	checkingPacketInterval time.Duration,
	bandClient band.Client,
	targetChainProvider chains.ChainProvider,
) TunnelRelayer {
	return TunnelRelayer{
		Log:                    log,
		TunnelID:               tunnelID,
		ContractAddress:        contractAddress,
		CheckingPacketInterval: checkingPacketInterval,
		BandClient:             bandClient,
		TargetChainProvider:    targetChainProvider,
		isExecuting:            false,
	}
}

// CheckAndRelay checks the tunnel and relays the packet
func (t *TunnelRelayer) CheckAndRelay(ctx context.Context) (err error) {
	t.isExecuting = true
	defer func() {
		t.isExecuting = false

		// Recover from panic
		if r := recover(); r != nil {
			newErr, ok := r.(error)
			if !ok {
				newErr = fmt.Errorf("%v", r)
			}
			err = newErr
		}
	}()

	// Query tunnel info from BandChain
	tunnelBandInfo, err := t.BandClient.GetTunnel(ctx, t.TunnelID)
	if err != nil {
		t.Log.Error(
			"failed to get tunnel",
			zap.Error(err),
			zap.Uint64("tunnel_id", t.TunnelID),
		)
		return err
	}

	// Query tunnel info from TargetChain
	tunnelChainInfo, err := t.TargetChainProvider.QueryTunnelInfo(ctx, t.TunnelID, t.ContractAddress)
	if err != nil {
		return err
	}

	if !tunnelChainInfo.IsActive {
		t.Log.Error("tunnel is not active on target chain", zap.Uint64("tunnel_id", t.TunnelID))
		return fmt.Errorf("tunnel is not active on target chain")
	}

	// end process if current packet is already relayed
	nextSeq := tunnelChainInfo.LatestSequence + 1
	if tunnelBandInfo.LatestSequence < nextSeq {
		t.Log.Info(
			"no new packet to relay",
			zap.Uint64("tunnel_id", t.TunnelID),
			zap.Uint64("sequence", tunnelChainInfo.LatestSequence),
		)
		return nil
	}

	// Relay packets
	for seq := nextSeq; seq <= tunnelBandInfo.LatestSequence; seq++ {
		t.Log.Info(
			"relaying packet",
			zap.Uint64("tunnel_id", t.TunnelID),
			zap.Uint64("sequence", seq),
		)

		packet, err := t.BandClient.GetTunnelPacket(ctx, t.TunnelID, seq)
		if err != nil {
			t.Log.Error(
				"failed to get packet",
				zap.Error(err),
				zap.Uint64("tunnel_id", t.TunnelID),
				zap.Uint64("sequence", seq),
			)
			return err
		}

		if err := t.TargetChainProvider.RelayPacket(ctx, packet); err != nil {
			t.Log.Error(
				"failed to relay packet",
				zap.Error(err),
				zap.Uint64("tunnel_id", t.TunnelID),
				zap.Uint64("sequence", seq),
			)
			return err
		}

		t.Log.Info(
			"successfully relayed packet",
			zap.Uint64("tunnel_id", t.TunnelID),
			zap.Uint64("sequence", seq),
		)
	}

	return nil
}

func (t *TunnelRelayer) IsExecuting() bool {
	return t.isExecuting
}
