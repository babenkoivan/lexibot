package telegram

import "gopkg.in/tucnak/telebot.v2"

func (b *bot) Translate(m *telebot.Message) {
	b.translate(m)
}
