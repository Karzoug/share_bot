package debt

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	apiModel "github.com/Karzoug/share_bot/internal/api/model"
	"github.com/Karzoug/share_bot/internal/model"
)

func (s Service) GetDebtsOwedToUser(ctx context.Context, userID int64, chat model.Chat) (string, error) {
	const op = "debt service: get debts owed to user"

	const (
		noDebtsMsg                 = "Вам никто не должен 😢"
		debtNotConfirmedMsg        = "Сумма не была подтверждена должником. Я не напоминаю ему о возврате."
		confirmReturnDebtButtonMsg = "Деньги получил"
	)

	if chat.Type != "private" {
		return "", fmt.Errorf("%s: %w", op, errors.New("called in not private chat"))
	}

	debts, err := s.debtRepo.GetByAuthorID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	getDebtMsg := func(debt model.Debt) string {
		msg := fmt.Sprintf("@%s \n%s: %d ₽ %s \n",
			debt.DebtorUsername,
			debt.Date.Format("02.01.06"),
			debt.Sum,
			debt.Comment)
		if !debt.Confirmed {
			msg += "\n" + s.t.Get(debtNotConfirmedMsg)
		}
		return msg
	}

	if len(debts) == 0 {
		return s.t.Get(noDebtsMsg), nil
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultBackgroundTaskTimeout)
		defer cancel()

		for i := range debts {
			_, err := s.api.SendMessage(ctx,
				getDebtMsg(model.Debt(debts[i])),
				chat.ID,
				apiModel.InlineKeyboardMarkup{
					InlineKeyboard: [][]apiModel.InlineKeyboardButton{{{
						Text: s.t.Get(confirmReturnDebtButtonMsg),
						CallbackData: fmt.Sprintf("confirm_return_debt:%s:%s",
							debts[i].RequestID,
							debts[i].DebtorUsername),
					}}},
				})
			if err != nil {
				s.logger.Error("error", zap.String("error message", err.Error()))
			}
		}
	}()

	return "", nil
}
