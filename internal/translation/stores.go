package translation

import (
	"gorm.io/gorm"
)

type TranslationStore interface {
	Save(translation *Translation) *Translation
	GetAuto(text, langFrom, langTo string) *Translation
}

type dbTranslationStore struct {
	db *gorm.DB
}

func (s *dbTranslationStore) Save(translation *Translation) *Translation {
	s.db.Create(translation)
	return translation
}

func (s *dbTranslationStore) GetAuto(text, langFrom, langTo string) *Translation {
	translation := &Translation{}
	conds := Translation{Text: text, LangFrom: langFrom, LangTo: langTo, Manual: false}

	if s.db.First(translation, conds).RowsAffected > 0 {
		return translation
	}

	return nil
}

func NewTranslationStore(db *gorm.DB) TranslationStore {
	return &dbTranslationStore{db}
}
