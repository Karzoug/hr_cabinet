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

func (s *storage) ListPassports(ctx context.Context, userID uint64) ([]model.Passport, error) {
	const op = "postrgresql user storage: list passports"

	rows, err := s.DB.Query(ctx, `SELECT 
	id, number, type, issued_date, issued_by, 
	(SELECT COUNT(*) FROM visas WHERE visas.passport_id = passports.id) AS visas_count 
	FROM passports
	WHERE passports.user_id = @user_id`,
		pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	psps, err := pgx.CollectRows[passport](rows, pgx.RowToStructByNameLax[passport])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	passports := make([]model.Passport, len(psps))
	for i, ed := range psps {
		passports[i] = convertPassportToModelPassport(ed)
	}

	return passports, nil
}

func (s *storage) GetPassport(ctx context.Context, userID, passportID uint64) (*model.Passport, error) {
	const op = "postrgresql user storage: get passport"

	rows, err := s.DB.Query(ctx,
		`SELECT id, number, type, issued_date, issued_by,
		(SELECT COUNT(*) FROM visas WHERE visas.passport_id = passports.id) AS visas_count 
		FROM passports
		WHERE id = @passport_id AND user_id = @user_id`,
		pgx.NamedArgs{
			"passport_id": passportID,
			"user_id":     userID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	p, err := pgx.CollectExactlyOneRow[passport](rows, pgx.RowToStructByNameLax[passport])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repoerr.ErrRecordNotFound
	}

	med := convertPassportToModelPassport(p)
	return &med, nil
}

func (s *storage) AddPassport(ctx context.Context, userID uint64, mp model.Passport) (uint64, error) {
	const op = "postrgresql user storage: add passport"

	p := convertModelPassportToPassport(mp)

	row := s.DB.QueryRow(ctx, `INSERT INTO passports
		("user_id", "number", "type", "issued_date", "issued_by")
		VALUES (@user_id, @number, @type, @issued_date, @issued_by)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":     userID,
			"number":      p.Number,
			"type":        p.Type,
			"issued_date": p.IssuedDate,
			"issued_by":   p.IssuedBy,
		})

	if err := row.Scan(&p.ID); err != nil {
		if strings.Contains(err.Error(), "23") && // Integrity Constraint Violation
			strings.Contains(err.Error(), "user_id") {
			return 0, fmt.Errorf("the user does not exist: %w", repoerr.ErrConflict)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return p.ID, nil
}

func (s *storage) UpdatePassport(ctx context.Context, userID uint64, mp model.Passport) error {
	const op = "postrgresql user storage: update passport"

	p := convertModelPassportToPassport(mp)

	tag, err := s.DB.Exec(ctx, `UPDATE passports
	SET number=@number, 
	type=@type, 
	issued_date=@issued_date, 
	issued_by=@issued_by
	WHERE id=@id AND user_id=@user_id`,
		pgx.NamedArgs{
			"user_id":     userID,
			"id":          p.ID,
			"number":      p.Number,
			"type":        p.Type,
			"issued_date": p.IssuedDate,
			"issued_by":   p.IssuedBy,
		})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 { // it's ok for pgx
		return repoerr.ErrRecordNotAffected
	}
	return nil
}
