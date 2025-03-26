package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mehmetalisavas/message-sender/config"
)

var ErrEmptyMessageID = errors.New("message ID cannot be empty")

type RedisCacheStore struct {
	client *redis.Client
}

// NewRedisCacheStore initializes a new Redis client
func NewRedisCacheStore(ctx context.Context, cfg config.Config) (*RedisCacheStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", cfg.RedisHost), // Redis port
		Password: cfg.RedisPassword,                     // Password if Redis is password protected
		DB:       0,                                     // Default database
	})

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	return &RedisCacheStore{client: rdb}, nil
}

func (r *RedisCacheStore) CacheMessage(ctx context.Context, messageId string, sendTime time.Time) error {
	if messageId == "" {
		return ErrEmptyMessageID
	}

	key := fmt.Sprintf("message:%s", messageId)
	value := sendTime.Format(time.RFC3339)

	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisCacheStore) GetMessageValue(ctx context.Context, messageId string) (time.Time, error) {
	key := fmt.Sprintf("message:%s", messageId)
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, value)
}
