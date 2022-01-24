package config

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/bot"
)

type SelectLangUIMessage struct{}

func (m *SelectLangUIMessage) Type() string {
	return "config.selectLangUI"
}

func (m *SelectLangUIMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	var captions []string
	for _, lang := range SupportedLangUI() {
		captions = append(captions, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "lang." + lang}))
	}

	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "config.langUI"})
	options = append(options, bot.WithReplyKeyboard(captions))

	return
}

type SelectLangDictMessage struct{}

func (m *SelectLangDictMessage) Type() string {
	return "config.selectLangDict"
}

func (m *SelectLangDictMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	var captions []string
	for _, lang := range SupportedLangDict() {
		captions = append(captions, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "lang." + lang}))
	}

	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "config.langDict"})
	options = append(options, bot.WithReplyKeyboard(captions))

	return
}

type NotSupportedMessage struct {
	Value string
}

func (m *NotSupportedMessage) Type() string {
	return "config.notSupported"
}

func (m *NotSupportedMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "config.notSupported",
		TemplateData: map[string]interface{}{
			"Value": m.Value,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())

	return
}
