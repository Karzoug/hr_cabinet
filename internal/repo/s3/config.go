package s3

import "time"

type Config struct {
	Host            string        `env:"HOST" env-default:"localhost"`
	Port            int           `env:"PORT,notEmpty" env-default:"9000"`
	AccessKeyID     string        `env:"ACCESS_KEY_ID" env-required:"true"`
	SecretAccessKey string        `env:"SECRET_ACCESS_KEY" env-required:"true"`
	UseSSL          bool          `env:"USE_SSL" env-default:"false"`
	URLExpires      time.Duration `env:"URL_EXPIRES" env-default:"12h"`
	// Location        string `env:"LOCATION"`
}
