package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID `json:"id"`
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductCreateRequest struct {
	CategoryID  string  `json:"category_id" validate:"required,uuid4"`
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type ProductUpdateRequest struct {
	CategoryID  string  `json:"category_id" validate:"required,uuid4"`
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}
