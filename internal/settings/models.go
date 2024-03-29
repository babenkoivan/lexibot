package settings

import "time"

type Settings struct {
	UserID           int `gorm:"primaryKey"`
	LangUI           string
	LangDict         string
	AutoTranslate    bool `gorm:"default:1"`
	WordsPerTraining int  `gorm:"default:10"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
