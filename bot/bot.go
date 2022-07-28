package bot

import (
	"fmt"
	"log"
	"share_bot/db"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

const rubleCode rune = '₽'
const receiptCode rune = '🧾'

const toWhomCommandText string = "Кому должен я?"
const whoMeCommandText string = "Кто должен мне?"

var keyboard = echotron.ReplyKeyboardMarkup{
	ResizeKeyboard: true,
	Keyboard: [][]echotron.KeyboardButton{
		{
			{Text: whoMeCommandText},
			{Text: toWhomCommandText},
		},
	},
}

type bot struct {
	chatID int64
	echotron.API
}

var (
	tokenBot    string
	usernameBot string
)

func Start(token string, username string) {
	tokenBot = token
	usernameBot = username
	dsp := echotron.NewDispatcher(tokenBot, newBot)
	log.Println(dsp.Poll())
}

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(tokenBot),
	}
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message.Text == "/start" && update.Message.Chat.Type == "private" {
		user := db.User{
			Username:   update.Message.From.Username,
			TelegramId: update.Message.From.ID,
		}
		db.InitUser(user)
		name := update.Message.From.FirstName
		if name == "" {
			name = user.Username
		}
		b.SendMessage(fmt.Sprintf(`Привет, %v!
Я бот, который поможет тебе и твоим знакомым не забыть об общих тратах друг друга.

* Чтобы добавить трату, отправь сообщение следующего вида:
/add @nickname_друга сумма комментарий_описывающий_трату
* Чтобы вернуть долг, отправь сообщение такого формата: 
/return @nickname_друга сумма
* Чтобы быстро узнать, кто и сколько должен воспользуйся кнопками снизу.

Начнем?`, name), b.chatID, &echotron.MessageOptions{
			// ReplyMarkup: echotron.ReplyKeyboardRemove{
			// 	RemoveKeyboard: true,
			// 	Selective:      false,
			// },
			ReplyMarkup: keyboard,
		})
	} else if update.Message != nil {

		switch {
		case strings.HasPrefix(update.Message.Text, "/add"):
			addCommand(b, update)
		case strings.HasPrefix(update.Message.Text, "/return"):
			//returnCommand(b, update)
		case strings.HasPrefix(update.Message.Text, whoMeCommandText):
			//whoMeCommand(b, update)
		case strings.HasPrefix(update.Message.Text, toWhomCommandText):
			toWhomCommand(b, update)
		default:
			return
		}
	}
}
