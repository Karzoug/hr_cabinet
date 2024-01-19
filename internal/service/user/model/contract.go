package model

import "time"

type Contract struct {
	ID              uint64
	Number          string
	HasScan         bool
	Type            contractType
	WorkTypeID      uint64
	WorkType        string
	ProbationPeriod *uint
	DateBegin       time.Time
	DateEnd         *time.Time
}

type contractType string

const (
	ContractTypePermanent contractType = "permanent"
	ContractTypeTemporary contractType = "temporary"
)
