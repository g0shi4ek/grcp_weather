package handlers

import (
	"context"

	pb "github.com/g0shi4ek/grcp_weather/gen/go"
	"github.com/g0shi4ek/grcp_weather/server/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WeatherHandler struct {
	service domain.IWeatherService
}

func NewWeatherHandler (s domain.IWeatherService) *WeatherHandler{
	return &WeatherHandler{
		service: s,
	}
}

func (h * WeatherHandler) GetHistoryWeather(ctx context.Context, req * pb.CityRequest) (*pb.WeatherHistoryData, error){
	city := req.GetCity()
	data, err := h.service.GetHistory(ctx, city)
	if err != nil{
		return nil, status.Errorf(codes.Internal, "failed to get: %v", err)
	}
	var weatherHistory []*pb.WeatherData
	for _, weather := range data{
		current := &pb.WeatherData{
			City: weather.City,
			Time: timestamppb.New(weather.Time),
			Temperature: weather.Temperature,
			FeelsLike: weather.FeelsLike,
			Humidity: weather.Humidity,
			Speed: weather.WindSpeed,
		} 
		weatherHistory = append(weatherHistory,current)
	}
	return &pb.WeatherHistoryData{Data: weatherHistory} ,nil
}
