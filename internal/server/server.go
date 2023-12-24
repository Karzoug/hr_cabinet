package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	srvErrors "github.com/Employee-s-file-cabinet/backend/internal/server/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/server/handlers"
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/server/middleware"
	"github.com/Employee-s-file-cabinet/backend/internal/utils/email"
	"github.com/Employee-s-file-cabinet/backend/internal/utils/password"
)

const baseURL = "/api/v1"

type server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func New(cfg Config,
	dbRepository handlers.DBRepository,
	s3FileRepository handlers.S3FileRepository,
	tokenManager handlers.TokenManager,
	keyRepository handlers.KeyRepository,
	mail email.Mail,
	logger *slog.Logger) *server {

	logger = logger.With(slog.String("from", "http-server"))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	s := &server{
		httpServer: srv,
		logger:     logger,
	}

	passwordVerification := password.New()
	handler := handlers.New(dbRepository, s3FileRepository, passwordVerification, tokenManager, keyRepository, mail, logger)

	mux := chi.NewRouter()
	mux.NotFound(srvErrors.NotFound)
	mux.MethodNotAllowed(srvErrors.MethodNotAllowed)
	mux.Use(middleware.LogAccess)
	mux.Use(middleware.RecoverPanic)

	srv.Handler = api.HandlerWithOptions(handler, api.ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: mux,
	})

	return s
}

func (s *server) Run(ctx context.Context) error {
	shutdownErrorChan := make(chan error)

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- s.httpServer.Shutdown(ctx)
	}()

	s.logger.Info("starting server", slog.String("addr", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	s.logger.Info("stopped server", slog.String("addr", s.httpServer.Addr))

	return nil
}
