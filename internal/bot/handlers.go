package bot

import (
	"gopkg.in/tucnak/telebot.v2"
)

type MessageHandler interface {
	Handle(b Bot, msg *telebot.Message)
}

type ReplyHandler interface {
	Handle(b Bot, re *telebot.Message, hm *HistoryMessage)
}

type MessageHandlerFunc func(b Bot, msg *telebot.Message)

func (h MessageHandlerFunc) Handle(b Bot, msg *telebot.Message) {
	h(b, msg)
}

type ReplyHandlerFunc func(b Bot, re *telebot.Message, hm *HistoryMessage)

func (h ReplyHandlerFunc) Handle(b Bot, re *telebot.Message, hm *HistoryMessage) {
	h(b, re, hm)
}
