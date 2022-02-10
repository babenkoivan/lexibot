package training

import (
	"gopkg.in/tucnak/telebot.v2"
	"lexibot/internal/bot"
)

const OnTraining = "/training"

type generateTaskHandler struct {
	taskGenerator *taskGenerator
}

func (h *generateTaskHandler) Handle(b bot.Bot, msg *telebot.Message) {
	task := h.taskGenerator.Next(msg.Sender.ID)
	b.Send(msg.Sender, &TaskMessage{task})
}

func NewGenerateTaskHandler(taskGenerator *taskGenerator) *generateTaskHandler {
	return &generateTaskHandler{taskGenerator}
}
