package translation_test

import (
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

		handler := translation.NewTranslateHandler(
			testkit.MockSettingsStore(t),
			testkit.MockTranslationStore(t),
			testkit.MockScoreStore(t),
			testkit.MockTranslator(t),
		)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: "bunt"})

		botSpy.AssertSent(user, &translation.LangErrorMessage{})
	})

	t.Run("translation already exists", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		langFrom, langTo := "de", "en"
		text, transl := "bunt", "colorful"

		settingsStoreMock := testkit.MockSettingsStore(t)
		settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
			return &settings.Settings{LangDict: langFrom, LangUI: langTo}
		})

		translationStoreMock := testkit.MockTranslationStore(t)
		translationStoreMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
			return &translation.Translation{Text: text, Translation: transl, LangFrom: langFrom, LangTo: langTo}
		})

		handler := translation.NewTranslateHandler(
			settingsStoreMock,
			translationStoreMock,
			testkit.MockScoreStore(t),
			testkit.MockTranslator(t),
		)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: text})

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
				return &settings.Settings{LangDict: langFrom, LangUI: langTo}
			})

			handler := translation.NewTranslateHandler(
				settingsStoreMock,
				testkit.MockTranslationStore(t),
				testkit.MockScoreStore(t),
				testkit.MockTranslator(t),
			)
			handler.Handle(botSpy, &telebot.Message{Sender: user, Text: text})

			botSpy.AssertSent(user, &translation.EnterTranslationMessage{Text: text})
		})
	}

	t.Run("auto-translation is disabled", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		langFrom, langTo := "de", "en"
		text, transl := "bunt", "colorful"

		settingsStoreMock := testkit.MockSettingsStore(t)
		settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
			return &settings.Settings{LangDict: langFrom, LangUI: langTo, AutoTranslate: false}
		})

		onFirstCounter := 0
		translationStoreMock := testkit.MockTranslationStore(t)
		translationStoreMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
			onFirstCounter++

			if onFirstCounter == 1 {
				return nil
			}

			return &translation.Translation{Text: text, Translation: transl, LangFrom: langFrom, LangTo: langTo}
		})

		scoreStoreMock := testkit.MockScoreStore(t)

		translatorMock := testkit.MockTranslator(t)
		translatorMock.OnTranslate(func(text, langFrom, langTo string) (string, error) {
			return transl, nil
		})

		handler := translation.NewTranslateHandler(
			settingsStoreMock,
			translationStoreMock,
			scoreStoreMock,
			translatorMock,
		)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: text})

		translationStoreMock.AssertNothingSaved()
		scoreStoreMock.AssertNothingSaved()
		botSpy.AssertSent(user, &translation.EnterTranslationMessage{Text: text, Suggestion: transl})
	})

	t.Run("auto-translation is enabled", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		langFrom, langTo := "de", "en"
		text, transl := "bunt", "colorful"

		settingsStoreMock := testkit.MockSettingsStore(t)
		settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
			return &settings.Settings{LangDict: langFrom, LangUI: langTo, AutoTranslate: true}
		})

		translationStoreMock := testkit.MockTranslationStore(t)
		scoreStoreMock := testkit.MockScoreStore(t)

		translatorMock := testkit.MockTranslator(t)
		translatorMock.OnTranslate(func(text, langFrom, langTo string) (string, error) {
			return transl, nil
		})

		handler := translation.NewTranslateHandler(
			settingsStoreMock,
			translationStoreMock,
			scoreStoreMock,
			translatorMock,
		)
		handler.Handle(botSpy, &telebot.Message{Sender: user, Text: text})

		newTransl := &translation.Translation{
			ID:          1,
			Text:        text,
			Translation: transl,
			LangFrom:    langFrom,
			LangTo:      langTo,
			Manual:      false,
		}

		translationStoreMock.AssertSaved(newTransl)
		scoreStoreMock.AssertSaved(newTransl.ID, user.ID)
		botSpy.AssertSent(user, &translation.AddedToDictionaryMessage{text, transl})
	})
}
