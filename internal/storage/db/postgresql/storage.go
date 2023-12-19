package postgresql

import "github.com/Employee-s-file-cabinet/backend/pkg/errhelper"

type storage struct {
	*DB
}

func NewStorage(cfg Config) (*storage, error) {
	const op = "create user storage"

	db, err := NewDB(cfg.DSN,
		MaxOpenConn(cfg.MaxOpenConns),
		MaxIdleConn(cfg.MaxIdleConns),
		ConnAttempts(cfg.ConnAttempts))
	if err != nil {
		return nil, errhelper.Wrap(op, err)
	}

	return &storage{db}, nil
}
