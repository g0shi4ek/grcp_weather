package server

import (
	"context"
	"encoding/json"
	"first_try/weather"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net"
	"net/http"
)

const apiKey = "cae0bfaacc6a42deb153581a503b95f7" // API-ключ

type weatherServer struct {
	weather.UnimplementedWeatherServiceServer
}

func (s *weatherServer) GetWeather(ctx context.Context, req *weather.CityRequest) (*weather.WeatherData, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", req.City, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read: %v", err)
	}

	var weatherData weather.WeatherData

	err = json.Unmarshal(body, &weatherData) //записывает данные в weatherdata
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal: %v", err)
	}

	return &weatherData, nil
}

func RunServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to listen: %v", err)
	}
	s := grpc.NewServer()
	weather.RegisterWeatherServiceServer(s, &weatherServer{})
	fmt.Printf("gRPC server started on port %s\n", port)
	return s.Serve(lis)
}
