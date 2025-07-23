package subscriber

import (
	"context"
	"fmt"
	"strconv"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"go.uber.org/zap"

	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
)

var _ Subscriber = &PacketSuccessSubscriber{}

// PacketSuccessSubscriber subscribes to the produce packet success event.
type PacketSuccessSubscriber struct {
	rpcClient  rpcclient.Client
	log        *zap.Logger
	eventCh    <-chan coretypes.ResultEvent
	tunnelIDCh chan<- uint64
}

// NewPacketSuccessSubscriber creates a new PacketSuccessSubscriber.
func NewPacketSuccessSubscriber(
	log *zap.Logger,
	tunnelIDCh chan<- uint64,
) *PacketSuccessSubscriber {
	return &PacketSuccessSubscriber{
		rpcClient:  nil,
		log:        log.With(zap.String("subscriber", "packet_success")),
		eventCh:    make(chan coretypes.ResultEvent),
		tunnelIDCh: tunnelIDCh,
	}
}

// Subscribe subscribes to the produce packet success event.
func (s *PacketSuccessSubscriber) Subscribe(ctx context.Context, endpoint string) error {
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
		"tm.event='NewBlock' AND %s.%s EXISTS AND %s.%s EXISTS",
		tunneltypes.EventTypeProducePacketSuccess,
		tunneltypes.AttributeKeyTunnelID,
		tunneltypes.EventTypeProducePacketSuccess,
		tunneltypes.AttributeKeySequence,
	)

	eventCh, err := s.rpcClient.Subscribe(ctx, "producePacketSuccess", subscriptionQuery, 1000)
	if err != nil {
		return err
	}
	s.eventCh = eventCh

	return nil
}

// HandleEvent handles the produce packet success event and
// forwards the received packet to the packet channel.
func (s *PacketSuccessSubscriber) HandleEvent(ctx context.Context) error {
	for msg := range s.eventCh {
		attrs := msg.Events

		// key for the tunnelID attribute
		key := fmt.Sprintf(
			"%s.%s",
			tunneltypes.EventTypeProducePacketSuccess,
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
			fmt.Println("tunnelID", tunnelID)
			s.tunnelIDCh <- tunnelID
		}
	}

	return nil
}
