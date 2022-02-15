package translation

import (
	"gorm.io/gorm"
	"lexibot/internal/utils"
	"time"
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
	Count(conds ...func(*translationQuery)) int64
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
	count := s.Count(conds...)
	offset := utils.NewRand().Intn(int(count) - 1)

	var transl []*Translation
	s.withQuery(query).Offset(offset).Find(&transl)
	return transl
}

func (s *dbTranslationStore) Count(conds ...func(*translationQuery)) int64 {
	query := makeTranslationQuery(conds)
	query.limit = nil

	var count int64
	s.withQuery(query).Model(&Translation{}).Count(&count)
	return count
}

func (s *dbTranslationStore) withQuery(query *translationQuery) *gorm.DB {
	tx := s.db

	if query.id != nil {
		tx = tx.Where("ID = ?", *query.id)
	}

	if query.notID != nil {
		tx = tx.Where("ID != ?", *query.notID)
	}

	if query.text != nil {
		tx = tx.Where("text = ?", *query.text)
	}

	if query.translation != nil {
		tx = tx.Where("translation = ?", *query.translation)
	}

	if query.textOrTranslation != nil {
		tx = tx.Where("text = ? or translation = ?", *query.textOrTranslation, *query.textOrTranslation)
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
		tx = tx.Joins("inner join scores on scores.translation_id = translations.id")
		tx = tx.Where("user_id = ?", *query.userID)
	}

	if query.limit != nil {
		tx = tx.Limit(*query.limit)
	}

	return tx
}

func NewTranslationStore(db *gorm.DB) TranslationStore {
	return &dbTranslationStore{db}
}

type ScoreStore interface {
	Save(translationID uint64, userID int) *Score
	Delete(translationID uint64, userID int)
	Increment(translationID uint64, userID int)
	Decrement(translationID uint64, userID int)
	AutoDecrement(after time.Duration)
	LowestNotTrained(userID int, langDict string) *Score
}

type dbScoreStore struct {
	db *gorm.DB
}

func (s *dbScoreStore) Save(translationID uint64, userID int) *Score {
	score := &Score{UserID: userID, TranslationID: translationID}
	s.db.Create(score)
	return score
}

func (s *dbScoreStore) Delete(translationID uint64, userID int) {
	s.db.Delete(&Score{}, "translation_id = ? and user_id = ?", translationID, userID)
}

func (s *dbScoreStore) Increment(translationID uint64, userID int) {
	s.db.Model(&Score{}).
		Where("translation_id = ? and user_id = ?", translationID, userID).
		Update("score", gorm.Expr("score + ?", 1))
}

func (s *dbScoreStore) Decrement(translationID uint64, userID int) {
	s.db.Model(&Score{}).
		Where("translation_id = ? and user_id = ?", translationID, userID).
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
		Joins("inner join translations on translations.id = scores.translation_id").
		Joins("left join tasks on tasks.translation_id = scores.translation_id and tasks.user_id = scores.user_id").
		Where("scores.user_id = ?", userID).
		Where("translations.lang_from = ?", langDict).
		Where("tasks.translation_id is null").
		Take(&score)

	if res.RowsAffected > 0 {
		return score
	}

	return nil
}

func NewScoreStore(db *gorm.DB) ScoreStore {
	return &dbScoreStore{db}
}
