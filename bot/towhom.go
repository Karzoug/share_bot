package bot

import (
	"fmt"
	"log"
	"share_bot/db"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

func toWhomCommand(b *bot, update *echotron.Update) {
	if update.Message.Chat.Type != "private" {
		resDel, err := b.DeleteMessage(b.chatID, update.Message.ID)
		if !resDel.Ok {
			b.SendMessage("The bot should be able to delete user messages", b.chatID, nil)
			return
		}
		if err != nil {
			log.Println(fmt.Errorf("toWhom Command delete message error: %w", err))
			return
		}
	}

	msgs, err := db.ShowExpensesByBorrower(update.Message.From.Username)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if len(msgs) == 0 {
		return
	}
	var bld strings.Builder
	for _, v := range msgs {
		fmt.Fprintf(&bld, "@%s \n%s: %d %c %s \n", v.Lender.Username, v.Request.Date.Format("02.01.06"), v.Sum, rubleCode, v.Request.Comment)
	}

	_, err = b.SendMessage(bld.String(), b.chatID, nil)
	if err != nil {
		log.Println(fmt.Errorf("toWhom Command send message error: %w", err))
	}
}
