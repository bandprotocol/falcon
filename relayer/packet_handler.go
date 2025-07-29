package relayer

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/band/subscriber"
	"github.com/bandprotocol/falcon/relayer/band/types"
)

const (
	cacheTTL = 24 * time.Hour
)

// CacheEntry represents a cache entry for packet information
type CacheEntry struct {
	TunnelID               uint64
	CurrentGroupSigningID  uint64
	IncomingGroupSigningID uint64
	CurrentGroupStatus     tsstypes.SigningStatus
	IncomingGroupStatus    tsstypes.SigningStatus
}

// NewCacheEntry creates a new cache entry.
func NewCacheEntry(
	tunnelID uint64,
	currentGroupSigningID uint64,
	incomingGroupSigningID uint64,
	currentGroupStatus tsstypes.SigningStatus,
	incomingGroupStatus tsstypes.SigningStatus,
) *CacheEntry {
	return &CacheEntry{
		TunnelID:               tunnelID,
		CurrentGroupSigningID:  currentGroupSigningID,
		IncomingGroupSigningID: incomingGroupSigningID,
		CurrentGroupStatus:     currentGroupStatus,
		IncomingGroupStatus:    incomingGroupStatus,
	}
}

// PacketHandler handles new signing IDs and packets.
type PacketHandler struct {
	Log              *zap.Logger
	BandClient       band.Client
	SigningResultCh  <-chan subscriber.SigningResult
	TunnelIDCh       <-chan uint64
	TriggerRelayerCh chan<- uint64
	ValidTunnelIDs   map[uint64]struct{}
	signingCache     *ttlcache.Cache[uint64, *CacheEntry]
}

// NewPacketHandler creates a new PacketHandler.
func NewPacketHandler(
	log *zap.Logger,
	bandClient band.Client,
	triggerRelayerCh chan<- uint64,
	signingResultCh <-chan subscriber.SigningResult,
	tunnelIDCh <-chan uint64,
) *PacketHandler {
	cache := ttlcache.New(
		ttlcache.WithTTL[uint64, *CacheEntry](cacheTTL),
		ttlcache.WithDisableTouchOnHit[uint64, *CacheEntry](),
	)

	cache.OnInsertion(
		func(ctx context.Context, item *ttlcache.Item[uint64, *CacheEntry]) {
			relayermetrics.IncreasePendingNewSigning()
		},
	)
	cache.OnEviction(
		func(
			ctx context.Context,
			reason ttlcache.EvictionReason,
			item *ttlcache.Item[uint64, *CacheEntry],
		) {
			relayermetrics.DecreasePendingNewSigning()
		},
	)

	return &PacketHandler{
		Log:              log,
		BandClient:       bandClient,
		SigningResultCh:  signingResultCh,
		TunnelIDCh:       tunnelIDCh,
		TriggerRelayerCh: triggerRelayerCh,
		ValidTunnelIDs:   make(map[uint64]struct{}),
		signingCache:     cache,
	}
}

// HandleSigningResult handles signingResult received from the channel.
func (h *PacketHandler) HandleSigningResult() {
	for signingResult := range h.SigningResultCh {
		if signingResult.IsSuccess {
			h.handleSigningSuccess(signingResult.SigningID)
		} else {
			h.handleSigningFailure(signingResult.SigningID)
		}
	}
}

