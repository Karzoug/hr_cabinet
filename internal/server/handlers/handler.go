package handlers

import (
	"log/slog"

	"github.com/Employee-s-file-cabinet/backend/internal/server/internal/api"
	"github.com/Employee-s-file-cabinet/backend/internal/utils/email"
)

var _ api.ServerInterface = (*handler)(nil)

type handler struct {
	fileRepository       S3FileRepository
	dbRepository         DBRepository
	passwordVerification PasswordVerification
	tokenManager         TokenManager
	keyRepository        KeyRepository
	mail                 *email.Mail
	logger               *slog.Logger
}

func New(dbRepository DBRepository,
	s3FileRepository S3FileRepository,
	passwordVerification PasswordVerification,
	tokenManager TokenManager,
	keyRepository KeyRepository,
	mail *email.Mail,
	logger *slog.Logger) *handler {
	logger = logger.With(slog.String("from", "handler"))

	h := &handler{
		fileRepository:       s3FileRepository,
		dbRepository:         dbRepository,
		passwordVerification: passwordVerification,
		tokenManager:         tokenManager,
		keyRepository:        keyRepository,
		mail:                 mail,
		logger:               logger,
	}

	return h
}
