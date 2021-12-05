package configs

import (
	"github.com/spf13/viper"
	"time"
)

const (
	DefaultConfigPath = "./configs"
)

type Telegram struct {
	Token   string
	Timeout time.Duration
}

type Azure struct {
	Endpoint string
	Key      string
	Region   string
}

type Config struct {
	Telegram Telegram
	Azure    Azure
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
