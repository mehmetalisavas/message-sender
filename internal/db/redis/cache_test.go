// File: internal/db/redis/cache_test.go
package redis

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mehmetalisavas/message-sender/config"
	"github.com/sethvargo/go-envconfig"
)

func setupTestRedis() (*RedisCacheStore, func(), error) {

	ctx := context.Background()

	var c config.Config
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	rds, err := NewRedisCacheStore(context.Background(), c)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Redis client: %w", err)
	}

	cleanup := func() {
		rds.client.FlushAll(context.Background())
	}

	return rds, cleanup, nil
}

func TestRedisCacheStore_CacheMessage(t *testing.T) {
	store, cleanup, err := setupTestRedis()
	if err != nil {
		t.Fatalf("failed to set up test Redis: %v", err)
	}
	defer cleanup()

	tests := []struct {
		name      string
		messageId string
		sendTime  time.Time
	}{
		{
			name:      "Valid message caching",
			messageId: "123",
			sendTime:  time.Now(),
		},
		{
			name:      "Another valid message caching",
			messageId: "456",
			sendTime:  time.Now().Add(1 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.CacheMessage(context.Background(), tt.messageId, tt.sendTime)
			if err != nil {
				t.Errorf("CacheMessage() error = %v, wantErr %v", err, false)
				return
			}

			// Verify the value is stored correctly
			cachedValue, err := store.GetMessageValue(context.Background(), tt.messageId)
			if err != nil {
				t.Errorf("GetMessageValue() error = %v, wantErr %v", err, false)
				return
			}
			if cachedValue.Format(time.RFC3339) != tt.sendTime.Format(time.RFC3339) {
				t.Errorf("GetMessageValue() = %v, want %v", cachedValue.Format(time.RFC3339), tt.sendTime.Format(time.RFC3339))
			}
		})
	}
}

func TestRedisCacheStore_CacheMessage_Overwrite(t *testing.T) {
	store, cleanup, err := setupTestRedis()
	if err != nil {
		t.Fatalf("failed to set up test Redis: %v", err)
	}
	defer cleanup()

	messageId := "123"
	initialTime := time.Now()
	updatedTime := initialTime.Add(2 * time.Hour)

	// Cache the initial message
	err = store.CacheMessage(context.Background(), messageId, initialTime)
	if err != nil {
		t.Errorf("CacheMessage() error = %v, wantErr %v", err, false)
		return
	}

	// Overwrite the message with a new time
	err = store.CacheMessage(context.Background(), messageId, updatedTime)
	if err != nil {
		t.Errorf("CacheMessage() error = %v, wantErr %v", err, false)
		return
	}

	// Verify the value is updated correctly
	cachedValue, err := store.GetMessageValue(context.Background(), messageId)
	if err != nil {
		t.Errorf("GetMessageValue() error = %v, wantErr %v", err, false)
		return
	}
	if cachedValue.Format(time.RFC3339) != updatedTime.Format(time.RFC3339) {
		t.Errorf("GetMessageValue() = %v, want %v", cachedValue.Format(time.RFC3339), updatedTime.Format(time.RFC3339))
	}
}

func TestRedisCacheStore_CacheMessage_EmptyMessageId(t *testing.T) {
	store, cleanup, err := setupTestRedis()
	if err != nil {
		t.Fatalf("failed to set up test Redis: %v", err)
	}
	defer cleanup()

	messageId := ""
	sendTime := time.Now()

	// Attempt to cache a message with an empty ID
	err = store.CacheMessage(context.Background(), messageId, sendTime)
	if err != ErrEmptyMessageID {
		t.Errorf("CacheMessage() error = %v, wantErr %v", err, ErrEmptyMessageID)
	}
}

func TestRedisCacheStore_GetMessageValue_InvalidFormat(t *testing.T) {
	store, cleanup, err := setupTestRedis()
	if err != nil {
		t.Fatalf("failed to set up test Redis: %v", err)
	}
	defer cleanup()

	messageId := "invalidFormat"
	invalidValue := "not-a-timestamp"

	// Manually set an invalid value in Redis
	err = store.client.Set(context.Background(), fmt.Sprintf("message:%s", messageId), invalidValue, 0).Err()
	if err != nil {
		t.Fatalf("failed to set invalid value in Redis: %v", err)
	}

	// Attempt to retrieve the message
	_, err = store.GetMessageValue(context.Background(), messageId)
	if err == nil {
		t.Errorf("GetMessageValue() error = %v, wantErr %v", err, true)
	}
}
