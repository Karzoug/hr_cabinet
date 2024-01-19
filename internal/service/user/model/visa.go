package model

import "time"

type Visa struct {
	ID          uint64
	Number      string
	Type        VisaType
	IssuedState *string
	ValidTo     time.Time
	ValidFrom   time.Time
}

type VisaType string
