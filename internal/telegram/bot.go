package telegram

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/configs"
	"time"
)

func NewBot(config configs.Telegram) (*telebot.Bot, error) {
	poller := &telebot.LongPoller{Timeout: config.Timeout * time.Second}
	settings := telebot.Settings{Token: config.Token, Poller: poller}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

type handlerRegister struct {
	bot *telebot.Bot
}

func (r *handlerRegister) Text(handler MessageHandler) {
	r.bot.Handle(telebot.OnText, handler.Handle)
}

func (r *handlerRegister) Callback(action string, handler CallbackHandler) {
	r.bot.Handle("\f"+action, handler.Handle)
}

func NewHandlerRegister(bot *telebot.Bot) *handlerRegister {
	return &handlerRegister{bot: bot}
}
