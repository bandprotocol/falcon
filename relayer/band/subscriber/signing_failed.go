package subscriber

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
)

var _ Subscriber = &SigningFailedSubscriber{}

// SigningFailedSubscriber is an object for handling the signing failed event.
type SigningFailedSubscriber struct {
	*Subscription
}

// NewSigningFailedSubscriber creates a new SigningFailedSubscriber.
func NewSigningFailedSubscriber(
	log *zap.Logger,
	signingResultCh chan<- SigningResult,
	timeout time.Duration,
) *SigningFailedSubscriber {
	name := "signing_failed"

	subscriptionQuery := fmt.Sprintf(
		"tm.event='NewBlock' AND %s.%s EXISTS",
		tsstypes.EventTypeSigningFailed,
		tsstypes.AttributeKeySigningID,
	)

	l := log.With(zap.String("subscriber", name))
	onEventReceived := onHandleSigningFailedEvent(signingResultCh, l)

	subscription := NewSubscription(
		name,
		subscriptionQuery,
		onEventReceived,
		timeout,
		l,
	)

	return &SigningFailedSubscriber{
		Subscription: subscription,
	}
}

// onHandleSigningFailedEvent handles the signing failed event.
func onHandleSigningFailedEvent(
	signingResultCh chan<- SigningResult,
	log *zap.Logger,
) func(ctx context.Context, msg coretypes.ResultEvent) {
	return func(ctx context.Context, msg coretypes.ResultEvent) {
		attrs := msg.Events

		// key for the signingID attribute
		key := fmt.Sprintf(
			"%s.%s",
			tsstypes.EventTypeSigningFailed,
			tsstypes.AttributeKeySigningID,
		)

		signingIDs := attrs[key]
		if len(signingIDs) == 0 {
			log.Error("Missing signing_id in event signing_failed")
			return
		}

		// handle each signingID in the event
		for _, idStr := range signingIDs {
			signingID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				log.Error(
					"Failed to parse signing_id in the event signing_failed",
					zap.String("signing_id", idStr),
					zap.Error(err),
				)
				continue
			}
			signingResultCh <- NewSigningResult(signingID, false)
		}
	}
}
