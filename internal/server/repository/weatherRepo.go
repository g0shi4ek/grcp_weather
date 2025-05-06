package repository

import (
	"context"

	"github.com/g0shi4ek/grcp_weather/internal/server/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WeatherRepo struct {
	pgx *pgxpool.Pool
}

func NewWeatherRepo(pgx *pgxpool.Pool) *WeatherRepo {
	return &WeatherRepo{
		pgx: pgx,
	}
}

func (r *WeatherRepo) Create(ctx context.Context, weather *domain.Weather) error {
	query := "INSERT INTO weather (city, time, temperature, feels_like, humidity, wind_speed) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.pgx.Exec(ctx, query, weather.City, weather.Time, weather.Temperature, weather.FeelsLike, weather.Humidity, weather.WindSpeed)
	return err
}

func (r *WeatherRepo) GetByCity(ctx context.Context, city string) ([]*domain.Weather, error) {
	var weather []*domain.Weather
	query := "SELECT city, time, temperature, feels_like, humidity, wind_speed FROM weather WHERE city = $1"
	rows, err := r.pgx.Query(ctx, query, city)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var w domain.Weather
		err = rows.Scan(&w.City, &w.Time, &w.Temperature, &w.FeelsLike, &w.Humidity, &w.WindSpeed)
		if err != nil {
			return nil, err
		}

		weather = append(weather, &w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return weather, nil
}
