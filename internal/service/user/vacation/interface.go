package vacation

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type dbRepository interface {
	Get(ctx context.Context, userID, vacationID uint64) (*model.Vacation, error)
	List(ctx context.Context, userID uint64) ([]model.Vacation, error)
	Add(ctx context.Context, userID uint64, v model.Vacation) (uint64, error)
	Update(ctx context.Context, userID uint64, v model.Vacation) error
}
