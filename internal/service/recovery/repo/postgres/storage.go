package postgres

import (
	pq "github.com/Employee-s-file-cabinet/backend/pkg/postgresql"
)

type storage struct {
	*pq.DB
}

func NewStorage(db *pq.DB) (*storage, error) {
	return &storage{db}, nil
}
