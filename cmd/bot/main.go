package main

import (
	"fmt"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	db, err := gorm.Open(mysql.Open(config.DB.DSN))
	if err != nil {
		panic(fmt.Errorf("cannot initiate database: %w", err))
	}

	textSanitizers := map[language.Tag]translation.TextSanitizer{
		language.German: translation.NewGermanTextSanitizer(),
	}

	translator := translation.NewAzureTranslator(config.Translator, textSanitizers)
	store := translation.NewDBStore(db)

	bot.OnText(translation.NewSuggestTranslationHandler(translator, store))
	bot.OnCallback(translation.OnCancelTranslation, translation.NewCancelTranslationHandler())
	bot.OnCallback(translation.OnSaveTranslation, translation.NewSaveTranslationHandler(store))
	bot.OnCallback(translation.OnDeleteTranslation, translation.NewDeleteTranslationHandler(store))

	bot.Start()
}
