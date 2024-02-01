package training

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type dbRepository interface {
	List(ctx context.Context, userID uint64) ([]model.Training, error)
	Get(ctx context.Context, userID, trainingID uint64) (*model.Training, error)
	Add(ctx context.Context, userID uint64, tr model.Training) (uint64, error)
	Update(ctx context.Context, userID uint64, tr model.Training) error
}
