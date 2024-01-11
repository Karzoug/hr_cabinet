package smtp

type Config struct {
	Name     string `env:"NAME" env-default:"Картотека сотрудника"`
	From     string `env:"FROM"`
	Login    string `env:"LOGIN"`
	Password string `env:"PASSWORD"`
	SMTPHost string `env:"SMTP_HOST" env-required:"true"`
	SMTPPort int    `env:"SMTP_PORT" env-required:"true"`
}
