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

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
)

var _ Subscriber = &SigningSuccessSubscriber{}

// SigningSuccessSubscriber is an object for handling the signing success event.
type SigningSuccessSubscriber struct {
	name              string
	subscriptionQuery string
	timeout           time.Duration
	rpcClient         rpcclient.Client
	log               *zap.Logger
	stopCh            chan struct{}
	eventCh           chan coretypes.ResultEvent
	signingResultCh   chan<- SigningResult
}

// NewSigningSuccessSubscriber creates a new SigningSuccessSubscriber.
func NewSigningSuccessSubscriber(
	log *zap.Logger,
	signingResultCh chan<- SigningResult,
	timeout time.Duration,
) *SigningSuccessSubscriber {
	name := "signing_success"

	return &SigningSuccessSubscriber{
		name: name,
		subscriptionQuery: fmt.Sprintf(
			"tm.event='NewBlock' AND %s.%s EXISTS",
			tsstypes.EventTypeSigningSuccess,
			tsstypes.AttributeKeySigningID,
		),
		timeout:         timeout,
		rpcClient:       nil,
		log:             log.With(zap.String("subscriber", name)),
		eventCh:         make(chan coretypes.ResultEvent, 1000),
		signingResultCh: signingResultCh,
	}
}

// Subscribe subscribes to the produce packet success event.
func (s *SigningSuccessSubscriber) Subscribe(ctx context.Context, endpoint string) error {
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
func (s *SigningSuccessSubscriber) HandleEvent(ctx context.Context) {
	for msg := range s.eventCh {
		attrs := msg.Events

		// key for the signingID attribute
		key := fmt.Sprintf(
			"%s.%s",
			tsstypes.EventTypeSigningSuccess,
			tsstypes.AttributeKeySigningID,
		)

		signingIDs := attrs[key]
		if len(signingIDs) == 0 {
			s.log.Error("Missing signing_id in event signing_success")
			continue
		}

		// handle each signingID in the event
		for _, idStr := range signingIDs {
			signingID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				s.log.Error(
					"Failed to parse signing_id in the event signing_success",
					zap.String("signing_id", idStr),
					zap.Error(err),
				)
				continue
			}
			s.signingResultCh <- NewSigningResult(signingID, true)
		}
	}
}

// unsubscribeAndStopPreviousClient unsubscribes from the previous RPC client if it exists.
// If error occurs (e.g. client is already stopped or timeout), it will be logged
// but not returned so that it doesn't block the subscription part.
func (s *SigningSuccessSubscriber) unsubscribeAndStopPreviousClient(ctx context.Context) {
	if s.rpcClient == nil {
		return
	}

	unsubCtx, unsubCtxCancel := context.WithTimeout(ctx, s.timeout)
	defer unsubCtxCancel()
	if err := s.rpcClient.Unsubscribe(unsubCtx, s.name, s.subscriptionQuery); err != nil {
		s.log.Debug(
			"Failed to unsubscribe from signing_success event",
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
