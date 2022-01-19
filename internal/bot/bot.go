package bot

import (
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

type Bot interface {
	OnText(handler MessageHandler)
	OnCommand(command string, handler MessageHandler)
	Send(recipient telebot.Recipient, msg Message)
	Start()
}

type handlerRegistry struct {
}

type bot struct {
	telebot *telebot.Bot
}

func (b *bot) OnText(handler MessageHandler) {
	b.telebot.Handle(telebot.OnText, func(m *telebot.Message) {
		handler.Handle(b, m)
	})
}

func (b *bot) OnCommand(command string, handler MessageHandler) {
	b.telebot.Handle(command, func(m *telebot.Message) {
		handler.Handle(b, m)
	})
}

func (b *bot) Send(recipient telebot.Recipient, msg Message) {
	text, options := msg.Render()
	b.telebot.Send(recipient, text, options...)
}

func (b *bot) Start() {
	b.telebot.Start()
}

func NewBot(token string, timout time.Duration) (Bot, error) {
	poller := &telebot.LongPoller{Timeout: timout * time.Second}
	settings := telebot.Settings{Token: token, Poller: poller}

	telebot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	return &bot{telebot}, nil
}
