// TODO: Структуры с тегами переменных окружения
package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	LogLevel string `env:"LOG_LEVEL" default:"debug"`
	HTTP     HTTP   `env-prefix:"HTTP_"`
	PG       PG     `env-prefix:"PG_"`
}

type HTTP struct {
	Host string `env:"HOST" default:"localhost"`
	Port int    `env-required:"true" env:"HTTP_PORT" default:"9990"`
}

type PG struct {
	MaxOpen      int    `env:"POOL_MAX"`
	ConnAttempts int    `env:"CONN_ATTEMPTS"`
	DSN          string `env-required:"true"  env:"DSN"`
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
