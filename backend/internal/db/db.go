package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"cau-used-goods-app/backend/internal/config"
)

var defaultDB *sql.DB

func Init(cfg config.DatabaseConfig) error {
	conn, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return fmt.Errorf("open mysql connection: %w", err)
	}

	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetConnMaxLifetime(cfg.ConnMaxLifetime())

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return fmt.Errorf("ping mysql: %w", err)
	}

	defaultDB = conn
	return nil
}

func DB() *sql.DB {
	return defaultDB
}

func Close() error {
	if defaultDB == nil {
		return nil
	}
	return defaultDB.Close()
}

func WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	if defaultDB == nil {
		return fmt.Errorf("database is not initialized")
	}

	tx, err := defaultDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx failed: %w; rollback failed: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
