package training

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"lexibot/internal/bot"
)

type TaskMessage struct {
	Task *Task
}

func (m *TaskMessage) Type() string {
	return "training.task"
}

func (m *TaskMessage) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(m.Task.Question)

	if len(m.Task.Hints) == 0 {
		options = append(options, bot.WithoutReplyKeyboard())
		return
	}

	options = append(options, bot.WithReplyKeyboard(m.Task.Hints))
	return
}

type NotEnoughWordsError struct {
	RequiredWords int64
}

func (m *NotEnoughWordsError) Type() string {
	return "training.notEnoughWordsError"
}

func (m *NotEnoughWordsError) Render(localizer *i18n.Localizer) (text string, options []interface{}) {
	text = localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "training.notEnoughWordsError",
		TemplateData: map[string]interface{}{
			"RequiredWords": m.RequiredWords,
		},
	})

	options = append(options, bot.WithoutReplyKeyboard())
	return
}
