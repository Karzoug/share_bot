package repo

import (
	"context"
	"fmt"

	"github.com/Karzoug/share_bot/internal/model"
)

func (r Repo) GetByID(ctx context.Context, id int64) (model.User, error) {
	const op = "user repo: get by id"

	const getByIDQuery = `SELECT * FROM users WHERE id = ?`

	u := model.User{}
	err := r.db.GetContext(ctx, &u, getByIDQuery, id)

	if err != nil {
		return u, fmt.Errorf("%s: %w", op, err)
	}
	return u, nil
}

func (r Repo) GetByUsername(ctx context.Context, username string) (model.User, error) {
	const (
		op                 = "user repo: get by username"
		getByUsernameQuery = `SELECT * FROM users WHERE username = ?`
	)

	u := model.User{}
	err := r.db.GetContext(ctx, &u, getByUsernameQuery, username)

	if err != nil {
		return u, fmt.Errorf("%s: %w", op, err)
	}
	return u, nil
}

func (r Repo) Save(ctx context.Context, user model.User) error {
	const op = "user repo: save"

	const saveQuery = `INSERT INTO users (id, username, first_name)
VALUES (:id, :username, :first_name)
ON CONFLICT(id) 
DO UPDATE SET username = excluded.username, first_name = excluded.first_name`

	_, err := r.db.NamedExecContext(ctx, saveQuery, user)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
