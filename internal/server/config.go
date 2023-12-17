package server

type Config struct {
	Host string `env:"HOST" default:"localhost"`
	Port int    `env-required:"true" env:"HTTP_PORT" default:"9990"`
}
