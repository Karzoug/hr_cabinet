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
}

type s3FileRepository interface {
	Upload(context.Context, s3.File) error
	Download(ctx context.Context, prefix, name, etag string) (file s3.File, closeFn func() error, err error)
}
