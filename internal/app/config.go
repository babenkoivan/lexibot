package app

import (
	"github.com/spf13/viper"
	"time"
)

const DefaultConfigPath = "./configs/.env"

type Config struct {
	Bot struct {
		Token   string
		Timeout time.Duration
	}
	Translator struct {
		Endpoint string
		Key      string
	}
	DB struct {
		DSN string
	}
}

func LoadConfig(path string) (config Config, err error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("_"))
	v.AutomaticEnv()
	v.SetConfigFile(path)
	_ = v.ReadInConfig()
	err = v.Unmarshal(&config)
	return
}
