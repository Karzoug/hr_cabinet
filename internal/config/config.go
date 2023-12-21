package config

import (
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/Employee-s-file-cabinet/backend/internal/server"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/db/postgresql"
	"github.com/Employee-s-file-cabinet/backend/internal/storage/s3"
	"github.com/Employee-s-file-cabinet/backend/internal/utils/email"
)

type Config struct {
	LogLevel string            `env:"LOG_LEVEL" env-default:"debug"`
	HTTP     server.Config     `env-prefix:"HTTP_"`
	PG       postgresql.Config `env-prefix:"PG_"`
	S3       s3.Config         `env-prefix:"S3_"`
	Mail     email.Config      `env-prefix:"MAIL_"`
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
