package pubsub

import "context"

// Producer represents a producer functionality that produces messages.
type Producer interface {
	Produce(ctx context.Context) error
}
