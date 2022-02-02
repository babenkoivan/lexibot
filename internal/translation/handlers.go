package translation

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
	"strings"
)

const (
	OnDelete             = "/delete"
	autoTranslationLimit = 50
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
		b.Send(msg.Sender, &SettingsErrorMessage{})
		return
	}

	text := strings.TrimSpace(msg.Text)

	// return error if the given text is already in the dictionary
	translation := h.translationStore.First(
		WithText(text),
		WithLangFrom(userSettings.LangDict),
		WithLangTo(userSettings.LangUI),
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
	translatedText, err := h.translator.Translate(text, userSettings.LangDict, userSettings.LangUI)
	if err != nil {
		b.Send(msg.Sender, &EnterTranslationMessage{Text: text})
		return
	}

	translation = h.translationStore.First(
		WithText(text),
		WithTranslation(translatedText),
		WithLangFrom(userSettings.LangDict),
		WithLangTo(userSettings.LangUI),
		WithManual(false),
	)

	if translation == nil {
		translation = h.translationStore.Save(&Translation{
			Text:        text,
			Translation: translatedText,
			LangFrom:    userSettings.LangDict,
			LangTo:      userSettings.LangUI,
			Manual:      false,
		})
	}

	// ask to enter a translation if auto-translation is disabled
	if !userSettings.AutoTranslate {
		b.Send(msg.Sender, &EnterTranslationMessage{text, translation.Translation})
		return
	}

	// save the translation otherwise
	h.translationStore.Attach(translation.ID, msg.Sender.ID)
	b.Send(msg.Sender, &AddedToDictionaryMessage{text, translation.Translation})
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

	translation := h.translationStore.First(
		WithText(text),
		WithTranslation(translatedText),
		WithLangFrom(userSettings.LangDict),
		WithLangTo(userSettings.LangUI),
	)

	if translation == nil {
		translation = h.translationStore.Save(&Translation{
			Text:        text,
			Translation: translatedText,
			LangFrom:    userSettings.LangDict,
			LangTo:      userSettings.LangUI,
			Manual:      true,
		})
	}

	h.translationStore.Attach(translation.ID, re.Sender.ID)
	b.Send(re.Sender, &AddedToDictionaryMessage{translation.Text, translation.Translation})
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

type deleteFromDictionaryIndirectHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
}

func (h *deleteFromDictionaryIndirectHandler) Handle(b bot.Bot, msg *telebot.Message) {
	text := strings.TrimSpace(msg.Payload)

	if text == "" {
		b.Send(msg.Sender, &WhatToDeleteMessage{})
		return
	}

	deleteFromDictionary(b, h.settingsStore, h.translationStore, msg.Sender, text)
}

func NewDeleteFromDictionaryIndirectHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
) *deleteFromDictionaryIndirectHandler {
	return &deleteFromDictionaryIndirectHandler{
		settingsStore:    settingsStore,
		translationStore: translationStore,
	}
}

type deleteFromDictionaryDirectHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
}

func (h *deleteFromDictionaryDirectHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	text := strings.TrimSpace(re.Text)
	deleteFromDictionary(b, h.settingsStore, h.translationStore, re.Sender, text)
}

func NewDeleteFromDictionaryDirectHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
) *deleteFromDictionaryDirectHandler {
	return &deleteFromDictionaryDirectHandler{
		settingsStore:    settingsStore,
		translationStore: translationStore,
	}
}

func deleteFromDictionary(
	b bot.Bot,
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	user *telebot.User,
	text string,
) {
	userSettings := settingsStore.FirstOrInit(user.ID)

	translation := translationStore.First(
		WithTextOrTranslation(text),
		WithLangFrom(userSettings.LangDict),
		WithLangTo(userSettings.LangUI),
		WithUserID(user.ID),
	)

	if translation == nil {
		b.Send(user, &NotFoundErrorMessage{text})
		return
	}

	translationStore.Detach(translation.ID, user.ID)
	b.Send(user, &DeletedFromDictionaryMessage{translation.Text, translation.Translation})
}
