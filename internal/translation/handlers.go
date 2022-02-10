package translation

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
	"lexibot/internal/training"
	"strings"
)

const (
	OnDelete             string = "/delete"
	autoTranslationLimit int    = 50
)

type translateHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
	scoreStore       training.ScoreStore
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
		WhereText(text),
		WhereLangFrom(userSettings.LangDict),
		WhereUserID(msg.Sender.ID),
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
		WhereText(text),
		WhereTranslation(translatedText),
		WhereLangFrom(userSettings.LangDict),
		WhereLangTo(userSettings.LangUI),
		WhereManual(false),
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

	// otherwise, save the translation
	h.scoreStore.Create(translation.ID, msg.Sender.ID)
	b.Send(msg.Sender, &AddedToDictionaryMessage{text, translation.Translation})
}

func NewTranslateHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	scoreStore training.ScoreStore,
	translator Translator,
) *translateHandler {
	return &translateHandler{
		settingsStore,
		translationStore,
		scoreStore,
		translator,
	}
}

type addToDictionaryHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
	scoreStore       training.ScoreStore
}

func (h *addToDictionaryHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	text := strings.TrimSpace(msg.(*EnterTranslationMessage).Text)
	translatedText := strings.TrimSpace(re.Text)
	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)

	translation := h.translationStore.First(
		WhereText(text),
		WhereTranslation(translatedText),
		WhereLangFrom(userSettings.LangDict),
		WhereLangTo(userSettings.LangUI),
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

	h.scoreStore.Create(translation.ID, re.Sender.ID)
	b.Send(re.Sender, &AddedToDictionaryMessage{translation.Text, translation.Translation})
}

func NewAddToDictionaryHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	scoreStore training.ScoreStore,
) *addToDictionaryHandler {
	return &addToDictionaryHandler{
		settingsStore,
		translationStore,
		scoreStore,
	}
}

func clarifyWhatToDeleteHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &WhatToDeleteMessage{})
}

func NewClarifyWhatToDeleteHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(clarifyWhatToDeleteHandler)
}

type deleteFromDictionaryHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
	scoreStore       training.ScoreStore
}

func (h *deleteFromDictionaryHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	text := strings.TrimSpace(re.Text)
	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)

	translation := h.translationStore.First(
		WhereTextOrTranslation(text),
		WhereLangFrom(userSettings.LangDict),
		WhereLangTo(userSettings.LangUI),
		WhereUserID(re.Sender.ID),
	)

	if translation == nil {
		b.Send(re.Sender, &NotFoundErrorMessage{text})
		return
	}

	h.scoreStore.Delete(translation.ID, re.Sender.ID)
	b.Send(re.Sender, &DeletedFromDictionaryMessage{translation.Text, translation.Translation})
}

func NewDeleteFromDictionaryHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	scoreStore training.ScoreStore,
) *deleteFromDictionaryHandler {
	return &deleteFromDictionaryHandler{
		settingsStore,
		translationStore,
		scoreStore,
	}
}
