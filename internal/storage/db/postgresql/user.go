package postgresql

import (
	"context"
)

func (s *storage) ExistUser(ctx context.Context, userID int) (bool, error) {
	const op = "postrgresql user storage: exist user"

	panic("not implemented")
}
