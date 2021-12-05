package telegram

import (
	"fmt"
	"golang.org/x/text/language"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/translations"
	"strings"
)

type MessageHandler interface {
	Handle(m *telebot.Message)
}

type CallbackHandler interface {
	Handle(c *telebot.Callback)
}

type translateHandler struct {
	translator translations.Translator
	bot        *telebot.Bot
}

func (h *translateHandler) Handle(m *telebot.Message) {
	// todo take from the user config
	from := language.German
	to := language.English

	res, err := h.translator.Translate(from, to, m.Text)
	if err != nil {
		h.bot.Send(m.Sender, &errorMessage{err})
		return
	}

	h.bot.Send(m.Sender, &selectTranslationMessage{m.Text, res})
}

func NewTranslateHandler(translator translations.Translator, bot *telebot.Bot) *translateHandler {
	return &translateHandler{translator: translator, bot: bot}
}

type cancelTranslationHandler struct {
	bot *telebot.Bot
}

func (h *cancelTranslationHandler) Handle(c *telebot.Callback) {
	h.bot.Delete(c.Message)
}

func NewCancelTranslationHandler(bot *telebot.Bot) *cancelTranslationHandler {
	return &cancelTranslationHandler{bot: bot}
}

type saveTranslationHandler struct {
	bot *telebot.Bot
}

func (h *saveTranslationHandler) Handle(c *telebot.Callback) {
	data := strings.Split(c.Data, "|")
	text := data[0]
	translation := data[1]

	// todo save in the db

	msg := fmt.Sprintf("%s → %s", text, translation)

	markup := &telebot.ReplyMarkup{}
	btn := markup.Data("✗", OnDeleteTranslation)
	row := telebot.Row{btn}
	markup.Inline(row)

	h.bot.Edit(c.Message, msg, markup)
}

func NewSaveTranslationHandler(bot *telebot.Bot) *saveTranslationHandler {
	return &saveTranslationHandler{bot: bot}
}

type deleteTranslationHandler struct {
	bot *telebot.Bot
}

func (h *deleteTranslationHandler) Handle(c *telebot.Callback) {
	// todo remove from the db
	h.bot.Delete(c.Message)
}

func NewDeleteTranslationHandler(bot *telebot.Bot) *deleteTranslationHandler {
	return &deleteTranslationHandler{bot: bot}
}
