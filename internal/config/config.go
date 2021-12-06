package config

import (
	"github.com/spf13/viper"
	"time"
)

const (
	DefaultConfigPath = "./configs"
)

type Bot struct {
	Token   string
	Timeout time.Duration
}

type Translator struct {
	Endpoint string
	Key      string
	Region   string
}

type Config struct {
	Bot        Bot
	Translator Translator
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
