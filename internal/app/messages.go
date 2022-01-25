package app

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/bot"
)

type StartMessage struct{}

func (m *StartMessage) Type() string {
	return "app.start"
}

func (m *StartMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "app.start"})
	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type HelpMessage struct{}

func (m *HelpMessage) Type() string {
	return "app.help"
}

func (m *HelpMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "app.help"})
	options = append(options, bot.WithoutReplyKeyboard())
	return
}
