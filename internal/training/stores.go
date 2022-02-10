package training

import "gorm.io/gorm"

type ScoreStore interface {
	Create(translationID uint64, userID int) *Score
	Delete(translationID uint64, userID int)
	LowestScore(userID int) *Score
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

func (s *dbScoreStore) LowestScore(userID int) *Score {
	score := &Score{}

	if s.db.Order("score asc").Where("user_id = ?", userID).First(&score).RowsAffected > 0 {
		return score
	}

	return nil
}

func NewScoreStore(db *gorm.DB) ScoreStore {
	return &dbScoreStore{db}
}
