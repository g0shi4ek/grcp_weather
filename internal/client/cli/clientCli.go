package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"


	"github.com/g0shi4ek/grcp_weather/config"
	"github.com/g0shi4ek/grcp_weather/internal/client/app"
)

type WeatherCli struct {
	client *app.GrpcWeatherClient
	cfg    *config.Config
}

func NewWeatherCli(client *app.GrpcWeatherClient, cfg *config.Config) *WeatherCli {
	return &WeatherCli{
		client: client,
		cfg:    cfg,
	}
}

func (c *WeatherCli) RunApp(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		select{
		case <-ctx.Done():
            fmt.Println("\nApplication shutdown requested...")
            c.client.CloseConnection()
            return

		default:
			fmt.Println("\n1. Get current weather")
			fmt.Println("2. Get weather history")
			fmt.Println("3. Exit")
			fmt.Print("Choose option: ")

			scanner.Scan()
			input := scanner.Text()

			choice, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("please enter a number")
				continue
			}

			switch choice {
			case 1:
				c.handleCurrentWeather(ctx)
			case 2:
				c.handleHistoryWeather(ctx)
			case 3:
				c.client.CloseConnection()
				os.Exit(0)
			default:
				fmt.Println("Invalid choice")
			}
		}
	}
}

func (c *WeatherCli) handleCurrentWeather(ctx context.Context) {
	fmt.Print("Enter cities (space separated): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	cities := strings.Fields(input)
	if len(cities) == 0 {
		fmt.Println("No cities provided")
		return
	}

	log.Print(cities)

	data, err := c.client.GetCurrentWeather(ctx, cities)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	log.Print("Get data")

	for _, weather := range data {
		fmt.Printf("\n\nWeather in %s:\n", weather.GetCity())
		fmt.Printf("Time: %s \n", weather.GetTime().AsTime().Local().Format("2006-01-02 15:04"))
		fmt.Printf("Temperature: %.1f°C\n", weather.GetTemperature())
		fmt.Printf("Feels like: %.1f°C\n", weather.GetFeelsLike())
		fmt.Printf("Humidity: %d%%\n", weather.GetHumidity())
		fmt.Printf("Wind speed: %.1f m/s\n", weather.GetSpeed())
	}
}

func (c *WeatherCli) handleHistoryWeather(ctx context.Context) {
	fmt.Print("Enter city: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	city := scanner.Text()

	if city == "" {
		fmt.Println("City name cannot be empty")
		return
	}

	data, err := c.client.GetHistoryWeather(ctx, city)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("\n\nWeather history in %s", city)
	for _, weather := range data.GetData() {
		fmt.Printf("\nTime: %s \n", weather.GetTime().AsTime().Local().Format("2006-01-02 15:04"))
		fmt.Printf("Temperature: %.1f°C\n", weather.GetTemperature())
		fmt.Printf("Feels like: %.1f°C\n", weather.GetFeelsLike())
		fmt.Printf("Humidity: %d%%\n", weather.GetHumidity())
		fmt.Printf("Wind speed: %.1f m/s\n", weather.GetSpeed())
	}
}
