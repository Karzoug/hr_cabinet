package postgresql

import (
	"context"

	pg "github.com/Employee-s-file-cabinet/backend/pkg/postgresql"
)

func New(ctx context.Context, cfg Config) (pg.DB, error) {
	return pg.NewDB(ctx, cfg.DSN,
		pg.MaxOpenConn(cfg.MaxOpenConns),
		pg.ConnAttempts(cfg.ConnAttempts))
}
