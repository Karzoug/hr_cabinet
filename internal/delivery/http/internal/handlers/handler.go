package handlers

import (
	"log/slog"

	"github.com/Employee-s-file-cabinet/backend/internal/config/env"
	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
)

var _ api.ServerInterface = (*handler)(nil)

type handler struct {
	userService             UserService
	authService             AuthService
	passwordRecoveryService PasswordRecoveryService
	envType                 env.Type
	logger                  *slog.Logger
}

func New(envType env.Type, userService UserService,
	authService AuthService,
	passwordRecoveryService PasswordRecoveryService,
	logger *slog.Logger) *handler {
	return &handler{
		envType:                 envType,
		logger:                  logger,
		userService:             userService,
		authService:             authService,
		passwordRecoveryService: passwordRecoveryService,
	}
}
