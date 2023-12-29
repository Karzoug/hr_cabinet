package auth

import (
	"context"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model"
	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/model/token"
)

type authRepository interface {
	Get(ctx context.Context, login string) (model.AuthnDAO, error)
}

// passwordVerification абстракция хеширования и проверки паролей.
type passwordVerificator interface {
	// Hash - хеширование пароля.
	Hash(password string) (string, error)

	// Check - проверка переданного пароля и оригинального хеша на соответствие.
	Check(password, hashedPassword string) error
}

// tokenManager абстракция для управления токенами.
type tokenManager interface {
	// Create создаёт токен для переданных данных и продолжительности действия.
	Create(data token.Data) (string, error)

	// Verify проверяет, является ли токен действительным.
	Verify(in string) (*token.Payload, error)

	// Expires возвращает время истечения срока годности токена (начиная с текущего момента времени).
	Expires() time.Time
}
