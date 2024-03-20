package debt

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Karzoug/share_bot/internal/usecase"
)

func (s Service) Confirm(ctx context.Context, reqID, debtorUsername string) (string, error) {
	const op = "debt service: confirm"

	const (
		confirmDebtMsg    = "Вы подтвердили долг. Спасибо!"
		notFoundDebtMsg   = "Похоже, что это не ваш долг 😉"
		needToRegisterMsg = "Пожалуйста, сначала зарегистрируйтесь в боте 🙏"
	)

	_, err := s.userRepo.GetByUsername(ctx, debtorUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", usecase.NewError(s.t.Get(needToRegisterMsg))
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	debt, err := s.debtRepo.Get(ctx, reqID, debtorUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", usecase.NewError(s.t.Get(notFoundDebtMsg))
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if debt.Confirmed {
		return "", nil
	}
	debt.Confirmed = true
	if err := s.debtRepo.Save(ctx, debt); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return s.t.Get(confirmDebtMsg), nil

}
