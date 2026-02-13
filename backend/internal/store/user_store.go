package store

import (
	"context"
	"errors"
	"mini-product-catalog/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, name, email, passwordHash, role string) (model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, role, created_at
	`, name, email, passwordHash, role).
		Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	return u, err
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, pgx.ErrNoRows
	}
	return u, err
}

func (s *UserStore) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, pgx.ErrNoRows
	}

	return u, err
}
