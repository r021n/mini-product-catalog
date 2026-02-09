package handlers

import (
	"mini-product-catalog/internal/response"
	"net/http"
)

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}
