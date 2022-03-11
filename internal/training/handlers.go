package training

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
	"lexibot/internal/settings"
	"lexibot/internal/translation"
	"strings"
)

const OnTraining = "/training"

type startTrainingHandler struct {
	settingsStore    settings.SettingsStore
	translationStore translation.TranslationStore
	taskStore        TaskStore
	taskGenerator    TaskGenerator
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

	h.taskStore.Cleanup(msg.Sender.ID)

	task := h.taskGenerator.Next(msg.Sender.ID)
	b.Send(msg.Sender, &TranslateTaskMessage{task, 1, requiredWordsPerTraining})
}

func NewStartTrainingHandler(
	settingsStore settings.SettingsStore,
	translationStore translation.TranslationStore,
	taskStore TaskStore,
	taskGenerator TaskGenerator,
) *startTrainingHandler {
	return &startTrainingHandler{
		settingsStore,
		translationStore,
		taskStore,
		taskGenerator,
	}
}

type checkAnswerHandler struct {
	scoreStore    translation.ScoreStore
	taskStore     TaskStore
	settingsStore settings.SettingsStore
	taskGenerator TaskGenerator
}

func (h *checkAnswerHandler) Handle(b bot.Bot, re *telebot.Message, msg bot.Message) {
	prevTask := msg.(*TranslateTaskMessage).Task

	correctAnswer := strings.TrimSpace(prevTask.Answer)
	givenAnswer := strings.TrimSpace(re.Text)

	if strings.EqualFold(correctAnswer, givenAnswer) {
		h.taskStore.IncrementScore(prevTask)
		h.scoreStore.Increment(prevTask.TranslationID, re.Sender.ID)
		b.Send(re.Sender, &CorrectAnswerMessage{})
	} else {
		h.taskStore.DecrementScore(prevTask)
		h.scoreStore.Decrement(prevTask.TranslationID, re.Sender.ID)
		b.Send(re.Sender, &IncorrectAnswerMessage{correctAnswer})
	}

	userSettings := h.settingsStore.FirstOrInit(re.Sender.ID)
	wordsPerTraining := int64(userSettings.WordsPerTraining)
	taskCount := h.taskStore.Count(re.Sender.ID)

	var nextTask *Task
	if taskCount < wordsPerTraining {
		nextTask = h.taskGenerator.Next(re.Sender.ID)
	}

	if nextTask == nil {
		correctAnswers := h.taskStore.TotalPositiveScore(re.Sender.ID)
		b.Send(re.Sender, &ResultsMessage{taskCount, correctAnswers})
		return
	}

	b.Send(re.Sender, &TranslateTaskMessage{nextTask, taskCount + 1, wordsPerTraining})
}

func NewCheckAnswerHandler(
	scoreStore translation.ScoreStore,
	taskStore TaskStore,
	settingsStore settings.SettingsStore,
	taskGenerator TaskGenerator,
) *checkAnswerHandler {
	return &checkAnswerHandler{
		scoreStore,
		taskStore,
		settingsStore,
		taskGenerator,
	}
}
