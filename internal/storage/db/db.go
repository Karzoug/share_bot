package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"share_bot/internal/logger"
	"share_bot/pkg/e"

	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

// New creates Storage instance, opens db, creates tables if not exists and returns close function
func New(dbPath string) (st Storage, closeFn func()) {
	st = Storage{}

	var err error
	defer func() {
		if err != nil {
			logger.Logger.Fatal("can't create new db storage", zap.Error(err))
		}
	}()

	err = os.MkdirAll(filepath.Dir(dbPath), 0750)
	if err != nil && !os.IsExist(err) {
		err = e.Wrap("cannot create directories to store database", err)
		return
	}

	st.db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		err = e.Wrap("cannot open database", err)
		return
	}

	closeFn = func() {
		err := st.db.Close()
		if err != nil {
			logger.Logger.Error("cannot close database", zap.Error(err))
		}
	}

	err = st.db.Ping()
	if err != nil {
		err = e.Wrap("cannot ping database", err)
		return
	}

	tx, err := st.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(createQuery)

	if err != nil {
		err = e.Wrap("error during initial creation of tables in database", err)
		return
	}

	if err = tx.Commit(); err != nil {
		err = e.Wrap("error during initial creation of tables in database", err)
		return
	}

	return
}

func (st Storage) mustOpenDb() {
	if st.db == nil {
		logger.Logger.Panic("database connection is not open")
	}
}
