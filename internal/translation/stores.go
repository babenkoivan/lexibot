package translation

import (
	"gorm.io/gorm"
)

type TranslationFilter struct {
	Text        string
	Translation string
	LangFrom    string
	LangTo      string
	Manual      bool
}

type TranslationStore interface {
	Create(translation *Translation)
	Get(filter TranslationFilter) *Translation
	GetOrInit(filter TranslationFilter) *Translation
	GetOrCreate(filter TranslationFilter) *Translation
	IsAttached(translationID uint64, userID int) bool
	Attach(translationID uint64, userID int)
	Detach(translationID uint64, userID int)
}

type dbTranslationStore struct {
	db *gorm.DB
}

func (s *dbTranslationStore) Create(translation *Translation) {
	s.db.Create(translation)
}

func (s *dbTranslationStore) Get(filter TranslationFilter) *Translation {
	translation := &Translation{}

	if s.db.First(translation, filter).RowsAffected > 0 {
		return translation
	}

	return nil
}

func (s *dbTranslationStore) GetOrInit(filter TranslationFilter) *Translation {
	translation := &Translation{}
	s.db.FirstOrInit(translation, filter)
	return translation
}

func (s *dbTranslationStore) GetOrCreate(filter TranslationFilter) *Translation {
	translation := &Translation{}
	s.db.FirstOrCreate(translation, filter)
	return translation
}

func (s *dbTranslationStore) IsAttached(translationID uint64, userID int) bool {
	return s.db.Take(&UserTranslation{userID, translationID}).RowsAffected > 0
}

func (s *dbTranslationStore) Attach(translationID uint64, userID int) {
	s.db.Create(&UserTranslation{userID, translationID})
}

func (s *dbTranslationStore) Detach(translationID uint64, userID int) {
	s.db.Delete(&UserTranslation{userID, translationID})
}

func NewTranslationStore(db *gorm.DB) TranslationStore {
	return &dbTranslationStore{db}
}
