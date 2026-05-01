package bot

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBot_HandleMessages_Command(t *testing.T) {
	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "/help",
		Chat: &tgbotapi.Chat{ID: 123},
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len("/help"),
			},
		},
	}

	b.handleMessages(context.Background(), msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), message.ChatID)
	assert.Contains(t, message.Text, "Доступные команды")
}

func TestBot_HandleMessages_NotCommand(t *testing.T) {
	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "Berlin",
		Chat: &tgbotapi.Chat{ID: 123},
	}

	b.handleMessages(context.Background(), msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), message.ChatID)
	assert.NotEmpty(t, message.Text)
}
