package bot

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

type MessageSig struct {
	messageId string
	chatId    int64
}

func (s *MessageSig) MessageSig() (messageID string, chatID int64) {
	return s.messageId, s.chatId
}

func ExtractMessageSig(msg *telebot.Message) *MessageSig {
	messageId, chatId := msg.MessageSig()
	return &MessageSig{messageId, chatId}
}

type Message interface {
	Text() string
	Options() (options []interface{})
}

type errorMessage struct {
	err error
}

func (m *errorMessage) Text() string {
	return fmt.Sprintf("❗️ %s", m.err)
}

func (m *errorMessage) Options() (options []interface{}) {
	return
}

func NewErrorMessage(err error) *errorMessage {
	return &errorMessage{err}
}

type infoMessage struct {
	text string
}

func (m *infoMessage) Text() string {
	return fmt.Sprintf("⚠️ %s", m.text)
}

func (m *infoMessage) Options() (options []interface{}) {
	return
}

func NewInfoMessage(text string) *infoMessage {
	return &infoMessage{text: text}
}
