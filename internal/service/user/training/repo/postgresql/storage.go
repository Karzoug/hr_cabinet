package postgresql

import (
	pq "github.com/Employee-s-file-cabinet/backend/pkg/postgresql"
)

type storage struct {
	*pq.DB
}

func NewUserStorage(db *pq.DB) (*storage, error) {
	return &storage{
		DB: db,
	}, nil
}
