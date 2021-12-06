package translation

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

type selectTranslationMessage struct {
	text         string
	translations []string
}

func (m *selectTranslationMessage) Text() string {
	return fmt.Sprintf("Select translation for %q", m.text)
}

func (m *selectTranslationMessage) Options() (options []interface{}) {
	markup := &telebot.ReplyMarkup{}

	var rows []telebot.Row
	for _, t := range m.translations {
		btn := markup.Data(t, OnSaveTranslation, m.text, t)
		rows = append(rows, telebot.Row{btn})
	}

	btn := markup.Data("✗", OnCancelTranslation)
	rows = append(rows, telebot.Row{btn})
	markup.Inline(rows...)

	return append(options, markup)
}

type savedTranslationMessage struct {
	text        string
	translation string
}

func (m *savedTranslationMessage) Text() string {
	return fmt.Sprintf("%s → %s", m.text, m.translation)
}

func (m *savedTranslationMessage) Options() (options []interface{}) {
	markup := &telebot.ReplyMarkup{}

	btn := markup.Data("✗", OnDeleteTranslation)
	row := telebot.Row{btn}
	markup.Inline(row)

	return append(options, markup)
}
