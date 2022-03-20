package translation

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
	"strings"
)

const (
	OnDelete             string = "/delete"
	AutoTranslationLimit int    = 50
)

type translateHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
	translator       Translator
}

func (h *translateHandler) Handle(b bot.Bot, msg *telebot.Message) {
	// if configuration hasn't been created ask to start the initial set up
	userSettings := h.settingsStore.FirstOrInit(msg.Sender.ID)

	if userSettings.LangUI == "" || userSettings.LangDict == "" {
		b.Send(msg.Sender, &LangErrorMessage{})
		return
	}

	text := strings.TrimSpace(msg.Text)

	// return error if the given text is already in the dictionary
	transl := h.translationStore.First(
		WithText(text),
		WithLangFrom(userSettings.LangDict),
		WithUserID(msg.Sender.ID),
	)

	if transl != nil {
		b.Send(msg.Sender, &ExistsErrorMessage{transl.Text, transl.Translation})
		return
	}

	// ask to enter a translation if the text is too long
	if len(text) > AutoTranslationLimit {
		b.Send(msg.Sender, &EnterTranslationMessage{Text: text})
		return
	}

	// translate the text and save the result
	translatedText, err := h.translator.Translate(text, userSettings.LangDict, userSettings.LangUI)
	if err != nil {
		b.Send(msg.Sender, &EnterTranslationMessage{Text: text})
		return
	}

	// ask to enter a translation if auto-translation is disabled
	if !userSettings.AutoTranslate {
		b.Send(msg.Sender, &EnterTranslationMessage{text, translatedText})
		return
	}

	// otherwise, save the translation
	h.translationStore.Save(&Translation{
		UserID:      msg.Sender.ID,
		Text:        text,
		Translation: translatedText,
		LangFrom:    userSettings.LangDict,
		LangTo:      userSettings.LangUI,
		Manual:      false,
	})

	b.Send(msg.Sender, &AddedToDictionaryMessage{text, translatedText})
}

func NewTranslateHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	translator Translator,
) *translateHandler {
	return &translateHandler{
		settingsStore,
		translationStore,
		translator,
	}
}

type addToDictionaryHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
}

func (h *addToDictionaryHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	text := strings.TrimSpace(msg.(*EnterTranslationMessage).Text)
	translatedText := strings.TrimSpace(re.Text)
	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)

	h.translationStore.Save(&Translation{
		UserID:      re.Sender.ID,
		Text:        text,
		Translation: translatedText,
		LangFrom:    userSettings.LangDict,
		LangTo:      userSettings.LangUI,
		Manual:      true,
	})

	b.Send(re.Sender, &AddedToDictionaryMessage{text, translatedText})
}

func NewAddToDictionaryHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
) *addToDictionaryHandler {
	return &addToDictionaryHandler{
		settingsStore,
		translationStore,
	}
}

func whatToDeleteHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &WhatToDeleteMessage{})
}

func NewWhatToDeleteHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(whatToDeleteHandler)
}

type deleteFromDictionaryHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
}

func (h *deleteFromDictionaryHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	text := strings.TrimSpace(re.Text)
	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)

	transl := h.translationStore.Find(
		WithTextOrTranslation(text),
		WithLangFrom(userSettings.LangDict),
		WithUserID(re.Sender.ID),
	)

	if len(transl) == 0 {
		b.Send(re.Sender, &NotFoundErrorMessage{text})
		return
	}

	for _, t := range transl {
		h.translationStore.Delete(t)
		b.Send(re.Sender, &DeletedFromDictionaryMessage{t.Text, t.Translation})
	}
}

func NewDeleteFromDictionaryHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
) *deleteFromDictionaryHandler {
	return &deleteFromDictionaryHandler{
		settingsStore,
		translationStore,
	}
}
