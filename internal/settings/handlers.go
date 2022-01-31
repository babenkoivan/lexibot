package settings

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/locale"
	"strconv"
	"strings"
)

const OnSettings string = "/settings"

func settingsHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &EnableAutoTranslateMessage{})
}

func NewSettingsHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(settingsHandler)
}

type saveAutoTranslateHandler struct {
	locale        locale.Locale
	settingsStore SettingsStore
}

func (h *saveAutoTranslateHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.locale.MakeLocalizer(re.Sender.ID)
	answer := matchLocalizedMessage(re.Text, []string{"yes", "no"}, localizer, "")

	if answer == "" {
		b.Send(re.Sender, &EnumErrorMessage{re.Text})
		b.Send(re.Sender, msg.(*EnableAutoTranslateMessage))
		return
	}

	s := h.settingsStore.GetOrInit(re.Sender.ID)
	s.AutoTranslate = answer == "yes"
	h.settingsStore.Save(s)

	b.Send(re.Sender, &EnterWordsPerTrainingMessage{})
}

func NewSaveAutoTranslateHandler(locale locale.Locale, settingsStore SettingsStore) *saveAutoTranslateHandler {
	return &saveAutoTranslateHandler{locale, settingsStore}
}

type saveWordsPerTrainingHandler struct {
	settingsStore SettingsStore
}

func (h *saveWordsPerTrainingHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	number, err := strconv.Atoi(re.Text)

	if err != nil {
		b.Send(re.Sender, &IntegerErrorMessage{re.Text})
		b.Send(re.Sender, msg.(*EnterWordsPerTrainingMessage))
		return
	}

	s := h.settingsStore.GetOrInit(re.Sender.ID)
	s.WordsPerTraining = number
	h.settingsStore.Save(s)

	b.Send(re.Sender, &SuccessMessage{})
}

func NewSaveWordsPerTrainingHandler(settingsStore SettingsStore) *saveWordsPerTrainingHandler {
	return &saveWordsPerTrainingHandler{settingsStore}
}

type saveLangUIHandler struct {
	locale        locale.Locale
	settingsStore SettingsStore
}

func (h *saveLangUIHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.locale.MakeLocalizer(re.Sender.ID)
	lang := matchLocalizedMessage(re.Text, SupportedLangUI(), localizer, "lang.")

	if lang == "" {
		b.Send(re.Sender, &EnumErrorMessage{re.Text})
		b.Send(re.Sender, msg.(*SelectLangUIMessage))
		return
	}

	s := h.settingsStore.GetOrInit(re.Sender.ID)
	s.LangUI = lang
	h.settingsStore.Save(s)

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
	lang := matchLocalizedMessage(re.Text, SupportedLangDict(), localizer, "lang.")

	if lang == "" {
		b.Send(re.Sender, &EnumErrorMessage{re.Text})
		b.Send(re.Sender, msg.(*SelectLangDictMessage))
		return
	}

	s := h.settingsStore.GetOrInit(re.Sender.ID)
	s.LangDict = lang
	h.settingsStore.Save(s)

	b.Send(re.Sender, &SuccessMessage{})
}

func NewSaveLangDictHandler(locale locale.Locale, settingsStore SettingsStore) *saveLangDictHandler {
	return &saveLangDictHandler{locale, settingsStore}
}

func matchLocalizedMessage(text string, messageIDs []string, localizer *i18n.Localizer, prefix string) (match string) {
	text = strings.TrimSpace(text)

	for _, ID := range messageIDs {
		localized := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: prefix + ID})

		if strings.EqualFold(text, localized) {
			match = ID
			break
		}
	}

	return
}
