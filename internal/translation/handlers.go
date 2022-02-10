package translation

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
	"strings"
)

const (
	OnDelete             string = "/delete"
	autoTranslationLimit int    = 50
)

type translateHandler struct {
	settingsStore    settings.SettingsStore
	translationStore TranslationStore
	scoreStore       ScoreStore
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

	transl = h.translationStore.First(
		WithText(text),
		WithTranslation(translatedText),
		WithLangFrom(userSettings.LangDict),
		WithLangTo(userSettings.LangUI),
		WithManual(false),
	)

	if transl == nil {
		transl = h.translationStore.Save(&Translation{
			Text:        text,
			Translation: translatedText,
			LangFrom:    userSettings.LangDict,
			LangTo:      userSettings.LangUI,
			Manual:      false,
		})
	}

	// ask to enter a translation if auto-translation is disabled
	if !userSettings.AutoTranslate {
		b.Send(msg.Sender, &EnterTranslationMessage{text, transl.Translation})
		return
	}

	// otherwise, save the translation
	h.scoreStore.Create(transl.ID, msg.Sender.ID)
	b.Send(msg.Sender, &AddedToDictionaryMessage{text, transl.Translation})
}

func NewTranslateHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	scoreStore ScoreStore,
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
	scoreStore       ScoreStore
}

func (h *addToDictionaryHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	text := strings.TrimSpace(msg.(*EnterTranslationMessage).Text)
	translatedText := strings.TrimSpace(re.Text)
	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)

	transl := h.translationStore.First(
		WithText(text),
		WithTranslation(translatedText),
		WithLangFrom(userSettings.LangDict),
		WithLangTo(userSettings.LangUI),
	)

	if transl == nil {
		transl = h.translationStore.Save(&Translation{
			Text:        text,
			Translation: translatedText,
			LangFrom:    userSettings.LangDict,
			LangTo:      userSettings.LangUI,
			Manual:      true,
		})
	}

	h.scoreStore.Create(transl.ID, re.Sender.ID)
	b.Send(re.Sender, &AddedToDictionaryMessage{transl.Text, transl.Translation})
}

func NewAddToDictionaryHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	scoreStore ScoreStore,
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
	scoreStore       ScoreStore
}

func (h *deleteFromDictionaryHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	text := strings.TrimSpace(re.Text)
	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)

	transl := h.translationStore.First(
		WithTextOrTranslation(text),
		WithLangFrom(userSettings.LangDict),
		WithLangTo(userSettings.LangUI),
		WithUserID(re.Sender.ID),
	)

	if transl == nil {
		b.Send(re.Sender, &NotFoundErrorMessage{text})
		return
	}

	h.scoreStore.Delete(transl.ID, re.Sender.ID)
	b.Send(re.Sender, &DeletedFromDictionaryMessage{transl.Text, transl.Translation})
}

func NewDeleteFromDictionaryHandler(
	settingsStore settings.SettingsStore,
	translationStore TranslationStore,
	scoreStore ScoreStore,
) *deleteFromDictionaryHandler {
	return &deleteFromDictionaryHandler{
		settingsStore,
		translationStore,
		scoreStore,
	}
}
