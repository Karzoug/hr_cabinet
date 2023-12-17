package app

import (
	"context"
	"log/slog"

	"github.com/Employee-s-file-cabinet/backend/internal/config"
	"github.com/Employee-s-file-cabinet/backend/internal/server"
	"golang.org/x/sync/errgroup"
)

func Run(pctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	db, err := postgresql.New(pctx, cfg.PG)
	if err != nil {
		return err
	}
	defer db.Close()

	srv := server.New(cfg.HTTP, db, nil, logger)

	eg, ctx := errgroup.WithContext(pctx)
	eg.Go(func() error {
		return srv.Run(ctx)
	})

	return eg.Wait()
}
