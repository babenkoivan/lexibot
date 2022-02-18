package bot

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/localization"
)

type Message interface {
	Type() string
	Render(localizer localization.Localizer) (text string, options []interface{})
}

func WithReplyKeyboard(captions []string) *telebot.ReplyMarkup {
	replyMarkup := &telebot.ReplyMarkup{
		ResizeReplyKeyboard: true,
		OneTimeKeyboard:     true,
	}

	var rows []telebot.Row
	for _, c := range captions {
		btn := replyMarkup.Text(c)
		rows = append(rows, telebot.Row{btn})
	}

	replyMarkup.Reply(rows...)
	return replyMarkup
}

func WithoutReplyKeyboard() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{ReplyKeyboardRemove: true}
}
