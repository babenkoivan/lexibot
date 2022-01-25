package settings

import (
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lexibot/internal/locale"
	"strconv"
	"time"
)

const (
	CacheExpiration = 5 * time.Minute
	CacheCleanup    = 10 * time.Minute
)

type SettingsStore interface {
	locale.LocaleStore
	Save(settings *Settings) *Settings
	Get(userID int) *Settings
}

type dbSettingsStore struct {
	db         *gorm.DB
	cacheStore *cache.Cache
}

func (s *dbSettingsStore) Save(settings *Settings) *Settings {
	s.cacheStore.Delete(strconv.Itoa(settings.UserID))

	s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(settings)

	return settings
}

func (s *dbSettingsStore) Get(userID int) *Settings {
	if cached, found := s.cacheStore.Get(strconv.Itoa(userID)); found {
		return cached.(*Settings)
	}

	settings := &Settings{}
	s.db.FirstOrInit(settings, Settings{UserID: userID})
	s.cacheStore.Set(strconv.Itoa(settings.UserID), settings, cache.DefaultExpiration)

	return settings
}

func (s *dbSettingsStore) GetLocale(userID int) string {
	settings := s.Get(userID)
	return settings.LangUI
}

func NewSettingsStore(db *gorm.DB) SettingsStore {
	cacheStore := cache.New(CacheExpiration, CacheCleanup)
	return &dbSettingsStore{db, cacheStore}
}
