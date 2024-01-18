package model

import "time"

type Training struct {
	ID                uint64
	Program           string
	IssuedInstitution string
	Cost              uint64
	DateTo            time.Time
	DateFrom          time.Time
	HasScan           bool
}
