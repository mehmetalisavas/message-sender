package pubsub

import (
	"context"
	"log"
	"time"

	"github.com/mehmetalisavas/message-sender/internal/models"
	"github.com/mehmetalisavas/message-sender/internal/service"
	"github.com/mehmetalisavas/message-sender/pkg/services/notification"
)

var _ Consumer = (*MessageConsumer)(nil)

type MessageConsumer struct {
	storageService      service.Storage
	messageBus          *MessageBus
	notificationService notification.NotificationSender
	cacheService        service.CacheStore
}

func NewMessageConsumer(storageService service.Storage, messageBus *MessageBus, notificationService notification.NotificationSender, cacheService service.CacheStore) *MessageConsumer {
	return &MessageConsumer{
		storageService:      storageService,
		messageBus:          messageBus,
		notificationService: notificationService,
		cacheService:        cacheService,
	}
}

func (mc *MessageConsumer) Consume(ctx context.Context, workerCount int) error {
	messageChannel, exists := mc.messageBus.GetChannel(MessageSenderTopic)
	if !exists {
		return ErrChannelNotFound
	}

	for i := 0; i < workerCount; i++ {
		go mc.worker(ctx, messageChannel)
	}

	return nil
}

func (mc *MessageConsumer) worker(ctx context.Context, ch chan interface{}) error {

	for {
		select {

		case msg := <-ch:
			message, ok := msg.(models.Message)
			if !ok {
				log.Println("invalid message type")
				continue
			}

			err := mc.processMessage(ctx, message)
			if err != nil {
				log.Printf("failed to process message: %v\n", err)
				continue
			}

		case <-ctx.Done():
			log.Println("message consumer is stopped")
			return nil
		}
	}
}

// processMessage simulates sending a message
func (mc *MessageConsumer) processMessage(ctx context.Context, msg models.Message) error {
	requestSendingTime := time.Now()
	resp, err := mc.notificationService.Send(ctx, msg.Recipient, msg.Content)
	if err != nil {
		log.Printf("failed to process message id:%d: %v\n", msg.ID, err)
		err := mc.storageService.UpdateMessageStatus(ctx, msg.ID, models.MessageStatusFailed)
		return err
	}

	err = mc.storageService.UpdateMessageStatus(ctx, msg.ID, models.MessageStatusSent)
	if err != nil {
		log.Printf("failed to update message status id:%d: %v\n", msg.ID, err)
		return err
	}
	err = mc.cacheService.CacheMessage(ctx, resp.MessageID, requestSendingTime)
	if err != nil {
		log.Printf("failed to cache message id:%s: %v\n", resp.MessageID, err)
		return err
	}

	log.Printf("message %d is marked as sent\n", msg.ID)

	return nil
}
