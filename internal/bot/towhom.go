package bot

import (
	"errors"
	"fmt"
	"share_bot/internal/storage"
	"share_bot/pkg/e"

	"github.com/NicoNex/echotron/v3"
)

func (b *bot) toWhom(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			b.logger.Println(e.Wrap("can't do to whom command", err))
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
		b.SendMessage(toWhomNoExpensesMsg, b.chatID, nil)
		return
	}

	for _, r := range req {
		for _, ex := range r.Exps {
			kb := echotron.InlineKeyboardMarkup{
				InlineKeyboard: [][]echotron.InlineKeyboardButton{
					{
						{
							Text:         returnExpenseButtonMsg,
							CallbackData: fmt.Sprintf("return_expense:%d", ex.Id),
						},
					}},
			}
			msg := fmt.Sprintf("@%s \n%s: %d ₽ %s \n", r.Lender, r.Date.Format("02.01.06"), ex.Sum, r.Comment)

			_, err = b.SendMessage(msg, b.chatID, &echotron.MessageOptions{ReplyMarkup: kb})
			if err != nil {
				return
			}
		}
	}

}
