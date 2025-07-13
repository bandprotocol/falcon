package band

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
	"github.com/bandprotocol/falcon/relayer/band/types"
)

// Subscribe subscribes events from BandChain.
func (c *client) Subscribe(ctx context.Context) error {
	if err := c.subscribeToProducePacketSuccess(ctx); err != nil {
		c.Log.Error(
			"Failed to subscribe to ProducePacketSuccess events",
			zap.String("rpcEndpoint", c.Context.NodeURI),
			zap.Error(err),
		)
		return err
	}

	if err := c.subscribeToSigningSuccess(ctx); err != nil {
		c.Log.Error(
			"Failed to subscribe to SigningSuccess events",
			zap.String("rpcEndpoint", c.Context.NodeURI),
			zap.Error(err),
		)
		return err
	}

	if err := c.subscribeToSigningFailed(ctx); err != nil {
		c.Log.Error(
			"Failed to subscribe to SigningFailed events",
			zap.String("rpcEndpoint", c.Context.NodeURI),
			zap.Error(err),
		)
		return err
	}

	c.Log.Info("Subscribed to BandChain", zap.String("rpcEndpoint", c.Context.NodeURI))
	return nil
}

// subscribeToProducePacketSuccess subscribes to BandChain that emits events
// whenever a packet is successfully produced on any tunnel.
func (c *client) subscribeToProducePacketSuccess(ctx context.Context) error {
	subscriptionQuery := fmt.Sprintf(
		"tm.event='NewBlock' AND %s.%s EXISTS AND %s.%s EXISTS",
		tunneltypes.EventTypeProducePacketSuccess,
		tunneltypes.AttributeKeyTunnelID,
		tunneltypes.EventTypeProducePacketSuccess,
		tunneltypes.AttributeKeySequence,
	)

	eventCh, err := c.rpcClient.Subscribe(ctx, "producePacketSuccess", subscriptionQuery)
	if err != nil {
		return err
	}
	c.producePacketEventCh = eventCh

	return nil
}

// subscribeToSigningSuccess subscribes to BandChain that emits events
// whenever a signature is successfully aggregated.
func (c *client) subscribeToSigningSuccess(ctx context.Context) error {
	subscriptionQuery := fmt.Sprintf(
		"tm.event='NewBlock' AND %s.%s EXISTS",
		tsstypes.EventTypeSigningSuccess,
		tsstypes.AttributeKeySigningID,
	)

	eventCh, err := c.rpcClient.Subscribe(ctx, "signingSuccess", subscriptionQuery)
	if err != nil {
		return err
	}
	c.signingSuccessEventCh = eventCh

	return nil
}

// subscribeToSigningFailed subscribes to BandChain that emits events
// whenever a signature is failed to aggregate.
func (c *client) subscribeToSigningFailed(ctx context.Context) error {
	subscriptionQuery := fmt.Sprintf(
		"tm.event='NewBlock' AND %s.%s EXISTS",
		tsstypes.EventTypeSigningFailed,
		tsstypes.AttributeKeySigningID,
	)

	eventCh, err := c.rpcClient.Subscribe(ctx, "signingFailed", subscriptionQuery)
	if err != nil {
		return err
	}
	c.signingFailureEventCh = eventCh

	return nil
}

// HandleProducePacketSuccess reads ProducePacketSuccess events from the channel
// and forwards the received tunnel IDs to the channel.
func (c *client) HandleProducePacketSuccess(packetCh chan<- *types.Packet) {
	for msg := range c.producePacketEventCh {
		attrs := msg.Events

		// key for the tunnelID attribute
		key := fmt.Sprintf("%s.%s",
			tunneltypes.EventTypeProducePacketSuccess,
			tunneltypes.AttributeKeyTunnelID,
		)

		emittedTunnelIDs := attrs[key]
		if len(emittedTunnelIDs) == 0 {
			c.Log.Error("Missing tunnel_id in event produce_packet_success")
			continue
		}

		// parse the tunnel IDs from the event
		var tunnelIDs []uint64
		for _, idStr := range emittedTunnelIDs {
			tunnelID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				c.Log.Error("Failed to parse tunnel_id in the event produce_packet_success",
					zap.String("tunnel_id", idStr),
					zap.Error(err),
				)
				continue
			}

			tunnelIDs = append(tunnelIDs, tunnelID)
		}

		// handle each tunnelID in the event
		for _, tunnelID := range tunnelIDs {
			go func(tunnelID uint64) {
				// get the latest packet from the tunnel and forward to the packet channel
				// only if the tunnel is tssRoute and produce a new packet
				packet, err := c.GetLatestPacket(context.Background(), tunnelID)
				if err != nil {
					c.Log.Error(
						"Failed to get latest packet",
						zap.Error(err),
						zap.Uint64("tunnel_id", tunnelID),
					)
					return
				} else if packet == nil {
					return
				}

				packetCh <- packet
			}(tunnelID)
		}
	}
}

// HandleSigningSuccess reads SigningSuccess events from the channel
// and forwards the received signing IDs to the channel.
func (c *client) HandleSigningSuccess(signingIDSuccessCh chan<- uint64) {
	for msg := range c.signingSuccessEventCh {
		attrs := msg.Events

		// key for the signingID attribute
		key := fmt.Sprintf("%s.%s",
			tsstypes.EventTypeSigningSuccess,
			tsstypes.AttributeKeySigningID,
		)

		signingIDs := attrs[key]
		if len(signingIDs) == 0 {
			c.Log.Error("Missing signing_id in event signing_success")
			continue
		}

		// handle each signingID in the event
		for _, idStr := range signingIDs {
			signingID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				c.Log.Error("Failed to parse signing_id in the event signing_success",
					zap.String("tunnel_id", idStr),
					zap.Error(err),
				)
				continue
			}
			signingIDSuccessCh <- signingID
		}
	}
}

// HandleSigningFailure reads SigningFailed events from the channel
// and forwards the received signing IDs to the channel.
func (c *client) HandleSigningFailure(signingIDFailureCh chan<- uint64) {
	for msg := range c.signingFailureEventCh {
		attrs := msg.Events

		// key for the signingID attribute
		key := fmt.Sprintf("%s.%s",
			tsstypes.EventTypeSigningFailed,
			tsstypes.AttributeKeySigningID,
		)

		signingIDs := attrs[key]
		if len(signingIDs) == 0 {
			c.Log.Error("Missing signing_id in event signing_failed")
			continue
		}

		// handle each signingID in the event
		for _, idStr := range signingIDs {
			signingID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				c.Log.Error("Failed to parse signing_id in the event signing_failed",
					zap.String("tunnel_id", idStr),
					zap.Error(err),
				)
				continue
			}
			signingIDFailureCh <- signingID
		}
	}
}
