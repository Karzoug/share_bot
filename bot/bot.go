package bot

import (
	"fmt"
	"share_bot/storage"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

var keyboard = echotron.ReplyKeyboardMarkup{
	ResizeKeyboard: true,
	Keyboard: [][]echotron.KeyboardButton{
		{
			{Text: whoMeButtonText},
			{Text: toWhomButtonText},
		},
	},
}

type bot struct {
	chatID int64
	echotron.API
	username string
	storage  storage.Storage
}

func NewDispatcher(token string, username string, storage storage.Storage) *echotron.Dispatcher {
	newBotFn := func(chatID int64) echotron.Bot {
		return &bot{
			chatID,
			echotron.NewAPI(token),
			username,
			storage,
		}
	}
	return echotron.NewDispatcher(token, newBotFn)
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message != nil {
		if update.Message.Text == "/start" && update.Message.Chat.Type == "private" {
			b.start(update)
			return
		}

		switch {
		case strings.HasPrefix(update.Message.Text, "/add"):
			b.add(update)
		case strings.HasPrefix(update.Message.Text, whoMeButtonText):
			b.whoMe(update)
		case strings.HasPrefix(update.Message.Text, toWhomButtonText):
			b.toWhom(update)
		}

		return
	}
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, "approve_in_request:") {
			b.approveRequest(update)
			return
		}
		if strings.HasPrefix(update.CallbackQuery.Data, "return_expense:") {
			b.returnExpense(update)
			return
		}

		if strings.HasPrefix(update.CallbackQuery.Data, "approve_return_expense:") {
			b.approveReturnExpense(update)
			return
		}
	}
}

func (b *bot) start(update *echotron.Update) {
	user := storage.User{
		Username: update.Message.From.Username,
		ChatId:   update.Message.From.ID,
	}
	b.storage.SaveUser(user)
	name := update.Message.From.FirstName
	if name == "" {
		name = user.Username
	}
	b.SendMessage(fmt.Sprintf(helloMsg, name), b.chatID, &echotron.MessageOptions{
		// ReplyMarkup: echotron.ReplyKeyboardRemove{
		// 	RemoveKeyboard: true,
		// 	Selective:      false,
		// },
		ReplyMarkup: keyboard,
	})
}
