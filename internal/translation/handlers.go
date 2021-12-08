package translation

import (
	"golang.org/x/text/language"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"strconv"
	"strings"
)

const (
	OnCancelTranslation string = "cancel_translation"
	OnSaveTranslation   string = "save_translation"
	OnDeleteTranslation string = "delete_translation"
)

type suggestTranslationHandler struct {
	translator Translator
	store      Store
}

func (h *suggestTranslationHandler) Handle(b bot.Bot, m *telebot.Message) {
	// todo take from the user config
	text := strings.TrimSpace(m.Text)

	if h.store.Exists(text) {
		b.Send(m.Sender, bot.NewInfoMessage("A translation already exists"))
		return
	}

	from := language.German
	to := language.English

	res, err := h.translator.Translate(from, to, text)
	if err != nil {
		b.Send(m.Sender, bot.NewErrorMessage(err))
		return
	}

	if len(res) > 1 {
		b.Send(m.Sender, &selectTranslationMessage{text, res})
		return
	}

	translation := h.store.Create(text, res[0])
	b.Send(m.Sender, &savedTranslationMessage{translation})
}

func NewSuggestTranslationHandler(translator Translator, store Store) *suggestTranslationHandler {
	return &suggestTranslationHandler{translator: translator, store: store}
}

func cancelTranslationHandler(b bot.Bot, c *telebot.Callback) {
	b.Edit(bot.ExtractMessageSig(c.Message), bot.NewInfoMessage("The translation has been canceled"))
}

func NewCancelTranslationHandler() bot.CallbackHandler {
	return bot.CallbackHandlerFunc(cancelTranslationHandler)
}

type saveTranslationHandler struct {
	store Store
}

func (h *saveTranslationHandler) Handle(b bot.Bot, c *telebot.Callback) {
	data := strings.Split(c.Data, "|")
	translation := h.store.Create(data[0], data[1])
	b.Edit(bot.ExtractMessageSig(c.Message), &savedTranslationMessage{translation})
}

func NewSaveTranslationHandler(store Store) *saveTranslationHandler {
	return &saveTranslationHandler{store: store}
}

type deleteTranslationHandler struct {
	store Store
}

func (h *deleteTranslationHandler) Handle(b bot.Bot, c *telebot.Callback) {
	ID, _ := strconv.ParseUint(c.Data, 10, 64)
	h.store.Delete(ID)
	b.Edit(bot.ExtractMessageSig(c.Message), bot.NewInfoMessage("The translation has been deleted"))
}

func NewDeleteTranslationHandler(store Store) *deleteTranslationHandler {
	return &deleteTranslationHandler{store: store}
}
