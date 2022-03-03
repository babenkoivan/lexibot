package testkit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"testing"
)

type botSpy struct {
	testing  *testing.T
	messages map[int][]bot.Message
}

func (s *botSpy) OnMessage(handler bot.MessageHandler) {}

func (s *botSpy) OnReply(msg bot.Message, handler bot.ReplyHandler) {}

func (s *botSpy) OnCommand(command string, handler bot.MessageHandler) {}

func (s *botSpy) Send(to *telebot.User, msg bot.Message) {
	s.messages[to.ID] = append(s.messages[to.ID], msg)
}

func (s *botSpy) AssertSent(to *telebot.User, msg bot.Message) {
	assert.Contains(s.testing, s.messages[to.ID], msg, fmt.Sprintf("Message %#v is not sent", msg))
}

func (s *botSpy) Start() {}

func NewBotSpy(t *testing.T) *botSpy {
	return &botSpy{t, map[int][]bot.Message{}}
}
