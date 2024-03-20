package model

import "time"

type Debt struct {
	AuthorID       int64  `db:"author_id"`
	RequestID      string `db:"request_id"`
	DebtorUsername string `db:"debtor_username"`
	Sum            int64
	Comment        string
	Date           time.Time
	Confirmed      bool
}
