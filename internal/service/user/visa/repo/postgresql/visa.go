package postgresql

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

func (s storage) List(ctx context.Context, userID uint64) ([]model.Visa, error) {
	const op = "postrgresql visa storage: list"

	rows, err := s.DB.Query(ctx, `SELECT 
	id, number, issued_state, 
	valid_to, valid_from, type 
	FROM visas
	WHERE visas.user_id = @user_id`,
		pgx.NamedArgs{
			"user_id": userID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	vs, err := pgx.CollectRows[visa](rows, pgx.RowToStructByNameLax[visa])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	visas := make([]model.Visa, len(vs))
	for i, ed := range vs {
		visas[i] = convertFromDBO(ed)
	}

	return visas, nil
}

func (s storage) Get(ctx context.Context, userID, visaID uint64) (*model.Visa, error) {
	const op = "postrgresql visa storage: get"

	rows, err := s.DB.Query(ctx,
		`SELECT id, number, issued_state, 
		valid_to, valid_from, type 
		FROM visas
		WHERE visas.id = @visa_id AND visas.user_id = @user_id`,
		pgx.NamedArgs{
			"visa_id": visaID,
			"user_id": userID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	p, err := pgx.CollectExactlyOneRow[visa](rows, pgx.RowToStructByNameLax[visa])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repoerr.ErrRecordNotFound
	}

	med := convertFromDBO(p)
	return &med, nil
}

func (s storage) Add(ctx context.Context, userID uint64, mv model.Visa) (uint64, error) {
	const op = "postrgresql visa storage: add"

	v := convertToDBO(mv)

	row := s.DB.QueryRow(ctx,
		`INSERT INTO visas
			("user_id", "number", 
			"issued_state", "valid_from", "valid_to", "type")
		VALUES (@user_id, @number, @issued_state, 
			@valid_from, @valid_to, @type)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":      userID,
			"number":       v.Number,
			"issued_state": v.IssuedState,
			"valid_to":     v.ValidTo,
			"valid_from":   v.ValidFrom,
			"type":         v.Type,
		})

	if err := row.Scan(&v.ID); err != nil {
		if strings.Contains(err.Error(), "23") { // Integrity Constraint Violation
			if strings.Contains(err.Error(), "user_id") {
				return 0, fmt.Errorf("the user does not exist: %w", repoerr.ErrConflict)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return v.ID, nil
}

func (s storage) Update(ctx context.Context, userID uint64, mv model.Visa) error {
	const op = "postrgresql visa storage: update"

	v := convertToDBO(mv)

	tag, err := s.DB.Exec(ctx, `UPDATE visas
	SET number = @number, issued_state = @issued_state, 
	valid_from = @valid_from, valid_to = @valid_to, type = @type
	WHERE id=@id AND user_id=@user_id`,
		pgx.NamedArgs{
			"user_id":      userID,
			"id":           v.ID,
			"number":       v.Number,
			"issued_state": v.IssuedState,
			"valid_to":     v.ValidTo,
			"valid_from":   v.ValidFrom,
			"type":         v.Type,
		})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 { // it's ok for pgx
		return repoerr.ErrRecordNotAffected
	}
	return nil
}
