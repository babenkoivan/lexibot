package training

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/bot"
	"lexibot/internal/localization"
)

type TranslateTaskMessage struct {
	Task *Task
}

func (m *TranslateTaskMessage) Type() string {
	return "training.task"
}

func (m *TranslateTaskMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "training.task",
		TemplateData: map[string]interface{}{
			"Question": m.Task.Question,
		},
	})

	if len(m.Task.Hints) == 0 {
		options = append(options, bot.WithoutReplyKeyboard())
		return
	}

	options = append(options, bot.WithReplyKeyboard(m.Task.Hints))
	return
}

type CorrectAnswerMessage struct{}

func (m *CorrectAnswerMessage) Type() string {
	return "training.correctAnswer"
}

func (m *CorrectAnswerMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "training.correctAnswer"})
	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type IncorrectAnswerMessage struct {
	CorrectAnswer string
}

func (m *IncorrectAnswerMessage) Type() string {
	return "training.incorrectAnswer"
}

func (m *IncorrectAnswerMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "training.incorrectAnswer",
		TemplateData: map[string]interface{}{
			"CorrectAnswer": m.CorrectAnswer,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type ResultsMessage struct {
	TaskCount      int64
	CorrectAnswers int64
}

func (m *ResultsMessage) Type() string {
	return "training.results"
}

func (m *ResultsMessage) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "training.results",
		TemplateData: map[string]interface{}{
			"TaskCount":      m.TaskCount,
			"CorrectAnswers": m.CorrectAnswers,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}

type NotEnoughWordsError struct {
	RequiredWords int64
}

func (m *NotEnoughWordsError) Type() string {
	return "training.notEnoughWordsError"
}

func (m *NotEnoughWordsError) Render(localizer localization.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "training.notEnoughWordsError",
		TemplateData: map[string]interface{}{
			"RequiredWords": m.RequiredWords,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}
