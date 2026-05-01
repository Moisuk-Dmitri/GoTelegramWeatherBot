package weather

import (
	"fmt"
	"strings"
	"time"
)

// Weather stamp

type currentUnits struct {
	Time          string
	Temperature   string
	Humidity      string
	CloudCover    string
	Pressure      string
	WindSpeed     string
	WindDirection string
	Precipitation string
}

type current struct {
	Time          time.Time
	Temperature   float64
	Humidity      int
	CloudCover    int
	Pressure      float64
	WindSpeed     float64
	WindDirection int
	Precipitation float64
}

type WeatherStamp struct {
	CurrentUnits currentUnits
	Current      current
}

func (w WeatherStamp) Format() string {
	return fmt.Sprintf(`🕖 Время - %v
	🌡️ Температура - %v%s
	💧 Влажность - %v%s
	☁️ Пасмурность - %v%s
	😨 Давление - %v%s
	💨 Скорость ветра - %v%s
	🧭 Направление ветра - %v%s
	🌧 Осадки - %v%s`,
		w.Current.Time,
		w.Current.Temperature, w.CurrentUnits.Temperature,
		w.Current.Humidity, w.CurrentUnits.Humidity,
		w.Current.CloudCover, w.CurrentUnits.CloudCover,
		w.Current.Pressure, w.CurrentUnits.Pressure,
		w.Current.WindSpeed, w.CurrentUnits.WindSpeed,
		w.Current.WindDirection, w.CurrentUnits.WindDirection,
		w.Current.Precipitation, w.CurrentUnits.Precipitation)
}

// Weather daily

type dailyUnits struct {
	Date             string
	TempMax          string
	TempMin          string
	SunriseTime      string
	SunsetTime       string
	UVIndexMax       string
	WindSpeedMax     string
	PrecipitationSum string
}

type daily struct {
	Date             time.Time
	TempMin          float64
	TempMax          float64
	SunriseTime      time.Time
	SunsetTime       time.Time
	UVIndexMax       float64
	WindSpeedMax     float64
	PrecipitationSum float64
}

type WeatherDaily struct {
	DailyUnits dailyUnits
	Daily      []daily
}

func (w WeatherDaily) Format() string {
	var s strings.Builder
	for _, wd := range w.Daily {
		s.WriteString(fmt.Sprintf(`📅 Дата - %v
		🥶 Минимальная температура - %v%s
		🔥 Максимальная температура - %v%s
		🌅 Время рассвета %v
		🌇 Время заката - %v
		☀️ Максимальный UV индекс - %v%s
		💨 Максимальная скорость ветра - %v%s
		🌧 Сумма осадков - %v%s`,
			wd.Date,
			wd.TempMin, w.DailyUnits.TempMin,
			wd.TempMax, w.DailyUnits.TempMax,
			wd.SunriseTime,
			wd.SunsetTime,
			wd.UVIndexMax, w.DailyUnits.UVIndexMax,
			wd.WindSpeedMax, w.DailyUnits.WindSpeedMax,
			wd.PrecipitationSum, w.DailyUnits.PrecipitationSum,
		))
		s.WriteString("\n\n")
	}

	return s.String()
}
