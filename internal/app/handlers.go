package app

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
)

const (
	OnStart string = "/start"
	OnHelp  string = "/help"
)

func startHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &StartMessage{})
	b.Send(msg.Sender, &settings.SelectLangUIMessage{})
}

func NewStartHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(startHandler)
}

func helpHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &HelpMessage{})
}

func NewHelpHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(helpHandler)
}
