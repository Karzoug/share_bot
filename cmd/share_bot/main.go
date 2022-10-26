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
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

var (
	configPath string
	debugLevel string
)

const (
	dbPath            string        = "data/db/share_bot.db"
	reminderFrequency time.Duration = time.Hour
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.yml", "path to config file")
	flag.StringVar(&debugLevel, "debug-level", "debug", "a debug level is a logging priority, there are debug, info, warn, error, dpanic, panic, fatal levels")
}

func main() {
	flag.Parse()
	var cfg config.Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()
	logger.InitializeLogger(debugLevel)

	storage, close := db.New(dbPath)
	defer close()

	worker := scheduler.NewScheduler()
	logger.Logger.Info("adding new worker to scheduler: reminder of debt")
	reminder := remind.New(cfg.Token, storage, cfg.Reminder)
	worker.Add(ctx, func(ctx context.Context) {
		reminder.Work(ctx)
	}, reminderFrequency)
	defer worker.Stop()

	for {
		dsp := bot.NewDispatcher(cfg.Token, storage)
		logger.Logger.Error("dispatcher poll error", zap.Error(dsp.Poll()))
		time.Sleep(time.Duration(cfg.RestartDurationSeconds))
	}
}
