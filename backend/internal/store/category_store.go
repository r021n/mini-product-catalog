package store

import (
	"context"
	"mini-product-catalog/internal/model"

	"github.com/google/uuid"
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

func (s *CategoryStore) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var ok bool
	err := s.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`, id).Scan(&ok)
	return ok, err
}

func (s *CategoryStore) Update(ctx context.Context, id uuid.UUID, name string) (model.Category, error) {
	var c model.Category
	err := s.db.QueryRow(ctx, `
		UPDATE categories
		SET name = $2
		WHERE id = $1
		RETURNING id, name, created_at
	`, id, name).Scan(&c.ID, &c.Name, &c.CreatedAt)

	if err != nil {
		return model.Category{}, err
	}

	return c, nil
}

func (s *CategoryStore) Delete(ctx context.Context, id uuid.UUID) (model.Category, error) {
	var c model.Category
	err := s.db.QueryRow(ctx, `
		DELETE FROM categories
		WHERE id = $1
		RETURNING id, name, created_at
	`, id).Scan(&c.ID, &c.Name, &c.CreatedAt)

	if err != nil {
		return model.Category{}, err
	}

	return c, nil
}
