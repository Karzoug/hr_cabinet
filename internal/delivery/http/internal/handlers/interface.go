package handlers

import (
	"context"
	"time"

	"github.com/casbin/casbin/v2"

	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/token"
	umodel "github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type UserService interface {
	List(ctx context.Context, params umodel.ListUsersParams) (users []umodel.User, totalCount int, err error)
	Get(ctx context.Context, userID uint64) (*umodel.User, error)
	DownloadPhoto(ctx context.Context, userID uint64, hash string) (f umodel.File, closeFn func() error, err error)
	UploadPhoto(ctx context.Context, userID uint64, f umodel.File) error

	ListEducations(ctx context.Context, userID uint64) ([]umodel.Education, error)
	GetEducation(ctx context.Context, userID, educationID uint64) (*umodel.Education, error)
	AddEducation(ctx context.Context, userID uint64, ed umodel.Education) (uint64, error)

	GetTraining(ctx context.Context, userID, trainingID uint64) (*umodel.Training, error)
	ListTrainings(ctx context.Context, userID uint64) ([]umodel.Training, error)
	AddTraining(ctx context.Context, userID uint64, ed umodel.Training) (uint64, error)

	GetPassport(ctx context.Context, userID, passportID uint64) (*umodel.Passport, error)
	ListPassports(ctx context.Context, userID uint64) ([]umodel.Passport, error)
	AddPassport(ctx context.Context, userID uint64, ed umodel.Passport) (uint64, error)

	GetVisa(ctx context.Context, userID, passportID, visaID uint64) (*umodel.Visa, error)
	ListVisas(ctx context.Context, userID, passportID uint64) ([]umodel.Visa, error)
	AddVisa(ctx context.Context, userID, passportID uint64, mv umodel.Visa) (uint64, error)
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Expires() time.Time
	Payload(token string) (*token.Payload, error)
	PolicyEnforcer() (*casbin.Enforcer, error)
}

type PasswordRecoveryService interface {
	InitChangePassword(ctx context.Context, login string) error
	ChangePassword(ctx context.Context, key, newPassword string) error
	Check(ctx context.Context, key string) error
}
