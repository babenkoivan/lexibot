package config

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/locale"
	"strings"
)

type saveLangUIHandler struct {
	locale      locale.UserLocale
	configStore ConfigStore
}

func (h *saveLangUIHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.locale.MakeLocalizer(re.Sender.ID)
	lang := matchLocalized(re.Text, SupportedLangUI(), localizer, "lang.")

	if lang == "" {
		b.Send(re.Sender, &NotSupportedMessage{re.Text})
		b.Send(re.Sender, msg.(*SelectLangUIMessage))
		return
	}

	config := h.configStore.Get(re.Sender.ID)
	config.LangUI = lang
	h.configStore.Save(config)

	b.Send(re.Sender, &SelectLangDictMessage{})
}

func NewSaveLangUIHandler(locale locale.UserLocale, configStore ConfigStore) *saveLangUIHandler {
	return &saveLangUIHandler{locale, configStore}
}

type saveLangDictHandler struct {
	locale      locale.UserLocale
	configStore ConfigStore
}

func (h *saveLangDictHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.locale.MakeLocalizer(re.Sender.ID)
	lang := matchLocalized(re.Text, SupportedLangDict(), localizer, "lang.")

	if lang == "" {
		b.Send(re.Sender, &NotSupportedMessage{re.Text})
		b.Send(re.Sender, msg.(*SelectLangDictMessage))
		return
	}

	config := h.configStore.Get(re.Sender.ID)
	config.LangDict = lang
	h.configStore.Save(config)

	b.Send(re.Sender, &bot.LocalizedTextMessage{"config.ok"})
}

func NewSaveLangDictHandler(locale locale.UserLocale, configStore ConfigStore) *saveLangDictHandler {
	return &saveLangDictHandler{locale, configStore}
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
