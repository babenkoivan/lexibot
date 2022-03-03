package translation

import (
	"gorm.io/gorm"
	"lexibot/internal/utils"
	"time"
)

type translationQuery struct {
	id                *int
	notID             *int
	text              *string
	translation       *string
	textOrTranslation *string
	langFrom          *string
	langTo            *string
	manual            *bool
	userID            *int
	limit             *int
}

type TranslationQueryCond func(*translationQuery)

func MakeTranslationQuery(conds []TranslationQueryCond) *translationQuery {
	query := &translationQuery{}

	for _, c := range conds {
		c(query)
	}

	return query
}

func WithID(ID int) TranslationQueryCond {
	return func(query *translationQuery) {
		query.id = &ID
	}
}

func WithoutID(ID int) TranslationQueryCond {
	return func(query *translationQuery) {
		query.notID = &ID
	}
}

func WithText(text string) TranslationQueryCond {
	return func(query *translationQuery) {
		query.text = &text
	}
}

func WithTranslation(translation string) TranslationQueryCond {
	return func(query *translationQuery) {
		query.translation = &translation
	}
}

func WithTextOrTranslation(textOrTranslation string) TranslationQueryCond {
	return func(query *translationQuery) {
		query.textOrTranslation = &textOrTranslation
	}
}

func WithLangFrom(langFrom string) TranslationQueryCond {
	return func(query *translationQuery) {
		query.langFrom = &langFrom
	}
}

func WithLangTo(langTo string) TranslationQueryCond {
	return func(query *translationQuery) {
		query.langTo = &langTo
	}
}

func WithManual(manual bool) TranslationQueryCond {
	return func(query *translationQuery) {
		query.manual = &manual
	}
}

func WithUserID(userID int) TranslationQueryCond {
	return func(query *translationQuery) {
		query.userID = &userID
	}
}

func WithLimit(limit int) TranslationQueryCond {
	return func(query *translationQuery) {
		query.limit = &limit
	}
}

type TranslationStore interface {
	Save(transl *Translation) *Translation
	First(conds ...TranslationQueryCond) *Translation
	Find(conds ...TranslationQueryCond) []*Translation
	Rand(conds ...TranslationQueryCond) []*Translation
	Count(conds ...TranslationQueryCond) int64
}

type dbTranslationStore struct {
	db *gorm.DB
}

func (s *dbTranslationStore) Save(transl *Translation) *Translation {
	s.db.Create(transl)
	return transl
}

func (s *dbTranslationStore) First(conds ...TranslationQueryCond) *Translation {
	transl := &Translation{}
	query := MakeTranslationQuery(conds)

	if s.withQuery(query).First(transl).RowsAffected > 0 {
		return transl
	}

	return nil
}

func (s *dbTranslationStore) Find(conds ...TranslationQueryCond) []*Translation {
	var transl []*Translation
	query := MakeTranslationQuery(conds)
	s.withQuery(query).Find(&transl)
	return transl
}

func (s *dbTranslationStore) Rand(conds ...TranslationQueryCond) []*Translation {
	query := MakeTranslationQuery(conds)
	count := s.Count(conds...)

	offset := 0
	if count > 1 {
		offset = utils.SourcedRand().Intn(int(count) - 1)
	}

	var transl []*Translation
	s.withQuery(query).Offset(offset).Find(&transl)
	return transl
}

func (s *dbTranslationStore) Count(conds ...TranslationQueryCond) int64 {
	query := MakeTranslationQuery(conds)
	query.limit = nil

	var count int64
	s.withQuery(query).Model(&Translation{}).Count(&count)
	return count
}

func (s *dbTranslationStore) withQuery(query *translationQuery) *gorm.DB {
	tx := s.db

	if query.id != nil {
		tx = tx.Where("id = ?", *query.id)
	}

	if query.notID != nil {
		tx = tx.Where("id != ?", *query.notID)
	}

	if query.text != nil {
		tx = tx.Where("text = ?", *query.text)
	}

	if query.translation != nil {
		tx = tx.Where("translation = ?", *query.translation)
	}

	if query.textOrTranslation != nil {
		tx = tx.Where("text = ? OR translation = ?", *query.textOrTranslation, *query.textOrTranslation)
	}

	if query.langFrom != nil {
		tx = tx.Where("lang_from = ?", *query.langFrom)
	}

	if query.langTo != nil {
		tx = tx.Where("lang_to = ?", *query.langTo)
	}

	if query.manual != nil {
		tx = tx.Where("manual = ?", *query.manual)
	}

	if query.userID != nil {
		tx = tx.Joins("INNER JOIN scores ON scores.translation_id = translations.id")
		tx = tx.Where("user_id = ?", *query.userID)
	}

	if query.limit != nil {
		tx = tx.Limit(*query.limit)
	}

	return tx
}

func NewDBTranslationStore(db *gorm.DB) *dbTranslationStore {
	return &dbTranslationStore{db}
}

type ScoreStore interface {
	Save(translationID, userID int) *Score
	Delete(translationID, userID int)
	Increment(translationID, userID int)
	Decrement(translationID, userID int)
	AutoDecrement(after time.Duration)
	LowestNotTrained(userID int, langDict string) *Score
}

type dbScoreStore struct {
	db *gorm.DB
}

func (s *dbScoreStore) Save(translationID, userID int) *Score {
	score := &Score{UserID: userID, TranslationID: translationID}
	s.db.Create(score)
	return score
}

func (s *dbScoreStore) Delete(translationID, userID int) {
	s.db.Delete(&Score{}, "translation_id = ? AND user_id = ?", translationID, userID)
}

func (s *dbScoreStore) Increment(translationID, userID int) {
	s.db.Model(&Score{}).
		Where("translation_id = ? AND user_id = ?", translationID, userID).
		Update("score", gorm.Expr("score + ?", 1))
}

func (s *dbScoreStore) Decrement(translationID, userID int) {
	s.db.Model(&Score{}).
		Where("translation_id = ? AND user_id = ?", translationID, userID).
		Update("score", gorm.Expr("score - ?", 1))
}

func (s *dbScoreStore) AutoDecrement(after time.Duration) {
	s.db.Model(&Score{}).
		Where("score > ?", 0).
		Where("updated_at < ?", time.Now().Add(-after)).
		Update("score", gorm.Expr("score - ?", 1))
}

func (s *dbScoreStore) LowestNotTrained(userID int, langDict string) *Score {
	score := &Score{}

	res := s.db.
		Order("scores.score asc").
		Joins("INNER JOIN translations ON translations.id = scores.translation_id").
		Joins("LEFT JOIN tasks ON tasks.translation_id = scores.translation_id AND tasks.user_id = scores.user_id").
		Where("scores.user_id = ?", userID).
		Where("translations.lang_from = ?", langDict).
		Where("tasks.translation_id IS NULL").
		Take(&score)

	if res.RowsAffected > 0 {
		return score
	}

	return nil
}

func NewDBScoreStore(db *gorm.DB) *dbScoreStore {
	return &dbScoreStore{db}
}
