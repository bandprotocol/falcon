package relayer

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/band"
	"github.com/bandprotocol/falcon/relayer/band/types"
)

const (
	cacheTTL = 24 * time.Hour
)

// CacheEntry represents a cache entry for packet information
type CacheEntry struct {
	TunnelID  uint64
	SigningID uint64
}

// NewCacheEntry creates a new cache entry.
func NewCacheEntry(
	tunnelID uint64,
	signingID uint64,
) *CacheEntry {
	return &CacheEntry{
		TunnelID:  tunnelID,
		SigningID: signingID,
	}
}

// PacketHandler handles new signing IDs and packets.
type PacketHandler struct {
	Log              *zap.Logger
	BandClient       band.Client
	SigningIDCh      <-chan uint64
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
	signingIDCh <-chan uint64,
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
		SigningIDCh:      signingIDCh,
		TunnelIDCh:       tunnelIDCh,
		TriggerRelayerCh: triggerRelayerCh,
		ValidTunnelIDs:   make(map[uint64]struct{}),
		signingCache:     cache,
	}
}

// HandleSigningResult handles signingResult received from the channel.
func (h *PacketHandler) HandleSigningResult() {
	for signingID := range h.SigningIDCh {
		item := h.signingCache.Get(signingID)
		if item == nil {
			return
		}

		cacheEntry := item.Value()
		h.Log.Debug("Found matching signing ID in cache",
			zap.Uint64("signing_id", signingID),
			zap.Uint64("tunnel_id", cacheEntry.TunnelID),
		)

		// Delete the signing ID from the cache.
		h.signingCache.Delete(signingID)

		h.Log.Debug("Triggered relayer for tunnel", zap.Uint64("tunnel_id", cacheEntry.TunnelID))
		h.TriggerRelayerCh <- cacheEntry.TunnelID
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

		if packet.CurrentGroupSigning != nil {
			h.cacheSigning(packet.TunnelID, packet.CurrentGroupSigning.ID)
		}

		if packet.IncomingGroupSigning != nil {
			h.cacheSigning(packet.TunnelID, packet.IncomingGroupSigning.ID)
		}
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

// cacheSigning caches the signing ID.
func (h *PacketHandler) cacheSigning(tunnelID uint64, signingID uint64) {
	entry := NewCacheEntry(
		tunnelID,
		signingID,
	)
	h.signingCache.Set(signingID, entry, cacheTTL)

	h.Log.Debug("Stored signing ID in cache",
		zap.Uint64("tunnel_id", tunnelID),
		zap.Uint64("signing_id", signingID),
	)
}
