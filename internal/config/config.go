package config

type Config struct {
	Reminder               `yaml:"reminder"`
	Token                  string `env:"SHARE_BOT_TELEGRAM_TOKEN"`
	RestartDurationSeconds int    `env:"SHARE_BOT_RESTART_DURATION_SECONDS" yaml:"restartDurationSeconds" env-default:"5"`
}

type Reminder struct {
	WaitInDays int `env:"SHARE_BOT_WAIT_IN_DAYS" yaml:"waitInDays" env-default:"3"`
	RunHour    int `env:"SHARE_BOT_RUN_HOUR" yaml:"runHour" env-default:"18"`
}
