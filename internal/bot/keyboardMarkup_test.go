package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBot_IntervalKeyboard_Succes(t *testing.T) {
	kb := intervalKeyboard()

	require.Len(t, kb.Keyboard, 1)
	require.Len(t, kb.Keyboard[0], 3)

	assert.Equal(t, "Сейчас", kb.Keyboard[0][0].Text)
	assert.Equal(t, "День", kb.Keyboard[0][1].Text)
	assert.Equal(t, "Неделя", kb.Keyboard[0][2].Text)
}

func TestBot_ReplyRemoveKeyboard_Succes(t *testing.T) {
	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	b.replyRemoveKeyboard(123, "test text")

	require.NotNil(t, api.sentMessage)

	msg, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), msg.ChatID)
	assert.Equal(t, "test text", msg.Text)

	_, ok = msg.ReplyMarkup.(tgbotapi.ReplyKeyboardRemove)
	require.True(t, ok)
}
