package bot

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type HistoryStore interface {
	Save(hm *HistoryMessage) *HistoryMessage
	Last(userID int) *HistoryMessage
}

type dbHistoryStore struct {
	db *gorm.DB
}

func (s *dbHistoryStore) Save(hm *HistoryMessage) *HistoryMessage {
	if hm.UpdatedAt.IsZero() {
		hm.UpdatedAt = time.Now()
	}

	s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"type":       hm.Type,
			"content":    hm.Content,
			"updated_at": hm.UpdatedAt,
		}),
	}).Create(hm)

	return hm
}

func (s *dbHistoryStore) Last(userID int) *HistoryMessage {
	hm := &HistoryMessage{}

	if s.db.Take(hm, userID).RowsAffected > 0 {
		return hm
	}

	return nil
}

func NewDBHistoryStore(db *gorm.DB) *dbHistoryStore {
	return &dbHistoryStore{db}
}
