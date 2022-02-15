package bot

import (
	"encoding/json"
	"time"
)

type HistoryMessage struct {
	UserID    int `gorm:"primaryKey"`
	Type      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (hm *HistoryMessage) TableName() string {
	return "history"
}

func newHistoryMessage(userID int, msg Message) *HistoryMessage {
	// todo error handling
	content, _ := json.Marshal(msg)
	return &HistoryMessage{UserID: userID, Type: msg.Type(), Content: string(content)}
}
