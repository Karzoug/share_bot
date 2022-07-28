package bot

import (
	"fmt"
	"log"
	"share_bot/db"
	"share_bot/parse"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
)

func addCommand(b *bot, update *echotron.Update) {
	message := strings.TrimPrefix(update.Message.Text, "/add@"+usernameBot)
	message = strings.TrimPrefix(message, "/add")
	if update.Message.Chat.Type != "private" {
		resDel, _ := b.DeleteMessage(b.chatID, update.Message.ID)
		if !resDel.Ok {
			b.SendMessage("The bot should be able to delete user messages", b.chatID, nil)
			return
		}
	}

	resParse, comment, err := parse.AddMessage(message)
	if err != nil {
		log.Println(fmt.Errorf("add Command parse message error: %w", err))
		return
	}
	lender := db.User{
		Username:   update.Message.From.Username,
		TelegramId: update.Message.From.ID,
	}
	req := db.Request{
		Comment: comment,
		ChatId:  update.Message.Chat.ID,
		Date:    time.Unix(int64(update.Message.Date), 0),
	}
	db.AddExpenses(lender, req, resParse)

	var bld strings.Builder
	fmt.Fprintf(&bld, "%s сообщил о тратах «%s»\n\n", lender.Username, req.Comment)
	for _, v := range resParse {
		fmt.Fprintf(&bld, "%c @%s: %d %c \n", receiptCode, v.Borrower, v.Sum, rubleCode)
	}

	_, err = b.SendMessage(bld.String(), b.chatID, nil)
	if err != nil {
		log.Println(fmt.Errorf("add Command send message error: %w", err))
	}
}
