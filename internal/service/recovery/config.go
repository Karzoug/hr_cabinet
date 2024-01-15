package recovery

import "time"

type Config struct {
	Domain           string        `env:"DOMAIN" env-required:"true"`
	CleanKeyInterval time.Duration `env:"CLEAN_KEY_INTERVAL" env-default:"10m"`
	KeyLifetime      time.Duration `env:"KEY_LIFETIME" env-default:"30m"`
}
