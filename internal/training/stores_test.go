package training_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lexibot/internal/training"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestDBTaskStore_Save(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	want := newDummyTask()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tasks` (`user_id`,`translation_id`,`question`,`answer`,`hints`,"+
		"`score`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(want.UserID, want.TranslationID, want.Question, want.Answer, strings.Join(want.Hints, ","),
			want.Score, want.CreatedAt, want.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	got := store.Save(want)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTaskStore_Cleanup(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	userID := 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `tasks` WHERE user_id = ?")).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store.Cleanup(userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTaskStore_Count(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	userID := 1
	want := int64(10)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `tasks` WHERE user_id = ?")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(10))

	got := store.Count(userID)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTaskStore_IncrementScore(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	stored := newDummyTask()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `tasks` SET `score`=score + ?,`updated_at`=? WHERE `user_id` = ? AND "+
		"`translation_id` = ?")).
		WithArgs(1, sqlmock.AnyArg(), stored.UserID, stored.TranslationID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store.IncrementScore(stored)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTaskStore_DecrementScore(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	stored := newDummyTask()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `tasks` SET `score`=score - ?,`updated_at`=? WHERE `user_id` = ? AND "+
		"`translation_id` = ?")).
		WithArgs(1, sqlmock.AnyArg(), stored.UserID, stored.TranslationID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store.DecrementScore(stored)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTaskStore_TotalPositiveScore(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	userID := 1
	want := int64(5)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `tasks` WHERE user_id = ? AND (score > ?)")).
		WithArgs(userID, 0).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(want))

	got := store.TotalPositiveScore(userID)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}))
	require.NoError(t, err)

	return conn, mock, db
}

func newDummyTask() *training.Task {
	return &training.Task{
		UserID:        1,
		TranslationID: 2,
		Question:      "bunt",
		Answer:        "colorful",
		Hints:         []string{"brot", "mund"},
		Score:         0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
