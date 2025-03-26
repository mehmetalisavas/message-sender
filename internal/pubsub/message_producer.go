package pubsub

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mehmetalisavas/message-sender/config"
	"github.com/mehmetalisavas/message-sender/internal/service"
)

// Make sure MessageProducer implements Producer interface.
var _ Producer = (*MessageProducer)(nil)

var ErrChannelNotFound = errors.New("channel not found")

const MessageSenderTopic = "message-sender"

type MessageProducer struct {
	cfg            *config.Config
	storageService service.Storage
	messageBus     *MessageBus
	// intervalInSec represents the interval in seconds to produce messages.
	intervalInSec int
}

// NewMessageProducer creates a new MessageProducer instance.
func NewMessageProducer(cfg *config.Config, storageService service.Storage, messageBus *MessageBus, interval int) *MessageProducer {
	return &MessageProducer{
		cfg:            cfg,
		storageService: storageService,
		messageBus:     messageBus,
		intervalInSec:  interval,
	}
}

// Produce produces messages to the message queue.
func (mp *MessageProducer) Produce(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(mp.intervalInSec) * time.Second)
	defer ticker.Stop()

	messageChannel, exists := mp.messageBus.GetChannel(MessageSenderTopic)
	if !exists {
		log.Printf("channel %s does not exist in the message bus\n", MessageSenderTopic)
		return ErrChannelNotFound
	}

	for {
		select {
		case <-ticker.C:
			if !mp.cfg.IsMessageProcessing {
				log.Printf("message processing is stopped\n")
				continue
			}
			// Get pending messages from storage.
			log.Printf("getting pending messages from storage\n")

			messages, err := mp.storageService.GetPendingMessages(ctx, 2)
			if err != nil {
				log.Printf("failed to get pending messages from storage: %v\n", err)
				continue
			}

			for _, message := range messages {
				// Publish message to the message queue.
				messageChannel <- message
			}
		case <-ctx.Done():
			log.Printf("message producer is stopped\n")
			return nil
		}
	}
}
