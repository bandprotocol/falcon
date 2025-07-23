package subscriber

import (
	"context"
	"fmt"
	"strconv"

	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"go.uber.org/zap"
)

var _ Subscriber = &ManualTriggerSubscriber{}

type ManualTriggerSubscriber struct {
	rpcClient  rpcclient.Client
	log        *zap.Logger
	eventCh    <-chan coretypes.ResultEvent
	tunnelIDCh chan<- uint64
}

var _ Subscriber = &ManualTriggerSubscriber{}

// NewManualTriggerSubscriber creates a new ManualTriggerSubscriber.
func NewManualTriggerSubscriber(
	log *zap.Logger,
	tunnelIDCh chan<- uint64,
) *ManualTriggerSubscriber {
	return &ManualTriggerSubscriber{
		rpcClient:  nil,
		log:        log.With(zap.String("subscriber", "manual_trigger")),
		eventCh:    make(chan coretypes.ResultEvent),
		tunnelIDCh: make(chan uint64),
	}
}

// Subscribe subscribes to the manual trigger event.
func (s *ManualTriggerSubscriber) Subscribe(ctx context.Context, endpoint string) error {
	client, err := httpclient.New(endpoint, "/websocket")
	if err != nil {
		return err
	}

	if err := client.Start(); err != nil {
		s.log.Error(
			"Failed to start HTTP client",
			zap.String("rpcEndpoint", endpoint),
			zap.Error(err),
		)
		return err
	}

	s.rpcClient = client

	subscriptionQuery := fmt.Sprintf(
		"tm.event='Tx' AND %s.%s EXISTS",
		tunneltypes.EventTypeTriggerTunnel,
		tunneltypes.AttributeKeyTunnelID,
	)

	eventCh, err := s.rpcClient.Subscribe(ctx, "manualTrigger", subscriptionQuery, 1000)
	if err != nil {
		return err
	}
	s.eventCh = eventCh

	return nil
}

// HandleEvent handles the produce packet success event and
// forwards the received packet to the packet channel.
func (s *ManualTriggerSubscriber) HandleEvent(ctx context.Context) error {
	for msg := range s.eventCh {
		attrs := msg.Events

		// key for the tunnelID attribute
		key := fmt.Sprintf(
			"%s.%s",
			tunneltypes.EventTypeTriggerTunnel,
			tunneltypes.AttributeKeyTunnelID,
		)

		emittedTunnelIDs := attrs[key]
		if len(emittedTunnelIDs) == 0 {
			s.log.Error("Missing tunnel_id in event produce_packet_success")
			continue
		}

		// parse the tunnel IDs from the event
		for _, idStr := range emittedTunnelIDs {
			tunnelID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				s.log.Error(
					"Failed to parse tunnel_id in the event produce_packet_success",
					zap.String("tunnel_id", idStr),
					zap.Error(err),
				)
				continue
			}
			s.tunnelIDCh <- tunnelID
		}
	}

	return nil
}
