package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Karzoug/share_bot/internal/model"
)

func (r Repo) GetDebts(ctx context.Context, before *time.Time) ([]model.Debt, error) {
	const op = "remind service: get debts"

	const (
		getDebtsQuery         = `SELECT * FROM debts`
		getDebtsQueryWithTime = `SELECT * FROM debts WHERE date < ?`
	)

	debts := []model.Debt{}

	var err error
	if before == nil {
		err = r.db.SelectContext(ctx, &debts, getDebtsQuery)
	} else {
		err = r.db.SelectContext(ctx, &debts, getDebtsQuery, before)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return debts, nil
}
