package db

import (
	"database/sql"
	"log"
	"share_bot/pkg/e"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db     *sql.DB
	logger *log.Logger
}

// New creates Storage instance, opens db, creates tables if not exists and returns close function
func New(dbPath string, logger *log.Logger) (st Storage, closeFn func()) {
	st = Storage{}

	var err error
	defer func() {
		if err != nil {
			logger.Panic(e.Wrap("can't create new db storage", err))
		}
	}()

	st.db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return
	}

	closeFn = func() {
		err := st.db.Close()
		if err != nil {
			logger.Print(e.Wrap("close db error", err))
		}
	}

	err = st.db.Ping()
	if err != nil {
		return
	}

	tx, err := st.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(createQuery)

	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (st Storage) mustOpenDb() {
	if st.db == nil {
		st.logger.Panic("db connection doesn't open")
	}
}
