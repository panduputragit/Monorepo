package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrMissingURL = errors.New("database URL is empty")

type Config struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func Connect(ctx context.Context, cfg Config) (*sql.DB, error) {
	if cfg.URL == "" {
		return nil, ErrMissingURL
	}

	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, err
	}

	maxOpenConns := cfg.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 25
	}

	maxIdleConns := cfg.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 25
	}

	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime == 0 {
		connMaxLifetime = 5 * time.Minute
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func ConnectOptional(ctx context.Context, cfg Config) (*sql.DB, error) {
	db, err := Connect(ctx, cfg)
	if errors.Is(err, ErrMissingURL) {
		return nil, nil
	}
	return db, err
}
