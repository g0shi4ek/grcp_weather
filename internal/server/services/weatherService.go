package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"net/http"

	"github.com/g0shi4ek/grcp_weather/internal/server/domain"
)

type WeatherService struct {
	repo domain.IWeatherRepository
}

func NewweatherService(repo domain.IWeatherRepository) *WeatherService {
	return &WeatherService{
		repo: repo,
	}
}

func ParseWeatherResponse(data []byte) (*domain.Weather, error) {
	var resp domain.WeatherResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if len(resp.Weather) == 0 {
		return nil, fmt.Errorf("no weather data available")
	}
	log.Print(resp.Wind)

	return &domain.Weather{
		City:        resp.Name,
		Time:        time.Unix(resp.Dt, 0),
		Temperature: resp.Main.Temp,
		FeelsLike:   resp.Main.FeelsLike,
		Humidity:    resp.Main.Humidity,
		WindSpeed:   resp.Wind.Speed,
	}, nil
}

func (s *WeatherService) GetCurrent(ctx context.Context, city string, apiKey string) (*domain.Weather, error) {
	log.Print(city)
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	weatherData, err := ParseWeatherResponse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}
	log.Print(weatherData.Temperature)

	err = s.repo.Create(ctx, weatherData)
	if err != nil {
		return nil, fmt.Errorf("failed to create: %v", err)
	}

	return weatherData, nil
}

func (s *WeatherService) GetHistory(ctx context.Context, city string) ([]*domain.Weather, error) {
	weatherData, err := s.repo.GetByCity(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("failed to get data: %v", err)
	}
	return weatherData, nil
}
