package main

import (
	"fmt"
	"lexibot/internal/bot"
	"lexibot/internal/config"
	"lexibot/internal/translation"
)

func main() {
	config, err := config.LoadConfig(config.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("cannot read from the config file: %w", err))
	}

	bot, err := bot.NewBot(config.Bot)
	if err != nil {
		panic(fmt.Errorf("cannot initiate telebot: %w", err))
	}

	translator := translation.NewAzureTranslator(config.Translator)

	bot.OnText(translation.NewSuggestTranslationHandler(translator))
	bot.OnCallback(translation.OnCancelTranslation, translation.NewCancelTranslationHandler())
	bot.OnCallback(translation.OnSaveTranslation, translation.NewSaveTranslationHandler())
	bot.OnCallback(translation.OnDeleteTranslation, translation.NewDeleteTranslationHandler())

	bot.Start()
}
