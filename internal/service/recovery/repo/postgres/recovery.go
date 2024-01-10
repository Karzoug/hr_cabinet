package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Employee-s-file-cabinet/backend/internal/service/recovery/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *storage) CheckAndReturnUser(ctx context.Context, login string) (*model.User, error) {
	const op = "postgresql recovery storage: check user exists"

	rows, err := s.Query(ctx,
		`SELECT users.id AS id, lastname, firstname, middlename, work_email
		FROM users
		JOIN authorizations a ON users.id = a.user_id
		WHERE work_email=@login`,
		pgx.NamedArgs{"login": login})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[user])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mu := convertUserToModelUser(u)
	return &mu, nil
}

func (s *storage) ChangePassword(ctx context.Context, userID int, hash string) error {
	const op = "postgresql recovery storage: change password"

	_, err := s.Exec(ctx,
		`UPDATE authorizations
		SET password_hash = @pass_hash
		WHERE user_id=@id`,
		pgx.NamedArgs{
			"pass_hash": hash,
			"id":        userID,
		})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
