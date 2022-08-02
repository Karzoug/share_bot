package db

import "time"

type tableUser struct {
	Id       int
	Username string
	ChatId   int64
}

type tableExpense struct {
	id       int
	Borrower int
	Lender   int
	Sum      int
	Request  int
	Returned bool
}

type tableRequest struct {
	id      int
	Comment string
	ChatId  int64
	Date    time.Time
}
