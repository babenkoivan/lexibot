package settings

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/bot"
)

type SelectLangUIMessage struct{}

func (m *SelectLangUIMessage) Type() string {
	return "settings.selectLangUI"
}

func (m *SelectLangUIMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	var captions []string
	for _, lang := range SupportedLangUI() {
		captions = append(captions, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "lang." + lang}))
	}

	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "settings.langUI"})
	options = append(options, bot.WithReplyKeyboard(captions))

	return
}

type SelectLangDictMessage struct{}

func (m *SelectLangDictMessage) Type() string {
	return "settings.selectLangDict"
}

func (m *SelectLangDictMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	var captions []string
	for _, lang := range SupportedLangDict() {
		captions = append(captions, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "lang." + lang}))
	}

	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "settings.langDict"})
	options = append(options, bot.WithReplyKeyboard(captions))

	return
}

type NotSupportedMessage struct {
	Value string
}

func (m *NotSupportedMessage) Type() string {
	return "settings.notSupported"
}

func (m *NotSupportedMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "settings.notSupported",
		TemplateData: map[string]interface{}{
			"Value": m.Value,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())

	return
}
