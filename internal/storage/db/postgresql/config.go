package postgresql

type Config struct {
	MaxOpenConns int    `env:"MAX_OPEN_CONNS" env-default:"4"`
	MaxIdleConns int    `env:"MAX_IDLE_CONNS" env-default:"4"`
	ConnAttempts int    `env:"CONN_ATTEMPTS" env-default:"10"`
	DSN          string `env-required:"true"  env:"DSN"`
}
