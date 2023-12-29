package postgresql

type Config struct {
	MaxOpenConns int    `env:"MAX_OPEN_CONNS" env-default:"4"`
	ConnAttempts int    `env:"CONN_ATTEMPTS" env-default:"10"`
	DSN          string `env:"DSN" env-required:"true"`
}
