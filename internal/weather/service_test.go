package weather

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_NewService_NewServiceSuccess(t *testing.T) {
	apiURL := "some-api"
	svc := NewService(apiURL)

	assert.NotNil(t, svc)

	s, ok := svc.(*service)
	require.True(t, ok)

	assert.Equal(t, apiURL, s.apiURL)
	require.NotNil(t, s.client)
}

type mockClient struct {
	resp []byte
	err  error
}

func (m *mockClient) get(ctx context.Context, baseURL string, params ...string) ([]byte, error) {
	return m.resp, m.err
}

func getJsonWeatherStampResponse() []byte {
	return []byte(`{
  "latitude": 52.52,
  "longitude": 13.419998,
  "generationtime_ms": 0.10073184967041,
  "utc_offset_seconds": 0,
  "timezone": "GMT",
  "timezone_abbreviation": "GMT",
  "elevation": 38,
  "current_units": {
    "time": "iso8601",
    "interval": "seconds",
    "temperature_2m": "°C",
    "relative_humidity_2m": "%",
    "cloud_cover": "%",
    "wind_speed_10m": "km/h",
    "wind_direction_10m": "°",
    "precipitation": "mm",
    "pressure_msl": "hPa"
  },
  "current": {
    "time": "2026-04-30T12:00",
    "interval": 900,
    "temperature_2m": 14.3,
    "relative_humidity_2m": 65,
    "cloud_cover": 40,
    "wind_speed_10m": 5.3,
    "wind_direction_10m": 332,
    "precipitation": 0,
    "pressure_msl": 1027.8
  }
}`)
}

