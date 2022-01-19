package bot

import (
	"encoding/json"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

type HistoryMessage struct {
	ChatID    int64 `gorm:"primaryKey"`
	Type      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (hm *HistoryMessage) TableName() string {
	return "history"
}

func NewHistoryMessage(chat *telebot.Chat, msg Message) *HistoryMessage {
	// todo error handling
	content, _ := json.Marshal(msg)
	return &HistoryMessage{ChatID: chat.ID, Type: msg.Type(), Content: string(content)}
}
