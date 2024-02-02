package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	pq "github.com/Employee-s-file-cabinet/backend/pkg/postgresql"
)

type storage struct {
	pq.DB
}

func New(db pq.DB) storage {
	return storage{DB: db}
}

func (s storage) Exist(ctx context.Context, userID uint64) (bool, error) {
	const op = "postrgresql user storage: exist user"

	row := s.DB.QueryRow(ctx, "SELECT COUNT(1) FROM users WHERE id = @user_id",
		pgx.NamedArgs{"user_id": userID})
	var count int
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}
