package contract

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type dbRepository interface {
	List(ctx context.Context, userID uint64) ([]model.Contract, error)
	Get(ctx context.Context, userID, contractID uint64) (*model.Contract, error)
	Add(ctx context.Context, userID uint64, tr model.Contract) (uint64, error)
	Update(ctx context.Context, userID uint64, mc model.Contract) error
}
