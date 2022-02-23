package settings_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"lexibot/internal/settings"
	"lexibot/internal/testkit"
	"regexp"
	"testing"
	"time"
)

func TestDBSettingsStore_Save(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := settings.NewDBSettingsStore(db)
	want := newDummySettings()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `settings` (`lang_ui`,`lang_dict`,`auto_translate`,"+
		"`words_per_training`,`created_at`,`updated_at`,`user_id`) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY "+
		"UPDATE `auto_translate`=?,`lang_dict`=?,`lang_ui`=?,`updated_at`=?,`words_per_training`=?")).
		WithArgs(want.LangUI, want.LangDict, want.AutoTranslate, want.WordsPerTraining, want.CreatedAt, want.UpdatedAt,
			want.UserID, want.AutoTranslate, want.LangDict, want.LangUI, want.UpdatedAt, want.WordsPerTraining).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	got := store.Save(want)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBSettingsStore_FirstOrInit(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := settings.NewDBSettingsStore(db)
	want := newDummySettings()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `settings` WHERE `settings`.`user_id` = ? ORDER BY " +
		"`settings`.`user_id` LIMIT 1")).
		WithArgs(want.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"lang_ui", "lang_dict", "auto_translate", "words_per_training",
			"created_at", "updated_at", "user_id"}).AddRow(want.LangUI, want.LangDict, want.AutoTranslate,
			want.WordsPerTraining, want.CreatedAt, want.UpdatedAt, want.UserID))

	got := store.FirstOrInit(want.UserID)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBSettingsStore_Locale(t *testing.T) {
	conn, mock, db := testkit.MockDB(t)
	defer conn.Close()

	store := settings.NewDBSettingsStore(db)
	stored := newDummySettings()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `settings` WHERE `settings`.`user_id` = ? ORDER BY " +
		"`settings`.`user_id` LIMIT 1")).
		WithArgs(stored.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"lang_ui", "lang_dict", "auto_translate", "words_per_training",
			"created_at", "updated_at", "user_id"}).AddRow(stored.LangUI, stored.LangDict, stored.AutoTranslate,
			stored.WordsPerTraining, stored.CreatedAt, stored.UpdatedAt, stored.UserID))

	want := stored.LangUI
	got := store.Locale(stored.UserID)

	assert.Equal(t, want, got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func newDummySettings() *settings.Settings {
	return &settings.Settings{
		UserID:           1,
		LangUI:           "en",
		LangDict:         "de",
		AutoTranslate:    true,
		WordsPerTraining: 10,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
