package http

import "time"

type Config struct {
	Host  string `env:"HOST" env-default:"localhost"` // not used
	Port  int    `env:"PORT" env-default:"9990" env-required:"true"`
	Token struct {
		Lifetime  time.Duration `env:"LIFETIME" env-default:"72h"`
		SecretKey string        `env:"SECRET_KEY" env-required:"true" env-description:"private key to sign token, 64 bits size"`
	} `env-prefix:"TOKEN_"`
}
