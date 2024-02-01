package user

import (
	"context"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type userRepository interface {
	Exist(ctx context.Context, userID uint64) (bool, error)
	ListShortUserInfo(ctx context.Context, pms model.ListUsersParams) ([]model.ShortUserInfo, int, error)
	Get(ctx context.Context, userID uint64) (*model.User, error)
	GetExpandedUser(ctx context.Context, userID uint64) (*model.ExpandedUser, error)
	Add(ctx context.Context, user model.User) (uint64, error)
	Update(ctx context.Context, user model.User) error
}
