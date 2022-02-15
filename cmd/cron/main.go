package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"lexibot/internal/app"
	"lexibot/internal/database"
	"lexibot/internal/translation"
)

func main() {
	config, err := app.LoadConfig(app.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("cannot read from the app file: %w", err))
	}

	db, err := database.NewConnection(config.DB.DSN)
	if err != nil {
		panic(fmt.Errorf("cannot initiate database: %w", err))
	}

	scoreStore := translation.NewScoreStore(db)

	c := cron.New()
	c.AddJob("* * * * *", translation.NewAutoDecrementScoreJob(scoreStore))
	c.Run()
}
