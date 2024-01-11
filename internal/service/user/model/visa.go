package model

import "time"

type Visa struct {
	ID            uint64
	Number        string
	IssuedState   string
	ValidTo       time.Time
	ValidFrom     time.Time
	NumberEntries VisaNumberEntries
}

type VisaNumberEntries string

const (
	VisaNumberEntriesMult VisaNumberEntries = "mult"
	VisaNumberEntriesN1   VisaNumberEntries = "1"
	VisaNumberEntriesN2   VisaNumberEntries = "2"
)
