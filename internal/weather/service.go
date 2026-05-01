package weather

import (
	"context"
	"encoding/json"
	"strconv"
)

type Service interface {
	GetWeatherStamp(context.Context, float64, float64) (WeatherStamp, error)
	GetWeatherDaily(context.Context, float64, float64) (WeatherDaily, error)
	GetWeatherWeekly(context.Context, float64, float64) (WeatherDaily, error)
}

type service struct {
	ApiURL string
	client clientI
}

func NewService(apiURL string) Service {
	return &service{
		ApiURL: apiURL,
		client: NewClient(nil),
	}
}

func (s *service) GetWeatherStamp(ctx context.Context, lat, lon float64) (WeatherStamp, error) {
	respJson, err := s.client.get(
		ctx,
		s.ApiURL,
		"?latitude="+strconv.FormatFloat(lat, 'f', -1, 64),
		"&longitude="+strconv.FormatFloat(lon, 'f', -1, 64),
		"&current=temperature_2m,relative_humidity_2m,cloud_cover,pressure_msl,wind_speed_10m,wind_direction_10m,precipitation",
	)
	if err != nil {
		return WeatherStamp{}, err
	}

	wsr := weatherStampResponse{}
	err = json.Unmarshal(respJson, &wsr)
	if err != nil {
		return WeatherStamp{}, err
	}

	return mapWeatherStamp(wsr)
}

func (s *service) GetWeatherDaily(ctx context.Context, lat, lon float64) (WeatherDaily, error) {
	respJson, err := s.client.get(
		ctx,
		s.ApiURL,
		"?latitude="+strconv.FormatFloat(lat, 'f', -1, 64),
		"&longitude="+strconv.FormatFloat(lon, 'f', -1, 64),
		"&daily=temperature_2m_max,temperature_2m_min,sunrise,sunset,uv_index_max,wind_speed_10m_max,precipitation_sum",
		"&past_days=0",
		"&forecast_days=1",
	)
	if err != nil {
		return WeatherDaily{}, err
	}

	wdr := weatherDailyResponse{}
	err = json.Unmarshal(respJson, &wdr)
	if err != nil {
		return WeatherDaily{}, err
	}

	return mapWeatherDaily(wdr)
}

func (s *service) GetWeatherWeekly(ctx context.Context, lat, lon float64) (WeatherDaily, error) {
	respJson, err := s.client.get(
		ctx,
		s.ApiURL,
		"?latitude="+strconv.FormatFloat(lat, 'f', -1, 64),
		"&longitude="+strconv.FormatFloat(lon, 'f', -1, 64),
		"&daily=temperature_2m_max,temperature_2m_min,sunrise,sunset,uv_index_max,wind_speed_10m_max,precipitation_sum",
		"&past_days=0",
		"&forecast_days=7",
	)
	if err != nil {
		return WeatherDaily{}, err
	}

	wdr := weatherDailyResponse{}
	err = json.Unmarshal(respJson, &wdr)
	if err != nil {
		return WeatherDaily{}, err
	}

	return mapWeatherDaily(wdr)
}
