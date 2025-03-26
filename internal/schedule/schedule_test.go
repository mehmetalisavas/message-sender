package schedule

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mehmetalisavas/message-sender/config"
	"github.com/mehmetalisavas/message-sender/internal/db/mysql"
	"github.com/mehmetalisavas/message-sender/internal/db/redis"
	"github.com/mehmetalisavas/message-sender/internal/models"
	"github.com/mehmetalisavas/message-sender/internal/pubsub"
	"github.com/mehmetalisavas/message-sender/pkg/services/notification"
	"github.com/sethvargo/go-envconfig"
)

type MockNotificationService struct{}

func (m *MockNotificationService) Send(ctx context.Context, to, content string) (*notification.NotificationResponse, error) {
	return &notification.NotificationResponse{
		Message:   "Accepted",
		MessageID: uuid.New().String(),
	}, nil
}

func TestScheduler(t *testing.T) {
	ctx := context.Background()

	c := config.New()
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
	client, err := mysql.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}
	store := mysql.NewSqlStore(client)
	// store := testStorage()
	// store := mysql.TestMysqlStorage()

	// Insert a test message
	now := time.Now()
	// Insert a message for testing

	cacheService, err := redis.NewRedisCacheStore(ctx, c)
	if err != nil {
		log.Fatalf("error while starting cache service: %v \n", err)
	}

	message1 := models.Message{
		Content:   "Test Message scheduled",
		Recipient: "+905555555555",
		Status:    models.MessageStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	message2 := models.Message{
		Content:   "Test Message scheduled 2",
		Recipient: "+905555555556",
		Status:    models.MessageStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	insertedMessage1, err := store.InsertTestMessages(ctx, message1)
	if err != nil {
		t.Fatalf("Failed to insert message: %v", err)
	}
	_, err = store.InsertTestMessages(ctx, message2)
	if err != nil {
		t.Fatalf("Failed to insert message: %v", err)
	}

	notificationService := &MockNotificationService{}
	scheduler := NewScheduler(store)
	messageProducer := pubsub.NewMessageProducer(&c, store, scheduler.MessageBus(), 1)
	scheduler.AddProducer(messageProducer)
	messageConsumer := pubsub.NewMessageConsumer(store, scheduler.MessageBus(), notificationService, cacheService)
	scheduler.AddConsumer(messageConsumer)

	go scheduler.Start(ctx, 2) // start with 2 workers

	// Wait for the message to be processed
	time.Sleep(3 * time.Second)

	fetchedMessage, err := store.GetTestMessage(ctx, insertedMessage1.ID)
	if err != nil {
		t.Fatalf("Failed to fetch message: %v", err)
	}
	if fetchedMessage.Status != models.MessageStatusSent {
		t.Errorf("UpdateMessageStatus() failed to update message status")
	}
}
