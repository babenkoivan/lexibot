package translation

import (
	"gorm.io/gorm"
	"lexibot/internal/utils"
)

type translationQuery struct {
	id                *uint64
	notID             *uint64
	text              *string
	translation       *string
	textOrTranslation *string
	langFrom          *string
	langTo            *string
	manual            *bool
	userID            *int
	limit             *int
}

func makeTranslationQuery(conds []func(*translationQuery)) *translationQuery {
	query := &translationQuery{}

	for _, c := range conds {
		c(query)
	}

	return query
}

func WithID(ID uint64) func(*translationQuery) {
	return func(query *translationQuery) {
		query.id = &ID
	}
}

func WithoutID(ID uint64) func(*translationQuery) {
	return func(query *translationQuery) {
		query.notID = &ID
	}
}

func WithText(text string) func(*translationQuery) {
	return func(query *translationQuery) {
		query.text = &text
	}
}

func WithTranslation(translation string) func(*translationQuery) {
	return func(query *translationQuery) {
		query.translation = &translation
	}
}

func WithTextOrTranslation(textOrTranslation string) func(*translationQuery) {
	return func(query *translationQuery) {
		query.textOrTranslation = &textOrTranslation
	}
}

func WithLangFrom(langFrom string) func(*translationQuery) {
	return func(query *translationQuery) {
		query.langFrom = &langFrom
	}
}

func WithLangTo(langTo string) func(*translationQuery) {
	return func(query *translationQuery) {
		query.langTo = &langTo
	}
}

func WithManual(manual bool) func(*translationQuery) {
	return func(query *translationQuery) {
		query.manual = &manual
	}
}

func WithUserID(userID int) func(*translationQuery) {
	return func(query *translationQuery) {
		query.userID = &userID
	}
}

func WithLimit(limit int) func(*translationQuery) {
	return func(query *translationQuery) {
		query.limit = &limit
	}
}

type TranslationStore interface {
	Save(transl *Translation) *Translation
	First(conds ...func(*translationQuery)) *Translation
	Rand(conds ...func(*translationQuery)) []*Translation
}

type dbTranslationStore struct {
	db *gorm.DB
}

func (s *dbTranslationStore) Save(transl *Translation) *Translation {
	s.db.Create(transl)
	return transl
}

func (s *dbTranslationStore) First(conds ...func(*translationQuery)) *Translation {
	transl := &Translation{}
	query := makeTranslationQuery(conds)

	if s.withQuery(query).First(transl).RowsAffected > 0 {
		return transl
	}

	return nil
}

func (s *dbTranslationStore) Rand(conds ...func(*translationQuery)) []*Translation {
	query := makeTranslationQuery(conds)
	db := s.withQuery(query)

	var count int64
	db.Model(&Translation{}).Count(&count)

	var transl []*Translation
	offset := utils.NewRand().Intn(int(count) - 1)
	db.Offset(offset).Find(&transl)

	return transl
}

func (s *dbTranslationStore) withQuery(query *translationQuery) *gorm.DB {
	db := s.db

	if query.id != nil {
		db = db.Where("ID = ?", *query.id)
	}

	if query.notID != nil {
		db = db.Where("ID != ?", *query.notID)
	}

	if query.text != nil {
		db = db.Where("text = ?", *query.text)
	}

	if query.translation != nil {
		db = db.Where("translation = ?", *query.translation)
	}

	if query.textOrTranslation != nil {
		db = db.Where("text = ? OR translation = ?", *query.textOrTranslation, *query.textOrTranslation)
	}

	if query.langFrom != nil {
		db = db.Where("lang_from = ?", *query.langFrom)
	}

	if query.langTo != nil {
		db = db.Where("lang_to = ?", *query.langTo)
	}

	if query.manual != nil {
		db = db.Where("manual = ?", *query.manual)
	}

	if query.userID != nil {
		db = db.Joins("inner join scores on scores.translation_id = translations.id")
		db = db.Where("user_id = ?", *query.userID)
	}

	if query.limit != nil {
		db = db.Limit(*query.limit)
	}

	return db
}

func NewTranslationStore(db *gorm.DB) TranslationStore {
	return &dbTranslationStore{db}
}

type ScoreStore interface {
	Create(translationID uint64, userID int) *Score
	Delete(translationID uint64, userID int)
	LowestScore(userID int, langDict string) *Score
}

type dbScoreStore struct {
	db *gorm.DB
}

func (s *dbScoreStore) Create(translationID uint64, userID int) *Score {
	score := &Score{UserID: userID, TranslationID: translationID}
	s.db.Create(score)
	return score
}

func (s *dbScoreStore) Delete(translationID uint64, userID int) {
	s.db.Delete(&Score{UserID: userID, TranslationID: translationID})
}

func (s *dbScoreStore) LowestScore(userID int, langDict string) *Score {
	score := &Score{}

	res := s.db.
		Order("score asc").
		Joins("inner join translations on translations.id = scores.translation_id").
		Where("user_id = ?", userID).
		Where("lang_from = ?", langDict).
		Take(&score)

	if res.RowsAffected > 0 {
		return score
	}

	return nil
}

func NewScoreStore(db *gorm.DB) ScoreStore {
	return &dbScoreStore{db}
}
