package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g0shi4ek/grcp_weather/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbConf.DbHost, cfg.DbConf.DbPort, cfg.DbConf.DbUser, cfg.DbConf.DbPassword, cfg.DbConf.DbName,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Cannot parse config, %v", err)
	}
	config.MaxConns = 10
	config.MaxConnLifetime = time.Hour * 2

	log.Printf("Connected to database, %s", cfg.DbConf.DbName)

	return pgxpool.NewWithConfig(ctx, config)
}
