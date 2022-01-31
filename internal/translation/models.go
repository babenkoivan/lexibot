package translation

import "time"

type Translation struct {
	ID          uint64 `gorm:"primaryKey"`
	Text        string
	Translation string
	LangFrom    string
	LangTo      string
	Manual      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserTranslation struct {
	UserID        int    `gorm:"primaryKey"`
	TranslationID uint64 `gorm:"primaryKey"`
}
