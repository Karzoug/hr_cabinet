package passport

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type dbRepository interface {
	List(ctx context.Context, userID uint64) ([]model.Passport, error)
	Get(ctx context.Context, userID, passportID uint64) (*model.Passport, error)
	Add(ctx context.Context, userID uint64, p model.Passport) (uint64, error)
	Update(ctx context.Context, userID uint64, p model.Passport) error
}
