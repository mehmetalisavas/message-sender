package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mehmetalisavas/message-sender/pkg/retry"
)

type NotificationSender interface {
	Send(ctx context.Context, recipient, content string) (*NotificationResponse, error)
}

// NotificationService provides methods to send notifications
type NotificationService struct {
	client  *http.Client
	baseURL string
}

// NewNotificationService initializes a new NotificationService instance
func NewNotificationService(baseURL string, timeout time.Duration) *NotificationService {
	return &NotificationService{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

// NotificationRequest represents the request payload
type NotificationRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

// NotificationResponse represents the response from the webhook
type NotificationResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

// Send sends a notification to the notification service
func (ns *NotificationService) Send(ctx context.Context, recipient, content string) (*NotificationResponse, error) {
	payload := NotificationRequest{
		To:      recipient,
		Content: content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	requestFn := func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, ns.baseURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		return ns.client.Do(req)
	}

	resp, err := retry.Retry(ctx, requestFn, retry.DefaultConfig)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response NotificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
