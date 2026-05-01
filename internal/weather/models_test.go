package weather

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeatherStamp_Format_Success(t *testing.T) {
	ws := WeatherStamp{
		CurrentUnits: currentUnits{
			Time:          "iso8601",
			Temperature:   "°C",
			Humidity:      "%",
			CloudCover:    "%",
			Pressure:      "hPa",
			WindSpeed:     "km/h",
			WindDirection: "°",
			Precipitation: "mm",
		},
		Current: current{
			Time:          time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC),
			Temperature:   14.3,
			Humidity:      65,
			CloudCover:    40,
			Pressure:      1013.2,
			WindSpeed:     12.5,
			WindDirection: 180,
			Precipitation: 0.0,
		},
	}

	result := ws.Format()

	assert.Contains(t, result, "🕖 Время - 2026-04-30 12:00:00 +0000 UTC")
	assert.Contains(t, result, "🌡️ Температура - 14.3°C")
	assert.Contains(t, result, "💧 Влажность - 65%")
	assert.Contains(t, result, "☁️ Пасмурность - 40%")
	assert.Contains(t, result, "😨 Давление - 1013.2hPa")
	assert.Contains(t, result, "💨 Скорость ветра - 12.5km/h")
	assert.Contains(t, result, "🧭 Направление ветра - 180°")
	assert.Contains(t, result, "🌧 Осадки - 0mm")
}

func TestWeatherDaily_Format_Success(t *testing.T) {
	wd := WeatherDaily{
		DailyUnits: dailyUnits{
			Date:             "iso8601",
			TempMin:          "°C",
			TempMax:          "°C",
			SunriseTime:      "iso8601",
			SunsetTime:       "iso8601",
			UVIndexMax:       "",
			WindSpeedMax:     "km/h",
			PrecipitationSum: "mm",
		},
		Daily: []daily{
			{
				Date:             time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
				TempMin:          8.2,
				TempMax:          15.6,
				SunriseTime:      time.Date(2026, 4, 30, 6, 12, 0, 0, time.UTC),
				SunsetTime:       time.Date(2026, 4, 30, 20, 45, 0, 0, time.UTC),
				UVIndexMax:       4.5,
				WindSpeedMax:     18.2,
				PrecipitationSum: 0.0,
			},
		},
	}

	result := wd.Format()

	assert.Contains(t, result, "📅 Дата - 2026-04-30 00:00:00 +0000 UTC")
	assert.Contains(t, result, "🥶 Минимальная температура - 8.2°C")
	assert.Contains(t, result, "🔥 Максимальная температура - 15.6°C")
	assert.Contains(t, result, "🌅 Время рассвета 2026-04-30 06:12:00 +0000 UTC")
	assert.Contains(t, result, "🌇 Время заката - 2026-04-30 20:45:00 +0000 UTC")
	assert.Contains(t, result, "☀️ Максимальный UV индекс - 4.5")
	assert.Contains(t, result, "💨 Максимальная скорость ветра - 18.2km/h")
	assert.Contains(t, result, "🌧 Сумма осадков - 0mm")
}
