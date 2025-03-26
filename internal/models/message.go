package models

import "time"

type MessageStatus string

const (
	MessageStatusPending    MessageStatus = "pending"
	MessageStatusProcessing MessageStatus = "processing"
	MessageStatusSent       MessageStatus = "sent"
	MessageStatusFailed     MessageStatus = "failed"
)

type Message struct {
	ID        int           `json:"id"`
	Recipient string        `json:"recipient"`
	Content   string        `json:"content"`
	Status    MessageStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
