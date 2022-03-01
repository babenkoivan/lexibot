package testkit

import (
	"lexibot/internal/translation"
)

type translationStoreMock struct {
	onFirst func(conds ...translation.TranslationQueryCond) *translation.Translation
}

func (s *translationStoreMock) Save(transl *translation.Translation) *translation.Translation {
	return transl
}

func (s *translationStoreMock) OnFirst(fn func(conds ...translation.TranslationQueryCond) *translation.Translation) {
	s.onFirst = fn
}

func (s *translationStoreMock) First(conds ...translation.TranslationQueryCond) *translation.Translation {
	return s.onFirst(conds...)
}

func (s *translationStoreMock) Rand(conds ...translation.TranslationQueryCond) []*translation.Translation {
	return []*translation.Translation{}
}

func (s *translationStoreMock) Count(conds ...translation.TranslationQueryCond) int64 {
	return 0
}

func MockTranslationStore() *translationStoreMock {
	return &translationStoreMock{}
}
