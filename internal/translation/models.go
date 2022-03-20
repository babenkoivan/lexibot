package translation

import "time"

type Translation struct {
	ID          int `gorm:"primaryKey"`
	UserID      int
	Text        string
	Translation string
	LangFrom    string
	LangTo      string
	Manual      bool
	Score       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
