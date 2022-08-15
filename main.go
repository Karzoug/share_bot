package main

import (
	"log"
	"os"
	"share_bot/bot"
	"share_bot/remind"
	"share_bot/storage/db"
)

const (
	dbPath                string = "data/db/share_bot.db"
	tokenEnv              string = "SHARE_BOT_TELEGRAM_TOKEN"
	waitInDayBeforeRemind int    = 3
	runReminderHour       int    = 18
)

func main() {
	token, exists := os.LookupEnv(tokenEnv)
	if !exists {
		log.Fatal("telegram token does not exist")
	}
	storage, close := db.New(dbPath)
	defer close()
	remainder := remind.NewReminder(token, storage, waitInDayBeforeRemind)
	stop := remainder.Start(runReminderHour)
	defer stop()

	dsp := bot.NewDispatcher(token, storage)
	log.Println(dsp.Poll())
}
