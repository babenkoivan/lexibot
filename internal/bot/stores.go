package bot

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HistoryStore interface {
	Save(hm *HistoryMessage) *HistoryMessage
	LastMessage(chatID int64) *HistoryMessage
}

type dbHistoryStore struct {
	db *gorm.DB
}

func (s *dbHistoryStore) Save(hm *HistoryMessage) *HistoryMessage {
	s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chat_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "content", "updated_at"}),
	}).Create(hm)

	return hm
}

func (s *dbHistoryStore) LastMessage(chatID int64) *HistoryMessage {
	hm := &HistoryMessage{}
	s.db.First(hm, chatID)
	return hm
}

func NewHistoryStore(db *gorm.DB) HistoryStore {
	return &dbHistoryStore{db: db}
}
