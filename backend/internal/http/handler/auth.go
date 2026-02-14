package handler

import (
	"mini-product-catalog/internal/model"
	"mini-product-catalog/internal/response"
	"mini-product-catalog/internal/store"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users     *store.UserStore
	validate  *validator.Validate
	jwtSecret string
}

func NewAuthHandler(users *store.UserStore, validate *validator.Validate, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		users:     users,
		validate:  validate,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := response.DecodeJSON(w, r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to hash password", nil)
		return
	}

	created, err := h.users.Create(r.Context(), req.Name, req.Email, string(hash), "user")
	if err != nil {
		if store.IsUniqueViolation(err) {
			response.WriteError(w, http.StatusConflict, "email already registered", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to create user", nil)
		return
	}

	response.WriteData(w, http.StatusCreated, created, nil)
}
