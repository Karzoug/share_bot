package debt

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	apiModel "github.com/Karzoug/share_bot/internal/api/model"
	"github.com/Karzoug/share_bot/internal/usecase"
)

func (s Service) RequestReturn(ctx context.Context, reqID, debtorUsername string) (string, error) {
	const op = "debt service: request return"

	const (
		returnMsg                  = "@%s сообщил, что отдал вам долг за «%s» в размере %d ₽"
		requestReturnSuccessMsg    = "Спасибо! Проверяем ..."
		notFoundDebtMsg            = "Что-то пошло не так ... Не могу найти это долг!"
		confirmReturnDebtButtonMsg = "Деньги получил"
	)

	debt, err := s.debtRepo.Get(ctx, reqID, debtorUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", usecase.NewError(s.t.Get(notFoundDebtMsg))
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultBackgroundTaskTimeout)
		defer cancel()

		_, err := s.api.SendMessage(ctx,
			s.t.Getf(returnMsg, debtorUsername, debt.Comment, debt.Sum),
			debt.AuthorID,
			apiModel.InlineKeyboardMarkup{
				InlineKeyboard: [][]apiModel.InlineKeyboardButton{{{
					Text:         s.t.Get(confirmReturnDebtButtonMsg),
					CallbackData: fmt.Sprintf("confirm_return_debt:%s:%s", reqID, debtorUsername),
				}}},
			})
		if err != nil {
			s.logger.Error("error", zap.String("error message", err.Error()))
		}
	}()

	return s.t.Get(requestReturnSuccessMsg), nil

}
