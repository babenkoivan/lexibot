package translation_test

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lexibot/internal/translation"
	"regexp"
	"testing"
	"time"
)

func TestDBTranslationStore_Save(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	want := newDummyTranslation(0)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `translations` (`text`,`translation`,`lang_from`,`lang_to`,`manual`,"+
		"`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(want.Text, want.Translation, want.LangFrom, want.LangTo, want.Manual, want.CreatedAt, want.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	got := store.Save(want)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTranslationStore_First(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	want := newDummyTranslation(1)

	simpleWhereCases := map[string]struct {
		cond  translation.TranslationQueryCond
		where string
		args  []driver.Value
	}{
		"with id":                  {translation.WithID(1), "id = ?", []driver.Value{1}},
		"without id":               {translation.WithoutID(1), "id != ?", []driver.Value{1}},
		"with text":                {translation.WithText("bunt"), "text = ?", []driver.Value{"bunt"}},
		"with translation":         {translation.WithTranslation("colorful"), "translation = ?", []driver.Value{"colorful"}},
		"with text or translation": {translation.WithTextOrTranslation("bunt"), "text = ? OR translation = ?", []driver.Value{"bunt", "bunt"}},
		"with lang from":           {translation.WithLangFrom("de"), "lang_from = ?", []driver.Value{"de"}},
		"with lang to":             {translation.WithLangTo("en"), "lang_to = ?", []driver.Value{"en"}},
		"with manual":              {translation.WithManual(true), "manual = ?", []driver.Value{true}},
	}

	for n, c := range simpleWhereCases {
		t.Run("first "+n, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `translations` WHERE " + c.where + " ORDER BY `translations`.`id` LIMIT 1")).
				WithArgs(c.args...).
				WillReturnRows(sqlmock.NewRows([]string{"id", "text", "translation", "lang_from", "lang_to", "manual",
					"created_at", "updated_at"}).AddRow(want.ID, want.Text, want.Translation, want.LangFrom, want.LangTo,
					want.Manual, want.CreatedAt, want.UpdatedAt))

			got := store.First(c.cond)

			assert.Equal(t, want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("first with user id", func(t *testing.T) {
		userID := 1

		mock.ExpectQuery(regexp.QuoteMeta("SELECT `translations`.`id`,`translations`.`text`,`translations`." +
			"`translation`,`translations`.`lang_from`,`translations`.`lang_to`,`translations`.`manual`,`translations`." +
			"`created_at`,`translations`.`updated_at` FROM `translations` INNER JOIN scores ON scores.translation_id = " +
			"translations.id WHERE user_id = ? ORDER BY `translations`.`id` LIMIT 1")).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "text", "translation", "lang_from", "lang_to", "manual",
				"created_at", "updated_at"}).AddRow(want.ID, want.Text, want.Translation, want.LangFrom, want.LangTo,
				want.Manual, want.CreatedAt, want.UpdatedAt))

		got := store.First(translation.WithUserID(userID))

		assert.Equal(t, want, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBTranslationStore_Rand(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBTranslationStore(db)
	langFrom := "de"
	want := newDummyTranslation(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `translations` WHERE lang_from = ?")).
		WithArgs(langFrom).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(10))

	mock.ExpectQuery("^SELECT \\* FROM `translations` WHERE lang_from = \\? LIMIT 1( OFFSET [0-9]+)?").
		WithArgs("de").
		WillReturnRows(sqlmock.NewRows([]string{"id", "text", "translation", "lang_from", "lang_to", "manual",
			"created_at", "updated_at"}).AddRow(want.ID, want.Text, want.Translation, want.LangFrom, want.LangTo,
			want.Manual, want.CreatedAt, want.UpdatedAt))

	got := store.Rand(translation.WithLangFrom(langFrom), translation.WithLimit(1))

	require.Len(t, got, 1)
	assert.Equal(t, want, got[0])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBTranslationStore_Count(t *testing.T) {
	conn, mock, db := setup(t)
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

func TestDBScoreStore_Save(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBScoreStore(db)
	userID := 1
	translationID := 2
	score := 0

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `scores` (`user_id`,`translation_id`,`score`,`created_at`,"+
		"`updated_at`) VALUES (?,?,?,?,?)")).
		WithArgs(userID, translationID, score, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	got := store.Save(translationID, userID)

	assert.Equal(t, userID, got.UserID)
	assert.Equal(t, translationID, got.TranslationID)
	assert.Equal(t, score, got.Score)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBScoreStore_Delete(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBScoreStore(db)
	userID := 1
	translationID := 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `scores` WHERE translation_id = ? AND user_id = ?")).
		WithArgs(translationID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store.Delete(translationID, userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBScoreStore_Increment(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBScoreStore(db)
	userID := 1
	translationID := 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `scores` SET `score`=score + ?,`updated_at`=? WHERE translation_id = ? "+
		"AND user_id = ?")).
		WithArgs(1, sqlmock.AnyArg(), translationID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store.Increment(translationID, userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBScoreStore_Decrement(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBScoreStore(db)
	userID := 1
	translationID := 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `scores` SET `score`=score - ?,`updated_at`=? WHERE translation_id = ? "+
		"AND user_id = ?")).
		WithArgs(1, sqlmock.AnyArg(), translationID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store.Decrement(translationID, userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBScoreStore_AutoDecrement(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBScoreStore(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `scores` SET `score`=score - ?,`updated_at`=? WHERE (score > ?) "+
		"AND updated_at < ?")).
		WithArgs(1, sqlmock.AnyArg(), 0, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store.AutoDecrement(time.Hour)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBScoreStore_LowestNotTrained(t *testing.T) {
	conn, mock, db := setup(t)
	defer conn.Close()

	store := translation.NewDBScoreStore(db)
	want := newDummyScore()
	userID := 1
	langDict := "de"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `scores`.`user_id`,`scores`.`translation_id`,`scores`.`score`,"+
		"`scores`.`created_at`,`scores`.`updated_at` FROM `scores` INNER JOIN translations ON translations.id = "+
		"scores.translation_id LEFT JOIN tasks ON tasks.translation_id = scores.translation_id AND tasks.user_id = "+
		"scores.user_id WHERE (scores.user_id = ?) AND translations.lang_from = ? AND tasks.translation_id IS NULL "+
		"ORDER BY scores.score asc LIMIT 1")).
		WithArgs(userID, langDict).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "translation_id", "score", "created_at", "updated_at"}).
			AddRow(want.UserID, want.TranslationID, want.Score, want.CreatedAt, want.UpdatedAt))

	got := store.LowestNotTrained(userID, langDict)

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

func newDummyScore() *translation.Score {
	return &translation.Score{
		UserID:        1,
		TranslationID: 2,
		Score:         0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func newDummyTranslation(id int) *translation.Translation {
	return &translation.Translation{
		ID:          id,
		Text:        "bunt",
		Translation: "colorful",
		LangFrom:    "de",
		LangTo:      "en",
		Manual:      false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
