package main

import (
	"github/g0shi4ek/grpc_weather/initClient"
	"bufio"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the city: ")
	city, _ := reader.ReadString('\n')
	city = city[:len(city)-1]

	weatherClient, err := client.NewWeatherClient(":50050")
	if err != nil {
		fmt.Println(err)
		return
	}

	weatherData, err := weatherClient.GetWeather(city)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The weather in %s: Temperature is %.1f°C, feels like %.1f°C, humidity is %d%%, wind speed is %.1f m/s\n",
		city, weatherData.Main.Temp, weatherData.Main.FeelsLike, weatherData.Main.Humidity, weatherData.Wind.Speed)

}
