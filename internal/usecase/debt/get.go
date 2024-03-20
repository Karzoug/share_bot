package debt

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	apiModel "github.com/Karzoug/share_bot/internal/api/model"
	"github.com/Karzoug/share_bot/internal/model"
)

func (s Service) GetUserDebts(ctx context.Context, username string, chat model.Chat) (string, error) {
	const op = "debt service: get user debts"

	const (
		noDebtsMsg                   = "Вы никому не должны 👍"
		debtReturnedRequestButtonMsg = "Я вернул этот долг"
	)

	if chat.Type != "private" {
		return "", fmt.Errorf("%s: %w", op, errors.New("called in not private chat"))
	}

	debts, err := s.debtRepo.GetByDebtorUsername(ctx, username, true)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	getDebtMsg := func(authorUsername string, debt model.Debt) string {
		msg := fmt.Sprintf("@%s \n%s: %d ₽ %s \n",
			authorUsername,
			debt.Date.Format("02.01.06"),
			debt.Sum,
			debt.Comment)
		return msg
	}

	if len(debts) == 0 {
		return s.t.Get(noDebtsMsg), nil
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultBackgroundTaskTimeout)
		defer cancel()

		for i := range debts {
			author, err := s.userRepo.GetByID(ctx, debts[i].AuthorID)
			if err != nil {
				s.logger.Error("error", zap.String("error message", err.Error()))
			}

			_, err = s.api.SendMessage(ctx, getDebtMsg(author.Username, debts[i]), chat.ID, apiModel.InlineKeyboardMarkup{
				InlineKeyboard: [][]apiModel.InlineKeyboardButton{{{
					Text:         s.t.Get(debtReturnedRequestButtonMsg),
					CallbackData: fmt.Sprintf("debt_returned_request:%s:%s", debts[i].RequestID, debts[i].DebtorUsername),
				}}},
			})
			if err != nil {
				s.logger.Error("error", zap.String("error message", err.Error()))
			}
		}
	}()

	return "", nil
}
