package s3

type Config struct {
	Host            string `env:"HOST" env-default:"localhost"`
	Port            int    `env:"PORT,notEmpty" env-default:"9000"`
	AccessKeyID     string `env:"ACCESS_KEY_ID" env-required:"true"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY" env-required:"true"`
	UseSSL          bool   `env:"USE_SSL" env-default:"false"`
	//Location        string `env:"LOCATION"`
}
