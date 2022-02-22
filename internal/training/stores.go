package training

import "gorm.io/gorm"

type TaskStore interface {
	Save(task *Task) *Task
	Cleanup(userID int)
	Count(userID int) int64
	IncrementScore(task *Task)
	DecrementScore(task *Task)
	TotalPositiveScore(userID int) int64
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

func (s *dbTaskStore) IncrementScore(task *Task) {
	s.db.Model(task).Update("score", gorm.Expr("score + ?", 1))
}

func (s *dbTaskStore) DecrementScore(task *Task) {
	s.db.Model(task).Update("score", gorm.Expr("score - ?", 1))
}

func (s *dbTaskStore) TotalPositiveScore(userID int) int64 {
	var count int64
	s.db.Model(&Task{}).Where("user_id = ?", userID).Where("score > ?", 0).Count(&count)
	return count
}

func NewDBTaskStore(db *gorm.DB) *dbTaskStore {
	return &dbTaskStore{db}
}
