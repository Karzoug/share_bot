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
			{Text: whoMeCommandText},
			{Text: toWhomCommandText},
		},
	},
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

type bot struct {
	chatID int64
	echotron.API
	username string
	storage  storage.Storage
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message.Text == "/start" && update.Message.Chat.Type == "private" {
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
	} else if update.Message != nil {

		switch {
		case strings.HasPrefix(update.Message.Text, "/add"):
			b.addCommand(update)
		case strings.HasPrefix(update.Message.Text, "/return"):
			//returnCommand(b, update)
		case strings.HasPrefix(update.Message.Text, whoMeCommandText):
			//whoMeCommand(b, update)
		case strings.HasPrefix(update.Message.Text, toWhomCommandText):
			b.toWhomCommand(update)
		default:
			return
		}
	}
}
