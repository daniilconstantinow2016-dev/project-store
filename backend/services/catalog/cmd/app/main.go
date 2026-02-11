package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/daniil-dev/project-store/backend/services/catalog/internal/repository"
	"github.com/daniil-dev/project-store/backend/services/catalog/internal/transport/rest"
	
	"github.com/jackc/pgx/v5"
)

func main() {
	conn := connectToDB()
	defer conn.Close(context.Background())

	// Инициализация слоев
	repo := repository.NewProductRepository(conn)
	handler := rest.NewHandler(repo)

	// Роутинг
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

	log.Println("🚀 CyberMarket (Clean Architecture) запущен на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func connectToDB() *pgx.Conn {
	databaseUrl := "postgres://cyber_user:cyber_password@localhost:5432/cyber_market_db"
	time.Sleep(2 * time.Second)
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v\n", err)
	}
	log.Println("Подключение к БД успешно! 📦")
	return conn
}
