package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/s3"
	"github.com/go-chi/chi/v5"
)

const baseURL = "/api/v1"

var _ api.ServerInterface = (*server)(nil)

type s3FileRepository interface {
	UploadFile(context.Context, s3.File) error
	DownloadFile(ctx context.Context, prefix, name string) (s3.File, func() error, error)
}

type userRepository interface {
	ExistUser(ctx context.Context, userID int) (bool, error)
}

type server struct {
	httpServer     *http.Server
	fileRepository s3FileRepository
	userRepository userRepository
	logger         *slog.Logger
}

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func New(cfg Config, userRepository userRepository, s3FileRepository s3FileRepository, logger *slog.Logger) *server {
	logger = logger.With(slog.String("from", "http-server"))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	s := &server{
		httpServer:     srv,
		fileRepository: s3FileRepository,
		userRepository: userRepository,
		logger:         logger,
	}

	mux := chi.NewRouter()
	mux.NotFound(s.notFound)
	mux.MethodNotAllowed(s.methodNotAllowed)
	mux.Use(s.logAccess)
	mux.Use(s.recoverPanic)

	srv.Handler = api.HandlerWithOptions(s, api.ChiServerOptions{
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
