package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	// 👇 ТВОИ ИМПОРТЫ
	"github.com/daniil-dev/project-store/backend/services/catalog/internal/config" // Наш новый конфиг
	"github.com/daniil-dev/project-store/backend/services/catalog/internal/repository"
	"github.com/daniil-dev/project-store/backend/services/catalog/internal/transport/rest"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	// 1. ⚙️ ЗАГРУЗКА КОНФИГУРАЦИИ
	cfg := config.MustLoad()

	fmt.Printf("Запуск конфига: %s\n", cfg.Env)

	// Формируем строку подключения, используя данные из конфига
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// 2. Подключение к БД
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
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

	// 3. 🦆 Миграции
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
	fmt.Println("Миграции успешно применены! 🦆")

	// 4. Инициализация слоев
	repo := repository.NewProductRepository(db)
	handler := rest.NewHandler(repo)

	// 5. Роутинг
	mux := http.NewServeMux()
	mux.HandleFunc("/products", handler.CreateProduct)
	mux.HandleFunc("/products/list", handler.GetAllProducts)

	// 6. Запуск сервера (с таймаутами из конфига)
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	fmt.Println("Сервер запущен на порту", cfg.HTTPServer.Address)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Ошибка сервера:", err)
	}
}
