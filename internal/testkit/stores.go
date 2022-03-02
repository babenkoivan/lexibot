package testkit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"lexibot/internal/settings"
	"lexibot/internal/translation"
	"testing"
)

type translationStoreMock struct {
	onFirst func(conds ...translation.TranslationQueryCond) *translation.Translation
}

func (m *translationStoreMock) Save(transl *translation.Translation) *translation.Translation {
	return transl
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

func (m *translationStoreMock) Rand(conds ...translation.TranslationQueryCond) []*translation.Translation {
	return []*translation.Translation{}
}

func (m *translationStoreMock) Count(conds ...translation.TranslationQueryCond) int64 {
	return 0
}

func MockTranslationStore() *translationStoreMock {
	return &translationStoreMock{}
}

type settingsStoreMock struct {
	testing  *testing.T
	onLocale func(userID int) string
	saved    []*settings.Settings
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

func (m settingsStoreMock) AssertSaved(settings *settings.Settings) {
	assert.Contains(m.testing, m.saved, settings, fmt.Sprintf("Settings %#v are not saved", settings))
}

func (m settingsStoreMock) AssertNothingSaved() {
	assert.Len(m.testing, m.saved, 0, "Settings are not supposed to be saved")
}

func (m *settingsStoreMock) FirstOrInit(userID int) *settings.Settings {
	return &settings.Settings{UserID: userID}
}

func MockSettingsStore(t *testing.T) *settingsStoreMock {
	return &settingsStoreMock{testing: t}
}
