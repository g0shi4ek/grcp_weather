package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/g0shi4ek/grcp_weather/config"
	pb "github.com/g0shi4ek/grcp_weather/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcWeatherClient struct {
	conn   *grpc.ClientConn
	client pb.WeatherServiceClient
}

func NewGrpcWeatherClient(cfg *config.Config) (*GrpcWeatherClient, error) {
	conn, err := grpc.NewClient(
		"localhost:"+cfg.ServerConf.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start client: %v", err)
	}

	return &GrpcWeatherClient{
		conn:   conn,
		client: pb.NewWeatherServiceClient(conn),
	}, nil
}

func (c *GrpcWeatherClient) GetCurrentWeather(ctx context.Context, cities []string) ([]*pb.WeatherData, error) {
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	stream, err := c.client.GetCurrentWeather(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	go func() {
		defer func() {
			if err := stream.CloseSend(); err != nil {
				log.Printf("Failed to close send: %v", err)
			}
		}()

		for _, city := range cities {
			req := &pb.CityRequest{City: city}
			if err := stream.Send(req); err != nil {
				log.Printf("Failed to send city %s: %v", city, err)
				return
			}
			time.Sleep(time.Second)
		}
	}()

	var results []*pb.WeatherData
	for {
		weather, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("receive error: %w", err)
		}
		results = append(results, weather)
	}

	return results, nil
}

func (c *GrpcWeatherClient) GetHistoryWeather(ctx context.Context, city string) (*pb.WeatherHistoryData, error) {
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	req := &pb.CityRequest{
		City: city,
	}
	data, err := c.client.GetHistoryWeather(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %v", err)
	}
	return data, nil
}

func (c *GrpcWeatherClient) CloseConnection() error {
	return c.conn.Close()
}
