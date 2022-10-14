package config

type Config struct {
	Database struct {
		Path     string `yaml:"path" env-default:"data/db/"`
		Filename string `yaml:"filename" env-default:"share_bot.db"`
	} `yaml:"database"`
	Log struct {
		Path     string `yaml:"path" env-default:"./"`
		Filename string `yaml:"filename" env-default:"log.txt"`
	} `yaml:"log"`
	Reminder struct {
		WaitInDay int `yaml:"waitInDay" env-default:"3"`
		RunHour   int `yaml:"runHour" env-default:"18"`
	} `yaml:"reminder"`
	Token                  string `env:"SHARE_BOT_TELEGRAM_TOKEN"`
	RestartDurationSeconds int    `yaml:"restartDurationSeconds" env-default:"5"`
}
