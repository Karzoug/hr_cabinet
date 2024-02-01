package postgresql

import (
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type vacation struct {
	ID        uint64    `db:"id"`
	DateBegin time.Time `db:"date_begin"`
	DateEnd   time.Time `db:"date_end"`
}

func convertVacationToModelVacation(v vacation) model.Vacation {
	return model.Vacation(v)
}
