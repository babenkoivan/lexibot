package settings

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/bot"
	"lexibot/internal/localization"
)

type SelectLangUIMessage struct{}

func (m *SelectLangUIMessage) Type() string {
	return "settings.selectLangUI"
}

func (m *SelectLangUIMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	var captions []string
	for _, lang := range SupportedLangUI() {
		captions = append(captions, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: lang}))
	}

	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "settings.langUI"})
	options = append(options, bot.WithReplyKeyboard(captions))

	return
}

type SelectLangDictMessage struct{}

func (m *SelectLangDictMessage) Type() string {
	return "settings.selectLangDict"
}

func (m *SelectLangDictMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	var captions []string
	for _, lang := range SupportedLangDict() {
		captions = append(captions, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: lang}))
	}

	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "settings.langDict"})
	options = append(options, bot.WithReplyKeyboard(captions))

	return
}

type EnableAutoTranslateMessage struct{}

func (m *EnableAutoTranslateMessage) Type() string {
	return "settings.enableAutoTranslate"
}

func (m *EnableAutoTranslateMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "settings.autoTranslate"})

	options = append(options, bot.WithReplyKeyboard([]string{
		localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "yes"}),
		localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "no"}),
	}))

	return
}

type EnterWordsPerTrainingMessage struct{}

func (m *EnterWordsPerTrainingMessage) Type() string {
	return "settings.enterWordsPerTraining"
}

func (m *EnterWordsPerTrainingMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "settings.wordsPerTraining"})
	options = append(options, bot.WithReplyKeyboard([]string{"5", "10", "20"}))
	return
}

type EnumErrorMessage struct {
	Value string
}

func (m *EnumErrorMessage) Type() string {
	return "settings.enumError"
}

func (m *EnumErrorMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "settings.enumError",
		TemplateData: map[string]interface{}{
			"Value": m.Value,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())

	return
}

type IntegerErrorMessage struct {
	MaxValue int
}

func (m *IntegerErrorMessage) Type() string {
	return "settings.integerError"
}

func (m *IntegerErrorMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "settings.integerError",
		TemplateData: map[string]interface{}{
			"MaxValue": m.MaxValue,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())

	return
}

type SuccessMessage struct{}

func (m *SuccessMessage) Type() string {
	return "settings.success"
}

func (m *SuccessMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "settings.success"})
	options = append(options, bot.WithoutReplyKeyboard())
	return
}
