package handler

import (
	"errors"
	"mini-product-catalog/internal/auth"
	"mini-product-catalog/internal/middleware"
	"mini-product-catalog/internal/model"
	"mini-product-catalog/internal/response"
	"mini-product-catalog/internal/store"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := response.DecodeJSON(w, r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if err := h.validate.Struct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "validation error", err.Error())
		return
	}

	u, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusUnauthorized, "invalid credentials", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to login", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		response.WriteError(w, http.StatusUnauthorized, "invalid credentials", nil)
		return
	}

	ttl := 1 * time.Hour
	token, exp, err := auth.GenerateAccessToken(u.ID, u.Role, h.jwtSecret, ttl)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to generate token", nil)
		return
	}

	resp := map[string]any{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_at":   exp,
	}

	response.WriteData(w, http.StatusOK, resp, nil)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cur, ok := middleware.CurrentUserFromContext(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	u, err := h.users.GetByID(r.Context(), cur.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, http.StatusUnauthorized, "user not found", nil)
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch profile", nil)
		return
	}

	response.WriteData(w, http.StatusOK, u, nil)
}
