package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lexibot/internal/app"
	"lexibot/internal/bot"
	"lexibot/internal/locale"
	"lexibot/internal/settings"
	"lexibot/internal/training"
	"lexibot/internal/translation"
)

func main() {
	config, err := app.LoadConfig(app.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("cannot read from the app file: %w", err))
	}

	db, err := gorm.Open(mysql.Open(config.DB.DSN), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic(fmt.Errorf("cannot initiate database: %w", err))
	}

	settingsStore := settings.NewSettingsStore(db)
	historyStore := bot.NewHistoryStore(db)
	translationStore := translation.NewTranslationStore(db)
	scoreStore := translation.NewScoreStore(db)

	translator := translation.NewTranslator(config.Translator.Endpoint, config.Translator.Key, translationStore)
	taskGenerator := training.NewTaskGenerator(settingsStore, translationStore, scoreStore)

	loc, err := locale.NewLocale(locale.DefaultPath, settingsStore)
	if err != nil {
		panic(fmt.Errorf("cannot create locale: %w", err))
	}

	b, err := bot.NewBot(config.Bot.Token, config.Bot.Timeout, loc, historyStore)
	if err != nil {
		panic(fmt.Errorf("cannot initiate telebot: %w", err))
	}

	b.OnMessage(translation.NewTranslateHandler(settingsStore, translationStore, scoreStore, translator))

	b.OnCommand(translation.OnDelete, translation.NewClarifyWhatToDeleteHandler())
	b.OnCommand(settings.OnSettings, settings.NewSettingsHandler())
	b.OnCommand(app.OnHelp, app.NewHelpHandler())
	b.OnCommand(app.OnStart, app.NewStartHandler())
	b.OnCommand(training.OnTraining, training.NewGenerateTaskHandler(taskGenerator))

	b.OnReply(&settings.SelectLangUIMessage{}, settings.NewSaveLangUIHandler(loc, settingsStore))
	b.OnReply(&settings.SelectLangDictMessage{}, settings.NewSaveLangDictHandler(loc, settingsStore))
	b.OnReply(&settings.EnableAutoTranslateMessage{}, settings.NewSaveAutoTranslateHandler(loc, settingsStore))
	b.OnReply(&settings.EnterWordsPerTrainingMessage{}, settings.NewSaveWordsPerTrainingHandler(settingsStore))
	b.OnReply(&translation.EnterTranslationMessage{}, translation.NewAddToDictionaryHandler(settingsStore, translationStore, scoreStore))
	b.OnReply(&translation.WhatToDeleteMessage{}, translation.NewDeleteFromDictionaryHandler(settingsStore, translationStore, scoreStore))

	b.Start()
}
