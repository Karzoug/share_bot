package bot

import (
	"errors"
	"fmt"
	"log"
	"share_bot/lib/e"
	"share_bot/storage"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

func (b *bot) toWhomCommand(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			log.Println(e.Wrap("can't do toWhom command", err))
		}
	}()

	if update.Message.Chat.Type != "private" {
		err = errors.New("call in not private chat")
		return
	}

	req, err := b.storage.GetRequestsByBorrower(update.Message.From.Username, true)
	if err != nil {
		if err == storage.ErrUserNotExist {
			b.SendMessage(somethingWrongTryToStartMsg, b.chatID, nil)
		}
		return
	}
	if len(req) == 0 {
		b.SendMessage(toWhomNoExpenses, b.chatID, nil)
		return
	}
	var bld strings.Builder
	for _, r := range req {
		for _, e := range r.Exps {
			fmt.Fprintf(&bld, "@%s \n%s: %d ₽ %s \n", e.Person, r.Date.Format("02.01.06"), e.Sum, r.Comment)
		}
	}

	_, err = b.SendMessage(bld.String(), b.chatID, nil)
	if err != nil {
		log.Println(e.Wrap("can't do toWhom command", err))
	}
}
