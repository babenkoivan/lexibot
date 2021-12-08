package translation

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
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
	translation *Translation
}

func (m *savedTranslationMessage) Text() string {
	return fmt.Sprintf("%s → %s", m.translation.Text, m.translation.Translation)
}

func (m *savedTranslationMessage) Options() (options []interface{}) {
	markup := &telebot.ReplyMarkup{}

	btn := markup.Data("✗", OnDeleteTranslation, strconv.FormatUint(m.translation.ID, 10))
	row := telebot.Row{btn}
	markup.Inline(row)

	return append(options, markup)
}
