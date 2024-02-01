package photo

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
)

type userRepository interface {
	Exist(ctx context.Context, userID uint64) (bool, error)
}

type s3FileRepository interface {
	Upload(context.Context, s3.File) error
	Download(ctx context.Context, prefix, name, etag string) (file s3.File, closeFn func() error, err error)
}
