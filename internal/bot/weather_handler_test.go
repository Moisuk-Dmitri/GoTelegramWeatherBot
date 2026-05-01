package bot

import (
	"context"
	"errors"
	"main/internal/geo"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetUserState() {
	usersState = make(map[int64]UserState)
	usersInfo = make(map[int64]*UserInfo)
}

func TestBot_HandleWeatherMessage_SuccessMoscow(t *testing.T) {
	resetUserState()

	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "Москва",
		Chat: &tgbotapi.Chat{ID: 123},
	}

	b.handleWeatherMessage(context.Background(), msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), message.ChatID)
	assert.Equal(t, "Выберите интервал:", message.Text)
	assert.NotNil(t, message.ReplyMarkup)

	assert.Equal(t, StateWaitingInterval, usersState[123])
	assert.Equal(t, geo.Moscow, usersInfo[123].City)
}

func TestBot_HandleWeatherMessage_InvalidCity(t *testing.T) {
	resetUserState()

	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "Лондон",
		Chat: &tgbotapi.Chat{ID: 123},
	}

	b.handleWeatherMessage(context.Background(), msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), message.ChatID)
	assert.Equal(t, "Выберите город:", message.Text)
	assert.NotNil(t, message.ReplyMarkup)

	assert.Equal(t, StateWaitingCity, usersState[123])
}

func TestBot_HandleWeatherMessage_SuccessStamp(t *testing.T) {
	resetUserState()

	chatID := int64(123)

	usersState[chatID] = StateWaitingInterval
	usersInfo[chatID] = &UserInfo{
		City: geo.Moscow,
	}

	api := &mockTelegramAPI{}
	service := &mockWeatherService{}

	b := NewBot(api, service)

	msg := &tgbotapi.Message{
		Text: "Сейчас",
		Chat: &tgbotapi.Chat{ID: chatID},
	}

	b.handleWeatherMessage(context.Background(), msg)

	assert.True(t, service.stampCalled)
	assert.False(t, service.dailyCalled)
	assert.False(t, service.weeklyCalled)

	require.Len(t, api.sentMessages, 2)

	firstMsg, ok := api.sentMessages[0].(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, chatID, firstMsg.ChatID)
	assert.NotEmpty(t, firstMsg.Text)

	_, ok = firstMsg.ReplyMarkup.(tgbotapi.ReplyKeyboardRemove)
	assert.True(t, ok)

	secondMsg, ok := api.sentMessages[1].(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, chatID, secondMsg.ChatID)
	assert.Equal(t, "Для продолжения работы введите /start", secondMsg.Text)

	assert.Equal(t, StateWaitingCity, usersState[chatID])
	assert.Nil(t, usersInfo[chatID])
}

func TestBot_HandleWeatherMessage_InvalidInterval(t *testing.T) {
	resetUserState()

	chatID := int64(123)

	usersState[chatID] = StateWaitingInterval
	usersInfo[chatID] = &UserInfo{
		City: geo.Moscow,
	}

	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "Месяц",
		Chat: &tgbotapi.Chat{ID: chatID},
	}

	b.handleWeatherMessage(context.Background(), msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, chatID, message.ChatID)
	assert.Equal(t, "Выберите интервал:", message.Text)
	assert.NotNil(t, message.ReplyMarkup)

	assert.Equal(t, StateWaitingInterval, usersState[chatID])
}

func TestBot_HandleWeatherMessage_SuccessDaily(t *testing.T) {
	resetUserState()

	chatID := int64(123)

	usersState[chatID] = StateWaitingInterval
	usersInfo[chatID] = &UserInfo{
		City: geo.Moscow,
	}

	api := &mockTelegramAPI{}
	service := &mockWeatherService{}

	b := NewBot(api, service)

	msg := &tgbotapi.Message{
		Text: "День",
		Chat: &tgbotapi.Chat{ID: chatID},
	}

	b.handleWeatherMessage(context.Background(), msg)

	assert.False(t, service.stampCalled)
	assert.True(t, service.dailyCalled)
	assert.False(t, service.weeklyCalled)

	require.Len(t, api.sentMessages, 2)

	assert.Equal(t, StateWaitingCity, usersState[chatID])
	assert.Nil(t, usersInfo[chatID])
}

func TestBot_HandleWeatherMessage_SuccessWeekly(t *testing.T) {
	resetUserState()

	chatID := int64(123)

	usersState[chatID] = StateWaitingInterval
	usersInfo[chatID] = &UserInfo{
		City: geo.Moscow,
	}

	api := &mockTelegramAPI{}
	service := &mockWeatherService{}

	b := NewBot(api, service)

	msg := &tgbotapi.Message{
		Text: "Неделя",
		Chat: &tgbotapi.Chat{ID: chatID},
	}

	b.handleWeatherMessage(context.Background(), msg)

	assert.False(t, service.stampCalled)
	assert.False(t, service.dailyCalled)
	assert.True(t, service.weeklyCalled)

	require.Len(t, api.sentMessages, 2)

	assert.Equal(t, StateWaitingCity, usersState[chatID])
	assert.Nil(t, usersInfo[chatID])
}

func TestBot_HandleWeatherMessage_WeatherServiceError(t *testing.T) {
	resetUserState()

	chatID := int64(123)

	usersState[chatID] = StateWaitingInterval
	usersInfo[chatID] = &UserInfo{
		City: geo.Moscow,
	}

	api := &mockTelegramAPI{}
	service := &mockWeatherService{
		err: errors.New("weather error"),
	}

	b := NewBot(api, service)

	msg := &tgbotapi.Message{
		Text: "Сейчас",
		Chat: &tgbotapi.Chat{ID: chatID},
	}

	b.handleWeatherMessage(context.Background(), msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, chatID, message.ChatID)
	assert.Equal(t, "Не удалось получить прогноз погоды", message.Text)
}

func TestBot_HandleWeatherMessage_SuccessStPetersburg(t *testing.T) {
	resetUserState()

	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "Санкт-Петербург",
		Chat: &tgbotapi.Chat{ID: 123},
	}

	b.handleWeatherMessage(context.Background(), msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), message.ChatID)
	assert.Equal(t, "Выберите интервал:", message.Text)
	assert.NotNil(t, message.ReplyMarkup)

	assert.Equal(t, StateWaitingInterval, usersState[123])
	assert.Equal(t, geo.StPetersburg, usersInfo[123].City)
}

func TestBot_HandleWeatherMessage_CoordinatesError(t *testing.T) {
	resetUserState()

	chatID := int64(123)

	usersState[chatID] = StateWaitingInterval
	usersInfo[chatID] = &UserInfo{
		City: geo.City("invalid-city"),
	}

	api := &mockTelegramAPI{}
	b := NewBot(api, &mockWeatherService{})

	msg := &tgbotapi.Message{
		Text: "Сейчас",
		Chat: &tgbotapi.Chat{ID: chatID},
	}

	b.handleWeatherMessage(context.Background(), msg)

	require.NotNil(t, api.sentMessage)

	message, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, chatID, message.ChatID)
	assert.Equal(t, "Не удалось получить координаты города", message.Text)
}
