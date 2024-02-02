package vacation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
	pq "github.com/Employee-s-file-cabinet/backend/pkg/postgresql"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

type storage struct {
	pq.DB
}

func New(db pq.DB) storage {
	return storage{DB: db}
}

func (s storage) Get(ctx context.Context, userID, vacationID uint64) (*model.Vacation, error) {
	const op = "postgresql vacation storage: get"

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
		return nil, repoerr.ErrRecordNotFound
	}

	mv := convertFromDBO(v)
	return &mv, nil
}

func (s storage) List(ctx context.Context, userID uint64) ([]model.Vacation, error) {
	const op = "postgresql vacation storage: list"

	rows, err := s.DB.Query(ctx, `SELECT 
	id, date_begin, date_end 
	FROM vacations
	WHERE user_id = @user_id`,
		pgx.NamedArgs{"user_id": userID})
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
		vacations[i] = convertFromDBO(v)
	}

	return vacations, nil
}

func (s storage) Add(ctx context.Context, userID uint64, v model.Vacation) (uint64, error) {
	const op = "postgresql vacation storage: add"

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
			return 0, fmt.Errorf("the user does not exist: %w", repoerr.ErrConflict)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return v.ID, nil
}

func (s storage) Update(ctx context.Context, userID uint64, v model.Vacation) error {
	const op = "postrgresql vacation storage: update"

	tag, err := s.DB.Exec(ctx, `UPDATE vacations
	SET date_begin = @date_begin, date_end = @date_end
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
		return repoerr.ErrRecordNotAffected
	}
	return nil
}
