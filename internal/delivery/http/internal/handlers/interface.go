package handlers

import (
	"context"
	"time"

	"github.com/casbin/casbin/v2"

	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/token"
	umodel "github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type UserService interface {
	ListShortUserInfo(ctx context.Context, params umodel.ListUsersParams) (users []umodel.ShortUserInfo, totalCount int, err error)
	Get(ctx context.Context, userID uint64) (*umodel.User, error)
	GetExpanded(ctx context.Context, userID uint64) (*umodel.ExpandedUser, error)
	Add(ctx context.Context, u umodel.User) (uint64, error)
	Update(ctx context.Context, user umodel.User) error
	DownloadPhoto(ctx context.Context, userID uint64, hash string) (f umodel.File, closeFn func() error, err error)
	UploadPhoto(ctx context.Context, userID uint64, f umodel.File) error

	ListEducations(ctx context.Context, userID uint64) ([]umodel.Education, error)
	GetEducation(ctx context.Context, userID, educationID uint64) (*umodel.Education, error)
	AddEducation(ctx context.Context, userID uint64, ed umodel.Education) (uint64, error)
	UpdateEducation(ctx context.Context, userID uint64, ed umodel.Education) error

	GetTraining(ctx context.Context, userID, trainingID uint64) (*umodel.Training, error)
	ListTrainings(ctx context.Context, userID uint64) ([]umodel.Training, error)
	AddTraining(ctx context.Context, userID uint64, ed umodel.Training) (uint64, error)
	UpdateTraining(ctx context.Context, userID uint64, tr umodel.Training) error

	GetPassport(ctx context.Context, userID, passportID uint64) (*umodel.Passport, error)
	ListPassports(ctx context.Context, userID uint64) ([]umodel.Passport, error)
	AddPassport(ctx context.Context, userID uint64, ed umodel.Passport) (uint64, error)
	UpdatePassport(ctx context.Context, userID uint64, p umodel.Passport) error

	GetVisa(ctx context.Context, userID, passportID, visaID uint64) (*umodel.Visa, error)
	ListVisas(ctx context.Context, userID, passportID uint64) ([]umodel.Visa, error)
	AddVisa(ctx context.Context, userID, passportID uint64, mv umodel.Visa) (uint64, error)
	UpdateVisa(ctx context.Context, userID, passportID uint64, v umodel.Visa) error

	GetVacation(ctx context.Context, userID, vacationID uint64) (*umodel.Vacation, error)
	ListVacations(ctx context.Context, userID uint64) ([]umodel.Vacation, error)
	AddVacation(ctx context.Context, userID uint64, v umodel.Vacation) (uint64, error)
	UpdateVacation(ctx context.Context, userID uint64, v umodel.Vacation) error

	GetScan(ctx context.Context, userID, scanID uint64) (*umodel.Scan, error)
	ListScans(ctx context.Context, userID uint64) ([]umodel.Scan, error)
	UploadScan(ctx context.Context, userID uint64, ms umodel.Scan, f umodel.File) (uint64, error)

	GetContract(ctx context.Context, userID, contractID uint64) (*umodel.Contract, error)
	ListContracts(ctx context.Context, userID uint64) ([]umodel.Contract, error)
	AddContract(ctx context.Context, userID uint64, ed umodel.Contract) (uint64, error)
	UpdateContract(ctx context.Context, userID uint64, c umodel.Contract) error
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, string, error)
	Expires() time.Time
	Payload(token, sign string) (*token.Payload, error)
	PolicyEnforcer() (*casbin.Enforcer, error)
}

type PasswordRecoveryService interface {
	InitChangePassword(ctx context.Context, login string) error
	ChangePassword(ctx context.Context, key, newPassword string) error
	Check(ctx context.Context, key string) error
}
