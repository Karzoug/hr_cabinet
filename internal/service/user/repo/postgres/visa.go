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

func (s *storage) ListVisas(ctx context.Context, userID, passportID uint64) ([]model.Visa, error) {
	const op = "postrgresql user storage: list visas"

	rows, err := s.DB.Query(ctx, `SELECT 
	id, number, issued_state, 
	valid_to, valid_from, number_entries 
	FROM visas
	WHERE visas.passport_id = @passport_id AND visas.user_id = @user_id`,
		pgx.NamedArgs{
			"passport_id": passportID,
			"user_id":     userID,
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
		visas[i] = convertVisaToModelVisa(ed)
	}

	return visas, nil
}

func (s *storage) GetVisa(ctx context.Context, userID, passportID, visaID uint64) (*model.Visa, error) {
	const op = "postrgresql user storage: get visa"

	rows, err := s.DB.Query(ctx,
		`SELECT id, number, issued_state, 
		valid_to, valid_from, number_entries 
		FROM visas
		WHERE visas.id = @visa_id AND visas.passport_id = @passport_id AND visas.user_id = @user_id`,
		pgx.NamedArgs{
			"visa_id":     visaID,
			"passport_id": passportID,
			"user_id":     userID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	p, err := pgx.CollectExactlyOneRow[visa](rows, pgx.RowToStructByNameLax[visa])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
	}

	med := convertVisaToModelVisa(p)
	return &med, nil
}

func (s *storage) AddVisa(ctx context.Context, userID, passportID uint64, mv model.Visa) (uint64, error) {
	const op = "postrgresql user storage: add visa"

	v := convertModelVisaToVisa(mv)

	row := s.DB.QueryRow(ctx,
		`INSERT INTO visas
			("user_id", "passport_id", "number", 
			"issued_state", "valid_from", "valid_to", "number_entries")
		VALUES (@user_id, @passport_id, @number, @issued_state, 
			@valid_from, @valid_to, @number_entries)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":        userID,
			"passport_id":    passportID,
			"number":         v.Number,
			"issued_state":   v.IssuedState,
			"valid_to":       v.ValidTo,
			"valid_from":     v.ValidFrom,
			"number_entries": v.NumberEntries,
		})

	if err := row.Scan(&v.ID); err != nil {
		if strings.Contains(err.Error(), "23") { // Integrity Constraint Violation
			if strings.Contains(err.Error(), "user_id") {
				return 0, fmt.Errorf("%s: the user does not exist: %w", op, repoerr.ErrRecordNotFound)
			}
			if strings.Contains(err.Error(), "passport_id") {
				return 0, fmt.Errorf("%s: the passport does not exist: %w", op, repoerr.ErrRecordNotFound)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return v.ID, nil
}

func (s *storage) UpdateVisa(ctx context.Context, userID, passportID uint64, mv model.Visa) error {
	const op = "postrgresql user storage: update visa"

	v := convertModelVisaToVisa(mv)

	tag, err := s.DB.Exec(ctx, `UPDATE visas
	SET number=@number, issued_state=@issued_state, 
	valid_from=@valid_from, valid_to=@valid_to, number_entries=@number_entries
	WHERE id=@id AND user_id=@user_id AND passport_id=@passport_id`,
		pgx.NamedArgs{
			"user_id":        userID,
			"passport_id":    passportID,
			"id":             v.ID,
			"number":         v.Number,
			"issued_state":   v.IssuedState,
			"valid_to":       v.ValidTo,
			"valid_from":     v.ValidFrom,
			"number_entries": v.NumberEntries,
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
