package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	pb "github.com/g0shi4ek/grcp_weather/gen/go"
	"github.com/g0shi4ek/grcp_weather/server/internal/domain"
	"net/http"
)

const apiKey = "cae0bfaacc6a42deb153581a503b95f7" // API-ключ

type WeatherService struct {
	repo domain.IWeatherRepository
	pb.UnimplementedWeatherServiceServer
}

func NewweatherService(repo domain.IWeatherRepository) *WeatherService{
	return &WeatherService{
		repo: repo,
	}
}

func (s* WeatherService) GetCurrent(ctx context.Context, city string) (*domain.Weather, error){
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %v", err)
	}

	var weatherData domain.Weather

	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}
	err = s.repo.Create(ctx, &weatherData)
	if err != nil{
		return nil,  fmt.Errorf("failed to create: %v", err)
	}
	
	return &weatherData, nil
}

func (s * WeatherService) GetHistory(ctx context.Context, city string) ([]*domain.Weather, error){
	weatherData, err := s.repo.GetByCity(ctx, city)
	if err != nil{
		return nil, fmt.Errorf("failed to get data: %v", err)
	}
	return weatherData, nil
}