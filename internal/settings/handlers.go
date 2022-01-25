package settings

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/locale"
	"strings"
)

type saveLangUIHandler struct {
	locale        locale.Locale
	settingsStore SettingsStore
}

func (h *saveLangUIHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.locale.MakeLocalizer(re.Sender.ID)
	lang := matchLocalized(re.Text, SupportedLangUI(), localizer, "lang.")

	if lang == "" {
		b.Send(re.Sender, &NotSupportedMessage{re.Text})
		b.Send(re.Sender, msg.(*SelectLangUIMessage))
		return
	}

	settings := h.settingsStore.Get(re.Sender.ID)
	settings.LangUI = lang
	h.settingsStore.Save(settings)

	b.Send(re.Sender, &SelectLangDictMessage{})
}

func NewSaveLangUIHandler(locale locale.Locale, settingsStore SettingsStore) *saveLangUIHandler {
	return &saveLangUIHandler{locale, settingsStore}
}

type saveLangDictHandler struct {
	locale        locale.Locale
	settingsStore SettingsStore
}

func (h *saveLangDictHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.locale.MakeLocalizer(re.Sender.ID)
	lang := matchLocalized(re.Text, SupportedLangDict(), localizer, "lang.")

	if lang == "" {
		b.Send(re.Sender, &NotSupportedMessage{re.Text})
		b.Send(re.Sender, msg.(*SelectLangDictMessage))
		return
	}

	settings := h.settingsStore.Get(re.Sender.ID)
	settings.LangDict = lang
	h.settingsStore.Save(settings)

	b.Send(re.Sender, &bot.LocalizedTextMessage{"settings.ok"})
}

func NewSaveLangDictHandler(locale locale.Locale, settingsStore SettingsStore) *saveLangDictHandler {
	return &saveLangDictHandler{locale, settingsStore}
}

func matchLocalized(str string, messageIDs []string, localizer *i18n.Localizer, prefix string) (match string) {
	str = strings.TrimSpace(str)

	for _, ID := range messageIDs {
		localized := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: prefix + ID})

		if strings.ToLower(str) == strings.ToLower(localized) {
			match = ID
			break
		}
	}

	return
}
