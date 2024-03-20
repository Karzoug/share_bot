package repo

import "github.com/jmoiron/sqlx"

func New(db *sqlx.DB) Repo {
	return Repo{db: db}
}

type Repo struct {
	db *sqlx.DB
}
