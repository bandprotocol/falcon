package subscriber

import (
	"context"
	"fmt"
	"strconv"
	"time"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"go.uber.org/zap"

	tsstypes "github.com/bandprotocol/falcon/internal/bandchain/tss"
)

var _ Subscriber = &SigningSuccessSubscriber{}

// SigningSuccessSubscriber is an object for handling the signing success event.
type SigningSuccessSubscriber struct {
	*Subscription
}

// NewSigningSuccessSubscriber creates a new SigningSuccessSubscriber.
func NewSigningSuccessSubscriber(
	log *zap.Logger,
	signingResultCh chan<- SigningResult,
	timeout time.Duration,
) *SigningSuccessSubscriber {
	name := "signing_success"

	subscriptionQuery := fmt.Sprintf(
		"tm.event='NewBlock' AND %s.%s EXISTS",
		tsstypes.EventTypeSigningSuccess,
		tsstypes.AttributeKeySigningID,
	)

	l := log.With(zap.String("subscriber", name))
	onEventReceived := onHandleSigningSuccessEvent(signingResultCh, l)

	subscription := NewSubscription(
		name,
		subscriptionQuery,
		onEventReceived,
		timeout,
		l,
	)

	return &SigningSuccessSubscriber{
		Subscription: subscription,
	}
}

// onHandleSigningSuccessEvent handles the signing success event.
func onHandleSigningSuccessEvent(
	signingResultCh chan<- SigningResult,
	log *zap.Logger,
) func(ctx context.Context, msg coretypes.ResultEvent) {
	return func(ctx context.Context, msg coretypes.ResultEvent) {
		attrs := msg.Events

		// key for the signingID attribute
		key := fmt.Sprintf(
			"%s.%s",
			tsstypes.EventTypeSigningSuccess,
			tsstypes.AttributeKeySigningID,
		)

		signingIDs := attrs[key]
		if len(signingIDs) == 0 {
			log.Error("Missing signing_id in event signing_success")
			return
		}

		// handle each signingID in the event
		for _, idStr := range signingIDs {
			signingID, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				log.Error(
					"Failed to parse signing_id in the event signing_success",
					zap.String("signing_id", idStr),
					zap.Error(err),
				)
				continue
			}
			signingResultCh <- NewSigningResult(signingID, true)
		}
	}
}
