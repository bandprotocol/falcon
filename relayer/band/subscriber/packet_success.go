package subscriber

import (
	"context"
	"fmt"
	"strconv"
	"time"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"go.uber.org/zap"

	tunneltypes "github.com/bandprotocol/falcon/internal/bandchain/tunnel"
)

var _ Subscriber = &PacketSuccessSubscriber{}

// PacketSuccessSubscriber is an object for handling the produce packet success event.
type PacketSuccessSubscriber struct {
	name              string
	subscriptionQuery string
	timeout           time.Duration
	rpcClient         rpcclient.Client
	log               *zap.Logger
	stopCh            chan struct{}
	eventCh           chan coretypes.ResultEvent
	tunnelIDCh        chan<- uint64
}

// NewPacketSuccessSubscriber creates a new PacketSuccessSubscriber.
func NewPacketSuccessSubscriber(
	log *zap.Logger,
	tunnelIDCh chan<- uint64,
	timeout time.Duration,
) *PacketSuccessSubscriber {
	name := "packet_success"

	return &PacketSuccessSubscriber{
		name: name,
		subscriptionQuery: fmt.Sprintf(
			"tm.event='NewBlock' AND %s.%s EXISTS",
			tunneltypes.EventTypeProducePacketSuccess,
			tunneltypes.AttributeKeyTunnelID,
		),
		timeout:    timeout,
		rpcClient:  nil,
		log:        log.With(zap.String("subscriber", name)),
		stopCh:     make(chan struct{}),
		eventCh:    make(chan coretypes.ResultEvent, 1000),
		tunnelIDCh: tunnelIDCh,
	}
}

// Subscribe subscribes to the produce packet success event.
func (s *PacketSuccessSubscriber) Subscribe(ctx context.Context, endpoint string) error {
	// unsubscribe from the previous RPC client if it exists.
	s.unsubscribeAndStopPreviousClient(ctx)

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

	eventCh, err := client.Subscribe(ctx, s.name, s.subscriptionQuery, 1000)
	if err != nil {
		return err
	}

	s.stopCh = make(chan struct{})
	go func() {
		for {
			select {
			case <-s.stopCh:
				return
			case msg := <-eventCh:
				s.eventCh <- msg
			}
		}
	}()

	return nil
}

// HandleEvent handles the produce packet success event and
// forwards the received packet to the packet channel.
func (s *PacketSuccessSubscriber) HandleEvent(ctx context.Context) {
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
			s.tunnelIDCh <- tunnelID
		}
	}
}

// unsubscribeAndStopPreviousClient unsubscribes from the previous RPC client if it exists.
// If error occurs (e.g. client is already stopped or timeout), it will be logged
// but not returned so that it doesn't block the subscription part.
func (s *PacketSuccessSubscriber) unsubscribeAndStopPreviousClient(ctx context.Context) {
	if s.rpcClient == nil {
		return
	}

	unsubCtx, unsubCtxCancel := context.WithTimeout(ctx, s.timeout)
	defer unsubCtxCancel()
	if err := s.rpcClient.Unsubscribe(unsubCtx, s.name, s.subscriptionQuery); err != nil {
		s.log.Debug(
			"Failed to unsubscribe from packet_success event",
			zap.Error(err),
		)
	}

	if err := s.rpcClient.Stop(); err != nil {
		s.log.Debug(
			"Failed to stop HTTP client",
			zap.Error(err),
		)
	}

	close(s.stopCh)

	s.log.Debug("Unsubscribe and stop HTTP client successfully")
}
