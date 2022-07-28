package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type User struct {
	Id         int
	Username   string
	TelegramId int64
}

type Expense struct {
	id       int
	Borrower User
	Lender   User
	Sum      int
	Request  Request
	Returned bool
}

type Request struct {
	id      int
	Comment string
	ChatId  int64
	Date    time.Time
}

var db *sql.DB

// Init opens db and returns close function
func Init(dbPath string) (err error, closeFn func()) {
	var e error
	db, e = sql.Open("sqlite", dbPath)
	if e != nil {
		return fmt.Errorf("open db error: %w", e), nil
	}

	closeFn = func() {
		err := db.Close()
		if err != nil {
			log.Print(fmt.Errorf("close db error: %w", err))
		}
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return fmt.Errorf("ping db error: %w", pingErr), closeFn
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(fmt.Errorf("create tables transaction begin error: %w", err))
	}
	defer tx.Rollback()

	// add request
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS "expenses" (
		"id"	INTEGER,
		"borrower_id"	INTEGER NOT NULL,
		"lender_id"	INTEGER NOT NULL,
		"sum"	INTEGER DEFAULT 0,
		"request_id"	INTEGER NOT NULL,
		"returned"	INTEGER DEFAULT 0,
		FOREIGN KEY("borrower_id") REFERENCES "users"("id"),
		FOREIGN KEY("request_id") REFERENCES "requests"("id"),
		FOREIGN KEY("lender_id") REFERENCES "users"("id"),
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "requests" (
		"id"	INTEGER,
		"comment"	TEXT,
		"date"	TEXT NOT NULL,
		"chat_id"	INTEGER,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "users" (
		"id"	INTEGER,
		"username"	TEXT NOT NULL UNIQUE,
		"telegram_id"	INTEGER,
		PRIMARY KEY("id" AUTOINCREMENT)
	);`)

	if err != nil {
		log.Fatal(fmt.Errorf("create tables query execute error: %w", err))
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(fmt.Errorf("create tables commit transaction error: %w", err))
	}

	return nil, closeFn
}

func mustOpenDb() {
	if db == nil {
		log.Panic("db connection doesn't open")
	}
}
