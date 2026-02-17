package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"mini-product-catalog/internal/model"
	"mini-product-catalog/internal/response"
	"mini-product-catalog/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (h *ProductsHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid product id", nil)
		return
	}

	p, err := h.products.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "product not found", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch product", nil)
		return
	}

	response.WriteData(w, http.StatusOK, p, nil)
}

func (h *ProductsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.ProductCreateRequest
	if err := response.DecodeJSON(w, r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	catID, _ := uuid.Parse(req.CategoryID)

	ok, err := h.categories.Exists(r.Context(), catID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to validate category", nil)
		return
	}
	if !ok {
		response.WriteError(w, http.StatusBadRequest, "category_id not found", nil)
		return
	}

	created, err := h.products.Create(r.Context(), catID, req.Name, req.Description, req.Price)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to create product", nil)
		return
	}

	response.WriteData(w, http.StatusCreated, created, nil)
}

func (h *ProductsHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid product id", nil)
		return
	}

	var req model.ProductUpdateRequest
	if err := response.DecodeJSON(w, r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	catID, _ := uuid.Parse(req.CategoryID)

	ok, err := h.categories.Exists(r.Context(), catID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to validate category", nil)
		return
	}
	if !ok {
		response.WriteError(w, http.StatusBadRequest, "category_id not found", nil)
		return
	}

	updated, err := h.products.Update(r.Context(), id, catID, req.Name, req.Description, req.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "product not found", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to update product", nil)
		return
	}

	response.WriteData(w, http.StatusOK, updated, nil)
}

func (h *ProductsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid product id", nil)
		return
	}

	deleted, err := h.products.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "product not found", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to delete product", nil)
		return
	}

	response.WriteData(w, http.StatusOK, deleted, nil)
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
