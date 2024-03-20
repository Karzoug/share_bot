package model

// User represents a Telegram user or bot.
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	ID        int64  `json:"id"`
}
