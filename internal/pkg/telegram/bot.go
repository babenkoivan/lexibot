package telegram

import (
	"context"
	"golang.org/x/text/language"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/pkg/translations"
)

type Telegram interface {
	Start()
	Send(to telebot.Recipient, what interface{}, options ...interface{}) (*telebot.Message, error)
	Handle(endpoint interface{}, handler interface{})
}

type bot struct {
	telegram   Telegram
	translator translations.Translator
}

func NewBot(telegram Telegram, translator translations.Translator) *bot {
	return &bot{telegram: telegram, translator: translator}
}

func (b *bot) Start() {
	b.telegram.Handle(b.telegram, b.translate)
	b.telegram.Start()
}

func (b *bot) translate(m *telebot.Message) {
	// todo pass context from main?
	ctx := context.Background()

	// todo take from user config
	from := language.German
	to := language.English

	// todo error handling
	res, _ := b.translator.Translate(ctx, from, to, m.Text)

	b.telegram.Send(m.Sender, TranslationMessage{m.Text, res})
}
