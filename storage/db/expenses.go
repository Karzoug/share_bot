package db

import (
	"share_bot/lib/e"
	"share_bot/storage"
	"time"
)

func (st Storage) AddRequest(req storage.Request) (err error) {
	st.mustOpenDb()

	if len(req.Exps) == 0 {
		return nil
	}

	defer func() { err = e.Wrap("can't add request", err) }()

	lender, exist := st.getUser(req.Lender)
	if !exist {
		st.addUser(&lender)
	}

	tableExpenses := make([]tableExpense, 0, len(req.Exps))

	for _, v := range req.Exps {
		borrower, exist := st.getUser(v.Person)
		if !exist {
			st.addUser(&borrower)
		}

		exp := tableExpense{
			Borrower: borrower.Id,
			Lender:   lender.Id,
			Sum:      v.Sum,
			Returned: false,
		}

		tableExpenses = append(tableExpenses, exp)
	}

	tx, err := st.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	// add request
	result, err := tx.Exec("INSERT INTO requests(comment, date, chat_id) VALUES($1, $2, $3)", req.Comment, req.Date.Unix(), req.ChatId)
	if err != nil {
		return
	}
	resID, err := result.LastInsertId()
	if err != nil {
		return
	}
	tableReqId := int(resID)

	// add expenses
	for _, exp := range tableExpenses {
		result, err = tx.Exec("INSERT INTO expenses(sum, lender_id, borrower_id, request_id) VALUES($1, $2, $3, $4)",
			exp.Sum, exp.Lender, exp.Borrower, tableReqId)
		if err != nil {
			return
		}
		_, err = result.LastInsertId()
		if err != nil {
			return
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return nil
}

func (st Storage) ReturnExpense(lender, borrower string, sum int) error {
	return nil
}

func (st Storage) GetRequestsByBorrower(borrower string, onlyNotReturned bool) (exps []storage.Request, err error) {
	st.mustOpenDb()

	defer func() { err = e.Wrap("can't get expense by borrower", err) }()

	exps = make([]storage.Request, 0)

	dbBorrower, exist := st.getUser(borrower)
	if !exist {
		return exps, storage.ErrUserNotExist
	}

	rows, err := st.db.Query(`SELECT expenses.sum, requests.id, requests.comment, requests.date, users.username
		FROM expenses 
		JOIN requests ON expenses.request_id = requests.id 
		JOIN users ON expenses.lender_id = users.id
		WHERE expenses.borrower_id = $1 AND expenses.returned = 0
		ORDER BY requests.id`, dbBorrower.Id)
	if err != nil {
		return
	}
	defer rows.Close()

	var (
		date      string
		reqId     int
		lastReqId int             = -1
		e         storage.Expense = storage.Expense{}
		r         storage.Request = storage.Request{}
		acc       storage.Request = storage.Request{}
	)

	for rows.Next() {
		if err = rows.Scan(&e.Sum, &reqId, &r.Comment, &date, &e.Person); err != nil {
			return
		}
		r.Date, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", date)

		if reqId != lastReqId {
			if len(acc.Exps) != 0 {
				exps = append(exps, acc)
			}
			lastReqId = reqId
			acc = r
			acc.Exps = make([]storage.Expense, 0, 1)
		}
		acc.Exps = append(acc.Exps, e)
	}
	if len(acc.Exps) != 0 {
		exps = append(exps, acc)
	}

	if err = rows.Err(); err != nil {
		return
	}
	return exps, nil
}

func (st Storage) GetRequestsByLender(lender string, onlyNotReturned bool) (exps []storage.Request, err error) {
	exps = make([]storage.Request, 0)
	return
}
