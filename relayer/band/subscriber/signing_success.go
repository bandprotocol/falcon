package subscriber

import (
	"context"
	"fmt"
	"strconv"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
)

var _ Subscriber = &SigningSuccessSubscriber{}

// SigningSuccessSubscriber subscribes to the signing success event.
type SigningSuccessSubscriber struct {
	rpcClient   rpcclient.Client
	log         *zap.Logger
	eventCh     <-chan coretypes.ResultEvent
	signingIDCh chan<- uint64
}

// NewSigningSuccessSubscriber creates a new SigningSuccessSubscriber.
func NewSigningSuccessSubscriber(
	log *zap.Logger,
	signingIDCh chan<- uint64,
) *SigningSuccessSubscriber {
	return &SigningSuccessSubscriber{
		rpcClient:   nil,
		log:         log.With(zap.String("subscriber", "signing_success")),
		eventCh:     make(chan coretypes.ResultEvent),
		signingIDCh: signingIDCh,
	}
}

// Subscribe subscribes to the produce packet success event.
func (s *SigningSuccessSubscriber) Subscribe(ctx context.Context, endpoint string) error {
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
		"tm.event='NewBlock' AND %s.%s EXISTS",
		tsstypes.EventTypeSigningSuccess,
		tsstypes.AttributeKeySigningID,
	)

	eventCh, err := s.rpcClient.Subscribe(ctx, "signingSuccess", subscriptionQuery, 1000)
	if err != nil {
		return err
	}
	s.eventCh = eventCh

	return nil
}

// HandleEvent handles the produce packet success event and
// forwards the received packet to the packet channel.
func (s *SigningSuccessSubscriber) HandleEvent(ctx context.Context) error {
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
					zap.String("tunnel_id", idStr),
					zap.Error(err),
				)
				continue
			}
			s.signingIDCh <- signingID
		}
	}

	return nil
}
