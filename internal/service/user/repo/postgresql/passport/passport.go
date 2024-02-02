package passport

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

func (s storage) List(ctx context.Context, userID uint64) ([]model.Passport, error) {
	const op = "postrgresql passport storage: list"

	rows, err := s.DB.Query(ctx,
		`SELECT id, number, citizenship, type, issued_date, issued_by, issued_by_code,		
		(SELECT COUNT(*)>0 FROM scans WHERE scans.document_id=passports.id AND scans.type='Паспорт') AS has_scan
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
		passports[i] = convertFromDBO(ed)
	}

	return passports, nil
}

func (s storage) Get(ctx context.Context, userID, passportID uint64) (*model.Passport, error) {
	const op = "postrgresql passport storage: get"

	rows, err := s.DB.Query(ctx,
		`SELECT id, number, citizenship, type, issued_date, issued_by, issued_by_code,
		(SELECT COUNT(*)>0 FROM scans WHERE user_id=@user_id AND scans.document_id=passports.id AND scans.type='Паспорт') AS has_scan
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

	med := convertFromDBO(p)
	return &med, nil
}

func (s storage) Add(ctx context.Context, userID uint64, mp model.Passport) (uint64, error) {
	const op = "postrgresql passport storage: add"

	p := convertToDBO(mp)

	row := s.DB.QueryRow(ctx, `INSERT INTO passports
		("user_id", "number", "citizenship" "type", "issued_date", "issued_by", "issued_by_code")
		VALUES (@user_id, @number, @type, @issued_date, @issued_by, @issued_by_code)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":        userID,
			"number":         p.Number,
			"citizenship":    p.Citizenship,
			"type":           p.Type,
			"issued_date":    p.IssuedDate,
			"issued_by":      p.IssuedBy,
			"issued_by_code": p.IssuedByCode,
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

func (s storage) Update(ctx context.Context, userID uint64, mp model.Passport) error {
	const op = "postrgresql passport storage: update"

	p := convertToDBO(mp)

	tag, err := s.DB.Exec(ctx, `UPDATE passports
	SET number = @number, 
	citizenship = @citizenship,
	type = @type, 
	issued_date = @issued_date, 
	issued_by = @issued_by, issued_by_code = @issued_by_code
	WHERE id=@id AND user_id=@user_id`,
		pgx.NamedArgs{
			"user_id":        userID,
			"id":             p.ID,
			"number":         p.Number,
			"citizenship":    p.Citizenship,
			"type":           p.Type,
			"issued_date":    p.IssuedDate,
			"issued_by":      p.IssuedBy,
			"issued_by_code": p.IssuedByCode,
		})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 { // it's ok for pgx
		return repoerr.ErrRecordNotAffected
	}
	return nil
}
