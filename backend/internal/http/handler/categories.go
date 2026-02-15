package handler

import (
	"errors"
	"mini-product-catalog/internal/model"
	"mini-product-catalog/internal/response"
	"mini-product-catalog/internal/store"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CategoriesHandler struct {
	store    *store.CategoryStore
	validate *validator.Validate
}

func NewCategoriesHandler(store *store.CategoryStore, validate *validator.Validate) *CategoriesHandler {
	return &CategoriesHandler{store: store, validate: validate}
}

func (h *CategoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch categories", nil)
		return
	}

	meta := map[string]any{
		"count": len(items),
	}
	response.WriteData(w, http.StatusOK, items, meta)
}

func (h *CategoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CategoryCreateRequest

	if err := response.DecodeJSON(w, r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	created, err := h.store.Create(r.Context(), req.Name)
	if err != nil {
		if store.IsUniqueViolation(err) {
			response.WriteError(w, http.StatusConflict, "category already exists", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to create category", nil)
		return
	}

	response.WriteData(w, http.StatusCreated, created, nil)
}

func (h *CategoriesHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid category id", nil)
		return
	}

	var req model.CategoryUpdateRequest
	if err := response.DecodeJSON(w, r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	updated, err := h.store.Update(r.Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "category not found", nil)
			return
		}
		if store.IsUniqueViolation(err) {
			response.WriteError(w, http.StatusConflict, "category already exists", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to update category", nil)
		return
	}

	response.WriteData(w, http.StatusOK, updated, nil)
}

func (h *CategoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid category id", nil)
		return
	}

	deleted, err := h.store.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "category not found", nil)
			return
		}
		if store.IsForeignKeyViolation(err) {
			response.WriteError(w, http.StatusConflict, "category is used by products", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to delete category", nil)
		return
	}

	response.WriteData(w, http.StatusOK, deleted, nil)
}
