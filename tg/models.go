package tg

import "io"

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type Message struct {
	Document    File                 `json:"document"`
	Text        string               `json:"text,omitempty"`
	Chat        Chat                 `json:"chat,omitempty"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	ID          int64                `json:"message_id"`
	From        User                 `json:"from,omitempty"`
}

type Update struct {
	Message  Message       `json:"message,omitempty"`
	CallBack CallBackQuery `json:"callback_query,omitempty"`
	ID       int64         `json:"update_id,omitempty"`
}

type User struct {
	ID int64 `json:"id,omitempty"`
}

type Chat struct {
	Type string `json:"type,omitempty"`
	ID   int64  `json:"id,omitempty"`
}

type CallBackQuery struct {
	ID      string  `json:"id,omitempty"`
	Data    string  `json:"data,omitempty"`
	Mensaje Message `json:"message,omitempty"`
	From    User    `json:"from,omitempty"`
}

type InlineKeyboardButton struct {
	Text     string `json:"text,omitempty"`
	CallBack string `json:"callback_data,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboardB [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type File struct {
	ID       string `json:"file_id"`
	UniqueID string `json:"file_unique_id"`
	Filename string `json:"file_name"`
}

type SendDocumentResponse struct {
	Result Message `json:"result"`
	Ok     bool    `json:"ok"`
}

type SendMessage struct {
	Text                     string               `json:"text,omitempty"`
	ReplyMarkup              InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	ChatID                   int64                `json:"chat_id,omitempty"`
	RelayMessageID           int64                `json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply bool                 `json:"allow_sending_without_reply,omitempty"`
	ParseMode                string               `json:"parse_mode"`
}

type EditMessage struct {
	Text        string               `json:"text,omitempty"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	ChatID      int64                `json:"chat_id,omitempty"`
	MessageID   int64                `json:"message_id,omitempty"`
}

type SendDocument struct {
	File           io.Reader
	Document       string `json:"document"`
	ChatID         int64  `json:"chat_id"`
	ProtectContent bool   `json:"protect_content"`
}

type DeleteMessage struct {
	ChatID    int64 `json:"chat_id"`
	MessageID int64 `json:"message_id"`
}

type TelegramError struct {
	ErrorCode   uint                   `json:"error_code"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

func (err TelegramError) Error() string {
	return err.Description
}
