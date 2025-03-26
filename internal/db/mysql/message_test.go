package mysql

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mehmetalisavas/message-sender/config"
	"github.com/mehmetalisavas/message-sender/internal/models"
	"github.com/sethvargo/go-envconfig"
)

func (s *SqlStore) insertTestMessages(ctx context.Context, message models.Message) (*models.Message, error) {
	query := `
			INSERT INTO messages (content, recipient, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
		`
	result, err := s.db.Exec(query, message.Content, message.Recipient, message.Status, message.CreatedAt, message.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert message: %w", err)
	}
	// Get the last inserted ID (auto-incremented value)
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the full message using the inserted ID
	var insertedMessage models.Message
	selectQuery := `
			SELECT id, content, recipient, status, created_at, updated_at
			FROM messages
			WHERE id = ?`

	// Query the inserted message
	err = s.db.QueryRowContext(ctx, selectQuery, lastInsertID).
		Scan(&insertedMessage.ID, &insertedMessage.Content, &insertedMessage.Recipient, &insertedMessage.Status, &insertedMessage.CreatedAt, &insertedMessage.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Return the inserted message struct
	return &insertedMessage, nil
}

func (s *SqlStore) getTestMessage(ctx context.Context, id int) (*models.Message, error) {
	var message models.Message
	query := `
		SELECT id, content, recipient, status, created_at, updated_at
		FROM messages
		WHERE id = ?
	`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&message.ID, &message.Content, &message.Recipient, &message.Status, &message.CreatedAt, &message.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func testStorage() *SqlStore {
	ctx := context.Background()

	c := config.New()
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
	client, err := NewClient(c)
	if err != nil {
		log.Fatal(err)
	}
	testSqlStorage := NewSqlStore(client)

	return testSqlStorage
}

// TestListSentMessages tests the ListSentMessages function.
func TestListSentMessages(t *testing.T) {
	// Set up the test database and SqlStore
	store := testStorage()
	now := time.Now()

	// Insert test messages
	messagesToInsert := []models.Message{
		{
			Content:   "Message 1",
			Recipient: "user1@example.com",
			Status:    models.MessageStatusSent,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Content:   "Message 2",
			Recipient: "user2@example.com",
			Status:    models.MessageStatusSent,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Insert messages into the database
	err := insertTestMessages(store, messagesToInsert)
	if err != nil {
		t.Fatalf("Failed to insert messages: %v", err)
	}

	// Fetch the messages back using ListSentMessages
	opts := models.ListOptions{
		Limit:  2,
		Offset: 0,
	}
	messages, err := store.ListSentMessages(context.Background(), opts)
	if err != nil {
		t.Fatalf("ListSentMessages() error = %v", err)
	}

	// Check if the result matches the expected number of messages
	if len(messages) != len(messagesToInsert) {
		t.Errorf("ListSentMessages() returned %d messages, expected %d", len(messages), len(messagesToInsert))
	}
}

// Insert messages into the database
func insertTestMessages(db *SqlStore, messages []models.Message) error {
	for _, message := range messages {
		query := `
			INSERT INTO messages (content, recipient, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
		`
		_, err := db.db.Exec(query, message.Content, message.Recipient, message.Status, message.CreatedAt, message.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert message: %w", err)
		}
	}
	return nil
}

func TestUpdateMessageStatus(t *testing.T) {
	ctx := context.Background()
	store := testStorage()

	// Insert a test message
	now := time.Now()
	// Insert a message for testing

	message := models.Message{
		Content:   "Test Message",
		Recipient: "user@example.com",
		Status:    models.MessageStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	insertedMessage, err := store.insertTestMessages(ctx, message)
	if err != nil {
		t.Fatalf("Failed to insert message: %v", err)
	}

	// Update the message status
	err = store.UpdateMessageStatus(context.Background(), insertedMessage.ID, models.MessageStatusProcessing)
	if err != nil {
		t.Errorf("UpdateMessageStatus() error = %v", err)
	}

	// Fetch the message back
	fetchedMessage, err := store.getTestMessage(ctx, insertedMessage.ID)
	if err != nil {
		t.Fatalf("Failed to fetch message: %v", err)
	}
	if fetchedMessage.Status != models.MessageStatusProcessing {
		t.Errorf("UpdateMessageStatus() failed to update message status")
	}
}
