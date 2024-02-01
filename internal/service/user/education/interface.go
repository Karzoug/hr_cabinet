package education

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type dbRepository interface {
	Get(ctx context.Context, userID, educationID uint64) (*model.Education, error)
	List(ctx context.Context, userID uint64) ([]model.Education, error)
	Add(ctx context.Context, userID uint64, ed model.Education) (uint64, error)
	Update(ctx context.Context, userID uint64, ed model.Education) error
}
