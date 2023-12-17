// TODO: Структуры с тегами переменных окружения
package config

import (
	"github.com/Employee-s-file-cabinet/backend/internal/server"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/db/postgresql"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string            `env:"LOG_LEVEL" default:"debug"`
	HTTP     server.Config     `env-prefix:"HTTP_"`
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
