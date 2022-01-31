package translation

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/bot"
)

type NewTranslationMessage struct {
	Text        string
	Translation string
}

func (m *NewTranslationMessage) Type() string {
	return "translation.new"
}

func (m *NewTranslationMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "translation.new",
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
	return "translation.enter"
}

func (m *EnterTranslationMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "translation.enter",
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

type SettingsErrorMessage struct{}

func (m *SettingsErrorMessage) Type() string {
	return "translation.settingsErr"
}

func (m *SettingsErrorMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "translation.settingsErr"})
	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type ExistsErrorMessage struct {
	Text        string
	Translation string
}

func (m *ExistsErrorMessage) Type() string {
	return "translation.existsErr"
}

func (m *ExistsErrorMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "translation.existsErr",
		TemplateData: map[string]interface{}{
			"Text":        m.Text,
			"Translation": m.Translation,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}
