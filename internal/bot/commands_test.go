package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBot_HandleCommands_Start(t *testing.T) {
	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "/start",
		Chat: &tgbotapi.Chat{ID: 123},
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len("/start"),
			},
		},
	}

	b.handleCommands(msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), message.ChatID)
	assert.Equal(t, "Выберите город:", message.Text)
	assert.NotNil(t, message.ReplyMarkup)
}

func TestBot_HandleCommands_Help(t *testing.T) {
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

	b.handleCommands(msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), message.ChatID)
	assert.Contains(t, message.Text, "Доступные команды")
	assert.Contains(t, message.Text, "/start")
	assert.Equal(t, "HTML", message.ParseMode)
}

func TestBot_HandleCommands_Unknown(t *testing.T) {
	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "/unknown",
		Chat: &tgbotapi.Chat{ID: 123},
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len("/unknown"),
			},
		},
	}

	b.handleCommands(msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), message.ChatID)
	assert.Contains(t, message.Text, "Неизвестная команда")
	assert.Equal(t, "HTML", message.ParseMode)
}
