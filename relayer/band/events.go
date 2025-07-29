package band

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
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

	eventCh, err := c.rpcClient.Subscribe(ctx, "", subscriptionQuery)
	if err != nil {
		c.Log.Error("failed to subscribe to packet success events")
		return err
	}

	c.eventCh = eventCh

	return nil
}

// HandleProducePacketSuccess reads ProducePacketSuccess events from the channel
// and invokes the given handler once for each tunnel ID found in an event.
func (c *client) HandleProducePacketSuccess(handler func(tunnelID uint64)) {
	for msg := range c.eventCh {
		attrs := msg.Events

		// key for the tunnelID attribute
		key := fmt.Sprintf("%s.%s",
			tunneltypes.EventTypeProducePacketSuccess,
			tunneltypes.AttributeKeyTunnelID,
		)

		tunnelIDs := attrs[key]
		if len(tunnelIDs) == 0 {
			c.Log.Debug("missing tunnel_id in event")
			continue
		}

		// handle *each* tunnelID in the event
		for _, idStr := range tunnelIDs {
			tunnelID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				c.Log.Debug("failed to parse tunnel_id",
					zap.String("tunnel_id", idStr),
					zap.Error(err),
				)
				continue
			}
			handler(tunnelID)
		}
	}
}
