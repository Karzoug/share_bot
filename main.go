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
const botUsernameEnv string = "SHARE_BOT_USERNAME"

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, dbFileName)

	token, exists := os.LookupEnv(tokenEnv)
	if !exists {
		log.Fatal("telegram token does not exist")
	}

	username := os.Getenv(botUsernameEnv)

	storage, close := db.New(dbPath)
	defer close()

	dsp := bot.NewDispatcher(token, username, storage)
	log.Println(dsp.Poll())

}
