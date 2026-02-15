package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryCreateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type CategoryUpdateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}
