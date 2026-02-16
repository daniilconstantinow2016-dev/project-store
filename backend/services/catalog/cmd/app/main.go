package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/daniil-dev/project-store/backend/services/catalog/internal/repository"
	"github.com/daniil-dev/project-store/backend/services/catalog/internal/transport/rest"

	_ "github.com/lib/pq"         // Стандартный драйвер (он у нас был изначально)
	"github.com/pressly/goose/v3" // Мигратор
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	// 1. Конфиг
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		"cyber_user", "cyber_password", dbHost, "cyber_market_db")

	// 2. Подключение (ОДНО, стандартное)
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr) // Используем стандартный sql.Open
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			fmt.Println("Подключение к БД успешно! 📦")
			break
		}
		fmt.Printf("Ждем БД... (%d/10)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer db.Close()

	// 3. 🦆 МИГРАЦИИ (Исправленный блок)
	// Вместо WithFS используем SetBaseFS - это работает во всех версиях
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
	fmt.Println("Миграции успешно применены! 🦆")

	// 4. Инициализация
	// Теперь db имеет тип *sql.DB, и твой репозиторий его примет без ошибок
	repo := repository.NewProductRepository(db)
	handler := rest.NewHandler(repo)

	// 5. Роутинг
	http.HandleFunc("/products", handler.CreateProduct)
	http.HandleFunc("/products/list", handler.GetAllProducts)

	// 6. Запуск
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = ":8080"
	}

	fmt.Println("Сервер запущен на порту", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
