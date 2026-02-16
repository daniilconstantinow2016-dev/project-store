package repository

import (
	"context"
	"database/sql"
	"fmt"

	// 👇 Убедись, что этот путь совпадает с твоим go.mod!
	"github.com/daniil-dev/project-store/backend/services/catalog/internal/domain"
)

type ProductRepository struct {
	db *sql.DB // 👈 Было *pgx.Conn, стало *sql.DB
}

// Конструктор теперь принимает стандартное подключение
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (name, description, price, material, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	// Используем QueryRowContext для поддержки контекста
	err := r.db.QueryRowContext(ctx, query,
		p.Name, p.Description, p.Price, p.Material, p.ImageURL).Scan(&p.ID, &p.CreatedAt)

	if err != nil {
		return fmt.Errorf("ошибка вставки: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	// Используем QueryContext
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, price, material, image_url, created_at FROM products")
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Material, &p.ImageURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	
	// Всегда проверяем ошибки после цикла rows.Next()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
