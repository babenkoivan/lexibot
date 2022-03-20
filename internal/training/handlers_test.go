package training_test

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/settings"
	"lexibot/internal/testkit"
	"lexibot/internal/training"
	"lexibot/internal/translation"
	"testing"
)

func TestStartTrainingHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	msg := &telebot.Message{Sender: user}
	langFrom, langTo := "de", "en"
	wordsPerTraining := 10

	settingsStoreMock := testkit.MockSettingsStore(t)
	settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
		assert.Equal(t, user.ID, userID)
		return &settings.Settings{LangDict: langFrom, LangUI: langTo, WordsPerTraining: wordsPerTraining}
	})

	t.Run("not enough words to start the training", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)

		training.NewStartTrainingHandler(
			settingsStoreMock,
			testkit.MockTranslationStore(t),
			testkit.MockTaskStore(t),
			testkit.MockTaskGenerator(),
		).Handle(botSpy, msg)

		botSpy.AssertSent(user, &training.NotEnoughWordsError{int64(wordsPerTraining)})
	})

	t.Run("enough words to start the training", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		taskStoreMock := testkit.MockTaskStore(t)

		translationStoreMock := testkit.MockTranslationStore(t)
		translationStoreMock.OnCount(func(conds ...translation.TranslationQueryCond) int64 {
			testkit.AssertEqualTranslationQuery(t, []translation.TranslationQueryCond{
				translation.WithUserID(user.ID),
				translation.WithLangFrom(langFrom),
			}, conds)

			return int64(wordsPerTraining + 1)
		})

		nextTask := &training.Task{
			UserID:        user.ID,
			TranslationID: 1,
			Question:      "bunt",
			Answer:        "colorful",
		}

		generatorMock := testkit.MockTaskGenerator()
		generatorMock.OnNext(func(userID int) *training.Task {
			assert.Equal(t, user.ID, userID)
			return nextTask
		})

		training.NewStartTrainingHandler(
			settingsStoreMock,
			translationStoreMock,
			taskStoreMock,
			generatorMock,
		).Handle(botSpy, msg)

		taskStoreMock.AssertCleaned(user.ID)
		botSpy.AssertSent(user, &training.TranslateTaskMessage{nextTask, 1, int64(wordsPerTraining)})
	})
}

func TestCheckAnswerHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	prevTask := &training.Task{UserID: user.ID, TranslationID: 1, Question: "mouth", Answer: "mund"}
	nextTask := &training.Task{UserID: user.ID, TranslationID: 2, Question: "bunt", Answer: "colorful"}
	taskCount := int64(5)
	wordsPerTraining := 10
	msg := &training.TranslateTaskMessage{prevTask, taskCount, int64(wordsPerTraining)}

	settingsStoreMock := testkit.MockSettingsStore(t)
	settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
		assert.Equal(t, user.ID, userID)
		return &settings.Settings{WordsPerTraining: wordsPerTraining}
	})

	taskGeneratorMock := testkit.MockTaskGenerator()
	taskGeneratorMock.OnNext(func(userID int) *training.Task {
		assert.Equal(t, user.ID, userID)
		return nextTask
	})

	t.Run("given incorrect answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		translationStoreMock := testkit.MockTranslationStore(t)

		taskStoreMock := testkit.MockTaskStore(t)
		taskStoreMock.OnCount(func(userID int) int64 {
			return taskCount
		})

		training.NewCheckAnswerHandler(
			taskStoreMock,
			translationStoreMock,
			settingsStoreMock,
			taskGeneratorMock,
		).Handle(botSpy, &telebot.Message{Sender: user, Text: "foo"}, msg)

		taskStoreMock.AssertScoreDecremented(prevTask.TranslationID, prevTask.UserID)
		translationStoreMock.AssertScoreDecremented(prevTask.TranslationID, prevTask.UserID)
		botSpy.AssertSent(user, &training.TranslateTaskMessage{nextTask, taskCount + 1, int64(wordsPerTraining)})
	})

	t.Run("given correct answer", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		translationStoreMock := testkit.MockTranslationStore(t)

		taskStoreMock := testkit.MockTaskStore(t)
		taskStoreMock.OnCount(func(userID int) int64 {
			return taskCount
		})

		training.NewCheckAnswerHandler(
			taskStoreMock,
			translationStoreMock,
			settingsStoreMock,
			taskGeneratorMock,
		).Handle(botSpy, &telebot.Message{Sender: user, Text: prevTask.Answer}, msg)

		taskStoreMock.AssertScoreIncremented(prevTask.TranslationID, prevTask.UserID)
		translationStoreMock.AssertScoreIncremented(prevTask.TranslationID, prevTask.UserID)
		botSpy.AssertSent(user, &training.TranslateTaskMessage{nextTask, taskCount + 1, int64(wordsPerTraining)})
	})

	t.Run("the training is completed", func(t *testing.T) {
		botSpy := testkit.NewBotSpy(t)
		correctAnswers := int64(3)

		taskStoreMock := testkit.MockTaskStore(t)
		taskStoreMock.OnCount(func(userID int) int64 {
			assert.Equal(t, user.ID, userID)
			return int64(wordsPerTraining)
		})
		taskStoreMock.OnCorrectCount(func(userID int) int64 {
			assert.Equal(t, user.ID, userID)
			return correctAnswers
		})

		training.NewCheckAnswerHandler(
			taskStoreMock,
			testkit.MockTranslationStore(t),
			settingsStoreMock,
			taskGeneratorMock,
		).Handle(botSpy, &telebot.Message{Sender: user, Text: prevTask.Answer}, msg)

		botSpy.AssertSent(user, &training.ResultsMessage{int64(wordsPerTraining), correctAnswers})
	})
}
