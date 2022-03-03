package translation

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/bot"
	"lexibot/internal/localization"
)

type AddedToDictionaryMessage struct {
	Text        string
	Translation string
}

func (m *AddedToDictionaryMessage) Type() string {
	return "translation.addedToDictionary"
}

func (m *AddedToDictionaryMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "translation.added",
		TemplateData: map[string]interface{}{
			"Text":        m.Text,
			"Translation": m.Translation,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type DeletedFromDictionaryMessage struct {
	Text        string
	Translation string
}

func (m *DeletedFromDictionaryMessage) Type() string {
	return "translation.deletedFromDictionary"
}

func (m *DeletedFromDictionaryMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "translation.deleted",
		TemplateData: map[string]interface{}{
			"Text":        m.Text,
			"Translation": m.Translation,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type EnterTranslationMessage struct {
	Text       string
	Suggestion string
}

func (m *EnterTranslationMessage) Type() string {
	return "translation.enterTranslation"
}

func (m *EnterTranslationMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "translation.enterTranslation",
		TemplateData: map[string]interface{}{
			"Text": m.Text,
		},
	})

	if m.Suggestion != "" {
		options = append(options, bot.WithReplyKeyboard([]string{m.Suggestion}))
	} else {
		options = append(options, bot.WithoutReplyKeyboard())
	}

	return
}

type WhatToDeleteMessage struct{}

func (m *WhatToDeleteMessage) Type() string {
	return "translation.whatToDelete"
}

func (m *WhatToDeleteMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "translation.whatToDelete"})
	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type LangErrorMessage struct{}

func (m *LangErrorMessage) Type() string {
	return "translation.langError"
}

func (m *LangErrorMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "translation.langError"})
	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type ExistsErrorMessage struct {
	Text        string
	Translation string
}

func (m *ExistsErrorMessage) Type() string {
	return "translation.existsError"
}

func (m *ExistsErrorMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "translation.existsError",
		TemplateData: map[string]interface{}{
			"Text":        m.Text,
			"Translation": m.Translation,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type NotFoundErrorMessage struct {
	TextOrTranslation string
}

func (m *NotFoundErrorMessage) Type() string {
	return "translation.notFoundError"
}

func (m *NotFoundErrorMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "translation.notFoundError",
		TemplateData: map[string]interface{}{
			"TextOrTranslation": m.TextOrTranslation,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}
