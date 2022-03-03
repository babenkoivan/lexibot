package testkit

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type translatorMock struct {
	testing     *testing.T
	onTranslate func(text, langFrom, langTo string) (string, error)
	translated  [][3]string
}

func (m *translatorMock) OnTranslate(callback func(text, langFrom, langTo string) (string, error)) {
	m.onTranslate = callback
}

func (m *translatorMock) Translate(text, langFrom, langTo string) (string, error) {
	m.translated = append(m.translated, [3]string{text, langFrom, langTo})

	if m.onTranslate == nil {
		return "", errors.New("not found")
	}

	return m.onTranslate(text, langFrom, langTo)
}

func (m *translatorMock) AssertTranslated(text, langFrom, langTo string) {
	assert.Contains(m.testing, m.translated, [3]string{text, langFrom, langTo})
}

func MockTranslator(t *testing.T) *translatorMock {
	return &translatorMock{testing: t}
}
