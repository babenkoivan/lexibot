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

	b, err := bot.NewBot(config.Bot)
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

	b.OnText(translation.NewTranslateTextHandler(translator, store))
	b.OnCallback(translation.OnCancelTranslation, translation.NewCancelTranslationHandler())
	b.OnCallback(translation.OnSaveTranslation, translation.NewSaveTranslationHandler(store))
	b.OnCallback(translation.OnDeleteTranslation, translation.NewDeleteTranslationHandler(store))
	b.OnCommand(bot.OnStart, bot.NewStartHandler())
	b.OnCommand(bot.OnHelp, bot.NewHelpHandler())
	b.OnCommand(bot.OnSettings, bot.NewSettingsHandler())

	b.Start()
}
