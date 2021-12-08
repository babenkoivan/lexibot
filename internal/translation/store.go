package translation

import (
	"gorm.io/gorm"
)

type Store interface {
	Create(text, translation string) *Translation
	Exists(text string) bool
	Delete(ID uint64)
}

type Translation struct {
	ID          uint64 `gorm:"primarykey"`
	Text        string `gorm:"uniqueIndex:idx_translation"`
	Translation string `gorm:"uniqueIndex:idx_translation"`
}

type dbStore struct {
	db *gorm.DB
}

func (s *dbStore) Create(text, translation string) *Translation {
	t := &Translation{Text: text, Translation: translation}
	s.db.Create(t)
	return t
}

func (s *dbStore) Exists(text string) bool {
	r := s.db.Where("text = ?", text).First(&Translation{})
	return r.RowsAffected > 0
}

func (s *dbStore) Delete(ID uint64) {
	s.db.Delete(&Translation{}, ID)
}

func NewDBStore(db *gorm.DB) *dbStore {
	return &dbStore{db: db}
}
