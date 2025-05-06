package domain

import (
	"time"
)

type Weather struct {
	City        string
	Time        time.Time
	Temperature float32
	FeelsLike   float32
	Humidity    int32
	WindSpeed   float32
}

// ответ weatherapi
type WeatherResponse struct{
	Name    string `json:"name"`
	Dt      int64  `json:"dt"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float32 `json:"temp"`
		FeelsLike float32 `json:"feels_like"`
		Humidity  int32   `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float32 `json:"speed"`
	} `json:"wind"`
}
