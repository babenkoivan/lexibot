package translation

import (
	"golang.org/x/text/language"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"strings"
)

const (
	OnCancelTranslation string = "cancel_translation"
	OnSaveTranslation   string = "save_translation"
	OnDeleteTranslation string = "delete_translation"
)

type suggestTranslationHandler struct {
	translator Translator
}

func (h *suggestTranslationHandler) Handle(b bot.Bot, m *telebot.Message) {
	// todo take from the user config
	from := language.German
	to := language.English

	res, err := h.translator.Translate(from, to, m.Text)
	if err != nil {
		b.Send(m.Sender, bot.NewErrorMessage(err))
		return
	}

	b.Send(m.Sender, &selectTranslationMessage{m.Text, res})
}

func NewSuggestTranslationHandler(translator Translator) *suggestTranslationHandler {
	return &suggestTranslationHandler{translator: translator}
}

func cancelTranslationHandler(b bot.Bot, c *telebot.Callback) {
	b.Edit(bot.ExtractMessageSig(c.Message), bot.NewInfoMessage("The translation has been canceled"))
}

func NewCancelTranslationHandler() bot.CallbackHandler {
	return bot.CallbackHandlerFunc(cancelTranslationHandler)
}

type saveTranslationHandler struct {
	// todo data store
}

func (h *saveTranslationHandler) Handle(b bot.Bot, c *telebot.Callback) {
	data := strings.Split(c.Data, "|")
	text := data[0]
	translation := data[1]

	// todo save in the db

	b.Edit(bot.ExtractMessageSig(c.Message), &savedTranslationMessage{text, translation})
}

func NewSaveTranslationHandler() *saveTranslationHandler {
	return &saveTranslationHandler{}
}

func deleteTranslationHandler(b bot.Bot, c *telebot.Callback) {
	b.Edit(bot.ExtractMessageSig(c.Message), bot.NewInfoMessage("The translation has been deleted"))
}

func NewDeleteTranslationHandler() bot.CallbackHandler {
	return bot.CallbackHandlerFunc(deleteTranslationHandler)
}
