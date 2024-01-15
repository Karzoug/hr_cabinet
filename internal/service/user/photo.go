package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	serr "github.com/Employee-s-file-cabinet/backend/internal/service"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) DownloadPhoto(ctx context.Context, userID uint64, hash string) (model.File, func() error, error) {
	const op = "user service: download photo"

	f, closeFn, err := s.fileRepository.Download(ctx, strconv.FormatUint(userID, 10), photoFileName, hash)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotModifiedSince):
			return model.File{}, nil, serr.NewError(serr.NotModified, "photo file not modified")
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return model.File{}, nil, serr.NewError(serr.NotFound, "photo file not found")
		default:
			return model.File{}, nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return model.File{
		ContentType: f.ContentType,
		Size:        f.Size,
		Hash:        f.ETag,
		Reader:      f.Reader,
	}, closeFn, nil
}

func (s *service) UploadPhoto(ctx context.Context, userID uint64, f model.File) error {
	const op = "user service: upload photo"

	if f.Size > MaxPhotoSize {
		return serr.NewError(serr.InvalidArgument, "photo file size too large")
	}

	if exist, err := s.userRepository.Exist(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	} else if !exist {
		return serr.NewError(serr.Conflict, "not uploaded: user not found")
	}

	if err := s.fileRepository.Upload(ctx, s3.File{
		Prefix:      strconv.FormatUint(userID, 10),
		Name:        photoFileName,
		Reader:      f.Reader,
		Size:        f.Size,
		ContentType: f.ContentType,
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
