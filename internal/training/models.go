package training

import (
	"lexibot/internal/database"
	"time"
)

type Task struct {
	UserID        int `gorm:"primaryKey"`
	TranslationID int `gorm:"primaryKey"`
	Question      string
	Answer        string
	Hints         database.StringArray
	Score         int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
