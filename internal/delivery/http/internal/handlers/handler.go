package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/delivery/http/internal/api"
	umodel "github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

var _ api.ServerInterface = (*handler)(nil)

type UserService interface {
	List(ctx context.Context, params umodel.ListUsersParams) (users []umodel.User, totalCount int, err error)
	Get(ctx context.Context, userID uint64) (*umodel.User, error)
	UploadPhoto(ctx context.Context, userID uint64, f umodel.File) error
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Expires() time.Time
}

type handler struct {
	userService UserService
	authService AuthService
	logger      *slog.Logger
}

func New(userService UserService,
	authService AuthService,
	logger *slog.Logger) *handler {
	return &handler{
		logger:      logger,
		userService: userService,
		authService: authService,
	}
}
