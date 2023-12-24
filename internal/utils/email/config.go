package email

type Config struct {
	Name     string `env:"NAME" env-default:"Картотека сотрудника"`
	FromMail string `env:"FROM_MAIL"`
	Login    string `env:"LOGIN"`
	AppPass  string `env:"APP_PASS"`
	SMTPHost string `env:"SMTP_HOST"`
}
