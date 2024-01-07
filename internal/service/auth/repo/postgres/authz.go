package postgresql

import (
	_ "github.com/jackc/pgx/stdlib" // use as driver for sqlx

	"github.com/Employee-s-file-cabinet/backend/internal/service/auth/repo/sqlxadapter"
)

func (s *storage) PolicyAdapter() *sqlxadapter.Adapter {
	opts := &sqlxadapter.AdapterOptions{
		DriverName:     "pgx",
		DataSourceName: s.Config().ConnString(),
		TableName:      "policies",
	}
	return sqlxadapter.NewAdapterFromOptions(opts)
}
