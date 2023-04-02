package config

import "time"

type Config struct {
	Reminder          `yaml:"reminder"`
	Token             string        `env:"SHARE_BOT_TELEGRAM_TOKEN"`
	DbPath            string        `env:"SHARE_BOT_DB_PATH" yaml:"dbPath" env-default:"data/db/share_bot.db"`
	ReminderFrequency time.Duration `env:"SHARE_BOT_REMINDER_FREQUENCY" yaml:"reminderFrequency" env-default:"1h"`
	SSLCertPath       string        `env:"SHARE_BOT_SSL_CERT_PATH" yaml:"sslCertPath"`
	DebugLevel        string        `env:"SHARE_BOT_DEBUG_LEVEL" yaml:"debugLevel" env-default:"debug"`
	Mode              string        `env:"SHARE_BOT_MODE" yaml:"mode" env-default:"poll"`
	ServerUrl         string        `env:"SHARE_BOT_SERVER_URL" yaml:"serverUrl"`
}

type Reminder struct {
	WaitInDays int `env:"SHARE_BOT_WAIT_IN_DAYS" yaml:"waitInDays" env-default:"3"`
	RunHour    int `env:"SHARE_BOT_RUN_HOUR" yaml:"runHour" env-default:"18"`
}
