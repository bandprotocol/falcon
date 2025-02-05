package relayer

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
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
		Log:                    log.With(zap.Uint64("tunnel_id", tunnelID)),
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

	for {
		// Query tunnel info from BandChain
		tunnelBandInfo, err := t.BandClient.GetTunnel(ctx, t.TunnelID)
		if err != nil {
			t.Log.Error("Failed to get tunnel", zap.Error(err))
			return err
		}

		// Query tunnel info from TargetChain
		tunnelChainInfo, err := t.TargetChainProvider.QueryTunnelInfo(ctx, t.TunnelID, t.ContractAddress)
		if err != nil {
			return err
		}
		if !tunnelChainInfo.IsActive {
			t.Log.Info("Tunnel is not active on target chain")
			return nil
		}

		// end process if current packet is already relayed
		seq := tunnelChainInfo.LatestSequence + 1
		if tunnelBandInfo.LatestSequence < seq {
			t.Log.Info("No new packet to relay", zap.Uint64("sequence", tunnelChainInfo.LatestSequence))
			return nil
		}

		t.Log.Info("Relaying packet", zap.Uint64("sequence", seq))

		// get packet of the sequence
		packet, err := t.BandClient.GetTunnelPacket(ctx, t.TunnelID, seq)
		if err != nil {
			t.Log.Error("Failed to get packet", zap.Error(err), zap.Uint64("sequence", seq))
			return err
		}

		// Check signing status; if it is waiting, wait for the completion of the EVM signature.
		// If it is not success (Failed or Undefined), return error.
		signing := packet.CurrentGroupSigning
		if signing == nil {
			signing = packet.IncomingGroupSigning
		}

		signingStatus := tsstypes.SigningStatus(tsstypes.SigningStatus_value[signing.Status])
		if signingStatus == tsstypes.SIGNING_STATUS_WAITING {
			t.Log.Info(
				"The current packet must wait for the completion of the EVM signature",
				zap.Uint64("sequence", seq),
			)
			return nil
		} else if signingStatus != tsstypes.SIGNING_STATUS_SUCCESS {
			err := fmt.Errorf("signing status is not success")
			t.Log.Error("Failed to relay packet", zap.Error(err), zap.Uint64("sequence", seq))
			return err
		}

		// Relay the packet to the target chain
		if err := t.TargetChainProvider.RelayPacket(ctx, packet); err != nil {
			t.Log.Error("Failed to relay packet", zap.Error(err), zap.Uint64("sequence", seq))
			return err
		}

		t.Log.Info("Successfully relayed packet", zap.Uint64("sequence", seq))
	}
}

func (t *TunnelRelayer) IsExecuting() bool {
	return t.isExecuting
}
