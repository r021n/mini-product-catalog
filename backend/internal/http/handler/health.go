package handler

import (
	"mini-product-catalog/internal/response"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.WriteData(w, http.StatusOK, map[string]string{
		"status": "ok",
	}, nil)
}
