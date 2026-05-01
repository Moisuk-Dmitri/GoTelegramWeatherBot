package bot

import (
	"context"
	"log"
	"main/internal/geo"
	"main/internal/weather"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState int

const (
	StateWaitingCity UserState = iota
	StateWaitingInterval
)

type TimeIntervalMode int

const (
	StampMode TimeIntervalMode = iota
	DailyMode
	WeeklyMode
)

type TimeIntervalName string

const (
	Stamp  TimeIntervalName = "Сейчас"
	Daily  TimeIntervalName = "День"
	Weekly TimeIntervalName = "Неделя"
)

type City string

const (
	Moscow       City = "Москва"
	StPetersburg City = "Санкт-Петербург"
)

var usersState = make(map[int64]UserState)

type UserInfo struct {
	City             geo.City
	TimeIntervalMode TimeIntervalMode
}

var usersInfo = make(map[int64]*UserInfo)

type WeatherResponse interface {
	Format() string
}

type handler func(weather.Service, context.Context, float64, float64) (WeatherResponse, error)

var handlers = map[TimeIntervalMode]handler{
	StampMode: func(s weather.Service, ctx context.Context, lat, lon float64) (WeatherResponse, error) {
		return s.GetWeatherStamp(ctx, lat, lon)
	},
	DailyMode: func(s weather.Service, ctx context.Context, lat, lon float64) (WeatherResponse, error) {
		return s.GetWeatherDaily(ctx, lat, lon)
	},
	WeeklyMode: func(s weather.Service, ctx context.Context, lat, lon float64) (WeatherResponse, error) {
		return s.GetWeatherWeekly(ctx, lat, lon)
	},
}

func (b *Bot) handleWeatherMessage(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	if usersInfo[chatID] == nil {
		usersInfo[chatID] = &UserInfo{}
	}

	state := usersState[chatID]
	text := strings.TrimSpace(msg.Text)

	switch state {
	case StateWaitingCity:
		switch text {
		case string(Moscow):
			usersInfo[chatID].City = geo.Moscow
		case string(StPetersburg):
			usersInfo[chatID].City = geo.StPetersburg
		default:
			b.replyWithKeyboard(chatID, "Выберите город:", cityKeyboard())
			return
		}

		b.replyWithKeyboard(chatID, "Выберите интервал:", intervalKeyboard())
		usersState[chatID] = StateWaitingInterval

	case StateWaitingInterval:
		switch text {
		case string(Stamp):
			usersInfo[chatID].TimeIntervalMode = StampMode
		case string(Daily):
			usersInfo[chatID].TimeIntervalMode = DailyMode
		case string(Weekly):
			usersInfo[chatID].TimeIntervalMode = WeeklyMode
		default:
			b.replyWithKeyboard(chatID, "Выберите интервал:", intervalKeyboard())
			return
		}

		coords, err := geo.CoordinatesByCity(usersInfo[chatID].City)
		if err != nil {
			log.Printf("failed to get coordinates: %v", err)
			b.reply(chatID, "Не удалось получить координаты города")
			return
		}

		resp, err := handlers[usersInfo[chatID].TimeIntervalMode](b.weather, ctx, coords.Latitude, coords.Longtitude)
		if err != nil {
			log.Printf("failed to get coordinates: %v", err)
			b.reply(chatID, "Не удалось получить прогноз погоды")
			return
		}

		b.replyRemoveKeyboard(chatID, resp.Format())
		usersState[chatID] = StateWaitingCity
		delete(usersInfo, chatID)

		b.reply(chatID, "Для продолжения работы введите /start")
	}
}
