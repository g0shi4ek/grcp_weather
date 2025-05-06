package handlers

import (
	"context"
	"io"
	"log"

	"github.com/g0shi4ek/grcp_weather/config"
	pb "github.com/g0shi4ek/grcp_weather/gen/go"
	"github.com/g0shi4ek/grcp_weather/internal/server/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WeatherHandler struct {
	service domain.IWeatherService
	cfg * config.Config
	pb.UnimplementedWeatherServiceServer
}

func NewWeatherHandler (s domain.IWeatherService, cfg * config.Config) *WeatherHandler{
	return &WeatherHandler{
		service: s,
		cfg: cfg,
	}
}

func (h * WeatherHandler) GetCurrentWeather(stream pb.WeatherService_GetCurrentWeatherServer) error{
	for {
		req, err := stream.Recv()
		if err == io.EOF{
			return nil
		}
		if err != nil{
			return err
		}
		log.Printf("city: %s", req.GetCity())

		api := h.cfg.ServerConf.ApiKey
		weather, err := h.service.GetCurrent(stream.Context(), req.GetCity(), api)
		if err != nil{
			log.Printf("Error getting weather for %s: %v", req.GetCity(), err)
            continue
		}
		
		log.Printf("get data: %s", weather.City)

		resp := &pb.WeatherData{
			City: weather.City,
			Time: timestamppb.New(weather.Time),
			Temperature: weather.Temperature,
			FeelsLike: weather.FeelsLike,
			Humidity: weather.Humidity,
			Speed: weather.WindSpeed,
		} 
		if err := stream.Send(resp); err != nil {
			return err
		}
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
