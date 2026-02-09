package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

// --- 1. НОВАЯ МОДЕЛЬ (МЕБЕЛЬ) ---
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`        // Название (Диван "Честер")
	Description string    `json:"description"` // Описание
	Price       float64   `json:"price"`       // Цена
	Material    string    `json:"material"`    // <--- НОВОЕ ПОЛЕ (Материал)
	ImageURL    string    `json:"image_url"`   // <--- НОВОЕ ПОЛЕ (Фото)
	CreatedAt   time.Time `json:"created_at"`
}

var dbConn *pgx.Conn

func main() {
	dbConn = connectToDB()
	defer dbConn.Close(context.Background())

	createTable(dbConn) // Создаст НОВУЮ таблицу с полями material и image_url

	startServer()
}

func connectToDB() *pgx.Conn {
	databaseUrl := "postgres://cyber_user:cyber_password@localhost:5432/cyber_market_db"
	time.Sleep(1 * time.Second) // Даем базе проснуться
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v\n", err)
	}
	log.Println("Успешное подключение! 🚀")
	return conn
}

func createTable(conn *pgx.Conn) {
	// SQL запрос изменился! Добавили material и image_url
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		material VARCHAR(50),
		image_url TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);`
	
	_, err := conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v\n", err)
	}
	log.Println("Таблица 'products' (Furniture Edition) готова! 🪑")
}

func startServer() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/products", createProductHandler) // И создать, и получить

	port := ":8080"
	log.Printf("Запускаю Мебельный Магазин на порту %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Designer Furniture Store API v1.0 🟢")
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	// --- МЕТОД GET (Получить список товаров) ---
	if r.Method == http.MethodGet {
		getProductsHandler(w, r)
		return
	}

	// --- МЕТОД POST (Создать товар) ---
	if r.Method == http.MethodPost {
		var p Product
		// Читаем JSON
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Кривой JSON", http.StatusBadRequest)
			return
		}

		// Вставляем в базу (теперь 5 полей вместо 3)
		sql := `
			INSERT INTO products (name, description, price, material, image_url)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at
		`
		
		err = dbConn.QueryRow(context.Background(), sql, 
			p.Name, p.Description, p.Price, p.Material, p.ImageURL).Scan(&p.ID, &p.CreatedAt)
		
		if err != nil {
			http.Error(w, "Ошибка БД: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
		
		log.Printf("Добавлен товар: %s (Материал: %s)", p.Name, p.Material)
		return
	}

	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}

// --- НОВАЯ ФУНКЦИЯ: Получить все товары ---
func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Запрашиваем все строки из таблицы
	rows, err := dbConn.Query(context.Background(), "SELECT id, name, description, price, material, image_url, created_at FROM products")
	if err != nil {
		http.Error(w, "Ошибка чтения БД", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Создаем пустой список (Слайс) товаров
	products := []Product{}

	// Бежим по каждой строке, которую вернула база
	for rows.Next() {
		var p Product
		// Сканируем данные из базы в структуру
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Material, &p.ImageURL, &p.CreatedAt)
		if err != nil {
			continue // Если ошибка в одной строке, пропускаем её
		}
		// Добавляем товар в список
		products = append(products, p)
	}

	// Отдаем список клиенту в виде JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
