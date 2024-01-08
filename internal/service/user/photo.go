package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *service) DownloadPhoto(ctx context.Context, userID uint64, hash string) (model.File, func() error, error) {
	const op = "user service: download photo"

	// if exist, err := s.userRepository.Exist(ctx, userID); err != nil {
	// 	return model.File{}, nil, fmt.Errorf("%s: %w", op, err)
	// } else if !exist {
	// 	return model.File{}, nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
	// }

	f, closeFn, err := s.fileRepository.Download(ctx, strconv.FormatUint(userID, 10), photoFileName, hash)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotModified):
			return model.File{}, nil, fmt.Errorf("%s: %w", op, ErrPhotoFileNotModified)
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return model.File{}, nil, fmt.Errorf("%s: %w", op, ErrPhotoFileNotFound)
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
		return fmt.Errorf("%s: %w", op, ErrPhotoFileSizeTooLarge)
	}

	if exist, err := s.userRepository.Exist(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	} else if !exist {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
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
