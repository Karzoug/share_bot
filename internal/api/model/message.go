package model

type Message struct {
	Text string `json:"text,omitempty"`
	ID   int64  `json:"message_id"`
	Chat Chat   `json:"chat"`
	From *User  `json:"from,omitempty"`
	Date int64  `json:"date"`
}
