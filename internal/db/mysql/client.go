package mysql

import (
	"database/sql"
	"fmt"

	"github.com/mehmetalisavas/message-sender/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewClient(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(cfg))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return db, nil
}

func dsn(cfg config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.MysqlUser, cfg.MysqlPassword, cfg.MysqlHost, cfg.MysqlDatabase)
}
