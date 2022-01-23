package config

import (
	"encoding/json"
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

func (h *saveLangUIHandler) Handle(b bot.Bot, re *telebot.Message, hm *bot.HistoryMessage) {
	input := strings.TrimSpace(re.Text)
	localizer := h.locale.MakeLocalizer(re.Sender.ID)

	msg := &SelectLangUIMessage{}
	json.Unmarshal([]byte(hm.Content), msg)

	var selectedLang string
	for _, lang := range msg.Lang {
		localizedLang := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "lang." + lang})

		if strings.ToLower(input) == strings.ToLower(localizedLang) {
			selectedLang = lang
			break
		}
	}

	if selectedLang == "" {
		b.Send(re.Sender, msg)
		return
	}

	config := h.configStore.Get(re.Sender.ID)
	config.LangUI = selectedLang
	h.configStore.Save(config)
}

func NewSaveLangUIHandler(locale locale.UserLocale, configStore ConfigStore) *saveLangUIHandler {
	return &saveLangUIHandler{locale, configStore}
}
