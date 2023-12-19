package config

import (
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/Employee-s-file-cabinet/backend/internal/server/handlers"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/db/postgresql"
)

type Config struct {
	LogLevel string            `env:"LOG_LEVEL" env-default:"debug"`
	HTTP     handlers.Config   `env-prefix:"HTTP_"`
	PG       postgresql.Config `env-prefix:"PG_"`
}

// New создаёт объект Config.
func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
