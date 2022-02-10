package translation

import (
	"gorm.io/gorm"
)

type translationFilter struct {
	id                *uint64
	text              *string
	translation       *string
	textOrTranslation *string
	langFrom          *string
	langTo            *string
	manual            *bool
	userID            *int
}

func makeTranslationFilter(conds []func(*translationFilter)) *translationFilter {
	filter := &translationFilter{}

	for _, c := range conds {
		c(filter)
	}

	return filter
}

func WhereID(ID uint64) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.id = &ID
	}
}

func WhereText(text string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.text = &text
	}
}

func WhereTranslation(translation string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.translation = &translation
	}
}

func WhereTextOrTranslation(textOrTranslation string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.textOrTranslation = &textOrTranslation
	}
}

func WhereLangFrom(langFrom string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.langFrom = &langFrom
	}
}

func WhereLangTo(langTo string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.langTo = &langTo
	}
}

func WhereManual(manual bool) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.manual = &manual
	}
}

func WhereUserID(userID int) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.userID = &userID
	}
}

type TranslationStore interface {
	Save(translation *Translation) *Translation
	First(conds ...func(*translationFilter)) *Translation
}

type dbTranslationStore struct {
	db *gorm.DB
}

func (s *dbTranslationStore) Save(translation *Translation) *Translation {
	s.db.Create(translation)
	return translation
}

func (s *dbTranslationStore) First(conds ...func(*translationFilter)) *Translation {
	filter := makeTranslationFilter(conds)
	translation := &Translation{}

	if s.applyFilter(filter).First(translation).RowsAffected > 0 {
		return translation
	}

	return nil
}

func (s *dbTranslationStore) applyFilter(filter *translationFilter) *gorm.DB {
	db := s.db

	if filter.id != nil {
		db = db.Where("ID = ?", filter.id)
	}

	if filter.text != nil {
		db = db.Where("text = ?", filter.text)
	}

	if filter.translation != nil {
		db = db.Where("translation = ?", filter.translation)
	}

	if filter.textOrTranslation != nil {
		db = db.Where("text = ? OR translation = ?", filter.textOrTranslation, filter.textOrTranslation)
	}

	if filter.langFrom != nil {
		db = db.Where("lang_from = ?", filter.langFrom)
	}

	if filter.langTo != nil {
		db = db.Where("lang_to = ?", filter.langTo)
	}

	if filter.manual != nil {
		db = db.Where("manual = ?", filter.manual)
	}

	if filter.userID != nil {
		db = db.Joins("inner join scores on scores.translation_id = translations.id")
		db = db.Where("user_id = ?", filter.userID)
	}

	return db
}

func NewTranslationStore(db *gorm.DB) TranslationStore {
	return &dbTranslationStore{db}
}
