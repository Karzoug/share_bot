package model

type User struct {
	ID        int64
	Username  string
	FirstName string `db:"first_name"`
}
