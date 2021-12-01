package configs

import "github.com/spf13/viper"

const (
	DefaultConfigPath = "./configs"
)

type Telegram struct {
	Token string
}

type Google struct {
	ApiKey string
}

type Config struct {
	Telegram Telegram
	Google   Google
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
