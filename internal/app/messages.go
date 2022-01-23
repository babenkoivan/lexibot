package app

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
)

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

type LocalizedTextMessage struct {
	MessageID string
}

func (m *LocalizedTextMessage) Type() string {
	return "localized_text"
}

func (m *LocalizedTextMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: m.MessageID})
	options = append(options, &telebot.ReplyMarkup{ReplyKeyboardRemove: true})
	return
}
