// TODO: Структуры с тегами переменных окружения
package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		HTTP
		PG
		LogLevel string `env:"LOG_LEVEL" default:"debug"`
	}

	// HTTP - настройка http-подключения.
	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT" default:"9990"`
	}

	// PG - настройки подключения к БД Postgres.
	PG struct {
		MaxOpen      int    `env:"PG_POOL_MAX"`
		ConnAttempts int    `env:"PG_CONN_ATTEMPTS"`
		DSN          string `env-required:"true"  env:"PG_DSN"`
	}
)

// New создаёт объект Config.
func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
