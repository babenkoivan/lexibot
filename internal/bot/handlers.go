package bot

import (
	"gopkg.in/tucnak/telebot.v2"
)

type MessageHandler interface {
	Handle(b Bot, m *telebot.Message)
}

type CallbackHandler interface {
	Handle(b Bot, c *telebot.Callback)
}

type CallbackHandlerFunc func(b Bot, c *telebot.Callback)

func (h CallbackHandlerFunc) Handle(b Bot, c *telebot.Callback) {
	h(b, c)
}
