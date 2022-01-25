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

	b.OnCommand(app.OnStart, app.NewStartHandler())
	b.OnCommand(app.OnHelp, app.NewHelpHandler())
	b.OnCommand(settings.OnSettings, settings.NewSettingsHandler())

	b.OnReply(&settings.SelectLangUIMessage{}, settings.NewSaveLangUIHandler(locale, settingsStore))
	b.OnReply(&settings.SelectLangDictMessage{}, settings.NewSaveLangDictHandler(locale, settingsStore))
	b.OnReply(&settings.EnableAutoTranslateMessage{}, settings.NewSaveAutoTranslateHandler(locale, settingsStore))
	b.OnReply(&settings.EnterWordsPerTrainingMessage{}, settings.NewSaveWordsPerTrainingHandler(settingsStore))

	b.Start()
}
