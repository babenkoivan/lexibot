package bot

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/user"
)

const (
	OnStart    string = "/start"
	OnHelp     string = "/help"
	OnSettings string = "/settings"
)

type MessageHandler interface {
	Handle(b Bot, m *telebot.Message)
}

type MessageHandlerFunc func(b Bot, m *telebot.Message)

func (h MessageHandlerFunc) Handle(b Bot, m *telebot.Message) {
	h(b, m)
}

type startHandler struct {
	bundle *i18n.Bundle
}

func (h *startHandler) Handle(b Bot, m *telebot.Message) {
	localizer := user.NewLocalizer(h.bundle, m.Sender.ID)
	localizeConfig := &i18n.LocalizeConfig{MessageID: "start"}
	b.Send(m.Sender, NewPlainTextMessage(localizer.MustLocalize(localizeConfig)))
}

func NewStartHandler(bundle *i18n.Bundle) *startHandler {
	return &startHandler{bundle}
}

func helpHandler(b Bot, m *telebot.Message) {
	b.Send(m.Sender, NewPlainTextMessage("#todo help"))
}

func NewHelpHandler() MessageHandler {
	return MessageHandlerFunc(helpHandler)
}

func settingsHandler(b Bot, m *telebot.Message) {
	b.Send(m.Sender, NewPlainTextMessage("#todo settings"))
}

func NewSettingsHandler() MessageHandler {
	return MessageHandlerFunc(settingsHandler)
}
