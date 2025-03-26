package mysql

import (
	"database/sql"
)

type SqlStore struct {
	db *sql.DB
}

func NewSqlStore(client *sql.DB) *SqlStore {
	return &SqlStore{
		db: client,
	}
}
