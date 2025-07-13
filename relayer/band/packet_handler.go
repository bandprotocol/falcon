package band

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	"github.com/bandprotocol/falcon/internal/relayermetrics"
	"github.com/bandprotocol/falcon/relayer/band/types"
)

const (
	cleanupCacheInterval = 1 * time.Hour
	cacheTTL             = 24 * time.Hour
)

// CacheEntry represents a cache entry for packet information
type CacheEntry struct {
	TunnelID               uint64
	CurrentGroupSigningID  uint64
	IncomingGroupSigningID uint64
}

// NewCacheEntry creates a new cache entry.
func NewCacheEntry(
	tunnelID uint64,
	currentGroupSigningID uint64,
	incomingGroupSigningID uint64,
) *CacheEntry {
	return &CacheEntry{
		TunnelID:               tunnelID,
		CurrentGroupSigningID:  currentGroupSigningID,
		IncomingGroupSigningID: incomingGroupSigningID,
	}
}

// PacketHandler handles new signing IDs and packets.
type PacketHandler struct {
	Log              *zap.Logger
	NewSigningIDCh   <-chan uint64
	NewPacketCh      <-chan *types.Packet
	TriggerRelayerCh chan<- uint64
	ValidTunnelIDs   map[uint64]struct{}
	signingCache     *ttlcache.Cache[uint64, *CacheEntry]
}

// NewPacketHandler creates a new PacketHandler.
func NewPacketHandler(
	log *zap.Logger,
	triggerRelayerCh chan<- uint64,
	newSigningIDCh <-chan uint64,
	newPacketCh <-chan *types.Packet,
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
		NewSigningIDCh:   newSigningIDCh,
		NewPacketCh:      newPacketCh,
		TriggerRelayerCh: triggerRelayerCh,
		ValidTunnelIDs:   make(map[uint64]struct{}),
		signingCache:     cache,
	}
}

// HandleNewSigning triggers relayer for new signing IDs that are in the cache.
func (h *PacketHandler) HandleNewSigning() {
	for signingID := range h.NewSigningIDCh {
		item := h.signingCache.Get(signingID)
		if item == nil {
			continue
		}

		cacheEntry := item.Value()
		h.Log.Debug("Found matching signing ID in cache",
			zap.Uint64("signing_id", signingID),
			zap.Uint64("tunnel_id", cacheEntry.TunnelID),
			zap.Uint64("current_group_signing_id", cacheEntry.CurrentGroupSigningID),
			zap.Uint64("incoming_group_signing_id", cacheEntry.IncomingGroupSigningID),
		)

		// Delete the cache entries for both signing IDs
		h.signingCache.Delete(cacheEntry.CurrentGroupSigningID)
		h.signingCache.Delete(cacheEntry.IncomingGroupSigningID)

		relayermetrics.DecreasePendingNewSigning()

		// Trigger relayer
		h.Log.Debug("Triggered relayer for tunnel", zap.Uint64("tunnel_id", cacheEntry.TunnelID))
		h.TriggerRelayerCh <- cacheEntry.TunnelID
	}
}

// HandleNewPacket stores new packets in the cache.
func (h *PacketHandler) HandleNewPacket() {
	for packet := range h.NewPacketCh {
		if _, ok := h.ValidTunnelIDs[packet.TunnelID]; !ok {
			continue
		}

		h.Log.Debug("Received new packet",
			zap.Uint64("tunnel_id", packet.TunnelID),
			zap.Uint64("sequence", packet.Sequence),
		)

		// trigger relayer if the current group signing and
		// incoming group signing status are successful.
		if packet.CurrentGroupSigning != nil &&
			packet.CurrentGroupSigning.SigningStatus == tsstypes.SIGNING_STATUS_SUCCESS {
			h.Log.Debug("Current group signing is successful",
				zap.Uint64("tunnel_id", packet.TunnelID),
				zap.Uint64("sequence", packet.Sequence),
			)

			h.TriggerRelayerCh <- packet.TunnelID
			continue
		}

		if packet.IncomingGroupSigning != nil &&
			packet.IncomingGroupSigning.SigningStatus == tsstypes.SIGNING_STATUS_SUCCESS {
			h.Log.Debug("Incoming group signing is successful",
				zap.Uint64("tunnel_id", packet.TunnelID),
				zap.Uint64("sequence", packet.Sequence),
			)

			h.TriggerRelayerCh <- packet.TunnelID
			continue
		}

		h.cacheSigningNewPacket(packet)
	}
}

// UpdateValidTunnelIDs updates the valid tunnel IDs.
func (h *PacketHandler) UpdateValidTunnelIDs(tunnelIDs []uint64) {
	for _, tunnelID := range tunnelIDs {
		h.ValidTunnelIDs[tunnelID] = struct{}{}
	}
}

// cacheSigningNewPacket caches the signing IDs of a new packet.
func (h *PacketHandler) cacheSigningNewPacket(packet *types.Packet) {
	// Get the signing IDs from the packet.
	currentSigningID, incomingSigningID := uint64(0), uint64(0)
	if packet.CurrentGroupSigning != nil {
		currentSigningID = packet.CurrentGroupSigning.ID
	}
	if packet.IncomingGroupSigning != nil {
		incomingSigningID = packet.IncomingGroupSigning.ID
	}

	cacheEntry := NewCacheEntry(packet.TunnelID, currentSigningID, incomingSigningID)

	// Cache the signing ID if it is not in the cache.
	if currentSigningID != 0 &&
		!h.signingCache.Has(currentSigningID) {
		h.signingCache.Set(currentSigningID, cacheEntry, cacheTTL)
	}
	if incomingSigningID != 0 &&
		!h.signingCache.Has(incomingSigningID) {
		h.signingCache.Set(incomingSigningID, cacheEntry, cacheTTL)
	}

	h.Log.Debug("Stored current group signing ID in cache",
		zap.Uint64("tunnel_id", packet.TunnelID),
		zap.Uint64("sequence", packet.Sequence),
		zap.Uint64("current_group_signing_id", cacheEntry.CurrentGroupSigningID),
		zap.Uint64("incoming_group_signing_id", cacheEntry.IncomingGroupSigningID),
	)
}
