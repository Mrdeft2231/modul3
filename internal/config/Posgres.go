package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var (
	user     string = "user=postgres"
	password string = "password=2231"
	host     string = "host=localhost"
	port     string = "port=5432"
	dbname   string = "dbname=postgres"
	sslmode  string = "sslmode=disable"
)

func DB() *sql.DB {
	connect := fmt.Sprintf("%s %s %s %s %s %s", user, password, host, port, dbname, sslmode)
	fmt.Println(connect)
	db, err := sql.Open("postgres", connect)
	if err != nil {
		log.Fatalf("Не удалось подлкючиться к бд: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка при проверке подключения: ", err)
	}
	fmt.Println("Успешное подключение к базе данных!")
	return db
}
