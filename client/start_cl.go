package main

import (
	"bufio"
	"database/sql"
	"first_try/client/cl"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
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

	connStr := "user=postgres password=11111111 dbname=weatherdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	result, err := db.Exec("insert into Weather (City, Time, Temperature, Feels, Humidity, Wind) values ($1,$2,$3,$4,$5,$6)", city, time.Now(), weatherData.Main.Temp, weatherData.Main.FeelsLike, weatherData.Main.Humidity, weatherData.Wind.Speed)
	if err != nil {
		panic(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Rows affected:", rowsAffected)
}
