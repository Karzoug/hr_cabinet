package server

type Config struct {
	Host string `env:"HOST" env-default:"localhost"` // not used
	Port int    `env-required:"true" env:"PORT" env-default:"9990"`
}
