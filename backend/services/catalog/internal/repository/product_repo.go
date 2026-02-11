package repository

import (
	"context"
	"fmt"

	"github.com/daniil-dev/project-store/backend/services/catalog/internal/domain"
	
	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	sql := `
		INSERT INTO products (name, description, price, material, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, sql, 
		p.Name, p.Description, p.Price, p.Material, p.ImageURL).Scan(&p.ID, &p.CreatedAt)
	
	if err != nil {
		return fmt.Errorf("ошибка вставки: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, description, price, material, image_url, created_at FROM products")
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
	return products, nil
}
