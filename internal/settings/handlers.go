package settings

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/localization"
	"strconv"
)

const (
	OnSettings          string = "/settings"
	MaxWordsPerTraining int    = 50
)

func settingsHandler(b bot.Bot, msg *telebot.Message) {
	b.Send(msg.Sender, &EnableAutoTranslateMessage{})
}

func NewSettingsHandler() bot.MessageHandler {
	return bot.MessageHandlerFunc(settingsHandler)
}

type saveAutoTranslateHandler struct {
	localizerFactory localization.LocalizerFactory
	settingsStore    SettingsStore
}

func (h *saveAutoTranslateHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.localizerFactory.New(re.Sender.ID)
	answer := localizer.MatchMessage(re.Text, []string{"yes", "no"})

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

func NewSaveAutoTranslateHandler(
	localizerFactory localization.LocalizerFactory,
	settingsStore SettingsStore,
) *saveAutoTranslateHandler {
	return &saveAutoTranslateHandler{localizerFactory, settingsStore}
}

type saveWordsPerTrainingHandler struct {
	settingsStore SettingsStore
}

func (h *saveWordsPerTrainingHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	number, err := strconv.Atoi(re.Text)

	if err != nil || number <= 0 || number > MaxWordsPerTraining {
		b.Send(re.Sender, &IntegerErrorMessage{MaxWordsPerTraining})
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
	localizerFactory localization.LocalizerFactory
	settingsStore    SettingsStore
}

func (h *saveLangUIHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.localizerFactory.New(re.Sender.ID)
	lang := localizer.MatchMessage(re.Text, SupportedLangUI())

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

func NewSaveLangUIHandler(localizerFactory localization.LocalizerFactory, settingsStore SettingsStore) *saveLangUIHandler {
	return &saveLangUIHandler{localizerFactory, settingsStore}
}

type saveLangDictHandler struct {
	localizerFactory localization.LocalizerFactory
	settingsStore    SettingsStore
}

func (h *saveLangDictHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	localizer := h.localizerFactory.New(re.Sender.ID)
	lang := localizer.MatchMessage(re.Text, SupportedLangDict())

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

func NewSaveLangDictHandler(
	localizerFactory localization.LocalizerFactory,
	settingsStore SettingsStore,
) *saveLangDictHandler {
	return &saveLangDictHandler{localizerFactory, settingsStore}
}
