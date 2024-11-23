package main

import (
	"bufio"
	"first_try/client/cl"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the city: ")
	city, _ := reader.ReadString('\n')
	city = city[:len(city)-1] // убирает символ перевода строки

	weatherClient, err := cl.NewWeatherClient(":5000")
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
