package postgresql

import (
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type contract struct {
	ID              uint64       `db:"id"`
	Number          string       `db:"number"`
	ContractType    contractType `db:"contract_type"`
	WorkTypeID      uint64       `db:"work_type_id"`
	ProbationPeriod *uint        `db:"probation_period"`
	DateBegin       time.Time    `db:"date_begin"`
	DateEnd         *time.Time   `db:"date_end"`
	HasScan         bool         `db:"has_scan"`
}

type contractType string

const (
	contractTypePermanent contractType = "Бессрочный"
	contractTypeTemporary contractType = "Срочный"
)

func convertFromDBO(c contract) model.Contract {
	mc := model.Contract{
		ID:              c.ID,
		Number:          c.Number,
		WorkTypeID:      c.WorkTypeID,
		ProbationPeriod: c.ProbationPeriod,
		DateBegin:       c.DateBegin,
		DateEnd:         c.DateEnd,
		HasScan:         c.HasScan,
	}

	switch c.ContractType {
	case contractTypePermanent:
		mc.Type = model.ContractTypePermanent
	case contractTypeTemporary:
		mc.Type = model.ContractTypeTemporary
	}

	return mc
}

func convertToDBO(mc model.Contract) contract {
	c := contract{
		ID:              mc.ID,
		Number:          mc.Number,
		WorkTypeID:      mc.WorkTypeID,
		ProbationPeriod: mc.ProbationPeriod,
		DateBegin:       mc.DateBegin,
		DateEnd:         mc.DateEnd,
	}

	switch mc.Type {
	case model.ContractTypePermanent:
		c.ContractType = contractTypePermanent
	case model.ContractTypeTemporary:
		c.ContractType = contractTypeTemporary
	}

	return c
}
