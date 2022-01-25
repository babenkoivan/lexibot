package bot

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
)

type Message interface {
	Type() string
	Render(localizer *i18n.Localizer) (text string, options []interface{})
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
