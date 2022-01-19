package bot

import "time"

type HistoryMessage struct {
	UserId    uint
	MessageId string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
