package translation

import (
	"gorm.io/gorm"
)

type translationFilter struct {
	Text        *string
	Translation *string
	LangFrom    *string
	LangTo      *string
	Manual      *bool
	UserID      *int
}

func makeTranslationFilter(conds []func(*translationFilter)) *translationFilter {
	filter := &translationFilter{}

	for _, c := range conds {
		c(filter)
	}

	return filter
}

func WithText(text string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.Text = &text
	}
}

func WithTranslation(translation string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.Translation = &translation
	}
}

func WithLangFrom(langFrom string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.LangFrom = &langFrom
	}
}

func WithLangTo(langTo string) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.LangTo = &langTo
	}
}

func WithManual(manual bool) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.Manual = &manual
	}
}

func WithUserID(userID int) func(*translationFilter) {
	return func(filter *translationFilter) {
		filter.UserID = &userID
	}
}

type TranslationStore interface {
	Save(translation *Translation) *Translation
	First(conds ...func(*translationFilter)) *Translation
	Attach(translationID uint64, userID int)
	Detach(translationID uint64, userID int)
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

func (s *dbTranslationStore) Attach(translationID uint64, userID int) {
	s.db.Create(&UserTranslation{userID, translationID})
}

func (s *dbTranslationStore) Detach(translationID uint64, userID int) {
	s.db.Delete(&UserTranslation{userID, translationID})
}

func (s *dbTranslationStore) applyFilter(filter *translationFilter) *gorm.DB {
	db := s.db

	if filter.Text != nil {
		db = db.Where("text = ?", filter.Text)
	}

	if filter.Translation != nil {
		db = db.Where("translation = ?", filter.Translation)
	}

	if filter.LangFrom != nil {
		db = db.Where("lang_from = ?", filter.LangFrom)
	}

	if filter.LangTo != nil {
		db = db.Where("lang_to = ?", filter.LangTo)
	}

	if filter.Manual != nil {
		db = db.Where("manual = ?", filter.Manual)
	}

	if filter.UserID != nil {
		db = db.Joins("inner join user_translations on user_translations.translation_id = translations.id")
		db = db.Where("user_id = ?", filter.UserID)
	}

	return db
}

func NewTranslationStore(db *gorm.DB) TranslationStore {
	return &dbTranslationStore{db}
}
