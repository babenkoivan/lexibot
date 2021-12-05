package main

import (
	"fmt"
	"lexibot/internal/configs"
	"lexibot/internal/telegram"
	"lexibot/internal/translations"
)

func main() {
	config, err := configs.LoadConfig(configs.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("cannot read from the config file: %w", err))
	}

	bot, err := telegram.NewBot(config.Telegram)
	if err != nil {
		panic(fmt.Errorf("cannot initiate telebot: %w", err))
	}

	translator := translations.NewAzureTranslator(config.Azure)

	register := telegram.NewHandlerRegister(bot)
	register.Text(telegram.NewTranslateHandler(translator, bot))
	register.Callback(telegram.OnCancelTranslation, telegram.NewCancelTranslationHandler(bot))
	register.Callback(telegram.OnSaveTranslation, telegram.NewSaveTranslationHandler(bot))
	register.Callback(telegram.OnDeleteTranslation, telegram.NewDeleteTranslationHandler(bot))

	bot.Start()
}
