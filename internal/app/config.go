package app

import (
	"github.com/spf13/viper"
	"time"
)

const DefaultConfigPath = "./configs"

type Config struct {
	Bot struct {
		Token   string
		Timeout time.Duration
	}
	Translator struct {
		Endpoint string
		Key      string
		Region   string
	}
	DB struct {
		DSN string
	}
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
