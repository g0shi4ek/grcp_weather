package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=11111111 dbname=weatherdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	err = db.Ping() // Проверяем подключение
	if err != nil {
		fmt.Println("Ошибка Ping:", err)
		return
	}

	rows, err := db.Query("SELECT version()") // Выполняем простой запрос
	if err != nil {
		fmt.Println("Ошибка выполнения запроса:", err)
		return
	}
	defer rows.Close()

	var version string
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			fmt.Println("Ошибка чтения данных:", err)
			return
		}
		fmt.Println("Версия PostgreSQL:", version)
	}

	fmt.Println("Подключение к базе данных успешно проверено!")

	//result, err := db.Exec("insert into Weather (City, Time, Temprature, Feels, Humidity, Wind) values ($1,$2,$3,$4,$5,$6)", city, time.Now(), weatherData.Main.Temp, weatherData.Main.FeelsLike, weatherData.Main.Humidity, weatherData.Wind.Speed)
	//if err != nil {
	//	panic(err)
	//}

	//rowsAffected, err := result.RowsAffected()
	//if err != nil {
	//	log.Fatal(err)

	//}
	//fmt.Println("Rows affected:", rowsAffected)
}
