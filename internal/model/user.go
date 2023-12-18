package model

import (
	"time"
)

const (
	DateLayout = "02.01.2006"
)

type User struct {
	ID         int    `db:"id"`
	Lastname   string `db:"lastname"`
	Firstname  string `db:"firstname"`
	Middlename string `db:"patronymic"`
	// etc...
}

// ConvertStringToDate использует формат DateLayout для преобразования входящей строки JSON в дату.
// Строка должна быть предварительно проверена на соответствие формату time.Time
func ConvertStringToDate(s string) time.Time {
	date, _ := time.Parse(DateLayout, s)
	return date
}
