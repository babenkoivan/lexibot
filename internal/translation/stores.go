package translation

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lexibot/internal/utils"
	"time"
)

const orderAsc = "asc"

type translationQuery struct {
	notIDs            *[]int
	userID            *int
	text              *string
	textOrTranslation *string
	langFrom          *string
	langTo            *string
	manual            *bool
	order             *[2]string
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

func WithoutIDs(IDs []int) TranslationQueryCond {
	return func(query *translationQuery) {
		query.notIDs = &IDs
	}
}

func WithUserID(userID int) TranslationQueryCond {
	return func(query *translationQuery) {
		query.userID = &userID
	}
}

func WithText(text string) TranslationQueryCond {
	return func(query *translationQuery) {
		query.text = &text
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

func WithLowestScore() TranslationQueryCond {
	return func(query *translationQuery) {
		query.order = &[2]string{"score", orderAsc}
	}
}

func WithLimit(limit int) TranslationQueryCond {
	return func(query *translationQuery) {
		query.limit = &limit
	}
}

type TranslationStore interface {
	Save(transl *Translation) *Translation
	Delete(transl *Translation)
	First(conds ...TranslationQueryCond) *Translation
	Find(conds ...TranslationQueryCond) []*Translation
	Rand(conds ...TranslationQueryCond) []*Translation
	Count(conds ...TranslationQueryCond) int64
	IncrementScore(id, userID int)
	DecrementScore(id, userID int)
	AutoDecrementScore(after time.Duration)
}

type dbTranslationStore struct {
	db *gorm.DB
}

func (s *dbTranslationStore) Save(transl *Translation) *Translation {
	s.db.Create(transl)
	return transl
}

func (s *dbTranslationStore) Delete(transl *Translation) {
	s.db.Delete(transl)
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

func (s *dbTranslationStore) IncrementScore(id, userID int) {
	s.db.Model(&Translation{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("score", gorm.Expr("score + ?", 1))
}

func (s *dbTranslationStore) DecrementScore(id, userID int) {
	s.db.Model(&Translation{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("score", gorm.Expr("score - ?", 1))
}

func (s *dbTranslationStore) AutoDecrementScore(after time.Duration) {
	s.db.Model(&Translation{}).
		Where("score > ?", 0).
		Where("updated_at < ?", time.Now().Add(-after)).
		Update("score", gorm.Expr("score - ?", 1))
}

func (s *dbTranslationStore) withQuery(query *translationQuery) *gorm.DB {
	tx := s.db

	if query.notIDs != nil {
		tx = tx.Where("id NOT IN ?", *query.notIDs)
	}

	if query.userID != nil {
		tx = tx.Where("user_id = ?", *query.userID)
	}

	if query.text != nil {
		tx = tx.Where("text = ?", *query.text)
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

	if query.order != nil {
		tx = tx.Order(clause.OrderByColumn{
			Column: clause.Column{Name: query.order[0]},
			Desc:   query.order[1] != orderAsc,
		})
	}

	if query.limit != nil {
		tx = tx.Limit(*query.limit)
	}

	return tx
}

func NewDBTranslationStore(db *gorm.DB) *dbTranslationStore {
	return &dbTranslationStore{db}
}
