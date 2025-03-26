package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/mehmetalisavas/message-sender/internal/models"
)

func (s *SqlStore) ListSentMessages(ctx context.Context, opts models.ListOptions) ([]models.Message, error) {
	options := models.InitWithDefaultListOptions(opts)

	query := `
		SELECT id, content, recipient, status, created_at, updated_at
		FROM messages
		WHERE status = ?
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, models.MessageStatusSent, options.Limit, options.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// initialize the slice with a length of 0 and a capacity of limit for better performance
	messages := make([]models.Message, 0, options.Limit)
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.Content, &m.Recipient, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

// GetPendingMessages returns pending messages from the storage in a given limit.
// func (s *SqlStore) GetPendingMessages(ctx context.Context, limit int) ([]models.Message, error) {
// 	query := `
// 		SELECT id, content, recipient, status, created_at, updated_at
// 		FROM messages
// 		WHERE status = ?
// 		ORDER BY created_at ASC
// 		LIMIT ?
// 	`

// 	rows, err := s.db.QueryContext(ctx, query, models.MessageStatusPending, limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// initialize the slice with a length of 0 and a capacity of limit for better performance
// 	messages := make([]models.Message, 0, limit)
// 	for rows.Next() {
// 		var m models.Message
// 		if err := rows.Scan(&m.ID, &m.Content, &m.Recipient, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		messages = append(messages, m)
// 	}

// 	return messages, nil
// }

// func (s *SqlStore) GetPendingMessages(ctx context.Context, limit int) ([]models.Message, error) {
// 	query := `
// 		UPDATE messages
// 		SET status = 'processing', updated_at = NOW()
// 		WHERE id IN (
// 			SELECT id
// 			FROM messages
// 			WHERE (status = 'pending' OR (status = 'processing' AND updated_at < NOW() - INTERVAL 5 MINUTE))
// 			ORDER BY created_at ASC
// 			LIMIT ?
// 			FOR UPDATE SKIP LOCKED
// 		)
// 		RETURNING id, content, recipient, status, created_at, updated_at
// 	`

// 	rows, err := s.db.QueryContext(ctx, query, limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// initialize the slice with a length of 0 and a capacity of limit for better performance
// 	messages := make([]models.Message, 0, limit)
// 	for rows.Next() {
// 		var m models.Message
// 		if err := rows.Scan(&m.ID, &m.Content, &m.Recipient, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		messages = append(messages, m)
// 	}

// 	return messages, nil
// }

func (s *SqlStore) GetPendingMessages(ctx context.Context, limit int) ([]models.Message, error) {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Ensure rollback in case of any error

	// Step 1: Select pending messages and lock them
	selectQuery := `
			SELECT id, content, recipient, status, created_at, updated_at
			FROM messages
			WHERE (status = 'pending' OR (status = 'processing' AND updated_at < NOW() - INTERVAL 5 MINUTE))
			ORDER BY created_at ASC
			LIMIT ?
			FOR UPDATE SKIP LOCKED
		`

	rows, err := tx.QueryContext(ctx, selectQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to hold the messages
	messages := make([]models.Message, 0, limit)
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.Content, &m.Recipient, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	// If no messages found, return early
	if len(messages) == 0 {
		return messages, nil
	}

	placeholder := "?" // Placeholder for each id
	placeholders := make([]string, len(messages))
	for i := range messages {
		placeholders[i] = placeholder
	}
	placeholderStr := strings.Join(placeholders, ",")

	updateQuery := fmt.Sprintf(`
		UPDATE messages 
		SET status = 'processing', updated_at = NOW() 
		WHERE id IN (%s)`, placeholderStr,
	)

	ids := make([]interface{}, len(messages))
	for i, msg := range messages {
		ids[i] = msg.ID
	}

	// Run the update query to mark the messages as processing
	_, err = tx.ExecContext(ctx, updateQuery, ids...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// UpdateMessageStatus updates the status of the message with the given ID.
func (s *SqlStore) UpdateMessageStatus(ctx context.Context, id int, status models.MessageStatus) error {
	query := `
		UPDATE messages
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query, status, id)
	return err
}

// InsertTestMessages inserts a test message into the database.
// Don't use this function in production code.
func (s *SqlStore) InsertTestMessages(ctx context.Context, message models.Message) (*models.Message, error) {
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

// GetTestMessage returns a test message from the database.
// Don't use this function in production code.
func (s *SqlStore) GetTestMessage(ctx context.Context, id int) (*models.Message, error) {
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
