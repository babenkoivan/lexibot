package translation_test

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/settings"
	"lexibot/internal/testkit"
	"lexibot/internal/translation"
	"strings"
	"testing"
)

func TestTranslateHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}

	t.Run("languages are not configured", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)

		translation.NewTranslateHandler(
			testkit.MockSettingsStore(t),
			testkit.MockTranslationStore(t),
			testkit.MockTranslator(t),
		).Handle(botSpy, &telebot.Message{Sender: user, Text: "bunt"})

		botSpy.AssertSent(user, &translation.LangErrorMessage{})
	})

	t.Run("translation already exists", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		langFrom, langTo := "de", "en"
		text, transl := "bunt", "colorful"

		settingsStoreMock := testkit.MockSettingsStore(t)
		settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
			assert.Equal(t, user.ID, userID)
			return &settings.Settings{LangDict: langFrom, LangUI: langTo}
		})

		translationStoreMock := testkit.MockTranslationStore(t)
		translationStoreMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
			testkit.AssertEqualTranslationQuery(t, []translation.TranslationQueryCond{
				translation.WithText(text),
				translation.WithLangFrom(langFrom),
				translation.WithUserID(user.ID),
			}, conds)

			return &translation.Translation{
				UserID:      user.ID,
				Text:        text,
				Translation: transl,
				LangFrom:    langFrom,
				LangTo:      langTo,
			}
		})

		translation.NewTranslateHandler(
			settingsStoreMock,
			translationStoreMock,
			testkit.MockTranslator(t),
		).Handle(botSpy, &telebot.Message{Sender: user, Text: text})

		botSpy.AssertSent(user, &translation.ExistsErrorMessage{text, transl})
	})

	for n, text := range map[string]string{
		"text is too long":           strings.Repeat("a", translation.AutoTranslationLimit+1),
		"text can not be translated": "bunt",
	} {
		t.Run(n, func(t *testing.T) {
			botSpy := testkit.NewBotSpy(t)
			langFrom, langTo := "de", "en"

			settingsStoreMock := testkit.MockSettingsStore(t)
			settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
				assert.Equal(t, user.ID, userID)
				return &settings.Settings{LangDict: langFrom, LangUI: langTo}
			})

			translation.NewTranslateHandler(
				settingsStoreMock,
				testkit.MockTranslationStore(t),
				testkit.MockTranslator(t),
			).Handle(botSpy, &telebot.Message{Sender: user, Text: text})

			botSpy.AssertSent(user, &translation.EnterTranslationMessage{Text: text})
		})
	}

	t.Run("auto-translation is disabled", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		langFrom, langTo := "de", "en"
		text, transl := "bunt", "colorful"

		settingsStoreMock := testkit.MockSettingsStore(t)
		settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
			assert.Equal(t, user.ID, userID)
			return &settings.Settings{LangDict: langFrom, LangUI: langTo, AutoTranslate: false}
		})

		translationStoreMock := testkit.MockTranslationStore(t)
		translationStoreMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
			testkit.AssertEqualTranslationQuery(t, []translation.TranslationQueryCond{
				translation.WithText(text),
				translation.WithLangFrom(langFrom),
				translation.WithUserID(user.ID),
			}, conds)

			return nil
		})

		translatorMock := testkit.MockTranslator(t)
		translatorMock.OnTranslate(func(text, langFrom, langTo string) (string, error) {
			return transl, nil
		})

		translation.NewTranslateHandler(
			settingsStoreMock,
			testkit.MockTranslationStore(t),
			translatorMock,
		).Handle(botSpy, &telebot.Message{Sender: user, Text: text})

		translationStoreMock.AssertNothingSaved()
		botSpy.AssertSent(user, &translation.EnterTranslationMessage{Text: text, Suggestion: transl})
	})

	t.Run("auto-translation is enabled", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		langFrom, langTo := "de", "en"
		text, transl := "bunt", "colorful"

		settingsStoreMock := testkit.MockSettingsStore(t)
		settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
			assert.Equal(t, user.ID, userID)
			return &settings.Settings{LangDict: langFrom, LangUI: langTo, AutoTranslate: true}
		})

		translationStoreMock := testkit.MockTranslationStore(t)
		translationStoreMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
			testkit.AssertEqualTranslationQuery(t, []translation.TranslationQueryCond{
				translation.WithText(text),
				translation.WithLangFrom(langFrom),
				translation.WithUserID(user.ID),
			}, conds)

			return nil
		})

		translatorMock := testkit.MockTranslator(t)
		translatorMock.OnTranslate(func(text, langFrom, langTo string) (string, error) {
			return transl, nil
		})

		translation.NewTranslateHandler(
			settingsStoreMock,
			translationStoreMock,
			translatorMock,
		).Handle(botSpy, &telebot.Message{Sender: user, Text: text})

		newTransl := &translation.Translation{
			ID:          1,
			UserID:      user.ID,
			Text:        text,
			Translation: transl,
			LangFrom:    langFrom,
			LangTo:      langTo,
			Manual:      false,
			Score:       0,
		}

		translationStoreMock.AssertSaved(newTransl)
		botSpy.AssertSent(user, &translation.AddedToDictionaryMessage{text, transl})
	})
}

func TestAddToDictionaryHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	msg := &translation.EnterTranslationMessage{"bunt", "colorful"}
	langFrom, langTo := "de", "en"

	settingsStoreMock := testkit.MockSettingsStore(t)
	settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
		assert.Equal(t, user.ID, userID)
		return &settings.Settings{LangDict: langFrom, LangUI: langTo}
	})

	translationStoreMock := testkit.MockTranslationStore(t)
	botSpy := testkit.NewBotSpy(t)

	translation.NewAddToDictionaryHandler(
		settingsStoreMock,
		translationStoreMock,
	).Handle(botSpy, &telebot.Message{Sender: user, Text: msg.Suggestion}, msg)

	newTransl := &translation.Translation{
		ID:          1,
		UserID:      user.ID,
		Text:        msg.Text,
		Translation: msg.Suggestion,
		LangFrom:    langFrom,
		LangTo:      langTo,
		Manual:      true,
		Score:       0,
	}

	translationStoreMock.AssertSaved(newTransl)
	botSpy.AssertSent(user, &translation.AddedToDictionaryMessage{newTransl.Text, newTransl.Translation})
}

func TestWhatToDeleteHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	botSpy := testkit.NewBotSpy(t)

	translation.NewWhatToDeleteHandler().Handle(botSpy, &telebot.Message{Sender: user})

	botSpy.AssertSent(user, &translation.WhatToDeleteMessage{})
}

func TestDeleteFromDictionaryHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	msg := &translation.WhatToDeleteMessage{}
	text := "bunt"

	settingsStoreMock := testkit.MockSettingsStore(t)
	settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
		assert.Equal(t, user.ID, userID)
		return &settings.Settings{LangDict: "de", LangUI: "en"}
	})

	t.Run("translations are not found", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)

		translation.NewDeleteFromDictionaryHandler(
			settingsStoreMock,
			testkit.MockTranslationStore(t),
		).Handle(botSpy, &telebot.Message{Sender: user, Text: "bunt"}, msg)

		botSpy.AssertSent(user, &translation.NotFoundErrorMessage{text})
	})

	t.Run("translations are found", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)

		existingTransl := []*translation.Translation{
			{
				ID:          1,
				UserID:      user.ID,
				Text:        "bunt",
				Translation: "colorful",
				LangFrom:    "de",
				LangTo:      "en",
				Manual:      false,
			},
			{
				ID:          2,
				UserID:      user.ID,
				Text:        "bunt",
				Translation: "яркий",
				LangFrom:    "de",
				LangTo:      "ru",
				Manual:      false,
			},
		}

		translationStoreMock := testkit.MockTranslationStore(t)
		translationStoreMock.OnFind(func(conds ...translation.TranslationQueryCond) []*translation.Translation {
			testkit.AssertEqualTranslationQuery(t, []translation.TranslationQueryCond{
				translation.WithTextOrTranslation(existingTransl[0].Text),
				translation.WithLangFrom(existingTransl[0].LangFrom),
				translation.WithUserID(user.ID),
			}, conds)

			return existingTransl
		})

		translation.NewDeleteFromDictionaryHandler(
			settingsStoreMock,
			translationStoreMock,
		).Handle(botSpy, &telebot.Message{Sender: user, Text: existingTransl[0].Text}, msg)

		for _, transl := range existingTransl {
			translationStoreMock.AssertDeleted(transl)
			botSpy.AssertSent(user, &translation.DeletedFromDictionaryMessage{transl.Text, transl.Translation})
		}
	})
}
