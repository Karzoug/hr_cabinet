package model

import "time"

type Passport struct {
	ID         uint64
	IssuedBy   string
	IssuedDate time.Time
	Number     string
	Type       PassportType
	VisasCount uint
}

type PassportType string

const (
	PassportTypeExternal   PassportType = "external"
	PassportTypeForeigners PassportType = "foreigners"
	PassportTypeInternal   PassportType = "internal"
)
