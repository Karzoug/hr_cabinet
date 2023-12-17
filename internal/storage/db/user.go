package db

import (
	"fmt"

	"github.com/Employee-s-file-cabinet/backend/internal/config"
	"github.com/Employee-s-file-cabinet/backend/internal/model"
)

type UserStorage struct {
	*model.DB
}

func NewUserStorage(cfg config.PG) (*UserStorage, error) {
	db, err := model.New(cfg.DSN,
		model.MaxOpenConn(cfg.MaxOpen),
		model.ConnAttempts(cfg.ConnAttempts))
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	return &UserStorage{db}, nil
}
