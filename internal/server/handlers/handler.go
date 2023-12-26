package handlers

import (
	"log/slog"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
)

var _ api.ServerInterface = (*handler)(nil)

type handler struct {
	fileRepository       S3FileRepository
	dbRepository         DBRepository
	passwordVerification PasswordVerification
	tokenManager         TokenManager
	logger               *slog.Logger
}

func New(dbRepository DBRepository,
	s3FileRepository S3FileRepository,
	passwordVerification PasswordVerification,
	tokenManager TokenManager,
	logger *slog.Logger) *handler {
	logger = logger.With(slog.String("from", "handler"))

	h := &handler{
		fileRepository:       s3FileRepository,
		dbRepository:         dbRepository,
		passwordVerification: passwordVerification,
		tokenManager:         tokenManager,
		logger:               logger,
	}

	return h
}
