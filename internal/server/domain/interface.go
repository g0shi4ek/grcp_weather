package domain

import "context"

type IWeatherRepository interface {
	Create(ctx context.Context, weather * Weather) error
	GetByCity(ctx context.Context, city string) ([]*Weather, error)
}

type IWeatherService interface{
	GetCurrent(ctx context.Context, city string, apiKey string) (*Weather, error)
	GetHistory(ctx context.Context, city string) ([]*Weather, error)
}