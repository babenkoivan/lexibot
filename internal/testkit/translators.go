package testkit

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type translatorSpy struct {
	testing    *testing.T
	translated [][3]string
}

func (s *translatorSpy) Translate(text, langFrom, langTo string) (string, error) {
	s.translated = append(s.translated, [3]string{text, langFrom, langTo})
	return "", errors.New("not found")
}

func (s *translatorSpy) AssertTranslated(text, langFrom, langTo string) {
	assert.Contains(s.testing, s.translated, [3]string{text, langFrom, langTo})
}

func NewTranslatorSpy(t *testing.T) *translatorSpy {
	return &translatorSpy{testing: t}
}
