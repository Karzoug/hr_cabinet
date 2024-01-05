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
	DownloadPhoto(ctx context.Context, userID uint64, hash string) (f umodel.File, closeFn func() error, err error)
	UploadPhoto(ctx context.Context, userID uint64, f umodel.File) error

	ListEducations(ctx context.Context, userID uint64) ([]umodel.Education, error)
	GetEducation(ctx context.Context, educationID uint64) (*umodel.Education, error)
	AddEducation(ctx context.Context, userID uint64, ed umodel.Education) (uint64, error)

	GetTraining(ctx context.Context, trainingID uint64) (*umodel.Training, error)
	ListTrainings(ctx context.Context, userID uint64) ([]umodel.Training, error)
	AddTraining(ctx context.Context, userID uint64, ed umodel.Training) (uint64, error)
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
