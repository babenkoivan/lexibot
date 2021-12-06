package bot

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/config"
	"time"
)

type Bot interface {
	OnText(handler MessageHandler)
	OnCallback(action string, handler CallbackHandler)
	Send(recipient telebot.Recipient, msg Message)
	Edit(sig *MessageSig, msg Message)
	Start()
}

type bot struct {
	telebot *telebot.Bot
}

func (b *bot) OnText(handler MessageHandler) {
	b.telebot.Handle(telebot.OnText, func(m *telebot.Message) {
		handler.Handle(b, m)
	})
}

func (b *bot) OnCallback(action string, handler CallbackHandler) {
	b.telebot.Handle("\f"+action, func(c *telebot.Callback) {
		handler.Handle(b, c)
	})
}

func (b *bot) Send(recipient telebot.Recipient, msg Message) {
	b.telebot.Send(recipient, msg.Text(), msg.Options()...)
}

func (b *bot) Edit(sig *MessageSig, msg Message) {
	b.telebot.Edit(sig, msg.Text(), msg.Options()...)
}

func (b *bot) Start() {
	b.telebot.Start()
}

func NewBot(config config.Bot) (Bot, error) {
	poller := &telebot.LongPoller{Timeout: config.Timeout * time.Second}
	settings := telebot.Settings{Token: config.Token, Poller: poller}

	telebot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	return &bot{telebot}, nil
}
