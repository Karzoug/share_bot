package db

import (
	"share_bot/lib/e"
	"share_bot/storage"
	"time"
)

func (st Storage) AddRequest(req *storage.Request) (err error) {
	st.mustOpenDb()

	if len(req.Exps) == 0 {
		return nil
	}

	defer func() { err = e.Wrap("can't add request", err) }()

	lender, exist := st.GetUserByUsername(req.Lender)
	if !exist {
		lender.Username = req.Lender
		st.addUser(&lender)
	}

	borrIds := make(map[string]int64, len(req.Exps))

	for _, v := range req.Exps {
		borrower, exist := st.GetUserByUsername(v.Borrower)
		if !exist {
			borrower.Username = v.Borrower
			st.addUser(&borrower)
		}

		borrIds[borrower.Username] = borrower.Id
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
	req.Id, err = result.LastInsertId()
	if err != nil {
		return
	}

	// add expenses
	for _, exp := range req.Exps {
		result, err = tx.Exec("INSERT INTO expenses(sum, lender_id, borrower_id, request_id) VALUES($1, $2, $3, $4)",
			exp.Sum, lender.Id, borrIds[exp.Borrower], req.Id)
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

func (st Storage) ApproveExpense(reqId int64, borrowerUsername string) (err error) {
	defer func() { err = e.Wrap("can't approve expense", err) }()

	user, exist := st.GetUserByUsername(borrowerUsername)
	if !exist {
		err = storage.ErrUserNotExist
		return
	}

	stmt, err := st.db.Prepare("UPDATE expenses SET approved = 1 WHERE request_id = $1 AND borrower_id=$2")
	if err != nil {
		return
	}
	r, err := stmt.Exec(reqId, user.Id)
	if err != nil {
		return
	}

	count, err := r.RowsAffected()
	if err != nil {
		return
	}
	if count == 0 {
		err = storage.ErrNoResult
	}
	return
}

func (st Storage) ApproveReturnExpense(expId int64, lenderUsername string) (err error) {
	defer func() { err = e.Wrap("can't approve return expense", err) }()

	user, exist := st.GetUserByUsername(lenderUsername)
	if !exist {
		err = storage.ErrUserNotExist
		return
	}

	stmt, err := st.db.Prepare("UPDATE expenses SET returned = 1 WHERE id = $1 AND lender_id=$2")
	if err != nil {
		return
	}
	r, err := stmt.Exec(expId, user.Id)
	if err != nil {
		return
	}

	count, err := r.RowsAffected()
	if err != nil {
		return
	}
	if count == 0 {
		err = storage.ErrNoResult
	}
	return
}

func (st Storage) ReturnExpense(lender, borrower string, sum int) error {
	return nil
}

func (st Storage) GetRequestsByBorrower(borrower string, onlyNotReturned bool) (exps []storage.Request, err error) {
	st.mustOpenDb()

	defer func() { err = e.Wrap("can't get expense by borrower", err) }()

	exps = make([]storage.Request, 0)

	dbBorrower, exist := st.GetUserByUsername(borrower)
	if !exist {
		return exps, storage.ErrUserNotExist
	}

	rows, err := st.db.Query(`SELECT expenses.id, expenses.sum, requests.id, requests.comment, requests.date, users.username
		FROM expenses 
		JOIN requests ON expenses.request_id = requests.id 
		JOIN users ON expenses.lender_id = users.id
		WHERE expenses.borrower_id = $1 AND expenses.returned = 0 AND expenses.approved = 1
		ORDER BY requests.id`, dbBorrower.Id)
	if err != nil {
		return
	}
	defer rows.Close()

	var (
		date      string
		lastReqId int64           = -1
		e         storage.Expense = storage.Expense{Borrower: dbBorrower.Username}
		r         storage.Request = storage.Request{}
		acc       storage.Request = storage.Request{}
	)

	for rows.Next() {
		if err = rows.Scan(&e.Id, &e.Sum, &r.Id, &r.Comment, &date, &r.Lender); err != nil {
			return
		}
		r.Date, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", date)

		if r.Id != lastReqId {
			if len(acc.Exps) != 0 {
				exps = append(exps, acc)
			}
			lastReqId = r.Id
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

func (st Storage) GetExpenseWithRequest(expId int) (req storage.Request, err error) {
	st.mustOpenDb()

	defer func() { err = e.Wrap("can't get expense with request", err) }()

	req = storage.Request{
		Exps: []storage.Expense{{}},
	}

	rows, err := st.db.Query(`SELECT expenses.id, expenses.sum, requests.id, requests.comment, requests.date, lu.username, bu.username
	FROM expenses 
	JOIN requests ON expenses.request_id = requests.id 
	JOIN users lu ON expenses.lender_id = lu.id
	JOIN users bu ON expenses.borrower_id = bu.id
	WHERE expenses.id = $1 AND expenses.returned = 0 AND expenses.approved = 1`, expId)
	if err != nil {
		return
	}
	defer rows.Close()

	var date string

	rows.Next()
	if err = rows.Scan(&req.Exps[0].Id, &req.Exps[0].Sum, &req.Id, &req.Comment, &date, &req.Lender, &req.Exps[0].Borrower); err != nil {
		return
	}
	req.Date, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", date)

	if err = rows.Err(); err != nil {
		return storage.Request{}, err
	}
	return req, nil
}

func (st Storage) GetRequestsByLender(lender string, onlyNotReturned bool) (exps []storage.Request, err error) {
	st.mustOpenDb()

	defer func() { err = e.Wrap("can't get expense by lender", err) }()

	exps = make([]storage.Request, 0)

	dbLender, exist := st.GetUserByUsername(lender)
	if !exist {
		return exps, storage.ErrUserNotExist
	}

	rows, err := st.db.Query(`SELECT expenses.id, expenses.sum, expenses.approved, requests.id, requests.comment, requests.date, users.username
		FROM expenses 
		JOIN requests ON expenses.request_id = requests.id 
		JOIN users ON expenses.borrower_id = users.id
		WHERE expenses.lender_id = $1 AND expenses.returned = 0
		ORDER BY requests.id`, dbLender.Id)
	if err != nil {
		return
	}
	defer rows.Close()

	var (
		date      string
		lastReqId int64           = -1
		e         storage.Expense = storage.Expense{}
		r         storage.Request = storage.Request{Lender: dbLender.Username}
		acc       storage.Request = storage.Request{}
	)

	for rows.Next() {
		if err = rows.Scan(&e.Id, &e.Sum, &e.Approved, &r.Id, &r.Comment, &date, &e.Borrower); err != nil {
			return
		}
		r.Date, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", date)

		if r.Id != lastReqId {
			if len(acc.Exps) != 0 {
				exps = append(exps, acc)
			}
			lastReqId = r.Id
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
