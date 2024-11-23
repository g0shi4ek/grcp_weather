package sv

import (
	"context"
	"database/sql"
	"encoding/json"
	"first_try/weather"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	_ "os"

	"net"
	"net/http"
	"time"
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

	connStr := "user=postgres password=11111111 dbname=weatherdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to open db: %v", err)
	}
	defer db.Close()

	result, err := db.Exec("insert into Weather (City, Time, Temperature, Feels, Humidity, Wind) values ($1,$2,$3,$4,$5,$6)", weatherData.Name, time.Now(), weatherData.Main.Temp, weatherData.Main.FeelsLike, weatherData.Main.Humidity, weatherData.Wind.Speed)
	if err != nil {
		panic(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add: %v", err)
	}
	fmt.Println("Rows affected:", rowsAffected)

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
