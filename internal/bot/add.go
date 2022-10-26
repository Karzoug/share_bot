package bot

import (
	"fmt"
	"share_bot/internal/logger"
	"share_bot/internal/parse"
	"share_bot/internal/storage"
	"strconv"
	"strings"
	"time"

	"github.com/NicoNex/echotron/v3"
	"go.uber.org/zap"
)

func (b *bot) add(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			logger.Logger.Error("can't do add command", zap.Error(err))
			b.API.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
				Text: somethingWrongMsg,
			})
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
		var exist bool
		exist, err = b.storage.IsUserExist(update.Message.From.Username)
		if err != nil {
			return
		}
		if !exist {
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
			b.SendMessage(fmt.Sprintf(needToRegisterMsg, u), b.chatID, &echotron.MessageOptions{ReplyMarkup: kb})
			return
		}
	}

	exps, comment, err := parse.AddMessage(message)
	if err != nil {
		return
	}
	if len(exps) == 0 {
		b.API.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
			Text: wrongAddMsg,
		})
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

	if update.Message.Chat.Type == "private" {
		for _, e := range req.Exps {
			var bld strings.Builder
			if user, exist, err := b.storage.GetUserByUsername(e.Borrower); exist && user.ChatId != 0 && err == nil {
				fmt.Fprintf(&bld, addMsg+"\n\n", req.Lender, req.Comment)
				fmt.Fprintf(&bld, "🧾 @%s: %d ₽ \n", e.Borrower, e.Sum)
				_, err = b.SendMessage(bld.String(), user.ChatId, &echotron.MessageOptions{ReplyMarkup: inlineKeyboard})
				if err != nil {
					return
				}
				_, err = b.SendMessage(fmt.Sprintf(sendAproveBorrowerFromPrivateChatMsg, user.Username), b.chatID, nil)
			} else if err != nil {
				fmt.Fprintf(&bld, "🧾 @%s: %d ₽ «%s»", e.Borrower, e.Sum, req.Comment)
				fmt.Fprint(&bld, "\n\n"+mentionedUserNotRegistered)
				_, err = b.SendMessage(bld.String(), b.chatID, nil)
			} else {
				return
			}
		}
	} else {
		var bld strings.Builder
		fmt.Fprintf(&bld, addMsg+"\n\n", req.Lender, req.Comment)
		for _, e := range req.Exps {
			fmt.Fprintf(&bld, "🧾 @%s: %d ₽ \n", e.Borrower, e.Sum)
		}

		_, err = b.SendMessage(bld.String(), b.chatID, &echotron.MessageOptions{ReplyMarkup: inlineKeyboard})
		if err != nil {
			return
		}
	}
}

func (b *bot) approveRequest(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			logger.Logger.Error("can't do command approve request", zap.Error(err))
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
	_, err = b.API.EditMessageReplyMarkup(
		echotron.NewMessageID(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.ID),
		&echotron.MessageReplyMarkup{})
	if err != nil {
		return
	}
}
