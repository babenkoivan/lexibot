package training

import "time"

type Score struct {
	UserID        int    `gorm:"primaryKey"`
	TranslationID uint64 `gorm:"primaryKey"`
	Score         int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TaskResult struct {
	UserID        int    `gorm:"primaryKey"`
	TranslationID uint64 `gorm:"primaryKey"`
	Success       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
