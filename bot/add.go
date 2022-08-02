package bot

import (
	"fmt"
	"log"
	"share_bot/lib/e"
	"share_bot/parse"
	"share_bot/storage"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
)

func (b *bot) addCommand(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			log.Println(e.Wrap("can't do addCommand", err))
		}
	}()

	message := strings.TrimPrefix(update.Message.Text, "/add@"+b.username)
	message = strings.TrimPrefix(message, "/add")
	if update.Message.Chat.Type != "private" {
		resDel, _ := b.DeleteMessage(b.chatID, update.Message.ID)
		if !resDel.Ok {
			b.SendMessage(needDeletePermissionMsg, b.chatID, nil)
			return
		}
	}

	exps, comment, err := parse.AddMessage(message)
	if err != nil {
		return
	}

	req := storage.Request{
		Lender:  update.Message.From.Username,
		Exps:    exps,
		Comment: comment,
		Date:    time.Unix(int64(update.Message.Date), 0),
		ChatId:  update.Message.Chat.ID,
	}
	err = b.storage.AddRequest(req)
	if err != nil {
		return
	}

	var bld strings.Builder
	fmt.Fprintf(&bld, "%s сообщил о тратах «%s»\n\n", req.Lender, req.Comment)
	for _, e := range req.Exps {
		fmt.Fprintf(&bld, "🧾 @%s: %d ₽ \n", e.Person, e.Sum)
	}

	_, err = b.SendMessage(bld.String(), b.chatID, nil)
	if err != nil {
		return
	}
}
