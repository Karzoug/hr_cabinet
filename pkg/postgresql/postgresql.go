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
	cfg config
}

// NewDB создаёт объект DB с заданными параметрами и подключается к БД.
func NewDB(ctx context.Context, dsn string, opts ...Option) (*DB, error) {
	const op = "postgresql: new db"

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := &DB{
		cfg: config{
			Config:              pgxCfg,
			connAttempts:        defaultConnAttempts,
			connAttemptsTimeout: defaultConnAttemptsTimeout,
		},
	}

	for _, opt := range opts {
		opt(&db.cfg)
	}

	for db.cfg.connAttempts > 0 {
		if db.Pool, err = pgxpool.NewWithConfig(ctx, db.cfg.Config); err == nil {
			break
		}

		slog.Info(fmt.Sprintf("trying to connect: attempts left %v", db.cfg.connAttempts))

		time.Sleep(db.cfg.connAttemptsTimeout)

		db.cfg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

// Close дожидается завершения запросов и закрывает все открытые соединения.
func (db *DB) Close() {
	db.Pool.Close()
}
