package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

var (
	user     string = "postgres"
	password string = "2231"
	host     string = "localhost"
	port     string = "5431"
	dbname   string = "postgres"
	sslmode  string = "disable"
)

func DB(ctx context.Context) *pgxpool.Pool {
	// Создаем контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Теперь отложим отмену контекста до конца работы с БД

	// Формируем строку подключения
	connect := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	fmt.Println("Подключаемся к БД с использованием строки:", connect)

	// Создаем пул подключений
	db, err := pgxpool.New(ctx, connect)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	// Создаем мигратор
	m, err := migrate.New(
		"file://../db/migrations", connect)
	if err != nil {
		log.Fatalf("не удалось создать миграцию %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("не удалось мигрировать данные %w", err)
	}
	// Возвращаем пул подключений (не закрываем его, т.к. он будет управляться в другом месте)
	return db
}
