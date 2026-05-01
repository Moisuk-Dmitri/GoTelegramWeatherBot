package weather

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validWeatherStampResponse() weatherStampResponse {
	return weatherStampResponse{
		CurrentUnits: struct {
			Time          string `json:"time"`
			Temperature   string `json:"temperature_2m"`
			Humidity      string `json:"relative_humidity_2m"`
			CloudCover    string `json:"cloud_cover"`
			Pressure      string `json:"pressure_msl"`
			WindSpeed     string `json:"wind_speed_10m"`
			WindDirection string `json:"wind_direction_10m"`
			Precipitation string `json:"precipitation"`
		}{
			Time:          "iso8601",
			Temperature:   "°C",
			Humidity:      "%",
			CloudCover:    "%",
			Pressure:      "hPa",
			WindSpeed:     "km/h",
			WindDirection: "°",
			Precipitation: "mm",
		},

		Current: struct {
			Time          string  `json:"time"`
			Temperature   float64 `json:"temperature_2m"`
			Humidity      int     `json:"relative_humidity_2m"`
			CloudCover    int     `json:"cloud_cover"`
			Pressure      float64 `json:"pressure_msl"`
			WindSpeed     float64 `json:"wind_speed_10m"`
			WindDirection int     `json:"wind_direction_10m"`
			Precipitation float64 `json:"precipitation"`
		}{
			Time:          "2026-04-30T12:00",
			Temperature:   14.3,
			Humidity:      65,
			CloudCover:    40,
			Pressure:      1013.2,
			WindSpeed:     12.5,
			WindDirection: 180,
			Precipitation: 0.0,
		},
	}
}

func validWeatherDailyResponse() weatherDailyResponse {
	return weatherDailyResponse{
		DailyUnits: struct {
			Date             string `json:"time"`
			TempMax          string `json:"temperature_2m_max"`
			TempMin          string `json:"temperature_2m_min"`
			SunriseTime      string `json:"sunrise"`
			SunsetTime       string `json:"sunset"`
			UVIndexMax       string `json:"uv_index_max"`
			WindSpeedMax     string `json:"wind_speed_10m_max"`
			PrecipitationSum string `json:"precipitation_sum"`
		}{
			Date:             "iso8601",
			TempMax:          "°C",
			TempMin:          "°C",
			SunriseTime:      "iso8601",
			SunsetTime:       "iso8601",
			UVIndexMax:       "index",
			WindSpeedMax:     "km/h",
			PrecipitationSum: "mm",
		},

		Daily: struct {
			Date             []string  `json:"time"`
			TempMin          []float64 `json:"temperature_2m_min"`
			TempMax          []float64 `json:"temperature_2m_max"`
			SunriseTime      []string  `json:"sunrise"`
			SunsetTime       []string  `json:"sunset"`
			UVIndexMax       []float64 `json:"uv_index_max"`
			WindSpeedMax     []float64 `json:"wind_speed_10m_max"`
			PrecipitationSum []float64 `json:"precipitation_sum"`
		}{
			Date: []string{
				"2026-04-30",
				"2026-05-01",
			},
			TempMin: []float64{
				8.2,
				9.1,
			},
			TempMax: []float64{
				15.6,
				16.3,
			},
			SunriseTime: []string{
				"2026-04-30T06:12",
				"2026-05-01T06:10",
			},
			SunsetTime: []string{
				"2026-04-30T20:45",
				"2026-05-01T20:47",
			},
			UVIndexMax: []float64{
				4.5,
				5.0,
			},
			WindSpeedMax: []float64{
				18.2,
				20.1,
			},
			PrecipitationSum: []float64{
				0.0,
				1.2,
			},
		},
	}
}

