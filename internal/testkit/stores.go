package testkit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"lexibot/internal/settings"
	"lexibot/internal/training"
	"lexibot/internal/translation"
	"testing"
	"time"
)

func AssertTranslationQuery(t *testing.T, expected, actual []translation.TranslationQueryCond) {
	expectedQuery := translation.MakeTranslationQuery(expected)
	actualQuery := translation.MakeTranslationQuery(actual)
	assert.Equal(t, expectedQuery, actualQuery)
}

type translationStoreMock struct {
	testing *testing.T
	onFirst func(conds ...translation.TranslationQueryCond) *translation.Translation
	onFind  func(conds ...translation.TranslationQueryCond) []*translation.Translation
	onRand  func(conds ...translation.TranslationQueryCond) []*translation.Translation
	saved   []*translation.Translation
}

func (m *translationStoreMock) Save(transl *translation.Translation) *translation.Translation {
	transl.ID = len(m.saved) + 1
	m.saved = append(m.saved, transl)
	return transl
}

func (m *translationStoreMock) AssertSaved(transl *translation.Translation) {
	assert.Contains(m.testing, m.saved, transl, fmt.Sprintf("Translation %#v is not saved", transl))
}

func (m *translationStoreMock) AssertNothingSaved() {
	assert.Len(m.testing, m.saved, 0, "Translations are not supposed to be saved")
}

func (m *translationStoreMock) OnFirst(callback func(conds ...translation.TranslationQueryCond) *translation.Translation) {
	m.onFirst = callback
}

func (m *translationStoreMock) First(conds ...translation.TranslationQueryCond) *translation.Translation {
	if m.onFirst == nil {
		return nil
	}

	return m.onFirst(conds...)
}

func (m *translationStoreMock) OnFind(callback func(conds ...translation.TranslationQueryCond) []*translation.Translation) {
	m.onFind = callback
}

func (m *translationStoreMock) Find(conds ...translation.TranslationQueryCond) []*translation.Translation {
	if m.onFind == nil {
		return []*translation.Translation{}
	}

	return m.onFind(conds...)
}

func (m *translationStoreMock) OnRand(callback func(conds ...translation.TranslationQueryCond) []*translation.Translation) {
	m.onRand = callback
}

func (m *translationStoreMock) Rand(conds ...translation.TranslationQueryCond) []*translation.Translation {
	if m.onRand == nil {
		return []*translation.Translation{}
	}

	return m.onRand(conds...)
}

func (m *translationStoreMock) Count(conds ...translation.TranslationQueryCond) int64 {
	return 0
}

func MockTranslationStore(t *testing.T) *translationStoreMock {
	return &translationStoreMock{testing: t}
}

type settingsStoreMock struct {
	testing       *testing.T
	onLocale      func(userID int) string
	onFirstOrInit func(userID int) *settings.Settings
	saved         []*settings.Settings
}

func (m *settingsStoreMock) OnLocale(callback func(userID int) string) {
	m.onLocale = callback
}

func (m *settingsStoreMock) Locale(userID int) string {
	if m.onLocale == nil {
		return ""
	}

	return m.onLocale(userID)
}

func (m *settingsStoreMock) Save(settings *settings.Settings) *settings.Settings {
	m.saved = append(m.saved, settings)
	return settings
}

func (m *settingsStoreMock) AssertSaved(settings *settings.Settings) {
	assert.Contains(m.testing, m.saved, settings, fmt.Sprintf("Settings %#v are not saved", settings))
}

func (m *settingsStoreMock) AssertNothingSaved() {
	assert.Len(m.testing, m.saved, 0, "Settings are not supposed to be saved")
}

func (m *settingsStoreMock) OnFirstOrInit(callback func(userID int) *settings.Settings) {
	m.onFirstOrInit = callback
}

func (m *settingsStoreMock) FirstOrInit(userID int) *settings.Settings {
	if m.onFirstOrInit == nil {
		return &settings.Settings{UserID: userID}
	}

	return m.onFirstOrInit(userID)
}

func MockSettingsStore(t *testing.T) *settingsStoreMock {
	return &settingsStoreMock{testing: t}
}

type scoreStoreMock struct {
	testing            *testing.T
	onLowestNotTrained func(userID int, langDict string) *translation.Score
	saved              [][2]int
	deleted            [][2]int
	autoDecremented    []time.Duration
}

func (m *scoreStoreMock) Save(translationID, userID int) *translation.Score {
	m.saved = append(m.saved, [2]int{translationID, userID})
	return &translation.Score{UserID: userID, TranslationID: translationID}
}

func (m *scoreStoreMock) AssertSaved(translationID, userID int) {
	msg := fmt.Sprintf("Score with translation %d for user %d is not saved", translationID, userID)
	assert.Contains(m.testing, m.saved, [2]int{translationID, userID}, msg)
}

func (m *scoreStoreMock) AssertNothingSaved() {
	assert.Len(m.testing, m.saved, 0, "Scores are not supposed to be saved")
}

func (m *scoreStoreMock) Delete(translationID, userID int) {
	m.deleted = append(m.deleted, [2]int{translationID, userID})
}

func (m *scoreStoreMock) AssertDeleted(translationID, userID int) {
	msg := fmt.Sprintf("Score with translation %d for user %d is not deleted", translationID, userID)
	assert.Contains(m.testing, m.deleted, [2]int{translationID, userID}, msg)
}

func (m *scoreStoreMock) Increment(translationID, userID int) {}

func (m *scoreStoreMock) Decrement(translationID, userID int) {}

func (m *scoreStoreMock) AutoDecrement(after time.Duration) {
	m.autoDecremented = append(m.autoDecremented, after)
}

func (m *scoreStoreMock) AssertAutoDecremented(after time.Duration) {
	msg := fmt.Sprintf("Scores are not auto-decremented after %v", after)
	assert.Contains(m.testing, m.autoDecremented, after, msg)
}

func (m *scoreStoreMock) OnLowestNotTrained(callback func(userID int, langDict string) *translation.Score) {
	m.onLowestNotTrained = callback
}

func (m *scoreStoreMock) LowestNotTrained(userID int, langDict string) *translation.Score {
	if m.onLowestNotTrained == nil {
		return nil
	}

	return m.onLowestNotTrained(userID, langDict)
}

func MockScoreStore(t *testing.T) *scoreStoreMock {
	return &scoreStoreMock{testing: t}
}

type taskStoreMock struct {
	testing *testing.T
	saved   []*training.Task
}

func (m *taskStoreMock) Save(task *training.Task) *training.Task {
	m.saved = append(m.saved, task)
	return task
}

func (m *taskStoreMock) AssertSaved(task *training.Task) {
	assert.Contains(m.testing, m.saved, task, fmt.Sprintf("Task %#v is not saved", task))
}

func (m *taskStoreMock) Cleanup(userID int) {}

func (m *taskStoreMock) Count(userID int) int64 {
	return 0
}

func (m *taskStoreMock) IncrementScore(task *training.Task) {}

func (m *taskStoreMock) DecrementScore(task *training.Task) {}

func (m *taskStoreMock) TotalPositiveScore(userID int) int64 {
	return 0
}

func MockTaskStore(t *testing.T) *taskStoreMock {
	return &taskStoreMock{testing: t}
}
