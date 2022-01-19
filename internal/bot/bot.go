package bot

import (
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

type Bot interface {
	OnMessage(handler MessageHandler)
	OnReply(msg Message, handler ReplyHandler)
	OnCommand(command string, handler MessageHandler)
	Send(chat *telebot.Chat, msg Message)
	Start()
}

type handlers struct {
	message  MessageHandler
	replies  map[Message]ReplyHandler
	commands map[string]MessageHandler
}

func newHandlers() *handlers {
	return &handlers{
		replies:  map[Message]ReplyHandler{},
		commands: map[string]MessageHandler{},
	}
}

type bot struct {
	telebot      *telebot.Bot
	handlers     *handlers
	historyStore HistoryStore
}

func (b *bot) OnMessage(handler MessageHandler) {
	b.handlers.message = handler
}

func (b *bot) OnReply(msg Message, handler ReplyHandler) {
	b.handlers.replies[msg] = handler
}

func (b *bot) OnCommand(command string, handler MessageHandler) {
	b.handlers.commands[command] = handler
}

func (b *bot) Send(chat *telebot.Chat, msg Message) {
	text, options := msg.Render()
	b.telebot.Send(chat, text, options...)

	hm := NewHistoryMessage(chat, msg)
	b.historyStore.Save(hm)
}

func (b *bot) Start() {
	b.telebot.Handle(telebot.OnText, func(re *telebot.Message) {
		hm := b.historyStore.LastMessage(re.Chat.ID)

		for msg, handler := range b.handlers.replies {
			if hm.Type != msg.Type() {
				continue
			}

			handler.Handle(b, re, hm)
			return
		}

		b.handlers.message.Handle(b, re)
	})

	for command, handler := range b.handlers.commands {
		b.telebot.Handle(command, func(msg *telebot.Message) {
			handler.Handle(b, msg)
		})
	}

	b.telebot.Start()
}

func NewBot(token string, timout time.Duration, historyStore HistoryStore) (Bot, error) {
	poller := &telebot.LongPoller{Timeout: timout * time.Second}
	settings := telebot.Settings{Token: token, Poller: poller}

	telebot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	handlers := newHandlers()

	return &bot{telebot, handlers, historyStore}, nil
}
