package storage

import (
	"errors"
	"time"
)

type Storage interface {
	AddRequest(req *Request) error
	GetRequestsByBorrower(borrower string, onlyNotReturned bool) ([]Request, error)
	GetRequestsByLender(lender string, onlyNotReturned bool) ([]Request, error)
	GetNotReturnedRequests() ([]Request, error)
	GetExpenseWithRequest(expId int) (Request, error)
	ApproveExpense(reqId int64, username string) error
	ApproveReturnExpense(reqId int64, username string) error
	GetUserByUsername(username string) (User, bool, error)
	GetUserById(id int) (User, bool, error)
	IsUserExist(username string) (bool, error)
	SaveUser(user User) (err error)
}

type Request struct {
	Id      int64
	Lender  string
	Exps    []Expense
	Comment string
	Date    time.Time
	ChatId  int64
}

type Expense struct {
	Id       int64
	Borrower string
	Sum      int
	Returned bool
	Approved bool
}

type User struct {
	Id       int64
	Username string
	ChatId   int64
}

var (
	ErrUserNotExist = errors.New("user doesn't exist")
	ErrNoResult     = errors.New("no data affected by command")
)
