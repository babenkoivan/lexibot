package telegram

import (
	"gopkg.in/tucnak/telebot.v2"
)

type TranslationMessage struct {
	Text   string
	Result []string
}

func (t *TranslationMessage) Send(b *telebot.Bot, recipient telebot.Recipient, options *telebot.SendOptions) (*telebot.Message, error) {
	markup := telebot.ReplyMarkup{}

	var btnRow telebot.Row

	for _, r := range t.Result {
		btnRow = append(btnRow, markup.Text(r))
	}

	markup.Inline(btnRow)

	return b.Send(recipient, t.Text, markup)
}
