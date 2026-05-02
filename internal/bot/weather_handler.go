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
	Stamp TimeIntervalMode = iota
	Daily
	Weekly
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
	Stamp: func(s weather.Service, ctx context.Context, lat, lon float64) (WeatherResponse, error) {
		return s.GetWeatherStamp(ctx, lat, lon)
	},
	Daily: func(s weather.Service, ctx context.Context, lat, lon float64) (WeatherResponse, error) {
		return s.GetWeatherDaily(ctx, lat, lon)
	},
	Weekly: func(s weather.Service, ctx context.Context, lat, lon float64) (WeatherResponse, error) {
		return s.GetWeatherWeekly(ctx, lat, lon)
	},
}

const (
	CityMoscowLabel       = "Москва"
	CityStPetersburgLabel = "Санкт-Петербург"

	IntervalNowLabel  = "Сейчас"
	IntervalDayLabel  = "День"
	IntervalWeekLabel = "Неделя"
)

var cityByLabel = map[string]geo.City{
	CityMoscowLabel:       geo.Moscow,
	CityStPetersburgLabel: geo.StPetersburg,
}

var intervalByLabel = map[string]TimeIntervalMode{
	IntervalNowLabel:  Stamp,
	IntervalDayLabel:  Daily,
	IntervalWeekLabel: Weekly,
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
		city, ok := cityByLabel[text]
		if !ok {
			b.replyWithKeyboard(chatID, "Выберите город:", cityKeyboard())
			return
		}

		usersInfo[chatID].City = city

		b.replyWithKeyboard(chatID, "Выберите интервал:", intervalKeyboard())
		usersState[chatID] = StateWaitingInterval

	case StateWaitingInterval:
		interval, ok := intervalByLabel[text]
		if !ok {
			b.replyWithKeyboard(chatID, "Выберите интервал:", intervalKeyboard())
			return
		}
		usersInfo[chatID].TimeIntervalMode = interval

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
