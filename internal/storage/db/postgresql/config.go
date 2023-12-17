package postgresql

type Config struct {
	MaxOpen      int    `env:"POOL_MAX"`
	ConnAttempts int    `env:"CONN_ATTEMPTS"`
	DSN          string `env-required:"true"  env:"DSN"`
}
