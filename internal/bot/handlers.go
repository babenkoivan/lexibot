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
	Handle(bot Bot, msg *telebot.Message)
}

type ReplyHandler interface {
	Handle(bot Bot, re *telebot.Message, hm *HistoryMessage)
}

type MessageHandlerFunc func(bot Bot, msg *telebot.Message)

func (h MessageHandlerFunc) Handle(bot Bot, msg *telebot.Message) {
	h(bot, msg)
}

type ReplyHandlerFunc func(bot Bot, re *telebot.Message, hm *HistoryMessage)

func (h ReplyHandlerFunc) Handle(bot Bot, re *telebot.Message, hm *HistoryMessage) {
	h(bot, re, hm)
}

type startHandler struct {
	bundle *i18n.Bundle
}

func (h *startHandler) Handle(bot Bot, msg *telebot.Message) {
	localizer := user.NewLocalizer(h.bundle, msg.Sender.ID)
	localizeConfig := &i18n.LocalizeConfig{MessageID: "start"}
	bot.Send(msg.Chat, &PlainTextMessage{localizer.MustLocalize(localizeConfig)})
}

func NewStartHandler(bundle *i18n.Bundle) *startHandler {
	return &startHandler{bundle}
}

//type ReplyTest struct {
//
//}
//
//func (r *ReplyTest) Handle(bot Bot, re *telebot.Message, hm *HistoryMessage) {
//	msg := &PlainTextMessage{}
//	json.Unmarshal([]byte(hm.Content), re)
//
//	msg.Text = "re: " + re.Text
//	bot.Send(re.Chat, msg)
//}
