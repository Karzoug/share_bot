package db

import (
	"fmt"
	"log"
	"share_bot/parse"
	"time"
)

// AddExpenses calls in group and private chat for add expenses
func AddExpenses(lender User, req Request, msgs []parse.ExpenseMessage) {
	mustOpenDb()

	if len(msgs) == 0 {
		return
	}

	dbLender, ok := getUser(lender.Username)
	if !ok {
		addUser(&lender)
	} else {
		lender = dbLender
	}

	expenses := make([]Expense, 0, len(msgs))

	for _, v := range msgs {
		exp := Expense{
			Borrower: User{
				Username: v.Borrower,
			},
			Lender:   lender,
			Sum:      v.Sum,
			Request:  req,
			Returned: false,
		}

		dbBorrower, ok := getUser(exp.Borrower.Username)
		if !ok {
			addUser(&exp.Borrower)
		} else {
			exp.Borrower = dbBorrower
		}

		expenses = append(expenses, exp)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println(fmt.Errorf("insert expenses and request transaction begin error: %w", err))
		return
	}
	defer tx.Rollback()

	// add request
	result, err := tx.Exec("INSERT INTO requests(comment, date, chat_id) VALUES($1, $2, $3)", req.Comment, req.Date.Unix(), req.ChatId)
	if err != nil {
		log.Println(fmt.Errorf("insert request query execute error: %w", err))
		return
	}
	resID, err := result.LastInsertId()
	if err != nil {
		log.Println(fmt.Errorf("insert request get last id error: %w", err))
		return
	}
	req.id = int(resID)

	// add expenses
	for i, exp := range expenses {
		expenses[i].Request.id = req.id
		result, err := tx.Exec("INSERT INTO expenses(sum, lender_id, borrower_id, request_id) VALUES($1, $2, $3, $4)",
			exp.Sum, exp.Lender.Id, exp.Borrower.Id, req.id)
		if err != nil {
			log.Println(fmt.Errorf("insert expenses query execute error: %w", err))
			return
		}
		resID, err := result.LastInsertId()
		if err != nil {
			log.Println(fmt.Errorf("insert expenses get last id error: %w", err))
			return
		}
		expenses[i].id = int(resID)
	}

	if err = tx.Commit(); err != nil {
		log.Println(fmt.Errorf("insert expenses and request transaction commit error: %w", err))
		return
	}
}

func ShowExpensesByBorrower(username string) ([]Expense, error) {
	mustOpenDb()

	exps := make([]Expense, 0)

	dbBorrower, exist := getUser(username)
	if !exist {
		return exps, nil
	}

	rows, err := db.Query(`SELECT expenses.sum, requests.comment, requests.date, users.username
		FROM expenses 
		JOIN requests ON expenses.request_id = requests.id 
		JOIN users ON expenses.lender_id = users.id
		WHERE expenses.borrower_id = $1 AND expenses.returned = 0`, dbBorrower.Id)
	if err != nil {
		return exps, fmt.Errorf("select expenses by borrower query execute error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		e := Expense{}
		var d string
		if err := rows.Scan(&e.Sum, &e.Request.Comment, &d, &e.Lender.Username); err != nil {
			return exps, fmt.Errorf("select expenses by borrower scan row error: %w", err)
		}
		e.Request.Date, _ = time.Parse("2006-01-02 15:04:05 -0700 MST", d)

		exps = append(exps, e)
	}
	if err = rows.Err(); err != nil {
		return exps, fmt.Errorf("select expenses by borrower row error: %w", err)
	}
	return exps, nil
}
