package bot_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"lexibot/internal/bot"
	"lexibot/internal/testkit"
	"regexp"
	"testing"
	"time"
)

func TestDBHistoryStore_Save(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := bot.NewDBHistoryStore(db)
	want := newDummyHistoryMessage()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `history` (`type`,`content`,`created_at`,`updated_at`,`user_id`) "+
		"VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE `content`=?,`type`=?,`updated_at`=?")).
		WithArgs(want.Type, want.Content, want.CreatedAt, want.UpdatedAt, want.UserID, want.Content, want.Type, want.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	got := store.Save(want)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBHistoryStore_Last(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := bot.NewDBHistoryStore(db)
	want := newDummyHistoryMessage()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `history` WHERE `history`.`user_id` = ? LIMIT 1")).
		WithArgs(want.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "type", "content", "created_at", "updated_at"}).
			AddRow(want.UserID, want.Type, want.Content, want.CreatedAt, want.UpdatedAt))

	got := store.Last(want.UserID)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func newDummyHistoryMessage() *bot.HistoryMessage {
	return &bot.HistoryMessage{
		UserID:    1,
		Type:      "app.start",
		Content:   "{}",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
