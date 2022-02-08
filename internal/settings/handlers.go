package settings

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/locale"
	"strconv"
	"strings"
)

const (
	OnSettings          string = "/settings"
	maxWordsPerTraining int    = 50
)

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

	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)
	userSettings.AutoTranslate = answer == "yes"
	h.settingsStore.Save(userSettings)

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

	if err != nil || number > maxWordsPerTraining {
		b.Send(re.Sender, &IntegerErrorMessage{maxWordsPerTraining})
		b.Send(re.Sender, msg.(*EnterWordsPerTrainingMessage))
		return
	}

	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)
	userSettings.WordsPerTraining = number
	h.settingsStore.Save(userSettings)

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

	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)
	userSettings.LangUI = lang
	h.settingsStore.Save(userSettings)

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

	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)
	userSettings.LangDict = lang
	h.settingsStore.Save(userSettings)

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
