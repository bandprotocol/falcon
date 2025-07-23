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

var _ Subscriber = &SigningFailedSubscriber{}

// SigningFailedSubscriber is an object for handling the signing failed event.
type SigningFailedSubscriber struct {
	rpcClient   rpcclient.Client
	log         *zap.Logger
	eventCh     <-chan coretypes.ResultEvent
	signingIDCh chan<- uint64
}

// NewSigningFailedSubscriber creates a new SigningFailedSubscriber.
func NewSigningFailedSubscriber(
	log *zap.Logger,
	signingIDCh chan<- uint64,
) *SigningFailedSubscriber {
	return &SigningFailedSubscriber{
		rpcClient:   nil,
		log:         log.With(zap.String("subscriber", "signing_failed")),
		eventCh:     make(chan coretypes.ResultEvent),
		signingIDCh: signingIDCh,
	}
}

// Subscribe subscribes to the signing failed event.
func (s *SigningFailedSubscriber) Subscribe(ctx context.Context, endpoint string) error {
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
		tsstypes.EventTypeSigningFailed,
		tsstypes.AttributeKeySigningID,
	)

	eventCh, err := s.rpcClient.Subscribe(ctx, "signingFailed", subscriptionQuery, 1000)
	if err != nil {
		return err
	}
	s.eventCh = eventCh

	return nil
}

// HandleEvent handles the signing failed event and
// forwards the received signing ID to the signing ID channel.
func (s *SigningFailedSubscriber) HandleEvent(ctx context.Context) error {
	for msg := range s.eventCh {
		attrs := msg.Events

		// key for the signingID attribute
		key := fmt.Sprintf(
			"%s.%s",
			tsstypes.EventTypeSigningFailed,
			tsstypes.AttributeKeySigningID,
		)

		signingIDs := attrs[key]
		if len(signingIDs) == 0 {
			s.log.Error("Missing signing_id in event signing_failed")
			continue
		}

		// handle each signingID in the event
		for _, idStr := range signingIDs {
			signingID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				s.log.Error(
					"Failed to parse signing_id in the event signing_failed",
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
