package service

import (
	"context"
	"time"

	"github.com/mehmetalisavas/message-sender/internal/db/mysql"
	"github.com/mehmetalisavas/message-sender/internal/models"
)

// Make sure SqlStore implements Storage interface.
var _ Storage = (*mysql.SqlStore)(nil)

// Storage represents the storage service.
type Storage interface {
	// ListSentMessages returns all sent messages according to given options.
	ListSentMessages(ctx context.Context, opts models.ListOptions) ([]models.Message, error)

	// GetPendingMessages returns pending messages from the storage in a given limit.
	GetPendingMessages(ctx context.Context, limit int) ([]models.Message, error)

	// UpdateMessageStatus updates the status of the message with the given id.
	UpdateMessageStatus(ctx context.Context, id int, status models.MessageStatus) error
}

// CacheStore represents the cache store service.
type CacheStore interface {
	CacheMessage(ctx context.Context, messageId string, sendTime time.Time) error
}
