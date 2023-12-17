package server

type Config struct {
	Host  string `env:"HOST" env-default:"localhost"` // not used
	Port  int    `env-required:"true" env:"PORT" env-default:"9990"`
	Token struct {
		Lifetime  int    `env:"LIFETIME" env-default:"43200"`
		SecretKey string `env-required:"true" env:"SECRET_KEY"`
	} `env-prefix:"TOKEN_"`
}
