package app

import (
	"golang.org/x/text/language"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/config"
)

const (
	OnStart    string = "/start"
	OnHelp     string = "/help"
	OnSettings string = "/settings"
)

func startHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &LocalizedTextMessage{"app.start"})

	langUI := []string{language.English.String(), language.Russian.String()}
	b.Send(msg.Sender, &config.SelectLangUIMessage{langUI})
}

func NewStartHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(startHandler)
}
