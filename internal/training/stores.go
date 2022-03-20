package training

import "gorm.io/gorm"

type TaskStore interface {
	Save(task *Task) *Task
	Cleanup(userID int)
	Count(userID int) int64
	CorrectCount(userID int) int64
	IncrementScore(translationId, userID int)
	DecrementScore(translationId, userID int)
	TranslationIDs(userID int) []int
}

type dbTaskStore struct {
	db *gorm.DB
}

func (s *dbTaskStore) Save(task *Task) *Task {
	s.db.Create(task)
	return task
}

func (s *dbTaskStore) Cleanup(userID int) {
	s.db.Delete(&Task{}, "user_id = ?", userID)
}

func (s *dbTaskStore) Count(userID int) int64 {
	var count int64
	s.db.Model(&Task{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (s *dbTaskStore) CorrectCount(userID int) int64 {
	var count int64
	s.db.Model(&Task{}).Where("user_id = ? AND score > 0", userID).Count(&count)
	return count
}

func (s *dbTaskStore) IncrementScore(translationId, userID int) {
	s.db.Model(&Task{}).
		Where("user_id = ? AND translation_id = ?", userID, translationId).
		Update("score", gorm.Expr("score + ?", 1))
}

func (s *dbTaskStore) DecrementScore(translationId, userID int) {
	s.db.Model(&Task{}).
		Where("user_id = ? AND translation_id = ?", userID, translationId).
		Update("score", gorm.Expr("score - ?", 1))
}

func (s *dbTaskStore) TranslationIDs(userID int) []int {
	var translationIDs []int

	s.db.Model(&Task{}).
		Select("translation_id").
		Where("user_id = ?", userID).
		Scan(&translationIDs)

	return translationIDs
}

func NewDBTaskStore(db *gorm.DB) *dbTaskStore {
	return &dbTaskStore{db}
}
