package training

import "time"

type TaskResult struct {
	UserID        int    `gorm:"primaryKey"`
	TranslationID uint64 `gorm:"primaryKey"`
	Success       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
