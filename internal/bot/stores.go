package bot

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HistoryStore interface {
	Save(hm *HistoryMessage) *HistoryMessage
	GetLastMessage(userID int) *HistoryMessage
}

type dbHistoryStore struct {
	db *gorm.DB
}

func (s *dbHistoryStore) Save(hm *HistoryMessage) *HistoryMessage {
	s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "content", "updated_at"}),
	}).Create(hm)

	return hm
}

func (s *dbHistoryStore) GetLastMessage(userID int) *HistoryMessage {
	hm := &HistoryMessage{}
	conds := HistoryMessage{UserID: userID}

	if s.db.First(hm, conds).RowsAffected > 0 {
		return hm
	}

	return nil
}

func NewHistoryStore(db *gorm.DB) HistoryStore {
	return &dbHistoryStore{db}
}
