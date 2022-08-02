package storage

import (
	"errors"
	"time"
)

type Storage interface {
	AddRequest(req Request) error
	ReturnExpense(lender, borrower string, sum int) error
	GetRequestsByBorrower(borrower string, onlyNotReturned bool) ([]Request, error)
	GetRequestsByLender(lender string, onlyNotReturned bool) ([]Request, error)
	SaveUser(user User)
}

type Request struct {
	Lender  string
	Exps    []Expense
	Comment string
	Date    time.Time
	ChatId  int64
}

type Expense struct {
	// Borrower or Lender
	Person string
	Sum    int
}

type User struct {
	Username string
	ChatId   int64
}

var (
	ErrUnknownMetaType = errors.New("unknown meta type")
	ErrUserNotExist    = errors.New("user doesn't exist")
)
