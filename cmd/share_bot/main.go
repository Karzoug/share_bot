package main

import (
	"context"
	"flag"
	"log"
	"share_bot/internal/bot"
	"share_bot/internal/config"
	"share_bot/internal/logger"
	"share_bot/internal/remind"
	"share_bot/internal/storage/db"
	"share_bot/pkg/scheduler"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

var (
	configPath string
)

const (
	webhookMode string = "webhook"
	poolMode    string = "pool"
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.yml", "path to config file")
}

func main() {
	flag.Parse()
	var cfg config.Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()
	logger.InitializeLogger(cfg.DebugLevel)

	storage, close := db.New(cfg.DbPath)
	defer close()

	worker := scheduler.NewScheduler()
	logger.Logger.Info("adding new worker to scheduler: reminder of debt")
	reminder := remind.New(cfg.Token, storage, cfg.Reminder)
	worker.Add(ctx, func(ctx context.Context) {
		reminder.Work(ctx)
	}, cfg.ReminderFrequency)
	defer worker.Stop()

	dsp := bot.NewDispatcher(cfg.Token, storage)
	if cfg.Mode == webhookMode {
		logger.Logger.Error("webhook listening error", zap.Error(dsp.ListenWebhook(cfg.ServerUrl, cfg.SSLCertPath)))
	} else {
		logger.Logger.Error("long poll listening error", zap.Error(dsp.Poll()))
	}
}
