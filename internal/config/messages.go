package config

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
)

type SelectLangUIMessage struct {
	Lang []string
}

func (m *SelectLangUIMessage) Type() string {
	return "select_lang_ui"
}

func (m *SelectLangUIMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "config.langUI"})
	replyMarkup := &telebot.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}

	var btnRows []telebot.Row
	for _, lang := range m.Lang {
		btnText := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "lang." + lang})
		btn := replyMarkup.Text(btnText)
		btnRows = append(btnRows, telebot.Row{btn})
	}

	replyMarkup.Reply(btnRows...)
	options = append(options, replyMarkup)

	return
}
