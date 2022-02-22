package settings

import (
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lexibot/internal/localization"
	"strconv"
	"time"
)

const (
	CacheExpiration = 5 * time.Minute
	CacheCleanup    = 10 * time.Minute
)

type SettingsStore interface {
	localization.LocaleStore
	Save(settings *Settings) *Settings
	FirstOrInit(userID int) *Settings
}

type dbSettingsStore struct {
	db         *gorm.DB
	cacheStore *cache.Cache
}

func (s *dbSettingsStore) Save(settings *Settings) *Settings {
	s.cacheStore.Delete(strconv.Itoa(settings.UserID))

	if settings.UpdatedAt.IsZero() {
		settings.UpdatedAt = time.Now()
	}

	s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"lang_ui":            settings.LangUI,
			"lang_dict":          settings.LangDict,
			"auto_translate":     settings.AutoTranslate,
			"words_per_training": settings.WordsPerTraining,
			"updated_at":         settings.UpdatedAt,
		}),
	}).Create(settings)

	return settings
}

func (s *dbSettingsStore) FirstOrInit(userID int) *Settings {
	if cached, found := s.cacheStore.Get(strconv.Itoa(userID)); found {
		return cached.(*Settings)
	}

	settings := &Settings{}
	s.db.FirstOrInit(settings, Settings{UserID: userID})
	s.cacheStore.Set(strconv.Itoa(settings.UserID), settings, cache.DefaultExpiration)

	return settings
}

func (s *dbSettingsStore) Locale(userID int) string {
	settings := s.FirstOrInit(userID)
	return settings.LangUI
}

func NewDBSettingsStore(db *gorm.DB) *dbSettingsStore {
	cacheStore := cache.New(CacheExpiration, CacheCleanup)
	return &dbSettingsStore{db, cacheStore}
}
