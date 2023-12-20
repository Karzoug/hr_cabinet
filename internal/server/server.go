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
	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/server/middleware"
)

const baseURL = "/api/v1"

type server struct {
	httpServer *http.Server
	handler    api.ServerInterface
	logger     *slog.Logger
}

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func New(cfg Config, handler api.ServerInterface, logger *slog.Logger) *server {
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
		handler:    handler,
		logger:     logger,
	}

	mux := chi.NewRouter()
	mux.NotFound(srvErrors.NotFound)
	mux.MethodNotAllowed(srvErrors.MethodNotAllowed)
	mux.Use(middleware.LogAccess)
	mux.Use(middleware.RecoverPanic)

	srv.Handler = api.HandlerWithOptions(s.handler, api.ChiServerOptions{
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
