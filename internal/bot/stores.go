package bot

import "gorm.io/gorm"

type HistoryStore interface {
	Save(h *HistoryMessage) *HistoryMessage
	GetLastMessage(userId int) *HistoryMessage
}

type dbHistoryStore struct {
	db *gorm.DB
}

func (s *dbHistoryStore) Save(h *HistoryMessage) *HistoryMessage {
	s.db.Create(h)
	return h
}

func (s *dbHistoryStore) GetLastMessage(userId int) *HistoryMessage {
	h := &HistoryMessage{}
	s.db.First(h, userId)
	return h
}

func NewHistoryStore(db *gorm.DB) HistoryStore {
	return &dbHistoryStore{db: db}
}
