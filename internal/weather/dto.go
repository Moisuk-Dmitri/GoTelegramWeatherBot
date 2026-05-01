package weather

type currentUnitsResponse struct {
	Time          string `json:"time"`
	Temperature   string `json:"temperature_2m"`
	Humidity      string `json:"relative_humidity_2m"`
	CloudCover    string `json:"cloud_cover"`
	Pressure      string `json:"pressure_msl"`
	WindSpeed     string `json:"wind_speed_10m"`
	WindDirection string `json:"wind_direction_10m"`
	Precipitation string `json:"precipitation"`
}

type currentResponse struct {
	Time          string  `json:"time"`
	Temperature   float64 `json:"temperature_2m"`
	Humidity      int     `json:"relative_humidity_2m"`
	CloudCover    int     `json:"cloud_cover"`
	Pressure      float64 `json:"pressure_msl"`
	WindSpeed     float64 `json:"wind_speed_10m"`
	WindDirection int     `json:"wind_direction_10m"`
	Precipitation float64 `json:"precipitation"`
}

type weatherStampResponse struct {
	CurrentUnits currentUnitsResponse `json:"current_units"`
	Current      currentResponse      `json:"current"`
}

type dailyUnitsResponse struct {
	Date             string `json:"time"`
	TempMax          string `json:"temperature_2m_max"`
	TempMin          string `json:"temperature_2m_min"`
	SunriseTime      string `json:"sunrise"`
	SunsetTime       string `json:"sunset"`
	UVIndexMax       string `json:"uv_index_max"`
	WindSpeedMax     string `json:"wind_speed_10m_max"`
	PrecipitationSum string `json:"precipitation_sum"`
}

type dailyResponse struct {
	Date             []string  `json:"time"`
	TempMin          []float64 `json:"temperature_2m_min"`
	TempMax          []float64 `json:"temperature_2m_max"`
	SunriseTime      []string  `json:"sunrise"`
	SunsetTime       []string  `json:"sunset"`
	UVIndexMax       []float64 `json:"uv_index_max"`
	WindSpeedMax     []float64 `json:"wind_speed_10m_max"`
	PrecipitationSum []float64 `json:"precipitation_sum"`
}

type weatherDailyResponse struct {
	DailyUnits dailyUnitsResponse `json:"daily_units"`
	Daily      dailyResponse      `json:"daily"`
}
