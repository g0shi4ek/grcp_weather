package client

import (
	"context"
	"first_try/weather"
	"fmt"
	"google.golang.org/grpc"
)

type WeatherClient struct {
	client weather.WeatherServiceClient
}

func NewWeatherClient(address string) (*WeatherClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	return &WeatherClient{
		client: weather.NewWeatherServiceClient(conn),
	}, nil
}

func (c *WeatherClient) GetWeather(city string) (*weather.WeatherData, error) {
	req := &weather.CityRequest{
		City: city,
	}
	res, err := c.client.GetWeather(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %v", err)
	}

	return res, nil
}
