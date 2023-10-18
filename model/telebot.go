package model

import "io"

type TeleBot interface {
	SetWebhook() error
	GetUpdate(r io.Reader) (*TeleBotUpdate, error)
	SendTextMessage(chatID int64, text string) (*TeleBotMessage, error)
}

type TeleBotMessage struct {
	MessageID int
	Chat      *TeleBotChat
	Text      string // optional

}

type TeleBotUpdate struct {
	UpdateID int
	Message  *TeleBotMessage
}

// TeleBotChat contains information about the place a message was sent.
type TeleBotChat struct {
	ID   int64
	Type string
}
