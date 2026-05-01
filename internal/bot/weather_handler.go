package bot

import (
	"context"
	"log"
	"main/internal/geo"
	"main/internal/weather"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TimeIntervalMode int

const (
	Stamp TimeIntervalMode = iota
	Daily
	Weekly
)

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

type UserState int

const (
	StateWaitingCity UserState = iota
	StateWaitingInterval
)

var usersState = make(map[int64]UserState)

type UserInfo struct {
	City             geo.City
	TimeIntervalMode TimeIntervalMode
}

var usersInfo = make(map[int64]*UserInfo)

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
		case "Москва":
			usersInfo[chatID].City = geo.Moscow
		case "Санкт-Петербург":
			usersInfo[chatID].City = geo.StPetersburg
		default:
			b.replyWithKeyboard(chatID, "Выберите город:", cityKeyboard())
			return
		}

		b.replyWithKeyboard(chatID, "Выберите интервал:", intervalKeyboard())
		usersState[chatID] = StateWaitingInterval

	case StateWaitingInterval:
		switch text {
		case "Сейчас":
			usersInfo[chatID].TimeIntervalMode = Stamp
		case "День":
			usersInfo[chatID].TimeIntervalMode = Daily
		case "Неделя":
			usersInfo[chatID].TimeIntervalMode = Weekly
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
