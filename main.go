package main

import (
	"log"
	"os"
	"path/filepath"
	"share_bot/bot"
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

	dsp := bot.NewDispatcher(token, storage)
	log.Println(dsp.Poll())

}
