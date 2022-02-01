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
	// if configuration hasn't been created ask to start the initial set up
	settings := h.settingsStore.FirstOrInit(msg.Sender.ID)

	if settings.LangUI == "" || settings.LangDict == "" {
		b.Send(msg.Sender, &SettingsErrorMessage{})
		return
	}

	text := strings.TrimSpace(msg.Text)

	// return error if the given text is already in the dictionary
	translation := h.translationStore.First(
		WithText(text),
		WithLangFrom(settings.LangDict),
		WithLangTo(settings.LangUI),
		WithUserID(msg.Sender.ID),
	)

	if translation != nil {
		b.Send(msg.Sender, &ExistsErrorMessage{translation.Text, translation.Translation})
		return
	}

	// ask to enter a translation if the text is too long
	if len(text) > autoTranslationLimit {
		b.Send(msg.Sender, &EnterTranslationMessage{Text: text})
		return
	}

	// translate the text and save the result
	translatedText, err := h.translator.Translate(text, settings.LangDict, settings.LangUI)
	if err != nil {
		b.Send(msg.Sender, &EnterTranslationMessage{Text: text})
		return
	}

	translation = h.translationStore.First(
		WithText(text),
		WithTranslation(translatedText),
		WithLangFrom(settings.LangDict),
		WithLangTo(settings.LangUI),
		WithManual(false),
	)

	if translation == nil {
		translation = h.translationStore.Save(&Translation{
			Text:        text,
			Translation: translatedText,
			LangFrom:    settings.LangDict,
			LangTo:      settings.LangUI,
			Manual:      false,
		})
	}

	// ask to enter a translation if auto-translation is disabled
	if !settings.AutoTranslate {
		b.Send(msg.Sender, &EnterTranslationMessage{text, translation.Translation})
		return
	}

	// save the translation otherwise
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
	settings := h.settingsStore.FirstOrInit(re.Sender.ID)

	translation := h.translationStore.First(
		WithText(text),
		WithTranslation(translatedText),
		WithLangFrom(settings.LangDict),
		WithLangTo(settings.LangUI),
	)

	if translation == nil {
		translation = h.translationStore.Save(&Translation{
			Text:        text,
			Translation: translatedText,
			LangFrom:    settings.LangDict,
			LangTo:      settings.LangUI,
			Manual:      true,
		})
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
