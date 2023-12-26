package app

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/Employee-s-file-cabinet/backend/internal/config"
	"github.com/Employee-s-file-cabinet/backend/internal/server"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/db/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/utils/token"
)

func Run(pctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	db, err := postgresql.NewStorage(cfg.PG)
	if err != nil {
		return err
	}
	defer db.Close()

	s3Storage, err := s3.New(pctx, cfg.S3)
	if err != nil {
		return err
	}

	tokenManager, err := token.NewPasetoMaker(cfg.HTTP.Token.SecretKey, cfg.HTTP.Token.Lifetime)
	if err != nil {
		return err
	}

	srv := server.New(cfg.HTTP, db, s3Storage, tokenManager, logger)

	eg, ctx := errgroup.WithContext(pctx)
	eg.Go(func() error {
		return srv.Run(ctx)
	})

	return eg.Wait()
}
