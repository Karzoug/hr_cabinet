package scan

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

func (s storage) Get(ctx context.Context, userID, scanID uint64) (*model.Scan, error) {
	const op = "postgresql user storage: get scan"

	rows, err := s.DB.Query(ctx,
		`SELECT
		id, "type", document_id, description, created_at
		FROM scans
		WHERE id = @scan_id AND user_id = @user_id`,
		pgx.NamedArgs{
			"user_id": userID,
			"scan_id": scanID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	sc, err := pgx.CollectExactlyOneRow[scan](rows, pgx.RowToStructByNameLax[scan])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
	}

	ms := convertFromDBO(sc)
	return &ms, nil
}

func (s storage) List(ctx context.Context, userID uint64) ([]model.Scan, error) {
	const op = "postgresql user storage: list scans"

	rows, err := s.DB.Query(ctx,
		`SELECT 
		id, "type", document_id, description, created_at
		FROM scans
		WHERE user_id = @user_id`,
		pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	ss, err := pgx.CollectRows[scan](rows, pgx.RowToStructByNameLax[scan])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	scans := make([]model.Scan, len(ss))
	for i, sc := range ss {
		scans[i] = convertFromDBO(sc)
	}

	return scans, nil
}

func (s storage) Add(ctx context.Context, userID uint64, ms model.Scan) (uint64, error) {
	const op = "postgresql user storage: add scan"

	sc := convertToDBO(ms)

	row := s.DB.QueryRow(ctx, `INSERT INTO scans
		("user_id", "document_id", "type", "description")
		VALUES (@user_id, @document_id, @type, @description)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":     userID,
			"document_id": sc.DocumentID,
			"type":        sc.Type,
			"description": sc.Description,
		})

	if err := row.Scan(&sc.ID); err != nil {
		if strings.Contains(err.Error(), "23") && // Integrity Constraint Violation
			strings.Contains(err.Error(), "user_id") {
			return 0, fmt.Errorf("%s: the user does not exist: %w", op, repoerr.ErrRecordNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return sc.ID, nil
}
