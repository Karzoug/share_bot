package model

type SendMessageRequest struct {
	ChatID      int64       `json:"chat_id"`
	Text        string      `json:"text"`
	ParseMode   string      `json:"parse_mode,omitempty"`
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type ReplyMarkup interface {
	ImplementsReplyMarkup()
}

// InlineKeyboardMarkup represents an inline keyboard.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

func (InlineKeyboardMarkup) ImplementsReplyMarkup() {}

// InlineKeyboardButton represents a button in an inline keyboard.
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
}

type ReplyKeyboardMarkup struct {
	// Keyboard is a slice of button rows, each represented by a slcie of KeyboardButton objects
	Keyboard [][]KeyboardButton `json:"keyboard"`
	// IsPersistent requests clients to always show the keyboard when the regular keyboard is hidden.
	// Defaults to false, in which case the custom keyboard can be hidden and opened with a keyboard icon.
	IsPersistent bool `json:"is_persistent,omitempty"`
	// ResizeKeyboard requests clients to resize the keyboard vertically for optimal fit
	// (e.g., make the keyboard smaller if there are just two rows of buttons).
	// Defaults to false, in which case the custom keyboard is always of the same height as the app's standard keyboard.
	ResizeKeyboard bool `json:"resize_keyboard,omitempty"`
}

func (ReplyKeyboardMarkup) ImplementsReplyMarkup() {}

type KeyboardButton struct {
	Text string `json:"text"`
}
