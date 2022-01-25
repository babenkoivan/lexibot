package app

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
)

const (
	OnStart    string = "/start"
	OnHelp     string = "/help"
	OnSettings string = "/settings"
)

func startHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &bot.LocalizedTextMessage{"app.start"})
	b.Send(msg.Sender, &settings.SelectLangUIMessage{})
}

func NewStartHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(startHandler)
}
