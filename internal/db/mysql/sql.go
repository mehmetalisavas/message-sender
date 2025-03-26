package mysql

import (
	"database/sql"
)

// SqlStore represents a MySQL store.
type SqlStore struct {
	db *sql.DB
}

// NewSqlStore creates a new SqlStore.
func NewSqlStore(client *sql.DB) *SqlStore {
	return &SqlStore{
		db: client,
	}
}
