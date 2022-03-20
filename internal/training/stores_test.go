package training_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"lexibot/internal/testkit"
	"lexibot/internal/training"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestDBTaskStore_Save(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
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
	conn, mock, db := testkit.MockDB(t)
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
	conn, mock, db := testkit.MockDB(t)
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

func TestDBTaskStore_CorrectCount(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	userID := 1
	want := int64(5)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `tasks` WHERE user_id = ? AND score > 0")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(want))

	got := store.CorrectCount(userID)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTaskStore_IncrementScore(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	userID := 1
	translationID := 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `tasks` SET `score`=score + ?,`updated_at`=? WHERE user_id = ? AND "+
		"translation_id = ?")).
		WithArgs(1, sqlmock.AnyArg(), userID, translationID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	store.IncrementScore(translationID, userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTaskStore_DecrementScore(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	userID := 1
	translationID := 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `tasks` SET `score`=score - ?,`updated_at`=? WHERE user_id = ? AND "+
		"translation_id = ?")).
		WithArgs(1, sqlmock.AnyArg(), userID, translationID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store.DecrementScore(translationID, userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTaskStore_TranslationIDs(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := training.NewDBTaskStore(db)
	stored := newDummyTask()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `translation_id` FROM `tasks` WHERE user_id = ?")).
		WithArgs(stored.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"translation_id"}).AddRow(stored.TranslationID))

	want := []int{stored.TranslationID}
	got := store.TranslationIDs(stored.UserID)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
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
