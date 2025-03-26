package pubsub

import "context"

// Consumer represents a consumer functionality that consumes messages.
// workerCount is the number of workers that will be used to consume messages.
type Consumer interface {
	Consume(ctx context.Context, workerCount int) error
}
