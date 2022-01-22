package app

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
)

const (
	OnStart    string = "/start"
	OnHelp     string = "/help"
	OnSettings string = "/settings"
)

func startHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &bot.PlainTextMessage{"app.start"})
}

func NewStartHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(startHandler)
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
