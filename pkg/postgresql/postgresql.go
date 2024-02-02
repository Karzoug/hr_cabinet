package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultConnAttemptsTimeout = time.Second
	defaultConnAttempts        = 5
)

type config struct {
	*pgxpool.Config
	connAttempts        int
	connAttemptsTimeout time.Duration
}

// DB структура с настройками подключения к БД и доступом к текущему соединению.
type DB struct {
	*pgxpool.Pool
}

// NewDB создаёт объект DB с заданными параметрами и подключается к БД.
func NewDB(ctx context.Context, dsn string, opts ...Option) (DB, error) {
	const op = "postgresql: new db"

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return DB{}, fmt.Errorf("%s: %w", op, err)
	}

	cfg := &config{
		Config:              pgxCfg,
		connAttempts:        defaultConnAttempts,
		connAttemptsTimeout: defaultConnAttemptsTimeout,
	}
	db := DB{}

	for _, opt := range opts {
		opt(cfg)
	}

	for cfg.connAttempts > 0 {
		if db.Pool, err = pgxpool.NewWithConfig(ctx, cfg.Config); err == nil {
			break
		}

		slog.Info(fmt.Sprintf("trying to connect: attempts left %v", cfg.connAttempts))

		time.Sleep(cfg.connAttemptsTimeout)

		cfg.connAttempts--
	}

	if err != nil {
		return DB{}, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

// Close дожидается завершения запросов и закрывает все открытые соединения.
func (db *DB) Close() {
	db.Pool.Close()
}
