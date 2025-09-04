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

var _ Subscriber = &ManualTriggerSubscriber{}

// ManualTriggerSubscriber is an object for handling the manual trigger event.
type ManualTriggerSubscriber struct {
	*Subscription
}

// NewManualTriggerSubscriber creates a new ManualTriggerSubscriber.
func NewManualTriggerSubscriber(
	log logger.Logger,
	tunnelIDCh chan<- uint64,
	timeout time.Duration,
) *ManualTriggerSubscriber {
	name := "manual_trigger"

	subscriptionQuery := fmt.Sprintf(
		"tm.event='Tx' AND %s.%s EXISTS",
		tunneltypes.EventTypeTriggerTunnel,
		tunneltypes.AttributeKeyTunnelID,
	)

	l := log.With("subscriber", name)
	onEventReceived := onHandleManualTriggeredEvent(tunnelIDCh, l)

	subscription := NewSubscription(
		name,
		subscriptionQuery,
		onEventReceived,
		timeout,
		l,
	)

	return &ManualTriggerSubscriber{
		Subscription: subscription,
	}
}

// onHandleManualTriggeredEvent handles the manual triggered event.
func onHandleManualTriggeredEvent(
	tunnelIDCh chan<- uint64,
	log logger.Logger,
) func(ctx context.Context, msg coretypes.ResultEvent) {
	return func(ctx context.Context, msg coretypes.ResultEvent) {
		attrs := msg.Events

		// key for the tunnelID attribute
		key := fmt.Sprintf(
			"%s.%s",
			tunneltypes.EventTypeTriggerTunnel,
			tunneltypes.AttributeKeyTunnelID,
		)

		emittedTunnelIDs := attrs[key]
		if len(emittedTunnelIDs) == 0 {
			log.Error("Missing tunnel_id in event manual_trigger")
			return
		}

		// parse the tunnel IDs from the event
		for _, idStr := range emittedTunnelIDs {
			tunnelID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				log.Error(
					"Failed to parse tunnel_id in the event manual_trigger",
					"tunnel_id", idStr,
					err,
				)
				continue
			}
			tunnelIDCh <- tunnelID
		}
	}
}
