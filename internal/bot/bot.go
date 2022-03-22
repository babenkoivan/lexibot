package bot

import (
	"encoding/json"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/localization"
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
	telebot          *telebot.Bot
	handlers         *handlers
	localizerFactory localization.LocalizerFactory
	historyStore     HistoryStore
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
	text, options := msg.Render(b.localizerFactory.New(to.ID))
	options = append(options, telebot.ModeHTML)
	// todo error handling
	b.telebot.Send(to, text, options...)

	hm := newHistoryMessage(to.ID, msg)
	b.historyStore.Save(hm)
}

func (b *bot) Start() {
	b.telebot.Handle(telebot.OnText, func(re *telebot.Message) {
		hm := b.historyStore.Last(re.Sender.ID)

		for msg, handler := range b.handlers.replies {
			if hm == nil || hm.Type != msg.Type() {
				continue
			}

			// todo error handling
			json.Unmarshal([]byte(hm.Content), msg)
			handler.Handle(b, re, msg)

			return
		}

		b.handlers.message.Handle(b, re)
	})

	for command, handler := range b.handlers.commands {
		h := handler

		b.telebot.Handle(command, func(msg *telebot.Message) {
			h.Handle(b, msg)
		})
	}

	b.telebot.Start()
}

func NewBot(
	token string,
	timout time.Duration,
	localizerFactory localization.LocalizerFactory,
	historyStore HistoryStore,
) (Bot, error) {
	poller := &telebot.LongPoller{Timeout: timout}
	settings := telebot.Settings{Token: token, Poller: poller}

	telebot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	handlers := newHandlers()

	return &bot{
		telebot,
		handlers,
		localizerFactory,
		historyStore,
	}, nil
}
