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

const listContractsQuery = `SELECT 
contracts.id as id, number, contract_type, work_type_id, probation_period, date_begin, date_end,
(SELECT COUNT(*)>0 FROM scans WHERE scans.document_id=contracts.id AND scans.type='Трудовой договор') AS has_scan
FROM contracts
WHERE user_id = @user_id`

func (s *storage) ListContracts(ctx context.Context, userID uint64) ([]model.Contract, error) {
	const op = "postgresql user storage: list contracts"

	rows, err := s.DB.Query(ctx, listContractsQuery, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	trs, err := pgx.CollectRows[contract](rows, pgx.RowToStructByNameLax[contract])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	contracts := make([]model.Contract, len(trs))
	for i, tr := range trs {
		contracts[i] = convertContractToModelContract(tr)
	}

	return contracts, nil
}

func (s *storage) GetContract(ctx context.Context, userID, contractID uint64) (*model.Contract, error) {
	const op = "postgresql user storage: get contract"

	//стр-ра
	rows, err := s.DB.Query(ctx,
		`SELECT 
		id, number, contract_type, work_type_id, probation_period, date_begin, date_end,
		(SELECT COUNT(*)>0 FROM scans WHERE user_id=@user_id AND scans.document_id=contracts.id AND scans.type='Трудовой договор') AS has_scan
		FROM contracts
		WHERE id = @contract_id AND user_id = @user_id`,
		pgx.NamedArgs{
			"contract_id": contractID,
			"user_id":     userID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	c, err := pgx.CollectExactlyOneRow[contract](rows, pgx.RowToStructByNameLax[contract])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repoerr.ErrRecordNotFound
	}

	mc := convertContractToModelContract(c)
	return &mc, nil
}

func (s *storage) AddContract(ctx context.Context, userID uint64, mc model.Contract) (uint64, error) {
	const op = "postgresql user storage: add contract"

	c := convertModelContractToContract(mc)

	row := s.DB.QueryRow(ctx, `INSERT INTO contracts
		("user_id", "number", "contract_type", "work_type_id", "probation_period", "date_begin", "date_end")
		VALUES (@user_id, @number, @contract_type, @work_type_id, @probation_period, @date_begin, @date_end)
		RETURNING "id"`,
		pgx.NamedArgs{
			"user_id":          userID,
			"number":           c.Number,
			"contract_type":    c.ContractType,
			"work_type_id":     c.WorkTypeID,
			"probation_period": c.ProbationPeriod,
			"date_begin":       c.DateBegin,
			"date_end":         c.DateEnd,
		})

	if err := row.Scan(&c.ID); err != nil {
		if strings.Contains(err.Error(), "23") { // Integrity Constraint Violation
			if strings.Contains(err.Error(), "user_id") {
				return 0, fmt.Errorf("the user does not exist: %w", repoerr.ErrConflict)
			}
			if strings.Contains(err.Error(), "work_type_id") {
				return 0, fmt.Errorf("the work type does not exist: %w", repoerr.ErrConflict)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return c.ID, nil
}

func (s *storage) UpdateContract(ctx context.Context, userID uint64, mc model.Contract) error {
	const op = "postgresql user storage: update contract"

	c := convertModelContractToContract(mc)

	tag, err := s.DB.Exec(ctx, `UPDATE contracts
		SET number = @number, contract_type = @contract_type, work_type_id = @work_type_id, 
		probation_period = @probation_period, date_begin = @date_begin, date_end = @date_end
		WHERE id=@id AND user_id=@user_id`,
		pgx.NamedArgs{
			"id":               c.ID,
			"user_id":          userID,
			"number":           c.Number,
			"contract_type":    c.ContractType,
			"work_type_id":     c.WorkTypeID,
			"probation_period": c.ProbationPeriod,
			"date_begin":       c.DateBegin,
			"date_end":         c.DateEnd,
		})

	if err != nil {
		if strings.Contains(err.Error(), "23") { // Integrity Constraint Violation
			if strings.Contains(err.Error(), "work_type_id") {
				return fmt.Errorf("the work type does not exist: %w", repoerr.ErrConflict)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 { // it's ok for pgx
		return repoerr.ErrRecordNotAffected
	}
	return nil
}
