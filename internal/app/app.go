package app

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/Employee-s-file-cabinet/backend/internal/config"
	httpsrv "github.com/Employee-s-file-cabinet/backend/internal/delivery/http"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/smtp"
	repopg "github.com/Employee-s-file-cabinet/backend/internal/repo/postgresql"
	repos3 "github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/service/auth"
	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/password"
	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/token"
	authdb "github.com/Employee-s-file-cabinet/backend/internal/service/auth/repo/postgres"
	"github.com/Employee-s-file-cabinet/backend/internal/service/recovery"
	recoverydb "github.com/Employee-s-file-cabinet/backend/internal/service/recovery/repo/postgres"
	recoverykv "github.com/Employee-s-file-cabinet/backend/internal/service/recovery/repo/ttlmap"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user"
	userdb "github.com/Employee-s-file-cabinet/backend/internal/service/user/repo/postgres"
	users3 "github.com/Employee-s-file-cabinet/backend/internal/service/user/repo/s3"
)

func Run(ctx context.Context, cfg *config.Config, logger *slog.Logger) error {
	db, err := repopg.New(ctx, cfg.PG)
	if err != nil {
		return err
	}
	defer db.Close()

	s3Client, err := repos3.NewClient(cfg.S3)
	if err != nil {
		return err
	}

	// create user service
	userDBRepo, err := userdb.NewUserStorage(db)
	if err != nil {
		return err
	}
	userFileRepo, err := users3.New(ctx, s3Client, cfg.S3)
	if err != nil {
		return err
	}
	userService := user.NewService(userDBRepo, userFileRepo)

	// create auth service
	tokenMng, err := token.NewPasetoMaker(cfg.HTTP.Token.SecretKey, cfg.HTTP.Token.Lifetime)
	if err != nil {
		return err
	}
	authDBRepo, err := authdb.NewStorage(db)
	if err != nil {
		return err
	}
	passVerification := password.New()
	authService := auth.NewService(authDBRepo, passVerification, tokenMng)

	// create recovery service
	recoveryKeyRepo := recoverykv.New[string, int](cfg.Recovery.CleanKeyInterval)
	defer recoveryKeyRepo.Close()
	recoveryDBRepo, err := recoverydb.NewStorage(db)
	if err != nil {
		return err
	}
	smtpClient := smtp.NewMock(cfg.Mail)
	recoveryService := recovery.NewService(recoveryDBRepo, recoveryKeyRepo, smtpClient, passVerification, cfg.Recovery)

	srv, err := httpsrv.New(cfg.HTTP, cfg.EnvType, userService, authService, recoveryService, logger)
	if err != nil {
		return err
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return srv.Run(ectx)
	})

	return eg.Wait()
}
