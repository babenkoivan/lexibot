package training

import "gorm.io/gorm"

type ScoreStore interface {
	Create(translationID uint64, userID int) *Score
	Delete(translationID uint64, userID int)
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

func NewScoreStore(db *gorm.DB) ScoreStore {
	return &dbScoreStore{db}
}
