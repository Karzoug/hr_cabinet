package handlers

import (
	"context"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/model"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/s3"
)

type S3FileRepository interface {
	UploadFile(context.Context, s3.File) error
	DownloadFile(ctx context.Context, prefix, name string) (file s3.File, closeFn func() error, err error)
}

type DBRepository interface {
	ExistUser(ctx context.Context, userID int) (bool, error)
	GetAuthnData(ctx context.Context, login string) (model.AuthnDAO, error)
	ExistEmployee(ctx context.Context, userID int) (bool, error)
	ChangePass(ctx context.Context, login, hash string) error
}

type KeyRepository interface {
	Set(ctx context.Context, key string, value string, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

// PasswordVerification абстракция хеширования и проверки паролей.
type PasswordVerification interface {
	// Hash - хеширование пароля.
	Hash(password string) (string, error)

	// Check - проверка переданного пароля и оригинального хеша на соответствие.
	Check(password, hashedPassword string) error
}

// TokenManager абстракция для управления токенами.
type TokenManager interface {
	// Create создаёт токен для переданных данных и продолжительности действия.
	Create(data model.TokenData) (string, error)

	// Verify проверяет, является ли токен действительным.
	Verify(in string) (Payload, error)
}

// Payload абстракция для полезной нагрузки токена
type Payload interface {
	Valid() error
	GetData() model.TokenData
}
