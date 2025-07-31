package subscriber

import (
	"context"
	"time"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"go.uber.org/zap"
)

// Subscription is an object for handling the subscription to the event.
type Subscription struct {
	name              string
	subscriptionQuery string
	timeout           time.Duration
	rpcClient         rpcclient.Client
	log               *zap.Logger
	stopCh            chan struct{}
	eventCh           chan coretypes.ResultEvent
	onEventReceived   func(ctx context.Context, msg coretypes.ResultEvent)
}

// NewSubscription creates a new Subscription object.
func NewSubscription(
	name string,
	subscriptionQuery string,
	onEventReceived func(ctx context.Context, msg coretypes.ResultEvent),
	timeout time.Duration,
	log *zap.Logger,
) *Subscription {
	return &Subscription{
		name:              name,
		subscriptionQuery: subscriptionQuery,
		timeout:           timeout,
		log:               log,
		stopCh:            make(chan struct{}),
		eventCh:           make(chan coretypes.ResultEvent, 1000),
		onEventReceived:   onEventReceived,
	}
}

// Subscribe subscribes to the manual trigger event.
func (s *Subscription) Subscribe(ctx context.Context, endpoint string) error {
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

	s.rpcClient = client

	return nil
}

// unsubscribeAndStopPreviousClient unsubscribes from the previous RPC client if it exists.
// If error occurs (e.g. client is already stopped or timeout), it will be logged
// but not returned so that it doesn't block the subscription part.
func (s *Subscription) unsubscribeAndStopPreviousClient(ctx context.Context) {
	if s.rpcClient == nil {
		return
	}

	unsubCtx, unsubCtxCancel := context.WithTimeout(ctx, s.timeout)
	defer unsubCtxCancel()
	if err := s.rpcClient.Unsubscribe(unsubCtx, s.name, s.subscriptionQuery); err != nil {
		s.log.Debug(
			"Failed to unsubscribe from event",
			zap.String("event_name", s.name),
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
	s.rpcClient = nil

	s.log.Debug("Unsubscribe and stop HTTP client successfully")
}

// HandleEvent handles the event from the subscribed channel.
func (s *Subscription) HandleEvent(ctx context.Context) {
	for msg := range s.eventCh {
		if s.onEventReceived != nil {
			s.onEventReceived(ctx, msg)
		}
	}
}
