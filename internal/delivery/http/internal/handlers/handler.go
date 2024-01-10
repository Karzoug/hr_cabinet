package handlers

import (
	"log/slog"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
)

var _ api.ServerInterface = (*handler)(nil)

type handler struct {
	userService             UserService
	authService             AuthService
	passwordRecoveryService PasswordRecoveryService
	logger                  *slog.Logger
}

func New(userService UserService,
	authService AuthService,
	passwordRecoveryService PasswordRecoveryService,
	logger *slog.Logger) *handler {
	return &handler{
		logger:                  logger,
		userService:             userService,
		authService:             authService,
		passwordRecoveryService: passwordRecoveryService,
	}
}
