package bot

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/user"
	"time"
)

type Bot interface {
	OnMessage(handler MessageHandler)
	OnReply(msg Message, handler ReplyHandler)
	OnCommand(command string, handler MessageHandler)
	Send(to *telebot.User, msg Message)
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
	localization *i18n.Bundle
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

func (b *bot) Send(to *telebot.User, msg Message) {
	localizeConfig, replyMarkup := msg.Render()

	localizer := user.NewLocalizer(b.localization, to.ID)
	text := localizer.MustLocalize(localizeConfig)

	if replyMarkup == nil {
		replyMarkup = &telebot.ReplyMarkup{ReplyKeyboardRemove: true}
	}

	b.telebot.Send(to, text, replyMarkup)

	hm := MakeHistoryMessage(to.ID, msg)
	b.historyStore.Save(hm)
}

func (b *bot) Start() {
	b.telebot.Handle(telebot.OnText, func(re *telebot.Message) {
		hm := b.historyStore.LastMessage(re.Sender.ID)

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

func NewBot(
	token string,
	timout time.Duration,
	localization *i18n.Bundle,
	historyStore HistoryStore,
) (Bot, error) {
	poller := &telebot.LongPoller{Timeout: timout * time.Second}
	settings := telebot.Settings{Token: token, Poller: poller}

	telebot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	handlers := newHandlers()

	return &bot{
		telebot,
		handlers,
		localization,
		historyStore,
	}, nil
}
