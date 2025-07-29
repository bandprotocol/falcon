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
