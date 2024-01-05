package model

import "time"

type Education struct {
	ID                uint64
	Number            string
	Program           string
	IssuedInstitution string
	DateTo            time.Time
	DateFrom          time.Time
}
