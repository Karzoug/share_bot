package main

import (
	"log"
	"os"
	"share_bot/bot"
	"share_bot/remind"
	"share_bot/storage/db"
)

const dbPath string = "data/db/share_bot.db"
const tokenEnv string = "SHARE_BOT_TELEGRAM_TOKEN"

func main() {
	token, exists := os.LookupEnv(tokenEnv)
	if !exists {
		log.Fatal("telegram token does not exist")
	}
	storage, close := db.New(dbPath)
	defer close()

	remainder := remind.NewReminder(token, storage, 3)
	stop := remainder.Start(18)
	defer stop()

	dsp := bot.NewDispatcher(token, storage)
	log.Println(dsp.Poll())
}
