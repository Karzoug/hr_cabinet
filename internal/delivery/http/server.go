package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/stdlib" // use as driver for sqlx

	"github.com/go-chi/chi/v5"
	"github.com/jub0bs/fcors"

	"github.com/Employee-s-file-cabinet/backend/internal/config/env"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	srverr "github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/errors"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/handler"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/middleware"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/response"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

type server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(cfg Config, envType env.Type,
	userService handler.UserService,
	authService handler.AuthService,
	passwordRecoveryService handler.PasswordRecoveryService,
	logger *slog.Logger) (*server, error) {
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

	handler := handler.New(envType, userService, authService, passwordRecoveryService, logger)

	mux := chi.NewRouter()
	mux.NotFound(srverr.NotFoundHandlerFn(logger))
	mux.MethodNotAllowed(srverr.MethodNotAllowedHandlerFn(logger))

	// Add middlewares
	mux.Use(middleware.LogAccessFn(s.logger))
	mux.Use(middleware.RecoverPanicFn(s.logger))

	// CORS middleware
	switch envType {
	case env.Development:
		cors, err := fcors.AllowAccessWithCredentials(
			fcors.FromOrigins(
				"https://localhost:*",
				"http://localhost:*"),
			fcors.WithAnyMethod(),
			fcors.WithAnyRequestHeaders(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create CORS middleware: %w", err)
		}
		mux.Use(cors)
	default:
	}

	// Authorization middleware
	e, err := authService.PolicyEnforcer()
	if err != nil {
		return nil, fmt.Errorf("failed to create authorization middleware: %w", err)
	}

	authz := middleware.Authorizer{
		TokenManager: authService,
		Enforcer:     e,
	}
	mux.Use(authz.AuthorizeMiddleware)

	srv.Handler = api.HandlerWithOptions(handler, api.ChiServerOptions{
		BaseURL:    api.BaseURL,
		BaseRouter: mux,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Debug("request error", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			var msg string
			var fmtErr *api.InvalidParamFormatError
			if errors.As(err, &fmtErr) {
				msg = "Invalid format for parameter: " + fmtErr.ParamName
			} else {
				msg = err.Error()
			}
			if err := response.JSON(w, http.StatusBadRequest, api.Error{Message: msg}); err != nil {
				srverr.LogError(r, err, false, s.logger)
				srverr.ResponseError(w, r,
					http.StatusInternalServerError,
					srverr.ErrInternalServerErrorMsg, s.logger)
			}
		},
	})

	return s, nil
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
