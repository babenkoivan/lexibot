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

func AssertEqualTranslationQuery(t *testing.T, expected, actual []translation.TranslationQueryCond) {
	expectedQuery := translation.MakeTranslationQuery(expected)
	actualQuery := translation.MakeTranslationQuery(actual)
	assert.Equal(t, expectedQuery, actualQuery)
}

type translationStoreMock struct {
	testing          *testing.T
	onFirst          func(conds ...translation.TranslationQueryCond) *translation.Translation
	onFind           func(conds ...translation.TranslationQueryCond) []*translation.Translation
	onRand           func(conds ...translation.TranslationQueryCond) []*translation.Translation
	onCount          func(conds ...translation.TranslationQueryCond) int64
	saved            []*translation.Translation
	deleted          []*translation.Translation
	incrementedScore [][2]int
	decrementedScore [][2]int
	autoDecremented  []time.Duration
}

func (m *translationStoreMock) Save(transl *translation.Translation) *translation.Translation {
	transl.ID = len(m.saved) + 1
	m.saved = append(m.saved, transl)
	return transl
}

func (m *translationStoreMock) Delete(transl *translation.Translation) {
	m.deleted = append(m.deleted, transl)
}

func (m *translationStoreMock) AssertDeleted(transl *translation.Translation) {
	assert.Contains(m.testing, m.deleted, transl, fmt.Sprintf("Translation %#v is not deleted", transl))
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

func (m *translationStoreMock) OnCount(callback func(conds ...translation.TranslationQueryCond) int64) {
	m.onCount = callback
}

func (m *translationStoreMock) Count(conds ...translation.TranslationQueryCond) int64 {
	if m.onCount == nil {
		return 0
	}

	return m.onCount(conds...)
}

func (m *translationStoreMock) IncrementScore(id, userID int) {
	m.incrementedScore = append(m.incrementedScore, [2]int{id, userID})
}

func (m *translationStoreMock) AssertScoreIncremented(id, userID int) {
	msg := fmt.Sprintf("Score is not incremented for a transltion with id %d and user id %d", id, userID)
	assert.Contains(m.testing, m.incrementedScore, [2]int{id, userID}, msg)
}

func (m *translationStoreMock) DecrementScore(id, userID int) {
	m.decrementedScore = append(m.decrementedScore, [2]int{id, userID})
}

func (m *translationStoreMock) AssertScoreDecremented(id, userID int) {
	msg := fmt.Sprintf("Score is not decremented for a translation with id %d and user id %d", id, userID)
	assert.Contains(m.testing, m.decrementedScore, [2]int{id, userID}, msg)
}

func (m *translationStoreMock) AutoDecrementScore(after time.Duration) {
	m.autoDecremented = append(m.autoDecremented, after)
}

func (m *translationStoreMock) AssertAutoDecremented(after time.Duration) {
	msg := fmt.Sprintf("Translation scores are not auto-decremented after %v", after)
	assert.Contains(m.testing, m.autoDecremented, after, msg)
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

type taskStoreMock struct {
	testing          *testing.T
	onCount          func(userID int) int64
	onCorrectCount   func(userID int) int64
	onTranslationIDs func(userID int) []int
	saved            []*training.Task
	cleaned          []int
	incrementedScore [][2]int
	decrementedScore [][2]int
}

func (m *taskStoreMock) Save(task *training.Task) *training.Task {
	m.saved = append(m.saved, task)
	return task
}

func (m *taskStoreMock) AssertSaved(task *training.Task) {
	assert.Contains(m.testing, m.saved, task, fmt.Sprintf("Task %#v is not saved", task))
}

func (m *taskStoreMock) Cleanup(userID int) {
	m.cleaned = append(m.cleaned, userID)
}

func (m *taskStoreMock) AssertCleaned(userID int) {
	assert.Contains(m.testing, m.cleaned, userID, fmt.Sprintf("Tasks are not cleaned up for user %d", userID))
}

func (m *taskStoreMock) OnCount(callback func(userID int) int64) {
	m.onCount = callback
}

func (m *taskStoreMock) Count(userID int) int64 {
	if m.onCount == nil {
		return 0
	}

	return m.onCount(userID)
}

func (m *taskStoreMock) OnCorrectCount(callback func(userID int) int64) {
	m.onCorrectCount = callback
}

func (m *taskStoreMock) CorrectCount(userID int) int64 {
	if m.onCorrectCount == nil {
		return 0
	}

	return m.onCorrectCount(userID)
}

func (m *taskStoreMock) IncrementScore(translationId, userID int) {
	m.incrementedScore = append(m.incrementedScore, [2]int{translationId, userID})
}

func (m *taskStoreMock) AssertScoreIncremented(translationId, userID int) {
	msg := fmt.Sprintf("Score is not incremented for a task with translation id %d and user id %d", translationId, userID)
	assert.Contains(m.testing, m.incrementedScore, [2]int{translationId, userID}, msg)
}

func (m *taskStoreMock) DecrementScore(translationId, userID int) {
	m.decrementedScore = append(m.decrementedScore, [2]int{translationId, userID})
}

func (m *taskStoreMock) AssertScoreDecremented(translationId, userID int) {
	msg := fmt.Sprintf("Score is not decremented for a task with translation id %d and user id %d", translationId, userID)
	assert.Contains(m.testing, m.decrementedScore, [2]int{translationId, userID}, msg)
}

func (m *taskStoreMock) OnTranslationIDs(callback func(userID int) []int) {
	m.onTranslationIDs = callback
}

func (m *taskStoreMock) TranslationIDs(userID int) []int {
	if m.onTranslationIDs == nil {
		return nil
	}

	return m.onTranslationIDs(userID)
}

func MockTaskStore(t *testing.T) *taskStoreMock {
	return &taskStoreMock{testing: t}
}
