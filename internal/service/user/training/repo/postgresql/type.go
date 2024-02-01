package postgresql

import (
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type training struct {
	ID                uint64    `db:"id"`
	Program           string    `db:"title_of_program"`
	IssuedInstitution string    `db:"title_of_institution"`
	Cost              uint64    `db:"cost"`
	DateTo            time.Time `db:"date_end"`
	DateFrom          time.Time `db:"date_begin"`
	HasScan           bool      `db:"has_scan"`
}

func convertTrainingToModelTraining(tr training) model.Training {
	return model.Training(tr)
}
