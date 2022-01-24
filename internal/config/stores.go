package config

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

type ConfigStore interface {
	locale.LocaleStore
	Save(config *Config) *Config
	Get(userID int) *Config
}

type dbConfigStore struct {
	db         *gorm.DB
	cacheStore *cache.Cache
}

func (s *dbConfigStore) Save(config *Config) *Config {
	s.cacheStore.Delete(strconv.Itoa(config.UserID))

	s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(config)

	return config
}

func (s *dbConfigStore) Get(userID int) *Config {
	if cached, found := s.cacheStore.Get(strconv.Itoa(userID)); found {
		return cached.(*Config)
	}

	config := &Config{}
	s.db.FirstOrInit(config, Config{UserID: userID})
	s.cacheStore.Set(strconv.Itoa(config.UserID), config, cache.DefaultExpiration)

	return config
}

func (s *dbConfigStore) GetLocale(userID int) string {
	config := s.Get(userID)
	return config.LangUI
}

func NewConfigStore(db *gorm.DB) ConfigStore {
	cacheStore := cache.New(CacheExpiration, CacheCleanup)
	return &dbConfigStore{db, cacheStore}
}
