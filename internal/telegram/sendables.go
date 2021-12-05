package telegram

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	OnCancelTranslation string = "cancel_translation"
	OnSaveTranslation   string = "save_translation"
	OnDeleteTranslation string = "delete_translation"
)

type errorMessage struct {
	err error
}

func (m *errorMessage) Send(bot *telebot.Bot, recipient telebot.Recipient, options *telebot.SendOptions) (*telebot.Message, error) {
	text := fmt.Sprintf("Error: %s", m.err)
	return bot.Send(recipient, text)
}

type selectTranslationMessage struct {
	text         string
	translations []string
}

func (m *selectTranslationMessage) Send(bot *telebot.Bot, recipient telebot.Recipient, options *telebot.SendOptions) (*telebot.Message, error) {
	text := fmt.Sprintf("Select translation for %q", m.text)
	markup := &telebot.ReplyMarkup{}

	var rows []telebot.Row
	for _, t := range m.translations {
		btn := markup.Data(t, OnSaveTranslation, m.text, t)
		rows = append(rows, telebot.Row{btn})
	}

	btn := markup.Data("✗", OnCancelTranslation)
	rows = append(rows, telebot.Row{btn})
	markup.Inline(rows...)

	return bot.Send(recipient, text, markup)
}
