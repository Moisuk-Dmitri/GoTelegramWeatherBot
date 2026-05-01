package weather

import (
	"fmt"
	"time"
)

func mapWeatherStamp(resp weatherStampResponse) (WeatherStamp, error) {
	t, err := time.Parse("2006-01-02T15:04", resp.Current.Time)
	if err != nil {
		return WeatherStamp{}, fmt.Errorf("parse weather time: %w", err)
	}

	return WeatherStamp{
		CurrentUnits: currentUnits{
			Time:          resp.CurrentUnits.Time,
			Temperature:   resp.CurrentUnits.Temperature,
			Humidity:      resp.CurrentUnits.Humidity,
			CloudCover:    resp.CurrentUnits.CloudCover,
			Pressure:      resp.CurrentUnits.Pressure,
			WindSpeed:     resp.CurrentUnits.WindSpeed,
			WindDirection: resp.CurrentUnits.WindDirection,
			Precipitation: resp.CurrentUnits.Precipitation,
		},
		Current: current{
			Time:          t,
			Temperature:   resp.Current.Temperature,
			Humidity:      resp.Current.Humidity,
			CloudCover:    resp.Current.CloudCover,
			Pressure:      resp.Current.Pressure,
			WindSpeed:     resp.Current.WindSpeed,
			WindDirection: resp.Current.WindDirection,
			Precipitation: resp.Current.Precipitation,
		},
	}, nil
}

func mapWeatherDaily(wdr weatherDailyResponse) (WeatherDaily, error) {
	wd := WeatherDaily{
		DailyUnits: dailyUnits{
			Date:             wdr.DailyUnits.Date,
			TempMax:          wdr.DailyUnits.TempMax,
			TempMin:          wdr.DailyUnits.TempMin,
			SunriseTime:      wdr.DailyUnits.SunriseTime,
			SunsetTime:       wdr.DailyUnits.SunsetTime,
			UVIndexMax:       wdr.DailyUnits.UVIndexMax,
			WindSpeedMax:     wdr.DailyUnits.WindSpeedMax,
			PrecipitationSum: wdr.DailyUnits.PrecipitationSum,
		},
	}
	for i := 0; i < len(wdr.Daily.Date); i++ {
		d, err := time.Parse("2006-01-02", wdr.Daily.Date[i])
		if err != nil {
			return WeatherDaily{}, fmt.Errorf("failed to parse date: %w", err)
		}

		sRise, err := time.Parse("2006-01-02T15:04", wdr.Daily.SunriseTime[i])
		if err != nil {
			return WeatherDaily{}, fmt.Errorf("failed to parse sunrise: %w", err)
		}

		sSet, err := time.Parse("2006-01-02T15:04", wdr.Daily.SunsetTime[i])
		if err != nil {
			return WeatherDaily{}, fmt.Errorf("failed to parse sunset: %w", err)
		}

		wd.Daily = append(wd.Daily, daily{
			Date:             d,
			TempMin:          wdr.Daily.TempMin[i],
			TempMax:          wdr.Daily.TempMax[i],
			SunriseTime:      sRise,
			SunsetTime:       sSet,
			UVIndexMax:       wdr.Daily.UVIndexMax[i],
			WindSpeedMax:     wdr.Daily.WindSpeedMax[i],
			PrecipitationSum: wdr.Daily.PrecipitationSum[i],
		})
	}

	return wd, nil
}
