package training

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
	"lexibot/internal/translation"
)

const OnTraining = "/training"

type startTrainingHandler struct {
	settingsStore    settings.SettingsStore
	translationStore translation.TranslationStore
	taskGenerator    *taskGenerator
}

func (h *startTrainingHandler) Handle(b bot.Bot, msg *telebot.Message) {
	userSettings := h.settingsStore.FirstOrInit(msg.Sender.ID)
	requiredWordsPerTraining := int64(userSettings.WordsPerTraining)

	actualWordsCount := h.translationStore.Count(
		translation.WithUserID(msg.Sender.ID),
		translation.WithLangFrom(userSettings.LangDict),
	)

	if actualWordsCount < requiredWordsPerTraining {
		b.Send(msg.Sender, &NotEnoughWordsError{requiredWordsPerTraining - actualWordsCount})
		return
	}

	task := h.taskGenerator.Next(msg.Sender.ID)
	b.Send(msg.Sender, &TaskMessage{task})
}

func NewStartTrainingHandler(
	settingsStore settings.SettingsStore,
	translationStore translation.TranslationStore,
	taskGenerator *taskGenerator,
) *startTrainingHandler {
	return &startTrainingHandler{
		settingsStore,
		translationStore,
		taskGenerator,
	}
}
