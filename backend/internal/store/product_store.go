package store

import (
	"context"
	"mini-product-catalog/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductStore struct {
	db *pgxpool.Pool
}

func NewProductStore(db *pgxpool.Pool) *ProductStore {
	return &ProductStore{db: db}
}

type ProductListOptions struct {
	Page  int
	Limit int

	CategoryId *uuid.UUID
	MinPrice   *float64
	MaxPrice   *float64
	Q          string

	Sort  string
	Order string
}

func (s *ProductStore) GetByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	var p model.Product
	err := s.db.QueryRow(ctx, `
		SELECT p.id, p.category_id, c.name, p.name, p.description, p.price::float8, p.created_at, p.updated_at
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.Name, &p.Description, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return model.Product{}, err
	}
	return p, nil
}

func (s *ProductStore) Create(ctx context.Context, categoryID uuid.UUID, name, description string, price float64) (model.Product, error) {
	var p model.Product
	err := s.db.QueryRow(ctx, `
		INSERT INTO products (category_id, name, description, price)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id, category_id,
			(SELECT name FROM categories WHERE id = $1) AS category_name,
			name, description, price::float8, created_at, updated_at
	`, categoryID, name, description, price).Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	return p, err
}
