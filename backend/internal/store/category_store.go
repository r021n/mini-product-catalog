package store

import (
	"context"
	"mini-product-catalog/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryStore struct {
	db *pgxpool.Pool
}

func NewCategoryStore(db *pgxpool.Pool) *CategoryStore {
	return &CategoryStore{db: db}
}

func (s *CategoryStore) List(ctx context.Context) ([]model.Category, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, name, created_at
		FROM categories
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.Category{}
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *CategoryStore) Create(ctx context.Context, name string) (model.Category, error) {
	var c model.Category
	err := s.db.QueryRow(ctx, `
		INSERT INTO categories (name)
		VALUES ($1)
		RETURNING id, name, created_at
	`, name).Scan(&c.ID, &c.Name, &c.CreatedAt)

	return c, err
}

func IsUniqueViolation(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	return ok && pgErr.Code == "23505"
}
