package store

import (
	"context"
	"fmt"
	"mini-product-catalog/internal/model"
	"strings"

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

	CategoryID *uuid.UUID
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
	`, id).Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)

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

func (s *ProductStore) Update(ctx context.Context, id uuid.UUID, categoryID uuid.UUID, name, description string, price float64) (model.Product, error) {
	var p model.Product
	err := s.db.QueryRow(ctx, `
		UPDATE products
		SET category_id = $2,
			name = $3,
			description = $4,
			price = $5,
			updated_at = now()
		WHERE id = $1
		RETURNING
			id, category_id,
			(SELECT name FROM categories WHERE id = $2) AS category_name,
			name, description, price::float8, created_at, updated_at
	`, id, categoryID, name, description, price).
		Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	return p, err
}

func (s *ProductStore) Delete(ctx context.Context, id uuid.UUID) (model.Product, error) {
	var p model.Product
	err := s.db.QueryRow(ctx, `
		DELETE FROM products
		WHERE id = $1
		RETURNING id, category_id, ''::text as category_name, name, description, price::float8, created_at, updated_at
	`, id).Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	return p, err
}

func (s *ProductStore) List(ctx context.Context, opt ProductListOptions) ([]model.Product, int, error) {
	if opt.Page < 1 {
		opt.Page = 1
	}
	if opt.Limit < 1 {
		opt.Limit = 10
	}
	if opt.Limit > 100 {
		opt.Limit = 100
	}
	offset := (opt.Page - 1) * opt.Limit

	sortCol := "p.created_at"
	if opt.Sort == "price" {
		sortCol = "p.price"
	}
	order := "DESC"
	if strings.ToLower(opt.Order) == "asc" {
		order = "ASC"
	}

	conds := []string{"1=1"}
	args := []any{}
	argN := 1

	if opt.CategoryID != nil {
		conds = append(conds, fmt.Sprintf("p.category_id = $%d", argN))
		args = append(args, *opt.CategoryID)
		argN++
	}
	if opt.MinPrice != nil {
		conds = append(conds, fmt.Sprintf("p.price >= $%d", argN))
		args = append(args, *opt.MinPrice)
		argN++
	}
	if opt.MaxPrice != nil {
		conds = append(conds, fmt.Sprintf("p.price <= $%d", argN))
		args = append(args, *opt.MaxPrice)
		argN++
	}
	if strings.TrimSpace(opt.Q) != "" {
		conds = append(conds, fmt.Sprintf("p.name ILIKE $%d", argN))
		args = append(args, "%"+strings.TrimSpace(opt.Q)+"%")
		argN++
	}

	whereSQL := strings.Join(conds, " AND ")

	var total int

	if err := s.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM products p
		WHERE `+whereSQL, args...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	argsList := append(args, opt.Limit, offset)
	limitPos := argN
	offsetPos := argN + 1

	rows, err := s.db.Query(ctx, `
		SELECT p.id, p.category_id, c.name, p.name, p.description, p.price::float8, p.created_at, p.updated_at
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE `+whereSQL+`
		ORDER BY `+sortCol+` `+order+`
		LIMIT $`+fmt.Sprint(limitPos)+` OFFSET $`+fmt.Sprint(offsetPos), argsList...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := []model.Product{}
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return out, total, nil
}
