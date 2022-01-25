package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lexibot/internal/app"
	"lexibot/internal/bot"
	"lexibot/internal/locale"
	"lexibot/internal/settings"
)

func main() {
	config, err := app.LoadConfig(app.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("cannot read from the app file: %w", err))
	}

	db, err := gorm.Open(mysql.Open(config.DB.DSN))
	if err != nil {
		panic(fmt.Errorf("cannot initiate database: %w", err))
	}

	settingsStore := settings.NewSettingsStore(db)
	historyStore := bot.NewHistoryStore(db)

	locale, err := locale.NewLocale(locale.DefaultPath, settingsStore)
	if err != nil {
		panic(fmt.Errorf("cannot create localization bundle: %w", err))
	}

	b, err := bot.NewBot(config.Bot.Token, config.Bot.Timeout, locale, historyStore)
	if err != nil {
		panic(fmt.Errorf("cannot initiate telebot: %w", err))
	}

	//textSanitizers := map[language.Tag]translation.TextSanitizer{
	//	language.German: translation.NewGermanTextSanitizer(),
	//}

	//translator := translation.NewAzureTranslator(settings.Translator.Endpoint, settings.Translator.Key, settings.Translator.Region, textSanitizers)
	//store := translation.NewDBStore(db)

	//b.OnText(translation.NewTranslateTextHandler(translator, store))
	//b.OnCallback(translation.OnCancelTranslation, translation.NewCancelTranslationHandler())
	//b.OnCallback(translation.OnSaveTranslation, translation.NewSaveTranslationHandler(store))
	//b.OnCallback(translation.OnDeleteTranslation, translation.NewDeleteTranslationHandler(store))
	b.OnCommand(app.OnStart, app.NewStartHandler())
	b.OnReply(&settings.SelectLangUIMessage{}, settings.NewSaveLangUIHandler(locale, settingsStore))
	b.OnReply(&settings.SelectLangDictMessage{}, settings.NewSaveLangDictHandler(locale, settingsStore))

	b.Start()
}
