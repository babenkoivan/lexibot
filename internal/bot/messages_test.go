package bot_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"lexibot/internal/bot"
	"testing"
)

func TestWithReplyKeyboard(t *testing.T) {
	replyMarkup := bot.WithReplyKeyboard([]string{"btn1", "btn2"})

	require.Len(t, replyMarkup.ReplyKeyboard, 2)
	assert.Equal(t, "btn1", replyMarkup.ReplyKeyboard[0][0].Text)
	assert.Equal(t, "btn2", replyMarkup.ReplyKeyboard[1][0].Text)
}

func TestWithoutReplyKeyboard(t *testing.T) {
	replyMarkup := bot.WithoutReplyKeyboard()

	assert.Len(t, replyMarkup.ReplyKeyboard, 0)
	assert.True(t, replyMarkup.ReplyKeyboardRemove)
}
