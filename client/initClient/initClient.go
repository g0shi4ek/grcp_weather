package cl

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "github.com/g0shi4ek/grcp_weather/gen/go"
)

type WeatherClient struct {
	client pb.WeatherServiceClient
}

func NewWeatherClient(address string) (*WeatherClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithDefaultCallOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	defer conn.Close()
	return &WeatherClient{
		client: pb.NewWeatherServiceClient(conn),
	}, nil
}

func (c *WeatherClient) GetWeather(city string) (*pb.WeatherData, error) {
	req := &pb.CityRequest{
		City: city,
	}
	res, err := c.client.GetWeather(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %v", err)
	}

	return res, nil
}
