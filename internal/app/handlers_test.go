package app_test

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/app"
	"lexibot/internal/settings"
	"lexibot/internal/testkit"
	"testing"
)

func TestStartHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	botSpy := testkit.NewBotSpy(t)

	handler := app.NewStartHandler()
	handler.Handle(botSpy, &telebot.Message{Sender: user})

	botSpy.AssertSent(user, &app.StartMessage{})
	botSpy.AssertSent(user, &settings.SelectLangUIMessage{})
}

func TestHelpHandler_Handle(t *testing.T) {
	user := &telebot.User{ID: 1}
	botSpy := testkit.NewBotSpy(t)

	handler := app.NewHelpHandler()
	handler.Handle(botSpy, &telebot.Message{Sender: user})

	botSpy.AssertSent(user, &app.HelpMessage{})
}
