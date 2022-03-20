package translation_test

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"lexibot/internal/testkit"
	"lexibot/internal/translation"
	"regexp"
	"testing"
	"time"
)

func TestDBTranslationStore_Save(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	want := newDummyTranslation(0)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `translations` (`user_id`,`text`,`translation`,`lang_from`,`lang_to`,"+
		"`manual`,`score`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)")).
		WithArgs(want.UserID, want.Text, want.Translation, want.LangFrom, want.LangTo, want.Manual, want.Score,
			want.CreatedAt, want.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	got := store.Save(want)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDbTranslationStore_Delete(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	want := newDummyTranslation(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `translations` WHERE `translations`.`id` = ?")).
		WithArgs(want.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	store.Delete(want)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTranslationStore_First(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	want := newDummyTranslation(1)

	for n, c := range map[string]struct {
		cond  translation.TranslationQueryCond
		where string
		args  []driver.Value
	}{
		"without ids":              {translation.WithoutIDs([]int{1}), "id NOT IN (?)", []driver.Value{1}},
		"with user id":             {translation.WithUserID(1), "user_id = ?", []driver.Value{1}},
		"with text":                {translation.WithText("bunt"), "text = ?", []driver.Value{"bunt"}},
		"with strict text":         {translation.WithTextStrict("bunt"), "BINARY text = ?", []driver.Value{"bunt"}},
		"with text or translation": {translation.WithTextOrTranslation("bunt"), "text = ? OR translation = ?", []driver.Value{"bunt", "bunt"}},
		"with lang from":           {translation.WithLangFrom("de"), "lang_from = ?", []driver.Value{"de"}},
		"with lang to":             {translation.WithLangTo("en"), "lang_to = ?", []driver.Value{"en"}},
		"with manual":              {translation.WithManual(true), "manual = ?", []driver.Value{true}},
	} {
		t.Run("first "+n, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `translations` WHERE " + c.where + " ORDER BY `translations`.`id` LIMIT 1")).
				WithArgs(c.args...).
				WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "text", "translation", "lang_from", "lang_to",
					"manual", "score", "created_at", "updated_at"}).AddRow(want.ID, want.UserID, want.Text, want.Translation,
					want.LangFrom, want.LangTo, want.Manual, want.Score, want.CreatedAt, want.UpdatedAt))

			got := store.First(c.cond)

			assert.Equal(t, want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBTranslationStore_Find(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	want := newDummyTranslation(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `translations` ORDER BY `score` LIMIT 1")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "text", "translation", "lang_from", "lang_to",
			"manual", "score", "created_at", "updated_at"}).AddRow(want.ID, want.UserID, want.Text, want.Translation,
			want.LangFrom, want.LangTo, want.Manual, want.Score, want.CreatedAt, want.UpdatedAt))

	got := store.Find(translation.WithLimit(1), translation.WithLowestScore())

	require.Len(t, got, 1)
	assert.Equal(t, want, got[0])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTranslationStore_Rand(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	langFrom := "de"
	want := newDummyTranslation(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `translations` WHERE lang_from = ?")).
		WithArgs(langFrom).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(10))

	mock.ExpectQuery("^SELECT \\* FROM `translations` WHERE lang_from = \\? LIMIT 1( OFFSET [0-9]+)?").
		WithArgs("de").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "text", "translation", "lang_from", "lang_to",
			"manual", "score", "created_at", "updated_at"}).AddRow(want.ID, want.UserID, want.Text, want.Translation,
			want.LangFrom, want.LangTo, want.Manual, want.Score, want.CreatedAt, want.UpdatedAt))

	got := store.Rand(translation.WithLangFrom(langFrom), translation.WithLimit(1))

	require.Len(t, got, 1)
	assert.Equal(t, want, got[0])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTranslationStore_Count(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	manual := false
	want := int64(10)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `translations` WHERE manual = ?")).
		WithArgs(manual).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(want))

	got := store.Count(translation.WithManual(manual))

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTranslationStore_IncrementScore(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	userID := 1
	translationID := 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `translations` SET `score`=score + ?,`updated_at`=? WHERE id = ? AND user_id = ?")).
		WithArgs(1, sqlmock.AnyArg(), translationID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	store.IncrementScore(translationID, userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTranslationStore_DecrementScore(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	userID := 1
	translationID := 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `translations` SET `score`=score - ?,`updated_at`=? WHERE id = ? AND user_id = ?")).
		WithArgs(1, sqlmock.AnyArg(), translationID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	store.DecrementScore(translationID, userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTranslationStore_AutoDecrementScore(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `translations` SET `score`=score - ?,`updated_at`=? WHERE (score > ?) "+
		"AND updated_at < ?")).
		WithArgs(1, sqlmock.AnyArg(), 0, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	store.AutoDecrementScore(time.Hour)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func newDummyTranslation(id int) *translation.Translation {
	return &translation.Translation{
		ID:          id,
		UserID:      1,
		Text:        "bunt",
		Translation: "colorful",
		LangFrom:    "de",
		LangTo:      "en",
		Manual:      false,
		Score:       0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
