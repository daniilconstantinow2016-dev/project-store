package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	// 1. Подключаемся к Базе Данных
	conn := connectToDB()
	// Гарантируем, что соединение закроется, когда программа остановится
	defer conn.Close(context.Background())

	// 2. Выполняем миграцию (Создаем таблицу, если ее нет)
	createTable(conn)

	// 3. Запускаем сайт (HTTP сервер)
	startServer(conn)
}

// --- ФУНКЦИЯ 1: Подключение к БД ---
func connectToDB() *pgx.Conn {
	// Строка подключения: протокол://логин:пароль@адрес:порт/имя_базы
	databaseUrl := "postgres://cyber_user:cyber_password@localhost:5432/cyber_market_db"
	
	// Небольшая пауза, чтобы Docker успел проснуться
	time.Sleep(1 * time.Second)

	log.Println("Подключаюсь к PostgreSQL...")

	// context.Background() — это стандартный контекст выполнения
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		// log.Fatalf печатает ошибку и сразу выключает программу (Exit code 1)
		log.Fatalf("Ошибка подключения к БД: %v\n", err)
	}
	
	log.Println("Успешное подключение! 🚀")
	return conn
}

// --- ФУНКЦИЯ 2: Создание таблицы (SQL) ---
func createTable(conn *pgx.Conn) {
	// SQL запрос: Создать таблицу products с полями id, name, description, price, created_at
	// IF NOT EXISTS гарантирует, что мы не сломаем базу, если запустим код второй раз
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	// Exec выполняет запрос без возврата строк (для CREATE, INSERT, UPDATE)
	_, err := conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatalf("Не удалось создать таблицу: %v\n", err)
	}

	log.Println("Таблица 'products' проверена и готова к работе! 📦")
}

// --- ФУНКЦИЯ 3: Запуск сервера ---
func startServer(conn *pgx.Conn) {
	// Обработчик главной страницы
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// При каждом заходе проверяем, жива ли база (Ping)
		err := conn.Ping(context.Background())
		if err != nil {
			fmt.Fprintf(w, "Ошибка связи с БД 🔴: %v", err)
			return
		}
		fmt.Fprintf(w, "CyberMarket v0.3. Таблица Products существует 🟢")
	})

	port := ":8080"
	log.Printf("Запускаю HTTP сервер на порту %s...", port)
	
	// ListenAndServe блокирует выполнение и слушает порт вечно
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
