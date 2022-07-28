package main

import (
	"log"
	"os"
	"path/filepath"
	"share_bot/bot"
	"share_bot/db"

	"github.com/mitchellh/go-homedir"
)

const dbFileName string = "share_bot.db"
const tokenEnv string = "SHARE_BOT_TELEGRAM_TOKEN"
const botUsernameEnv string = "SHARE_BOT_USERNAME"

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, dbFileName)
	err, close := db.Init(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	token, exists := os.LookupEnv(tokenEnv)
	if !exists {
		log.Fatal("Telegram token does not exist")
	}

	username := os.Getenv(botUsernameEnv)

	bot.Start(token, username)
}
