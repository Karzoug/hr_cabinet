package postgresql

import (
	"time"

	"github.com/Employee-s-file-cabinet/backend/internal/service/user/model"
)

type education struct {
	ID                uint64    `db:"id"`
	Number            string    `db:"document_number"`
	Program           string    `db:"title_of_program"`
	IssuedInstitution string    `db:"title_of_institution"`
	DateTo            time.Time `db:"year_of_end"`
	DateFrom          time.Time `db:"year_of_begin"`
	HasScan           bool      `db:"has_scan"`
}

func convertFromDBO(ed education) model.Education {
	return model.Education(ed)
}
