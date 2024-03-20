package repo

import (
	"context"
	"fmt"

	"github.com/Karzoug/share_bot/internal/model"
)

func (r Repo) Get(ctx context.Context, reqID, debtorUsername string) (model.Debt, error) {
	const op = "debt repo: get"

	const getQuery = `SELECT * FROM debts WHERE request_id = ? AND debtor_username = ?`

	var debt model.Debt
	err := r.db.GetContext(ctx, &debt, getQuery, reqID, debtorUsername)

	if err != nil {
		return debt, fmt.Errorf("%s: %w", op, err)
	}
	return debt, nil
}

func (r Repo) GetByAuthorID(ctx context.Context, authorID int64) ([]model.Debt, error) {
	const op = "debt repo: get by author id"

	const getByAuthorQuery = `SELECT * FROM debts WHERE author_id = ?`

	debts := []model.Debt{}
	err := r.db.SelectContext(ctx, &debts, getByAuthorQuery, authorID)

	if err != nil {
		return debts, fmt.Errorf("%s: %w", op, err)
	}
	return debts, nil
}

func (r Repo) GetByDebtorUsername(ctx context.Context, debtorUsername string, onlyConfirmed bool) ([]model.Debt, error) {
	const op = "debt repo: get by debtor username"

	const (
		getByDebtorQuery          = `SELECT * FROM debts WHERE debtor_username = ?`
		getConfirmedByDebtorQuery = `SELECT * FROM debts WHERE debtor_username = ? AND confirmed = true`
	)

	debts := []model.Debt{}
	var err error
	if onlyConfirmed {
		err = r.db.SelectContext(ctx, &debts, getConfirmedByDebtorQuery, debtorUsername)
	} else {
		err = r.db.SelectContext(ctx, &debts, getByDebtorQuery, debtorUsername)
	}

	if err != nil {
		return debts, fmt.Errorf("%s: %w", op, err)
	}
	return debts, nil
}

func (r Repo) Save(ctx context.Context, debt model.Debt) error {
	const op = "debt repo: save"

	const (
		saveQuery = `INSERT INTO debts (author_id, request_id, debtor_username, sum, comment, date, confirmed)
VALUES (:author_id, :request_id, :debtor_username, :sum, :comment, :date, :confirmed)
ON CONFLICT(request_id, debtor_username) 
DO UPDATE SET author_id = excluded.author_id, sum = excluded.sum, comment = excluded.comment, 
date = excluded.date, confirmed = excluded.confirmed`
	)

	_, err := r.db.NamedExecContext(ctx, saveQuery, debt)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r Repo) Delete(ctx context.Context, authorID int64, reqID, debtorUsername string) error {
	const op = "debt repo: delete"

	const deleteQuery = `DELETE FROM debts WHERE request_id = ? AND debtor_username = ? AND author_id = ?`

	_, err := r.db.ExecContext(ctx, deleteQuery, reqID, debtorUsername, authorID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
