package handlers

import (
	"context"
	"log/slog"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/s3"
)

var _ api.ServerInterface = (*handler)(nil)

type S3FileUploader interface {
	UploadFile(context.Context, s3.UploadableFile) error
}

type UserRepository interface {
	ExistUser(ctx context.Context, userID int) (bool, error)
}

type handler struct {
	fileUploader   S3FileUploader
	userRepository UserRepository
	logger         *slog.Logger
}

func New(userRepository UserRepository, s3FileUploader S3FileUploader, logger *slog.Logger) *handler {
	logger = logger.With(slog.String("from", "handler"))

	h := &handler{
		fileUploader:   s3FileUploader,
		userRepository: userRepository,
		logger:         logger,
	}

	return h
}
