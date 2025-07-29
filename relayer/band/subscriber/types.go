package subscriber

import "context"

type Subscriber interface {
	// Subscribe subscribes to the event by initialize the rpc client from the given endpoint.
	// New client is created per subscriber due to the issue related to the limit number
	//  of channel when resubscription.
	Subscribe(ctx context.Context, endpoint string) error

	// HandleEvent handles the event from the subscribed channel.
	HandleEvent(ctx context.Context)
}

// SigningResult is a struct for handling the signing result from BandChain.
// It is used to group signingSuccess and signingFailed into the same channel to
// avoid race conditions when handling the signing result.
type SigningResult struct {
	SigningID uint64
	IsSuccess bool
}

// NewSigningResult creates a new SigningResult.
func NewSigningResult(signingID uint64, isSuccess bool) SigningResult {
	return SigningResult{
		SigningID: signingID,
		IsSuccess: isSuccess,
	}
}
