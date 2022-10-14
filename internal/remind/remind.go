package remind

import (
	"fmt"
	"log"
	"share_bot/internal/storage"
	"share_bot/pkg/e"
	"strconv"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/robfig/cron"
)

type Reminder struct {
	api               echotron.API
	storage           storage.Storage
	logger            *log.Logger
	waitDurationInDay int
}

func NewReminder(token string, storage storage.Storage, logger *log.Logger, waitInDay int) *Reminder {
	if waitInDay < 1 {
		waitInDay = 1
	}
	return &Reminder{
		echotron.NewAPI(token),
		storage,
		logger,
		waitInDay,
	}
}

func (r *Reminder) Start(runHour int) (stop func()) {
	c := cron.New()
	c.AddFunc("0 0 "+strconv.Itoa(runHour)+" * * *", func() {
		r.work()
	})
	c.Start()
	return func() { c.Stop() }
}

func (r *Reminder) work() {
	reqs, err := r.storage.GetNotReturnedRequests()
	if err != nil {
		err = e.Wrap("can't work remainder", err)
		r.logger.Println(err)
		return
	}

	for _, req := range reqs {
		if diff := time.Now().Sub(req.Date); int(diff.Hours())%(24*r.waitDurationInDay) > 24 {
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
				borr, exist := r.storage.GetUserByUsername(exp.Borrower)
				if !exist {
					continue
				}
				r.api.SendMessage(
					fmt.Sprintf(remindToBorrower, req.Lender, exp.Sum, req.Comment),
					borr.ChatId,
					&echotron.MessageOptions{ReplyMarkup: kb},
				)
			} else {
				lend, exist := r.storage.GetUserByUsername(req.Lender)
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
