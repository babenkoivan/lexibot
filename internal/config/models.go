package config

import "time"

type Config struct {
	UserID        int `gorm:"primaryKey"`
	LangUI        string
	LangDict      string
	AutoTranslate bool `gorm:"default:1"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