// handleSigningSuccess handles signingID that is successfully signed.
func (h *PacketHandler) handleSigningSuccess(signingID uint64) {
	item := h.signingCache.Get(signingID)
	if item == nil {
		return
	}

	cacheEntry := item.Value()
	incomingSigningID := cacheEntry.IncomingGroupSigningID
	currentSigningID := cacheEntry.CurrentGroupSigningID
	canSkipCurrentSigning := currentSigningID == 0 ||
		cacheEntry.CurrentGroupStatus == tsstypes.SIGNING_STATUS_FALLEN

	h.Log.Debug("Found matching signing ID in cache",
		zap.Uint64("signing_id", signingID),
		zap.Uint64("tunnel_id", cacheEntry.TunnelID),
		zap.Uint64("current_group_signing_id", currentSigningID),
		zap.Uint64("incoming_group_signing_id", incomingSigningID),
	)

	// Delete the signing ID from the cache.
	h.signingCache.Delete(signingID)

	if signingID == currentSigningID {
		// Delete the incoming group SigningID if it exists and trigger relayer.
		if incomingSigningID != 0 && h.signingCache.Has(incomingSigningID) {
			h.signingCache.Delete(incomingSigningID)
		}

		h.Log.Debug("Triggered relayer for tunnel", zap.Uint64("tunnel_id", cacheEntry.TunnelID))
		h.TriggerRelayerCh <- cacheEntry.TunnelID
	} else if signingID == incomingSigningID && canSkipCurrentSigning {
		// currentGroupSigning ID shouldn't exist, just trigger relayer.
		h.Log.Debug("Triggered relayer for tunnel", zap.Uint64("tunnel_id", cacheEntry.TunnelID))
		h.TriggerRelayerCh <- cacheEntry.TunnelID
	} else {
		// have to wait for the current group signing to be successful.
		// update the incoming group signing status to success on currentGroupSigning ID key.
		cacheEntry.IncomingGroupStatus = tsstypes.SIGNING_STATUS_SUCCESS
		h.signingCache.Set(currentSigningID, cacheEntry, cacheTTL)

		h.Log.Debug("Updated incoming group signing status to success",
			zap.Uint64("tunnel_id", cacheEntry.TunnelID),
			zap.Uint64("incoming_group_signing_id", incomingSigningID),
			zap.Uint64("current_group_signing_id", currentSigningID),
		)
	}
}

// handleSigningFailure handles the signingID that is failed to sign.
func (h *PacketHandler) handleSigningFailure(signingID uint64) {
	item := h.signingCache.Get(signingID)
	if item == nil {
		return
	}

	cacheEntry := item.Value()
	incomingSigningID := cacheEntry.IncomingGroupSigningID
	currentSigningID := cacheEntry.CurrentGroupSigningID
	incomingSigningStatus := cacheEntry.IncomingGroupStatus

	// Delete the signing ID from the cache.
	h.signingCache.Delete(signingID)

	if signingID == currentSigningID && incomingSigningStatus == tsstypes.SIGNING_STATUS_SUCCESS {
		// if it is the current group signing ID and the incoming group signing is successful,
		// trigger relayer.
		h.Log.Debug("Triggered relayer for tunnel", zap.Uint64("tunnel_id", cacheEntry.TunnelID))
		h.TriggerRelayerCh <- cacheEntry.TunnelID
	} else if signingID == currentSigningID && h.signingCache.Has(incomingSigningID) {
		// if it is the current group signing ID and the incoming group signing ID exists,
		// update the current group signing status to failed in incomingGroupSigning ID key.
		cacheEntry.CurrentGroupStatus = tsstypes.SIGNING_STATUS_FALLEN
		h.signingCache.Set(incomingSigningID, cacheEntry, cacheTTL)

		h.Log.Debug("Updated current group signing status to failed",
			zap.Uint64("tunnel_id", cacheEntry.TunnelID),
			zap.Uint64("current_group_signing_id", currentSigningID),
		)
	}
}

// HandleNewPacket handles new packets received from the channel. If the packet is valid,
// the requested signing IDs of the packet are cached.
func (h *PacketHandler) HandleNewPacket(ctx context.Context) {
	for tunnelID := range h.TunnelIDCh {
		if _, ok := h.ValidTunnelIDs[tunnelID]; !ok {
			continue
		}

		tunnel, err := h.BandClient.GetTunnel(ctx, tunnelID)
		if err != nil {
			h.Log.Error(
				"Failed to get tunnel",
				zap.Error(err),
				zap.Uint64("tunnel_id", tunnelID),
			)
			continue
		}

		// skip tunnel that doesn't produce a packet.
		if tunnel.LatestSequence == 0 {
			h.Log.Debug("Tunnel doesn't produce a packet", zap.Uint64("tunnel_id", tunnelID))
			continue
		}

		packet, err := h.BandClient.GetTunnelPacket(ctx, tunnelID, tunnel.LatestSequence)
		if err != nil {
			h.Log.Error(
				"Failed to get latest packet",
				zap.Error(err),
				zap.Uint64("tunnel_id", tunnelID),
			)
			continue
		}

		h.Log.Debug("Received new packet",
			zap.Uint64("tunnel_id", packet.TunnelID),
			zap.Uint64("sequence", packet.Sequence),
		)

		if h.isSigningCompleted(packet) {
			h.TriggerRelayerCh <- packet.TunnelID
			continue
		}

		h.cacheSigningNewPacket(packet)
	}
}

