package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Karzoug/share_bot/internal/app"
	"github.com/Karzoug/share_bot/internal/config"
)

func main() {
	cfg := &config.Config{}
	if err := env.ParseWithOptions(cfg,
		env.Options{Prefix: "SHARE_BOT_"}); err != nil {
		log.Fatal(err)
	}

	logger, err := buildLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("starting app")

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if err := app.Run(ctx, cfg, logger); err != nil {
		logger.Error("app stopped with error", zap.Error(err))
	}

}

func buildLogger(level zapcore.Level) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableCaller = true
	zapConfig.Level.SetLevel(level)
	return zapConfig.Build()
}
