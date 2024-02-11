package model

// AnswerCallbackQuery is an answer to callback queries sent from inline keyboards.
// The answer will be displayed to the user as a notification
// at the top of the chat screen or as an alert.
type AnswerCallbackQuery struct {
	// Unique identifier for the query to be answered
	CallbackQueryID string `json:"callback_query_id"`
	// Text of the notification. If not specified, nothing will be shown to the user, 0-200 characters
	Text string `json:"text"`
}
