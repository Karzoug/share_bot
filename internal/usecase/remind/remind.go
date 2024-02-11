package remind

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	apiModel "github.com/Karzoug/share_bot/internal/api/model"
)

func (s Service) remind(ctx context.Context) {
	const op = "remind service: remind"

	const (
		remindDebtMsg                       = "Напоминаю:\nвы должны @%s %d ₽ за «%s»"
		remindDebtNotConfirmedMsg    string = "Напоминаю:\nвам должен(-на) @%s %d ₽ за «%s».\n\nСумма не была подтверждена должником. Я не напоминаю ему о возврате."
		debtReturnedRequestButtonMsg        = "Я вернул этот долг"
		confirmReturnDebtButtonMsg          = "Деньги получил"
	)

	debts, err := s.remindRepo.GetDebts(ctx, nil)
	if err != nil {
		s.logger.Info("error",
			zap.Error(fmt.Errorf("%s: %w", op, err)))
		return
	}

	for i := range debts {
		if diff := time.Since(debts[i].Date); diff < s.cfg.InitDelay {
			continue
		} else if (diff.Microseconds() % s.cfg.Frequency.Microseconds()) > s.cfg.RunFrequency.Microseconds() {
			continue
		}

		if debts[i].Confirmed {
			author, err := s.userRepo.GetByID(ctx, debts[i].AuthorID)
			if err != nil {
				s.logger.Info("error",
					zap.Error(fmt.Errorf("%s: %w", op, err)))
				continue
			}
			debtor, err := s.userRepo.GetByUsername(ctx, debts[i].DebtorUsername)
			if err != nil {
				s.logger.Info("error",
					zap.Error(fmt.Errorf("%s: %w", op, err)))
				continue
			}

			if _, err := s.api.SendMessage(ctx, s.t.Getf(remindDebtMsg,
				author.Username,
				debts[i].Sum,
				debts[i].Comment),
				debtor.ID, apiModel.InlineKeyboardMarkup{
					InlineKeyboard: [][]apiModel.InlineKeyboardButton{{{
						Text: s.t.Get(debtReturnedRequestButtonMsg),
						CallbackData: fmt.Sprintf("debt_returned_request:%s:%s",
							debts[i].RequestID,
							debts[i].DebtorUsername),
					}}},
				}); err != nil {
				s.logger.Error("error", zap.String("error message", err.Error()))
			}
		} else {
			if _, err := s.api.SendMessage(ctx, s.t.Getf(remindDebtNotConfirmedMsg,
				debts[i].DebtorUsername,
				debts[i].Sum,
				debts[i].Comment),
				debts[i].AuthorID, apiModel.InlineKeyboardMarkup{
					InlineKeyboard: [][]apiModel.InlineKeyboardButton{{{
						Text: s.t.Get(confirmReturnDebtButtonMsg),
						CallbackData: fmt.Sprintf("confirm_return_debt:%s:%s",
							debts[i].RequestID,
							debts[i].DebtorUsername),
					}}},
				}); err != nil {
				s.logger.Error("error", zap.String("error message", err.Error()))
			}
		}
	}
}
