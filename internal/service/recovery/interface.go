package recovery

import (
	"context"
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/recovery/model"
)

type notificationDeliverer interface {
	SendMessage(recipient, subject, msg string) error
}

type recoveryRepository interface {
	CheckAndReturnUser(ctx context.Context, login string) (*model.User, error)
	ChangePassword(ctx context.Context, userID int, hash string) error
}

type keyRepository interface {
	Set(key string, value int, duration time.Duration) error
	Get(key string) (int, bool)
	Delete(key string) error
}

// passwordVerification абстракция хеширования паролей.
type passwordVerificator interface {
	// Hash - хеширование пароля.
	Hash(password string) (string, error)
}
