package telegram_test

import (
	"context"
	"errors"
	"golang.org/x/text/language"
	"gopkg.in/tucnak/telebot.v2"
	telegram2 "lexibot/internal/telegram"
	"lexibot/internal/translations"
	"reflect"
	"testing"
)

func TestBotTranslate(t *testing.T) {
	telebotSpy := &telebotSpy{}
	translatorStub := &translatorStub{map[string]string{"bunt": "colorful"}}

	bot := telegram2.NewBot(telebotSpy, translatorStub)

	t.Run("sends translations when found", func(t *testing.T) {
		incoming := &telebot.Message{Text: "bunt"}
		outgoing := telegram2.TranslationMessage{Text: incoming.Text, Result: []string{"colorful"}}

		bot.Translate(incoming)
		telebotSpy.assertSent(t, outgoing)
	})
}

type telebotSpy struct {
	sent []interface{}
}

func (t *telebotSpy) Start() {

}

func (t *telebotSpy) Send(to telebot.Recipient, what interface{}, options ...interface{}) (*telebot.Message, error) {
	t.sent = append(t.sent, what)
	return &telebot.Message{}, nil
}

func (t *telebotSpy) Handle(endpoint interface{}, handler interface{}) {

}

func (t *telebotSpy) assertSent(testing testing.TB, want interface{}) {
	testing.Helper()

	for _, got := range t.sent {
		if reflect.DeepEqual(got, want) {
			return
		}
	}

	testing.Errorf("%v was not sent", want)
}

type translatorStub struct {
	translations map[string]string
}

func (t *translatorStub) Translate(ctx context.Context, from, to language.Tag, text string) ([]string, error) {
	translation, ok := t.translations[text]

	if !ok {
		return nil, errors.New(translations.NoTranslationsErr)
	}

	return []string{translation}, nil
}
