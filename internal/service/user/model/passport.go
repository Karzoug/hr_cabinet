package model

import "time"

type Passport struct {
	ID         uint64
	IssuedBy   string
	IssuedDate time.Time
	Number     string
	Type       PassportType
	VisasCount uint
	HasScan    bool
}

type PassportType string

const (
	PassportTypeExternal   PassportType = "external"
	PassportTypeForeigners PassportType = "foreigners"
	PassportTypeInternal   PassportType = "internal"
)

type ExpandedPassport struct {
	Passport
	Visas []Visa
}