func TestService_GetWeatherStamp_Success(t *testing.T) {
	mock := &mockClient{
		resp: getJsonWeatherStampResponse(),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	result, err := s.GetWeatherStamp(t.Context(), 0, 0)

	require.NoError(t, err)

	assert.Equal(t, 14.3, result.Current.Temperature)
	assert.Equal(t, 65, result.Current.Humidity)
	assert.Equal(t, 40, result.Current.CloudCover)
	assert.Equal(t, 1027.8, result.Current.Pressure)
	assert.Equal(t, 5.3, result.Current.WindSpeed)
	assert.Equal(t, 332, result.Current.WindDirection)
	assert.Equal(t, 0.0, result.Current.Precipitation)
	assert.Equal(
		t,
		"2026-04-30T12:00",
		result.Current.Time.Format("2006-01-02T15:04"),
	)

	assert.Equal(t, "°C", result.CurrentUnits.Temperature)
	assert.Equal(t, "%", result.CurrentUnits.Humidity)
	assert.Equal(t, "iso8601", result.CurrentUnits.Time)
	assert.Equal(t, "%", result.CurrentUnits.CloudCover)
	assert.Equal(t, "hPa", result.CurrentUnits.Pressure)
	assert.Equal(t, "km/h", result.CurrentUnits.WindSpeed)
	assert.Equal(t, "°", result.CurrentUnits.WindDirection)
	assert.Equal(t, "mm", result.CurrentUnits.Precipitation)
}

func TestService_GetWeatherStamp_ResponseFail(t *testing.T) {
	mock := &mockClient{
		err: errors.New("some-error"),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	_, err := s.GetWeatherStamp(t.Context(), 0, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "some-error")
}

func TestService_GetWeatherStamp_JsonUnmarshalFail(t *testing.T) {
	mock := &mockClient{
		resp: []byte(`invalid-json`),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	_, err := s.GetWeatherStamp(t.Context(), 0, 0)

	require.Error(t, err)
}

func getJsonWeatherDailyResponse() []byte {
	return []byte(`{
  "latitude": 52.52,
  "longitude": 13.419998,
  "generationtime_ms": 0.0922679901123047,
  "utc_offset_seconds": 0,
  "timezone": "GMT",
  "timezone_abbreviation": "GMT",
  "elevation": 38,
  "daily_units": {
    "time": "iso8601",
    "temperature_2m_max": "°C",
    "temperature_2m_min": "°C",
    "sunrise": "iso8601",
    "sunset": "iso8601",
    "uv_index_max": "",
    "wind_speed_10m_max": "km/h",
    "precipitation_sum": "mm"
  },
  "daily": {
    "time": [
      "2026-05-01"
    ],
    "temperature_2m_max": [23.3],
    "temperature_2m_min": [9.2],
    "sunrise": [
      "2026-05-01T03:34"
    ],
    "sunset": [
      "2026-05-01T18:32"
    ],
    "uv_index_max": [5.85],
    "wind_speed_10m_max": [8.7],
    "precipitation_sum": [0]
  }
}`)
}

func TestService_GetWeatherDaily_Success(t *testing.T) {
	mock := &mockClient{
		resp: getJsonWeatherDailyResponse(),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	result, err := s.GetWeatherDaily(t.Context(), 0, 0)

	require.NoError(t, err)

	assert.Equal(t, "iso8601", result.DailyUnits.Date)
	assert.Equal(t, "°C", result.DailyUnits.TempMax)
	assert.Equal(t, "°C", result.DailyUnits.TempMin)
	assert.Equal(t, "iso8601", result.DailyUnits.SunriseTime)
	assert.Equal(t, "iso8601", result.DailyUnits.SunsetTime)
	assert.Equal(t, "", result.DailyUnits.UVIndexMax)
	assert.Equal(t, "km/h", result.DailyUnits.WindSpeedMax)
	assert.Equal(t, "mm", result.DailyUnits.PrecipitationSum)

	require.Len(t, result.Daily, 1)

	day := result.Daily[0]
	assert.Equal(t, 23.3, day.TempMax)
	assert.Equal(t, 9.2, day.TempMin)
	assert.Equal(t, 5.85, day.UVIndexMax)
	assert.Equal(t, 8.7, day.WindSpeedMax)
	assert.Equal(t, 0.0, day.PrecipitationSum)
	assert.Equal(
		t,
		"2026-05-01",
		day.Date.Format("2006-01-02"),
	)
	assert.Equal(
		t,
		"2026-05-01T03:34",
		day.SunriseTime.Format("2006-01-02T15:04"),
	)
	assert.Equal(
		t,
		"2026-05-01T18:32",
		day.SunsetTime.Format("2006-01-02T15:04"),
	)
}

func TestService_GetWeatherDaily_ResponseFail(t *testing.T) {
	mock := &mockClient{
		err: errors.New("some-error"),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	_, err := s.GetWeatherDaily(t.Context(), 0, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "some-error")
}

func TestService_GetWeatherDaily_JsonUnmarshalFail(t *testing.T) {
	mock := &mockClient{
		resp: []byte(`invalid-json`),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	_, err := s.GetWeatherDaily(t.Context(), 0, 0)

	require.Error(t, err)
}

func getJsonWeatherWeeklyResponse() []byte {
	return []byte(`{
  "latitude": 52.52,
  "longitude": 13.419998,
  "generationtime_ms": 0.246286392211914,
  "utc_offset_seconds": 0,
  "timezone": "GMT",
  "timezone_abbreviation": "GMT",
  "elevation": 38,
  "daily_units": {
    "time": "iso8601",
    "temperature_2m_max": "°C",
    "temperature_2m_min": "°C",
    "sunrise": "iso8601",
    "sunset": "iso8601",
    "uv_index_max": "",
    "wind_speed_10m_max": "km/h",
    "precipitation_sum": "mm"
  },
  "daily": {
    "time": [
      "2026-05-01",
      "2026-05-02",
      "2026-05-03",
      "2026-05-04",
      "2026-05-05",
      "2026-05-06",
      "2026-05-07"
    ],
    "temperature_2m_max": [23.3, 26.2, 28.1, 20.7, 16.7, 13.2, 16],
    "temperature_2m_min": [9.2, 11.7, 12.5, 13.7, 7.2, 4.2, 2.7],
    "sunrise": [
      "2026-05-01T03:34",
      "2026-05-02T03:32",
      "2026-05-03T03:30",
      "2026-05-04T03:28",
      "2026-05-05T03:27",
      "2026-05-06T03:25",
      "2026-05-07T03:23"
    ],
    "sunset": [
      "2026-05-01T18:32",
      "2026-05-02T18:33",
      "2026-05-03T18:35",
      "2026-05-04T18:37",
      "2026-05-05T18:39",
      "2026-05-06T18:40",
      "2026-05-07T18:42"
    ],
    "uv_index_max": [5.85, 5.9, 5.1, 4.3, 0.55, 1.45, 1.35],
    "wind_speed_10m_max": [8.7, 15.4, 14.7, 13.2, 11.8, 13.7, 8.6],
    "precipitation_sum": [0, 0, 0, 0.6, 0, 0, 0]
  }
}`)
}

func TestService_GetWeatherWeekly_Success(t *testing.T) {
	mock := &mockClient{
		resp: getJsonWeatherWeeklyResponse(),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	result, err := s.GetWeatherWeekly(t.Context(), 0, 0)

	require.NoError(t, err)

	assert.Equal(t, "iso8601", result.DailyUnits.Date)
	assert.Equal(t, "°C", result.DailyUnits.TempMax)
	assert.Equal(t, "°C", result.DailyUnits.TempMin)
	assert.Equal(t, "iso8601", result.DailyUnits.SunriseTime)
	assert.Equal(t, "iso8601", result.DailyUnits.SunsetTime)
	assert.Equal(t, "", result.DailyUnits.UVIndexMax)
	assert.Equal(t, "km/h", result.DailyUnits.WindSpeedMax)
	assert.Equal(t, "mm", result.DailyUnits.PrecipitationSum)

	require.Len(t, result.Daily, 7)

	d0 := result.Daily[0]

	assert.Equal(t, 23.3, d0.TempMax)
	assert.Equal(t, 9.2, d0.TempMin)
	assert.Equal(t, 5.85, d0.UVIndexMax)
	assert.Equal(t, 8.7, d0.WindSpeedMax)
	assert.Equal(t, 0.0, d0.PrecipitationSum)
	assert.Equal(t, "2026-05-01", d0.Date.Format("2006-01-02"))
	assert.Equal(t, "2026-05-01T03:34", d0.SunriseTime.Format("2006-01-02T15:04"))
	assert.Equal(t, "2026-05-01T18:32", d0.SunsetTime.Format("2006-01-02T15:04"))

	dLast := result.Daily[6]

	assert.Equal(t, 16.0, dLast.TempMax)
	assert.Equal(t, 2.7, dLast.TempMin)
	assert.Equal(t, 1.35, dLast.UVIndexMax)
	assert.Equal(t, 8.6, dLast.WindSpeedMax)
	assert.Equal(t, 0.0, dLast.PrecipitationSum)
	assert.Equal(t, "2026-05-07", dLast.Date.Format("2006-01-02"))
	assert.Equal(t, "2026-05-07T03:23", dLast.SunriseTime.Format("2006-01-02T15:04"))
	assert.Equal(t, "2026-05-07T18:42", dLast.SunsetTime.Format("2006-01-02T15:04"))

	for _, d := range result.Daily {
		assert.False(t, d.Date.IsZero())
		assert.False(t, d.SunriseTime.IsZero())
		assert.False(t, d.SunsetTime.IsZero())
	}
}

func TestService_GetWeatherWeekly_ResponseFail(t *testing.T) {
	mock := &mockClient{
		err: errors.New("some-error"),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	_, err := s.GetWeatherWeekly(t.Context(), 0, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "some-error")
}

func TestService_GetWeatherWeekly_JsonUnmarshalFail(t *testing.T) {
	mock := &mockClient{
		resp: []byte(`invalid-json`),
	}

	s := service{
		apiURL: "some-api",
		client: mock,
	}

	_, err := s.GetWeatherWeekly(t.Context(), 0, 0)

	require.Error(t, err)
}
