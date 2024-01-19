package user

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type userRepository interface {
	Exist(ctx context.Context, userID uint64) (bool, error)
	ListShortUserInfo(ctx context.Context, pms model.ListUsersParams) ([]model.ShortUserInfo, int, error)
	Get(ctx context.Context, userID uint64) (*model.User, error)
	GetExpandedUser(ctx context.Context, userID uint64) (*model.ExpandedUser, error)
	Add(ctx context.Context, user model.User) (uint64, error)
	Update(ctx context.Context, user model.User) error

	GetEducation(ctx context.Context, userID, educationID uint64) (*model.Education, error)
	ListEducations(ctx context.Context, userID uint64) ([]model.Education, error)
	AddEducation(ctx context.Context, userID uint64, ed model.Education) (uint64, error)
	UpdateEducation(ctx context.Context, userID uint64, ed model.Education) error

	ListTrainings(ctx context.Context, userID uint64) ([]model.Training, error)
	GetTraining(ctx context.Context, userID, trainingID uint64) (*model.Training, error)
	AddTraining(ctx context.Context, userID uint64, tr model.Training) (uint64, error)
	UpdateTraining(ctx context.Context, userID uint64, tr model.Training) error

	ListPassports(ctx context.Context, userID uint64) ([]model.Passport, error)
	GetPassport(ctx context.Context, userID, passportID uint64) (*model.Passport, error)
	AddPassport(ctx context.Context, userID uint64, p model.Passport) (uint64, error)
	UpdatePassport(ctx context.Context, userID uint64, p model.Passport) error

	ListVisas(ctx context.Context, userID uint64) ([]model.Visa, error)
	GetVisa(ctx context.Context, userID, visaID uint64) (*model.Visa, error)
	AddVisa(ctx context.Context, userID uint64, mv model.Visa) (uint64, error)
	UpdateVisa(ctx context.Context, userID uint64, v model.Visa) error

	GetVacation(ctx context.Context, userID, vacationID uint64) (*model.Vacation, error)
	ListVacations(ctx context.Context, userID uint64) ([]model.Vacation, error)
	AddVacation(ctx context.Context, userID uint64, v model.Vacation) (uint64, error)
	UpdateVacation(ctx context.Context, userID uint64, v model.Vacation) error

	GetScan(ctx context.Context, userID, scanID uint64) (*model.Scan, error)
	ListScans(ctx context.Context, userID uint64) ([]model.Scan, error)
	AddScan(ctx context.Context, userID uint64, ms model.Scan) (uint64, error)

	ListContracts(ctx context.Context, userID uint64) ([]model.Contract, error)
	GetContract(ctx context.Context, userID, contractID uint64) (*model.Contract, error)
	AddContract(ctx context.Context, userID uint64, tr model.Contract) (uint64, error)
	UpdateContract(ctx context.Context, userID uint64, mc model.Contract) error
}

type s3FileRepository interface {
	Upload(context.Context, s3.File) error
	Download(ctx context.Context, prefix, name, etag string) (file s3.File, closeFn func() error, err error)
	PresignedURL(ctx context.Context, prefix, name string) (string, error)
}
