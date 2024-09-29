package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func GetPostgres() (*pgx.Conn, error) {
	dsn, err := buildDSN()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Невозможно подключиться к базе данных: %v\n", err)
	}

	return conn, err
}

func buildDSN() (string, error) {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()

	// Получение данных для подключения из переменных окружения
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	return dsn, err
}
