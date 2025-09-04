package subscriber

import (
	"context"
	"fmt"
	"strconv"
	"time"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"

	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
	"github.com/bandprotocol/falcon/relayer/logger"
)

var _ Subscriber = &PacketSuccessSubscriber{}

// PacketSuccessSubscriber is an object for handling the produce packet success event.
type PacketSuccessSubscriber struct {
	*Subscription
}

// NewPacketSuccessSubscriber creates a new PacketSuccessSubscriber.
func NewPacketSuccessSubscriber(
	log logger.Logger,
	tunnelIDCh chan<- uint64,
	timeout time.Duration,
) *PacketSuccessSubscriber {
	name := "packet_success"

	subscriptionQuery := fmt.Sprintf(
		"tm.event='NewBlock' AND %s.%s EXISTS",
		tunneltypes.EventTypeProducePacketSuccess,
		tunneltypes.AttributeKeyTunnelID,
	)

	l := log.With("subscriber", name)
	onEventReceived := onHandlePacketSuccessEvent(tunnelIDCh, l)

	subscription := NewSubscription(
		name,
		subscriptionQuery,
		onEventReceived,
		timeout,
		l,
	)

	return &PacketSuccessSubscriber{
		Subscription: subscription,
	}
}

// onHandlePacketSuccessEvent handles the produce packet success event.
func onHandlePacketSuccessEvent(
	tunnelIDCh chan<- uint64,
	log logger.Logger,
) func(ctx context.Context, msg coretypes.ResultEvent) {
	return func(ctx context.Context, msg coretypes.ResultEvent) {
		attrs := msg.Events

		// key for the tunnelID attribute
		key := fmt.Sprintf(
			"%s.%s",
			tunneltypes.EventTypeProducePacketSuccess,
			tunneltypes.AttributeKeyTunnelID,
		)

		emittedTunnelIDs := attrs[key]
		if len(emittedTunnelIDs) == 0 {
			log.Error("Missing tunnel_id in event produce_packet_success")
			return
		}

		// parse the tunnel IDs from the event
		for _, idStr := range emittedTunnelIDs {
			tunnelID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				log.Error(
					"Failed to parse tunnel_id in the event produce_packet_success",
					"tunnel_id", idStr,
					err,
				)
				continue
			}
			tunnelIDCh <- tunnelID
		}
	}
}
