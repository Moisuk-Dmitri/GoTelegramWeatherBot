package bot

import (
	"context"
	"errors"
	"main/internal/weather"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTelegramAPI struct {
	sentMessage  tgbotapi.Chattable
	sentMessages []tgbotapi.Chattable
	sendErr      error

	stopCalled  bool
	updatesChan chan tgbotapi.Update
}

func (m *mockTelegramAPI) GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return m.updatesChan
}

func (m *mockTelegramAPI) StopReceivingUpdates() {
	m.stopCalled = true
}

func (m *mockTelegramAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.sentMessage = c
	m.sentMessages = append(m.sentMessages, c)
	return tgbotapi.Message{}, m.sendErr
}

type mockWeatherService struct {
	stampCalled  bool
	dailyCalled  bool
	weeklyCalled bool

	err error
}

func (m *mockWeatherService) GetWeatherStamp(ctx context.Context, lat, lon float64) (weather.WeatherStamp, error) {
	m.stampCalled = true
	return weather.WeatherStamp{}, m.err
}

func (m *mockWeatherService) GetWeatherDaily(ctx context.Context, lat, lon float64) (weather.WeatherDaily, error) {
	m.dailyCalled = true
	return weather.WeatherDaily{}, m.err
}

func (m *mockWeatherService) GetWeatherWeekly(ctx context.Context, lat, lon float64) (weather.WeatherDaily, error) {
	m.weeklyCalled = true
	return weather.WeatherDaily{}, m.err
}

func TestBot_NewBot_Success(t *testing.T) {
	api := &mockTelegramAPI{}
	service := &mockWeatherService{}

	b := NewBot(api, service)

	require.NotNil(t, b)
	assert.Equal(t, api, b.api)
	assert.Equal(t, service, b.weather)
}

func TestBot_Run_ContextCanceled(t *testing.T) {
	api := &mockTelegramAPI{
		updatesChan: make(chan tgbotapi.Update),
	}

	b := NewBot(api, &mockWeatherService{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := b.Run(ctx)

	require.ErrorIs(t, err, context.Canceled)
	assert.True(t, api.stopCalled)
}

func TestBot_Run_MessageUpdateSendsReply(t *testing.T) {
	api := &mockTelegramAPI{
		updatesChan: make(chan tgbotapi.Update, 1),
	}

	b := NewBot(api, &mockWeatherService{})

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- b.Run(ctx)
	}()

	api.updatesChan <- tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: 1,
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			Text: "/start",
		},
	}

	require.Eventually(t, func() bool {
		return api.sentMessage != nil
	}, time.Second, 10*time.Millisecond)

	cancel()

	err := <-done
	require.ErrorIs(t, err, context.Canceled)
	assert.True(t, api.stopCalled)

	msg, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), msg.ChatID)
	assert.NotEmpty(t, msg.Text)
}

func TestBot_Reply_Success(t *testing.T) {
	api := &mockTelegramAPI{}

	b := NewBot(api, &mockWeatherService{})

	b.reply(123, "hello")

	require.NotNil(t, api.sentMessage)

	msg, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), msg.ChatID)
	assert.Equal(t, "hello", msg.Text)
	assert.Equal(t, "HTML", msg.ParseMode)
}

func TestBot_Reply_SendError(t *testing.T) {
	api := &mockTelegramAPI{
		sendErr: errors.New("send failed"),
	}

	b := NewBot(api, &mockWeatherService{})

	require.NotPanics(t, func() {
		b.reply(123, "hello")
	})

	require.NotNil(t, api.sentMessage)

	msg, ok := api.sentMessage.(tgbotapi.MessageConfig)
	require.True(t, ok)

	assert.Equal(t, int64(123), msg.ChatID)
	assert.Equal(t, "hello", msg.Text)
	assert.Equal(t, "HTML", msg.ParseMode)
}
