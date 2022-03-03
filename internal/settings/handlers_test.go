package settings_test

import (
	"fmt"
	"golang.org/x/text/language"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/settings"
	"lexibot/internal/testkit"
	"strconv"
	"testing"
)

func TestSettingsHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	botSpy := testkit.NewBotSpy(t)

	handler := settings.NewSettingsHandler()
	handler.Handle(botSpy, &telebot.Message{Sender: user})

	botSpy.AssertSent(user, &settings.EnableAutoTranslateMessage{})
}

func TestSaveAutoTranslateHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	msg := &settings.EnableAutoTranslateMessage{}
	localizerFactory := testkit.MockLocalizerFactory(t, language.English)

	t.Run("unexpected answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		settingsStoreMock := testkit.MockSettingsStore(t)

		handler := settings.NewSaveAutoTranslateHandler(localizerFactory, settingsStoreMock)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: "foo"}, msg)

		settingsStoreMock.AssertNothingSaved()
		botSpy.AssertSent(user, &settings.EnumErrorMessage{"foo"})
		botSpy.AssertSent(user, msg)
	})

	t.Run("expected answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		settingsStoreMock := testkit.MockSettingsStore(t)

		handler := settings.NewSaveAutoTranslateHandler(localizerFactory, settingsStoreMock)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: "yes"}, msg)

		settingsStoreMock.AssertSaved(&settings.Settings{UserID: user.ID, AutoTranslate: true})
		botSpy.AssertSent(user, &settings.EnterWordsPerTrainingMessage{})
	})
}

func TestSaveWordsPerTrainingHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	msg := &settings.EnterWordsPerTrainingMessage{}

	for _, answer := range []string{"", "foo", "-1", strconv.Itoa(settings.MaxWordsPerTraining + 1)} {
		t.Run(fmt.Sprintf("unexpected answer: %q", answer), func(t *testing.T) {
			botSpy := testkit.NewBotSpy(t)
			settingsStoreMock := testkit.MockSettingsStore(t)

			handler := settings.NewSaveWordsPerTrainingHandler(settingsStoreMock)
			handler.Handle(botSpy, &telebot.Message{Sender: user, Text: answer}, msg)

			settingsStoreMock.AssertNothingSaved()
			botSpy.AssertSent(user, &settings.IntegerErrorMessage{settings.MaxWordsPerTraining})
			botSpy.AssertSent(user, msg)
		})
	}

	t.Run("expected answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		settingsStoreMock := testkit.MockSettingsStore(t)

		handler := settings.NewSaveWordsPerTrainingHandler(settingsStoreMock)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: "10"}, msg)

		settingsStoreMock.AssertSaved(&settings.Settings{UserID: user.ID, WordsPerTraining: 10})
		botSpy.AssertSent(user, &settings.SuccessMessage{})
	})
}

func TestSaveLangUIHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	msg := &settings.SelectLangUIMessage{}
	localizerFactory := testkit.MockLocalizerFactory(t, language.English)

	t.Run("unexpected answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		settingsStoreMock := testkit.MockSettingsStore(t)

		handler := settings.NewSaveLangUIHandler(localizerFactory, settingsStoreMock)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: "foo"}, msg)

		settingsStoreMock.AssertNothingSaved()
		botSpy.AssertSent(user, &settings.EnumErrorMessage{"foo"})
		botSpy.AssertSent(user, msg)
	})

	t.Run("expected answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		settingsStoreMock := testkit.MockSettingsStore(t)

		handler := settings.NewSaveLangUIHandler(localizerFactory, settingsStoreMock)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: "english"}, msg)

		settingsStoreMock.AssertSaved(&settings.Settings{UserID: user.ID, LangUI: language.English.String()})
		botSpy.AssertSent(user, &settings.SelectLangDictMessage{})
	})
}

func TestSaveLangDictHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	msg := &settings.SelectLangDictMessage{}
	localizerFactory := testkit.MockLocalizerFactory(t, language.English)

	t.Run("unexpected answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		settingsStoreMock := testkit.MockSettingsStore(t)

		handler := settings.NewSaveLangDictHandler(localizerFactory, settingsStoreMock)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: "foo"}, msg)

		settingsStoreMock.AssertNothingSaved()
		botSpy.AssertSent(user, &settings.EnumErrorMessage{"foo"})
		botSpy.AssertSent(user, msg)
	})

	t.Run("expected answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		settingsStoreMock := testkit.MockSettingsStore(t)

		handler := settings.NewSaveLangDictHandler(localizerFactory, settingsStoreMock)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: "english"}, msg)

		settingsStoreMock.AssertSaved(&settings.Settings{UserID: user.ID, LangDict: language.English.String()})
		botSpy.AssertSent(user, &settings.SuccessMessage{})
	})
}
