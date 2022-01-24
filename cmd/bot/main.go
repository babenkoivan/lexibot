package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lexibot/internal/app"
	"lexibot/internal/bot"
	"lexibot/internal/config"
	"lexibot/internal/locale"
)

func main() {
	conf, err := app.LoadConfig(app.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("cannot read from the app file: %w", err))
	}

	db, err := gorm.Open(mysql.Open(conf.DB.DSN))
	if err != nil {
		panic(fmt.Errorf("cannot initiate database: %w", err))
	}

	configStore := config.NewConfigStore(db)
	historyStore := bot.NewHistoryStore(db)

	locale, err := locale.NewLocale(locale.DefaultPath, configStore)
	if err != nil {
		panic(fmt.Errorf("cannot create localization bundle: %w", err))
	}

	b, err := bot.NewBot(conf.Bot.Token, conf.Bot.Timeout, locale, historyStore)
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
	b.OnCommand(app.OnStart, app.NewStartHandler())
	b.OnReply(&config.SelectLangUIMessage{}, config.NewSaveLangUIHandler(locale, configStore))
	b.OnReply(&config.SelectLangDictMessage{}, config.NewSaveLangDictHandler(locale, configStore))

	b.Start()
}
