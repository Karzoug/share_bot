package config

import (
	"go.uber.org/zap/zapcore"

	"github.com/Karzoug/share_bot/internal/api"
	"github.com/Karzoug/share_bot/internal/delivery/http"
	"github.com/Karzoug/share_bot/internal/usecase/remind"
)

type Config struct {
	LogLevel zapcore.Level `env:"LOG_LEVEL" envDefault:"INFO"`
	API      api.Config    `envPrefix:"API_"`
	HTTP     http.Config   `envPrefix:"HTTP_"`
	DB       DB            `envPrefix:"DB_"`
	Remind   remind.Config `envPrefix:"REMIND_"`
}

type DB struct {
	DSN string `env:"DSN,notEmpty"`
}
