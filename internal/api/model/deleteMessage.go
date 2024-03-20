package model

// DeleteMessageRequest represents a request to delete a message.
type DeleteMessageRequest struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int64 `json:"message_id"`
}
