package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lexibot/internal/app"
	"lexibot/internal/bot"
)

func main() {
	config, err := app.LoadConfig(app.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("cannot read from the app file: %w", err))
	}

	bundle, err := app.NewBundle(app.DefaultLocalePath)
	if err != nil {
		panic(fmt.Errorf("cannot create localization bundle: %w", err))
	}

	db, err := gorm.Open(mysql.Open(config.DB.DSN))
	if err != nil {
		panic(fmt.Errorf("cannot initiate database: %w", err))
	}

	historyStore := bot.NewHistoryStore(db)

	b, err := bot.NewBot(config.Bot.Token, config.Bot.Timeout, historyStore)
	if err != nil {
		panic(fmt.Errorf("cannot initiate telebot: %w", err))
	}

	//textSanitizers := map[language.Tag]translation.TextSanitizer{
	//	language.German: translation.NewGermanTextSanitizer(),
	//}

	//translator := translation.NewAzureTranslator(config.Translator.Endpoint, config.Translator.Key, config.Translator.Region, textSanitizers)
	//store := translation.NewDBStore(db)

	//b.OnText(translation.NewTranslateTextHandler(translator, store))
	//b.OnCallback(translation.OnCancelTranslation, translation.NewCancelTranslationHandler())
	//b.OnCallback(translation.OnSaveTranslation, translation.NewSaveTranslationHandler(store))
	//b.OnCallback(translation.OnDeleteTranslation, translation.NewDeleteTranslationHandler(store))
	b.OnCommand(bot.OnStart, bot.NewStartHandler(bundle))

	b.Start()
}
