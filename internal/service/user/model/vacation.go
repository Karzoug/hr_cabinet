package model

import "time"

type Vacation struct {
	ID        uint64
	DateBegin time.Time
	DateEnd   time.Time
}
