package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

func (s *storage) GetVacation(ctx context.Context, userID, vacationID uint64) (*model.Vacation, error) {
	const op = "postgresql user storage: get vacation"

	rows, err := s.DB.Query(ctx,
		`SELECT 
		id, date_begin, date_end 
		FROM vacations
		WHERE id = @id AND user_id = @user_id`,
		pgx.NamedArgs{
			"id":      vacationID,
			"user_id": userID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	v, err := pgx.CollectExactlyOneRow[vacation](rows, pgx.RowToStructByNameLax[vacation])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
	}

	mv := convertVacationToModelVacation(v)
	return &mv, nil
}

var listVacationsQuery = `SELECT 
		id, date_begin, date_end 
		FROM vacations
		WHERE user_id = @user_id`

func (s *storage) ListVacations(ctx context.Context, userID uint64) ([]model.Vacation, error) {
	const op = "postgresql user storage: list vacations"

	rows, err := s.DB.Query(ctx, listVacationsQuery, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	vs, err := pgx.CollectRows[vacation](rows, pgx.RowToStructByNameLax[vacation])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	vacations := make([]model.Vacation, len(vs))
	for i, v := range vs {
		vacations[i] = convertVacationToModelVacation(v)
	}

	return vacations, nil
}

func (s *storage) AddVacation(ctx context.Context, userID uint64, v model.Vacation) (uint64, error) {
	const op = "postgresql user storage: add vacation"

	row := s.DB.QueryRow(ctx, `INSERT INTO vacations
		("user_id", "date_begin", "date_end")
		VALUES (@user_id, @date_begin, @date_end)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":    userID,
			"date_begin": v.DateBegin,
			"date_end":   v.DateEnd,
		})

	if err := row.Scan(&v.ID); err != nil {
		if strings.Contains(err.Error(), "23") && // Integrity Constraint Violation
			strings.Contains(err.Error(), "user_id") {
			return 0, fmt.Errorf("%s: the user does not exist: %w", op, repoerr.ErrRecordNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return v.ID, nil
}

func (s *storage) UpdateVacation(ctx context.Context, userID uint64, v model.Vacation) error {
	const op = "postrgresql user storage: update vacation"

	tag, err := s.DB.Exec(ctx, `UPDATE vacations
	SET date_begin=@date_begin, date_end=@date_end
	WHERE id=@id AND user_id=@user_id`,
		pgx.NamedArgs{
			"user_id":    userID,
			"id":         v.ID,
			"date_begin": v.DateBegin,
			"date_end":   v.DateEnd,
		})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 { // it's ok for pgx
		return fmt.Errorf("%s: %w and %w", op,
			repoerr.ErrRecordNotModified,
			repoerr.ErrRecordNotFound)
	}
	return nil
}
