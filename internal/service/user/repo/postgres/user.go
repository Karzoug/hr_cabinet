package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/henvic/pgq"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const LimitListUsers = 10

func (s *storage) Exist(ctx context.Context, userID uint64) (bool, error) {
	const op = "postrgresql user storage: exist user"

	var count int
	if err := pgxscan.Get(ctx, s.DB, &count,
		"SELECT COUNT(1) FROM users WHERE id = $1", userID); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (s *storage) Get(ctx context.Context, userID uint64) (*model.User, error) {
	const op = "postrgresql user storage: get user"

	u := new(user)
	if err := pgxscan.Get(ctx, s.DB, u,
		`SELECT 
		users.id AS id,lastname,firstname,middlename,gender,
		date_of_birth,place_of_birth,grade,phone_numbers,
		work_email,registration_address,residential_address,nationality,
		insurance_number,taxpayer_number, 
		positions.title AS position, 
		departments.title AS department 
		FROM users		 
		JOIN departments ON users.department_id = departments.id
		JOIN positions ON users.position_id = positions.id
		WHERE users.id = $1`, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mu := convertUserToModelUser(u)
	return &mu, nil
}

func (s *storage) List(ctx context.Context, pms model.ListUsersParams) ([]model.User, int, error) {
	const op = "postrgresql user storage: list users"

	sb := pgq.
		Select(`users.id AS id,lastname,firstname,middlename,gender,
		date_of_birth,place_of_birth,grade,phone_numbers,
		work_email,registration_address,residential_address,nationality,
		insurance_number,taxpayer_number, 
		positions.title AS position, 
		departments.title AS department, count(*) OVER() AS total_count`).
		From("users").
		Join("departments ON users.department_id = departments.id").
		Join("positions ON users.position_id = positions.id")
	if pms.Query != "" {
		sb = sb.Where(pgq.ILike{"lastname": pms.Query + "%"})
	}
	// nolint:exhaustive
	switch pms.SortBy {
	case model.ListUsersParamsSortByDepartment:
		sb = sb.OrderBy("department, lastname, firstname")
	default:
		sb = sb.OrderBy("lastname, firstname")
	}
	sb = sb.
		Limit(uint64(pms.Limit)).
		Offset(uint64((pms.Page - 1) * pms.Limit))
	query, args, err := sb.SQL()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	users := make([]model.User, 0)
	u := new(listUser)
	for rows.Next() {
		if err := pgxscan.ScanRow(u, rows); err != nil {
			return nil, 0, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, converListUserToModelUser(u))
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return users, u.TotalCount, nil
}
