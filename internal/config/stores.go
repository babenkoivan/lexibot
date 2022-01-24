package config

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lexibot/internal/locale"
)

type ConfigStore interface {
	locale.UserLocaleStore
	Save(config *Config) *Config
	Get(userID int) *Config
}

type dbConfigStore struct {
	db *gorm.DB
}

func (s *dbConfigStore) Save(config *Config) *Config {
	s.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(config)

	return config
}

func (s *dbConfigStore) Get(userID int) *Config {
	config := &Config{}
	s.db.FirstOrInit(config, Config{UserID: userID})
	return config
}

func (s *dbConfigStore) GetLocale(userID int) string {
	config := s.Get(userID)
	return config.LangUI
}

func NewConfigStore(db *gorm.DB) ConfigStore {
	return &dbConfigStore{db: db}
}
