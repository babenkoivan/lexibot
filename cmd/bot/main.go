package main

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/configs"
	"lexibot/internal/telegram"
	"lexibot/internal/translations"
	"time"
)

func main() {
	config, err := configs.LoadConfig(configs.DefaultConfigPath)

	if err != nil {
		panic(fmt.Errorf("cannot read from the config file: %w", err))
	}

	poller := &telebot.LongPoller{Timeout: 10 * time.Second}
	settings := telebot.Settings{Token: config.Telegram.Token, Poller: poller}
	telebot, err := telebot.NewBot(settings)

	if err != nil {
		panic(fmt.Errorf("cannot initiate telebot: %w", err))
	}

	ctx := context.Background()
	auth := option.WithAPIKey(config.Google.ApiKey)
	client, err := translate.NewClient(ctx, auth)

	if err != nil {
		panic(fmt.Errorf("cannot initiate google translate client: %w", err))
	}

	translator := translations.NewGoogleTranslator(client)

	bot := telegram.NewBot(telebot, translator)
	bot.Start()
}
