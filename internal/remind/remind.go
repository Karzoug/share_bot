package remind

import (
	"context"
	"fmt"
	"share_bot/internal/config"
	"share_bot/internal/logger"
	"share_bot/internal/storage"
	"time"

	"github.com/NicoNex/echotron/v3"
	"go.uber.org/zap"
)

type Reminder struct {
	api     echotron.API
	storage storage.Storage
	config  config.Reminder
}

func New(token string, storage storage.Storage, cfg config.Reminder) *Reminder {
	if token == "" {
		logger.Logger.Fatal("telegram token does not exist")
	}
	if cfg.WaitInDays < 1 {
		cfg.WaitInDays = 1
	}
	return &Reminder{
		echotron.NewAPI(token),
		storage,
		cfg,
	}
}

func (r *Reminder) Work(ctx context.Context) {
	if time.Now().Hour() != r.config.RunHour {
		return
	}
	reqs, err := r.storage.GetNotReturnedRequests()
	if err != nil {
		logger.Logger.Info("can't work reminder: storage not returned requests", zap.Error(err))
		return
	}

	for _, req := range reqs {
		if diff := time.Since(req.Date); int(diff.Hours())%(24*r.config.WaitInDays) > 24 || int(diff.Hours())/(24*r.config.WaitInDays) == 0 {
			continue
		}
		for _, exp := range req.Exps {
			if exp.Approved {
				kb := echotron.InlineKeyboardMarkup{
					InlineKeyboard: [][]echotron.InlineKeyboardButton{
						{
							{
								Text:         returnExpenseButtonMsg,
								CallbackData: fmt.Sprintf("return_expense:%d", exp.Id),
							},
						}},
				}
				borr, exist, err := r.storage.GetUserByUsername(exp.Borrower)
				if err != nil {
					logger.Logger.Error("reminder: can't process request: storage error", zap.Error(err))
					continue
				}
				if !exist {
					continue
				}
				r.api.SendMessage(
					fmt.Sprintf(remindToBorrower, req.Lender, exp.Sum, req.Comment),
					borr.ChatId,
					&echotron.MessageOptions{ReplyMarkup: kb},
				)
			} else {
				lend, exist, err := r.storage.GetUserByUsername(req.Lender)
				if err != nil {
					logger.Logger.Error("reminder: can't process request: storage error", zap.Error(err))
					continue
				}
				if !exist {
					continue
				}
				r.api.SendMessage(
					fmt.Sprintf(remindToLender, exp.Borrower, exp.Sum, req.Comment),
					lend.ChatId,
					nil,
				)
			}

		}
	}
}
