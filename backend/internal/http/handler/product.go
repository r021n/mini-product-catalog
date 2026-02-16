package handler

import (
	"net/http"
	"strconv"
	"strings"

	"mini-product-catalog/internal/response"
	"mini-product-catalog/internal/store"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProductsHandler struct {
	products   *store.ProductStore
	categories *store.CategoryStore
	validate   *validator.Validate
}

func NewProductsHandler(products *store.ProductStore, categories *store.CategoryStore, validate *validator.Validate) *ProductsHandler {
	return &ProductsHandler{products: products, categories: categories, validate: validate}
}

func (h *ProductsHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page := parseInt(q.Get("page"), 1)
	limit := parseInt(q.Get("limit"), 10)

	var categoryID *uuid.UUID
	if v := strings.TrimSpace(q.Get("category_id")); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid category_id", nil)
			return
		}
		categoryID = &id
	}

	var minPrice *float64
	if v := strings.TrimSpace(q.Get("min_price")); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid min_price", nil)
			return
		}
		minPrice = &f
	}

	var maxPrice *float64
	if v := strings.TrimSpace(q.Get("max_price")); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid max_price", nil)
			return
		}
		maxPrice = &f
	}

	opt := store.ProductListOptions{
		Page:       page,
		Limit:      limit,
		CategoryID: categoryID,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Q:          q.Get("q"),
		Sort:       q.Get("sort"),
		Order:      q.Get("order"),
	}

	items, total, err := h.products.List(r.Context(), opt)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch products", nil)
		return
	}

	meta := map[string]any{
		"page":  opt.Page,
		"limit": opt.Limit,
		"total": total,
	}

	response.WriteData(w, http.StatusOK, items, meta)
}

func parseInt(s string, def int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
