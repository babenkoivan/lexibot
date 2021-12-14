package bot

import (
	"gopkg.in/tucnak/telebot.v2"
)

const (
	OnStart    string = "/start"
	OnHelp     string = "/help"
	OnSettings string = "/settings"
)

type MessageHandler interface {
	Handle(b Bot, m *telebot.Message)
}

type CallbackHandler interface {
	Handle(b Bot, c *telebot.Callback)
}

type MessageHandlerFunc func(b Bot, m *telebot.Message)

func (h MessageHandlerFunc) Handle(b Bot, m *telebot.Message) {
	h(b, m)
}

type CallbackHandlerFunc func(b Bot, c *telebot.Callback)

func (h CallbackHandlerFunc) Handle(b Bot, c *telebot.Callback) {
	h(b, c)
}

func startHandler(b Bot, m *telebot.Message) {
	b.Send(m.Sender, &infoMessage{"#todo start"})
}

func NewStartHandler() MessageHandler {
	return MessageHandlerFunc(startHandler)
}

func helpHandler(b Bot, m *telebot.Message) {
	b.Send(m.Sender, &infoMessage{"#todo help"})
}

func NewHelpHandler() MessageHandler {
	return MessageHandlerFunc(helpHandler)
}

func settingsHandler(b Bot, m *telebot.Message) {
	b.Send(m.Sender, &infoMessage{"#todo settings"})
}

func NewSettingsHandler() MessageHandler {
	return MessageHandlerFunc(settingsHandler)
}
