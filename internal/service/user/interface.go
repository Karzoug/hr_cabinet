package user

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type userRepository interface {
	Exist(ctx context.Context, userID uint64) (bool, error)
	List(ctx context.Context, params model.ListUsersParams) (users []model.User, totalCount int, err error)
	Get(ctx context.Context, userID uint64) (*model.User, error)

	GetEducation(ctx context.Context, educationID uint64) (*model.Education, error)
	ListEducations(ctx context.Context, userID uint64) ([]model.Education, error)
	AddEducation(ctx context.Context, userID uint64, ed model.Education) (uint64, error)

	ListTrainings(ctx context.Context, userID uint64) ([]model.Training, error)
	GetTraining(ctx context.Context, trainingID uint64) (*model.Training, error)
	AddTraining(ctx context.Context, userID uint64, tr model.Training) (uint64, error)
}

type s3FileRepository interface {
	Upload(context.Context, s3.File) error
	Download(ctx context.Context, prefix, name, etag string) (file s3.File, closeFn func() error, err error)
}