// isSigningCompleted checks if the signing is completed.
func (h *PacketHandler) isSigningCompleted(packet *types.Packet) bool {
	currentSigning := packet.CurrentGroupSigning
	incomingSigning := packet.IncomingGroupSigning
	canSkipCurrentSigning := currentSigning == nil ||
		currentSigning.SigningStatus == tsstypes.SIGNING_STATUS_FALLEN

	// trigger relayer if the current group signing is successful or
	// incoming group signing is successful and there is only an incoming signing.
	if currentSigning != nil &&
		currentSigning.SigningStatus == tsstypes.SIGNING_STATUS_SUCCESS {
		h.Log.Debug("Current group signing is successful",
			zap.Uint64("tunnel_id", packet.TunnelID),
			zap.Uint64("sequence", packet.Sequence),
		)

		return true
	} else if incomingSigning != nil &&
		canSkipCurrentSigning &&
		incomingSigning.SigningStatus == tsstypes.SIGNING_STATUS_SUCCESS {
		h.Log.Debug("Incoming group signing is successful",
			zap.Uint64("tunnel_id", packet.TunnelID),
			zap.Uint64("sequence", packet.Sequence),
		)

		return true
	}

	return false
}

// UpdateValidTunnelIDs updates the valid tunnel IDs.
func (h *PacketHandler) UpdateValidTunnelIDs(tunnelIDs []uint64) {
	for _, tunnelID := range tunnelIDs {
		h.ValidTunnelIDs[tunnelID] = struct{}{}
	}
}

// cacheSigningNewPacket caches the signing IDs of the given packet.
func (h *PacketHandler) cacheSigningNewPacket(packet *types.Packet) {
	currentSigningID, incomingSigningID := uint64(0), uint64(0)
	currentSigningStatus := tsstypes.SIGNING_STATUS_UNSPECIFIED
	incomingSigningStatus := tsstypes.SIGNING_STATUS_UNSPECIFIED

	if packet.CurrentGroupSigning != nil {
		currentSigningID = packet.CurrentGroupSigning.ID
		currentSigningStatus = packet.CurrentGroupSigning.SigningStatus
	}
	if packet.IncomingGroupSigning != nil {
		incomingSigningID = packet.IncomingGroupSigning.ID
		incomingSigningStatus = packet.IncomingGroupSigning.SigningStatus
	}

	cacheEntry := NewCacheEntry(
		packet.TunnelID,
		currentSigningID,
		incomingSigningID,
		currentSigningStatus,
		incomingSigningStatus,
	)

	// Cache the signing ID if it is not in the cache.
	if currentSigningID != 0 && !h.signingCache.Has(currentSigningID) {
		h.signingCache.Set(currentSigningID, cacheEntry, cacheTTL)
	}
	if incomingSigningID != 0 && !h.signingCache.Has(incomingSigningID) {
		h.signingCache.Set(incomingSigningID, cacheEntry, cacheTTL)
	}

	h.Log.Debug("Stored current group signing ID in cache",
		zap.Uint64("tunnel_id", packet.TunnelID),
		zap.Uint64("sequence", packet.Sequence),
		zap.Uint64("current_group_signing_id", cacheEntry.CurrentGroupSigningID),
		zap.Uint64("incoming_group_signing_id", cacheEntry.IncomingGroupSigningID),
	)
}
