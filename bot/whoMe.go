package bot

import (
	"errors"
	"fmt"
	"log"
	"share_bot/lib/e"
	"share_bot/storage"

	"github.com/NicoNex/echotron/v3"
)

func (b *bot) whoMe(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			log.Println(e.Wrap("can't do who me command", err))
		}
	}()

	if update.Message.Chat.Type != "private" {
		err = errors.New("call in not private chat")
		return
	}

	req, err := b.storage.GetRequestsByLender(update.Message.From.Username, true)
	if err != nil {
		if err == storage.ErrUserNotExist {
			b.SendMessage(somethingWrongTryToStartMsg, b.chatID, nil)
		}
		return
	}
	if len(req) == 0 {
		b.SendMessage(whoMeNoExpensesMsg, b.chatID, nil)
		return
	}

	for _, r := range req {
		for _, ex := range r.Exps {
			kb := echotron.InlineKeyboardMarkup{
				InlineKeyboard: [][]echotron.InlineKeyboardButton{
					{
						{
							Text:         approveReturnExpenseButtonMsg,
							CallbackData: fmt.Sprintf("approve_return_expense:%d", ex.Id),
						},
					}},
			}
			msg := fmt.Sprintf("@%s \n%s: %d ₽ %s \n", ex.Borrower, r.Date.Format("02.01.06"), ex.Sum, r.Comment)

			if !ex.Approved {
				msg += "\n" + notApprovedMsg
			}

			_, err = b.SendMessage(msg, b.chatID, &echotron.MessageOptions{ReplyMarkup: kb})
			if err != nil {
				return
			}
		}
	}
}
