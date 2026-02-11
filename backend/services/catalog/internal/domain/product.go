package domain

import "time"

// Product — это наша бизнес-сущность.
// Она описывает, что такое "Товар" для всего нашего приложения.
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Material    string    `json:"material"`    // Материал (новое поле)
	ImageURL    string    `json:"image_url"`   // Ссылка на фото (новое поле)
	CreatedAt   time.Time `json:"created_at"`
}
