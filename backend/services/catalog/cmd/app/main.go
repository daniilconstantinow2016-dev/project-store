package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/daniil-dev/project-store/backend/services/catalog/internal/repository"
	"github.com/daniil-dev/project-store/backend/services/catalog/internal/transport/rest"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Загружаем переменные из .env файла
	// Если файла нет (например, в проде в Kubernetes), он просто пропустит ошибку
	if err := godotenv.Load(); err != nil {
		log.Println("Info: Файл .env не найден, используем системные переменные")
	}

	// 2. Подключаемся к БД
	conn := connectToDB()
	defer conn.Close(context.Background())

	// 3. Инициализация слоев
	repo := repository.NewProductRepository(conn)
	handler := rest.NewHandler(repo)

	// 4. Роутинг
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetAllProducts(w, r)
		case http.MethodPost:
			handler.CreateProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Читаем порт из конфига (или ставим :8080 по умолчанию)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = ":8080"
	}

	log.Printf("🚀 CyberMarket запущен на порту %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func connectToDB() *pgx.Conn {
	// 👇 ТЕПЕРЬ МЫ ЧИТАЕМ URL ИЗ ПЕРЕМЕННОЙ ОКРУЖЕНИЯ
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("ОШИБКА: Не задана переменная DATABASE_URL")
	}

	// Ждем немного (для Docker Compose в будущем)
	time.Sleep(3 * time.Second)
	
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v\n", err)
	}
	log.Println("Подключение к БД успешно! 📦")
	return conn
}
