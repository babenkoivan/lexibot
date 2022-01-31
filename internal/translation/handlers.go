package translation

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
	"strings"
)

const autoTranslationLimit = 50

type enterTranslationHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
	translator       Translator
}

func (h *enterTranslationHandler) Handle(b bot.Bot, msg *telebot.Message) {
	settings := h.settingsStore.GetOrInit(msg.Sender.ID)

	if settings.LangUI == "" || settings.LangDict == "" {
		b.Send(msg.Sender, &SettingsErrorMessage{})
		return
	}

	text := strings.TrimSpace(msg.Text)

	if len(text) > autoTranslationLimit {
		b.Send(msg.Sender, &EnterTranslationMessage{Text: text})
		return
	}

	translatedText, err := h.translator.Translate(text, settings.LangDict, settings.LangUI)
	if err != nil {
		b.Send(msg.Sender, &EnterTranslationMessage{Text: text})
		return
	}

	translation := h.translationStore.GetOrCreate(TranslationFilter{
		Text:        text,
		Translation: translatedText,
		LangFrom:    settings.LangDict,
		LangTo:      settings.LangUI,
		Manual:      false,
	})

	if !settings.AutoTranslate {
		b.Send(msg.Sender, &EnterTranslationMessage{text, translation.Translation})
		return
	}

	if h.translationStore.IsAttached(translation.ID, msg.Sender.ID) {
		b.Send(msg.Sender, &ExistsErrorMessage{translation.Text, translation.Translation})
		return
	}

	h.translationStore.Attach(translation.ID, msg.Sender.ID)
	b.Send(msg.Sender, &NewTranslationMessage{text, translation.Translation})
}

func NewEnterTranslationHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	translator Translator,
) *enterTranslationHandler {
	return &enterTranslationHandler{
		settingsStore,
		translationStore,
		translator,
	}
}

type saveTranslationHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
}

func (h *saveTranslationHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	text := strings.TrimSpace(msg.(*EnterTranslationMessage).Text)
	translatedText := strings.TrimSpace(re.Text)
	settings := h.settingsStore.GetOrInit(re.Sender.ID)

	translation := h.translationStore.GetOrInit(TranslationFilter{
		Text:        text,
		Translation: translatedText,
		LangFrom:    settings.LangDict,
		LangTo:      settings.LangUI,
	})

	if translation.ID == 0 {
		translation.Manual = true
		h.translationStore.Create(translation)
	} else if h.translationStore.IsAttached(translation.ID, re.Sender.ID) {
		b.Send(re.Sender, &ExistsErrorMessage{translation.Text, translation.Translation})
		return
	}

	h.translationStore.Attach(translation.ID, re.Sender.ID)
	b.Send(re.Sender, &NewTranslationMessage{translation.Text, translation.Translation})
}

func NewSaveTranslationHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
) *saveTranslationHandler {
	return &saveTranslationHandler{
		settingsStore,
		translationStore,
	}
}
