package main

import (
	"log"
	"os"
	"path/filepath"
	"share_bot/bot"
	"share_bot/remind"
	"share_bot/storage/db"

	"github.com/mitchellh/go-homedir"
)

const dbFileName string = "share_bot.db"
const tokenEnv string = "SHARE_BOT_TELEGRAM_TOKEN"

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, dbFileName)

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
