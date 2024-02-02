package visa

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type visaService interface {
	Get(ctx context.Context, userID, visaID uint64) (*model.Visa, error)
	List(ctx context.Context, userID uint64) ([]model.Visa, error)
	Add(ctx context.Context, userID uint64, mv model.Visa) (uint64, error)
	Update(ctx context.Context, userID uint64, v model.Visa) error
}
