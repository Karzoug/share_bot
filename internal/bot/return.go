package bot

import (
	"fmt"
	"share_bot/internal/logger"
	"share_bot/internal/storage"
	"strconv"
	"strings"

	"github.com/NicoNex/echotron/v3"
	"go.uber.org/zap"
)

func (b *bot) returnExpense(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			logger.Logger.Error("can't do command return expense", zap.Error(err))
			b.API.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
				Text: somethingWrongMsg,
			})
		}
	}()

	s := strings.TrimPrefix(update.CallbackQuery.Data, "return_expense:")
	expId, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	req, err := b.storage.GetExpenseWithRequest(expId)
	if err != nil {
		return
	}
	borr, exist, err := b.storage.GetUserByUsername(req.Lender)
	if err != nil {
		return
	}
	if !exist {
		err = storage.ErrUserNotExist
		return
	}

	_, err = b.API.EditMessageReplyMarkup(
		echotron.NewMessageID(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.ID),
		&echotron.MessageReplyMarkup{})
	if err != nil {
		return
	}

	kb := echotron.InlineKeyboardMarkup{
		InlineKeyboard: [][]echotron.InlineKeyboardButton{
			{
				{
					Text:         approveReturnExpenseButtonMsg,
					CallbackData: fmt.Sprintf("approve_return_expense:%d", req.Exps[0].Id),
				},
			}},
	}
	msg := fmt.Sprintf(returnMsg, req.Exps[0].Borrower, req.Comment, req.Exps[0].Sum)

	_, err = b.SendMessage(msg, borr.ChatId, &echotron.MessageOptions{ReplyMarkup: kb})
	if err != nil {
		return
	}

	b.API.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
		Text: sendReturnExpenseMsg,
	})
}

func (b *bot) approveReturnExpense(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			logger.Logger.Error("can't do command approve return expense", zap.Error(err))
			b.API.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
				Text: somethingWrongMsg,
			})
		}
	}()

	s := strings.TrimPrefix(update.CallbackQuery.Data, "approve_return_expense:")
	expId, err := strconv.Atoi(s)
	if err != nil {
		return
	}

	err = b.storage.ApproveReturnExpense(int64(expId), update.CallbackQuery.From.Username)
	if err != nil {
		return
	}
	_, err = b.API.EditMessageReplyMarkup(
		echotron.NewMessageID(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.ID),
		&echotron.MessageReplyMarkup{})
	if err != nil {
		return
	}
	_, err = b.API.AnswerCallbackQuery(update.CallbackQuery.ID, &echotron.CallbackQueryOptions{
		Text: thanksMsg,
	})
	if err != nil {
		return
	}
	_, err = b.SendMessage(approveReturnExpenseMsg, update.CallbackQuery.Message.Chat.ID, nil)
	if err != nil {
		return
	}
}
