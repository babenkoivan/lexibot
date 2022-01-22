package bot

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
)

type Message interface {
	Type() string
	Render() (*i18n.LocalizeConfig, *telebot.ReplyMarkup)
}

//type ErrorMessage struct {
//	Err error
//}
//
//func (m *ErrorMessage) Type() string {
//	return "error"
//}
//
//func (m *ErrorMessage) Render() (text string, options []interface{}) {
//	text = fmt.Sprintf("❗️ %s", m.Err)
//	return
//}

type PlainTextMessage struct {
	MessageID string
}

func (m *PlainTextMessage) Type() string {
	return "plain_text"
}

func (m *PlainTextMessage) Render() (*i18n.LocalizeConfig, *telebot.ReplyMarkup) {
	return &i18n.LocalizeConfig{MessageID: m.MessageID}, nil
}