func TestMapper_mapWeatherStamp_Success(t *testing.T) {
	result, err := mapWeatherStamp(validWeatherStampResponse())

	require.NoError(t, err)

	assert.Equal(t, "iso8601", result.CurrentUnits.Time)
	assert.Equal(t, "°C", result.CurrentUnits.Temperature)
	assert.Equal(t, "%", result.CurrentUnits.Humidity)
	assert.Equal(t, "%", result.CurrentUnits.CloudCover)
	assert.Equal(t, "hPa", result.CurrentUnits.Pressure)
	assert.Equal(t, "km/h", result.CurrentUnits.WindSpeed)
	assert.Equal(t, "°", result.CurrentUnits.WindDirection)
	assert.Equal(t, "mm", result.CurrentUnits.Precipitation)

	assert.Equal(t, 14.3, result.Current.Temperature)
	assert.Equal(t, 65, result.Current.Humidity)
	assert.Equal(t, 40, result.Current.CloudCover)
	assert.Equal(t, 1013.2, result.Current.Pressure)
	assert.Equal(t, 12.5, result.Current.WindSpeed)
	assert.Equal(t, 180, result.Current.WindDirection)
	assert.Equal(t, 0.0, result.Current.Precipitation)

	assert.Equal(
		t,
		"2026-04-30T12:00",
		result.Current.Time.Format("2006-01-02T15:04"),
	)
}

func TestMapper_mapWeatherStamp_ParseTimeFail(t *testing.T) {
	wsrTest := validWeatherStampResponse()
	wsrTest.Current.Time = "bad-time"
	_, err := mapWeatherStamp(wsrTest)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse weather time")
}

func TestMapper_mapWeatherDaily_Success(t *testing.T) {
	result, err := mapWeatherDaily(validWeatherDailyResponse())
	require.NoError(t, err)

	assert.Equal(t, "iso8601", result.DailyUnits.Date)
	assert.Equal(t, "°C", result.DailyUnits.TempMax)
	assert.Equal(t, "°C", result.DailyUnits.TempMin)
	assert.Equal(t, "iso8601", result.DailyUnits.SunriseTime)
	assert.Equal(t, "iso8601", result.DailyUnits.SunsetTime)
	assert.Equal(t, "index", result.DailyUnits.UVIndexMax)
	assert.Equal(t, "km/h", result.DailyUnits.WindSpeedMax)
	assert.Equal(t, "mm", result.DailyUnits.PrecipitationSum)

	require.Len(t, result.Daily, 2)

	assert.Equal(t, "2026-04-30", result.Daily[0].Date.Format("2006-01-02"))
	assert.Equal(t, 8.2, result.Daily[0].TempMin)
	assert.Equal(t, 15.6, result.Daily[0].TempMax)
	assert.Equal(t, "2026-04-30T06:12", result.Daily[0].SunriseTime.Format("2006-01-02T15:04"))
	assert.Equal(t, "2026-04-30T20:45", result.Daily[0].SunsetTime.Format("2006-01-02T15:04"))
	assert.Equal(t, 4.5, result.Daily[0].UVIndexMax)
	assert.Equal(t, 18.2, result.Daily[0].WindSpeedMax)
	assert.Equal(t, 0.0, result.Daily[0].PrecipitationSum)
}

func TestMapper_mapWeatherDaily_ParseTimeFail(t *testing.T) {
	wdrTest1 := validWeatherDailyResponse()
	wdrTest1.Daily.Date[0] = "bad-date"

	_, err1 := mapWeatherDaily(wdrTest1)

	require.Error(t, err1)
	assert.Contains(t, err1.Error(), "failed to parse date")

	wdrTest2 := validWeatherDailyResponse()
	wdrTest2.Daily.SunriseTime[0] = "bad-time"

	_, err2 := mapWeatherDaily(wdrTest2)

	require.Error(t, err2)
	assert.Contains(t, err2.Error(), "failed to parse sunrise")

	wdrTest3 := validWeatherDailyResponse()
	wdrTest3.Daily.SunsetTime[0] = "bad-time"

	_, err3 := mapWeatherDaily(wdrTest3)

	require.Error(t, err3)
	assert.Contains(t, err3.Error(), "failed to parse sunset")
}
