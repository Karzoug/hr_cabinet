package handlers

type Config struct {
	Host  string `env:"HOST" env-default:"localhost"` // not used
	Port  int    `env:"PORT" env-default:"9990" env-required:"true"`
	Token struct {
		Lifetime  int    `env:"LIFETIME" env-default:"43200"`
		SecretKey string `env:"SECRET_KEY" env-required:"true"`
	} `env-prefix:"TOKEN_"`
}
