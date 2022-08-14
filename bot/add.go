package bot

import (
	"fmt"
	"log"
	"share_bot/lib/e"
	"share_bot/parse"
	"share_bot/storage"
	"strconv"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
)

func (b *bot) add(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			log.Println(e.Wrap("can't do add command", err))
		}
	}()

	botName, err := b.getUsername()
	if err != nil {
		return
	}
	message := strings.TrimPrefix(update.Message.Text, "/add@"+botName)
	message = strings.TrimPrefix(message, "/add")

	resDel, _ := b.DeleteMessage(b.chatID, update.Message.ID)
	if update.Message.Chat.Type != "private" {
		if !resDel.Ok {
			b.SendMessage(needDeletePermissionMsg, b.chatID, nil)
			return
		}

		u := update.Message.From.FirstName
		if u == "" {
			u = update.Message.From.Username
		}

		kb := echotron.InlineKeyboardMarkup{
			InlineKeyboard: [][]echotron.InlineKeyboardButton{
				{
					{
						Text: registerBotButtonMsg,
						URL:  "https://t.me/" + botName,
					},
				}},
		}

		if !b.storage.IsUserExist(update.Message.From.Username) {
			b.SendMessage(fmt.Sprintf(needToRegisterMsg, u), b.chatID, &echotron.MessageOptions{ReplyMarkup: kb})
			return
		}
	}

	exps, comment, err := parse.AddMessage(message)
	if err != nil {
		return
	}
	if len(exps) == 0 {
		return
	}

	req := storage.Request{
		Lender:  update.Message.From.Username,
		Exps:    exps,
		Comment: comment,
		Date:    time.Unix(int64(update.Message.Date), 0),
		ChatId:  update.Message.Chat.ID,
	}

	err = b.storage.AddRequest(&req)
	if err != nil {
		return
	}

	var bld strings.Builder
	fmt.Fprintf(&bld, addMsg+"\n\n", req.Lender, req.Comment)
	for _, e := range req.Exps {
		fmt.Fprintf(&bld, "🧾 @%s: %d ₽ \n", e.Borrower, e.Sum)
	}

	var inlineKeyboard = echotron.InlineKeyboardMarkup{
		InlineKeyboard: [][]echotron.InlineKeyboardButton{
			{
				{
					Text:         approveButtonMsg,
					CallbackData: fmt.Sprintf("approve_in_request:%d", req.Id),
				},
			},
		},
	}

	_, err = b.SendMessage(bld.String(), b.chatID, &echotron.MessageOptions{ReplyMarkup: inlineKeyboard})
	if err != nil {
		return
	}
}

func (b *bot) approveRequest(update *echotron.Update) {
	var err error
	defer func() {
		err = e.Wrap("can't do command approve request", err)
		if err != nil {
			log.Println(err)
			b.API.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
				Text: somethingWrongMsg,
			})
		}
	}()
	s := strings.TrimPrefix(update.CallbackQuery.Data, "approve_in_request:")
	reqId, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	err = b.storage.ApproveExpense(int64(reqId), update.CallbackQuery.From.Username)
	if err != nil {
		return
	}
	b.API.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
		Text: approveExpenseMsg,
	})
}
