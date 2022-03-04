package training_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/settings"
	"lexibot/internal/testkit"
	"lexibot/internal/training"
	"lexibot/internal/translation"
	"testing"
)

func TestTranslateTaskGenerator_Next(t *testing.T) {
	user := telebot.User{ID: 1}
	langFrom, langTo := "de", "en"

	transl := &translation.Translation{ID: 1, Text: "bunt", Translation: "colorful", LangFrom: langFrom, LangTo: langTo}
	randTransl := &translation.Translation{ID: 2, Text: "mund", Translation: "mouth", LangFrom: langFrom, LangTo: langTo}

	settingsStoreMock := testkit.MockSettingsStore(t)
	settingsStoreMock.OnFirstOrInit(func(userID int) *settings.Settings {
		assert.Equal(t, user.ID, userID)
		return &settings.Settings{LangDict: langFrom, LangUI: langTo}
	})

	translationStoreMock := testkit.MockTranslationStore(t)
	translationStoreMock.OnFirst(func(conds ...translation.TranslationQueryCond) *translation.Translation {
		testkit.AssertTranslationQuery(t, []translation.TranslationQueryCond{
			translation.WithID(transl.ID),
		}, conds)

		return transl
	})
	translationStoreMock.OnRand(func(conds ...translation.TranslationQueryCond) []*translation.Translation {
		testkit.AssertTranslationQuery(t, []translation.TranslationQueryCond{
			translation.WithoutID(transl.ID),
			translation.WithUserID(user.ID),
			translation.WithLangFrom(langFrom),
			translation.WithLimit(training.HintsLimit - 1),
		}, conds)

		return []*translation.Translation{randTransl}
	})

	t.Run("score is not found", func(t *testing.T) {
		task := training.NewTaskGenerator(
			settingsStoreMock,
			translationStoreMock,
			testkit.MockScoreStore(t),
			testkit.MockTaskStore(t),
		).Next(user.ID)

		assert.Nil(t, task)
	})

	for _, score := range []int{0, training.FamiliarWordScore} {
		t.Run(fmt.Sprintf("score is %q", score), func(t *testing.T) {
			taskStoreMock := testkit.MockTaskStore(t)

			scoreStoreMock := testkit.MockScoreStore(t)
			scoreStoreMock.OnLowestNotTrained(func(userID int, langDict string) *translation.Score {
				assert.Equal(t, user.ID, userID)
				assert.Equal(t, langFrom, langDict)
				return &translation.Score{UserID: user.ID, TranslationID: transl.ID, Score: training.FamiliarWordScore}
			})

			task := training.NewTaskGenerator(
				settingsStoreMock,
				translationStoreMock,
				scoreStoreMock,
				taskStoreMock,
			).Next(user.ID)

			assert.Equal(t, user.ID, task.UserID)
			assert.Equal(t, transl.ID, task.TranslationID)
			assert.Equal(t, 0, task.Score)

			expectedQuestionOrAnswer := []string{transl.Text, transl.Translation}
			assert.Contains(t, expectedQuestionOrAnswer, task.Question)
			assert.Contains(t, expectedQuestionOrAnswer, task.Answer)

			if score >= training.FamiliarWordScore && len(task.Hints) > 0 {
				expectedHints := []string{transl.Text, transl.Translation, randTransl.Text, randTransl.Translation}
				assert.Contains(t, expectedHints, task.Hints[0])
				assert.Contains(t, expectedHints, task.Hints[1])
			}

			taskStoreMock.AssertSaved(task)
		})
	}
}
