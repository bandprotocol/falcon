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

var _ Subscriber = &ManualTriggerSubscriber{}

// ManualTriggerSubscriber is an object for handling the manual trigger event.
type ManualTriggerSubscriber struct {
	name              string
	subscriptionQuery string
	timeout           time.Duration
	rpcClient         rpcclient.Client
	log               *zap.Logger
	stopCh            chan struct{}
	eventCh           chan coretypes.ResultEvent
	tunnelIDCh        chan<- uint64
}

// NewManualTriggerSubscriber creates a new ManualTriggerSubscriber.
func NewManualTriggerSubscriber(
	log *zap.Logger,
	tunnelIDCh chan<- uint64,
	timeout time.Duration,
) *ManualTriggerSubscriber {
	name := "manual_trigger"

	return &ManualTriggerSubscriber{
		name: name,
		subscriptionQuery: fmt.Sprintf(
			"tm.event='Tx' AND %s.%s EXISTS",
			tunneltypes.EventTypeTriggerTunnel,
			tunneltypes.AttributeKeyTunnelID,
		),
		timeout:    timeout,
		rpcClient:  nil,
		log:        log.With(zap.String("subscriber", name)),
		stopCh:     make(chan struct{}),
		eventCh:    make(chan coretypes.ResultEvent, 1000),
		tunnelIDCh: make(chan uint64),
	}
}

// Subscribe subscribes to the manual trigger event.
func (s *ManualTriggerSubscriber) Subscribe(ctx context.Context, endpoint string) error {
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

	subCtx, subCtxCancel := context.WithTimeout(ctx, s.timeout)
	defer subCtxCancel()

	eventCh, err := client.Subscribe(subCtx, s.name, s.subscriptionQuery, 1000)
	if err != nil {
		return err
	}

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

	s.rpcClient = client

	return nil
}

// HandleEvent handles the produce packet success event and
// forwards the received packet to the packet channel.
func (s *ManualTriggerSubscriber) HandleEvent(ctx context.Context) {
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
			s.log.Error("Missing tunnel_id in event manual_trigger")
			continue
		}

		// parse the tunnel IDs from the event
		for _, idStr := range emittedTunnelIDs {
			tunnelID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				s.log.Error(
					"Failed to parse tunnel_id in the event manual_trigger",
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
func (s *ManualTriggerSubscriber) unsubscribeAndStopPreviousClient(ctx context.Context) {
	if s.rpcClient == nil {
		return
	}

	unsubCtx, unsubCtxCancel := context.WithTimeout(ctx, s.timeout)
	defer unsubCtxCancel()
	if err := s.rpcClient.Unsubscribe(unsubCtx, s.name, s.subscriptionQuery); err != nil {
		s.log.Debug(
			"Failed to unsubscribe from manual_trigger event",
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
