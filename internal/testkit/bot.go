package testkit

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"testing"
)

type botSpy struct {
	testing  *testing.T
	messages map[int][]bot.Message
}

func (b *botSpy) OnMessage(handler bot.MessageHandler) {}

func (b *botSpy) OnReply(msg bot.Message, handler bot.ReplyHandler) {}

func (b *botSpy) OnCommand(command string, handler bot.MessageHandler) {}

func (b *botSpy) Send(to *telebot.User, msg bot.Message) {
	b.messages[to.ID] = append(b.messages[to.ID], msg)
}

func (b *botSpy) AssertSent(to *telebot.User, msg bot.Message) {
	assert.Contains(b.testing, b.messages[to.ID], msg)
}

func (b *botSpy) Start() {}

func NewBotSpy(t *testing.T) *botSpy {
	return &botSpy{t, map[int][]bot.Message{}}
}
