package debt

import (
	"context"
	"fmt"
)

func (s Service) ConfirmReturn(ctx context.Context, authorID int64, reqID, debtorUsername string) (string, error) {
	const op = "debt service: confirm return"

	const confirmReturnDebtMsg string = "Вы подтвердили, что получили деньги. Спасибо!"

	if err := s.debtRepo.Delete(ctx, authorID, reqID, debtorUsername); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return s.t.Get(confirmReturnDebtMsg), nil

}
