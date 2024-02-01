package scan

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type userRepository interface {
	Exist(ctx context.Context, userID uint64) (bool, error)
}

type scanRepository interface {
	Get(ctx context.Context, userID, scanID uint64) (*model.Scan, error)
	List(ctx context.Context, userID uint64) ([]model.Scan, error)
	Add(ctx context.Context, userID uint64, ms model.Scan) (uint64, error)
}

type s3FileRepository interface {
	Upload(context.Context, s3.File) error
	PresignedURL(ctx context.Context, prefix, name string) (string, error)
}
