package app

import "github.com/Employee-s-file-cabinet/backend/internal/config"

type App struct{}

func New(cfg *config.Config) (*App, error) {
	return &App{}, nil

	// Инициализация БД
	// Инициализация сервера/хэндлеров
}

func (a *App) Run() error {
	return nil
}
