package middleware

import (
	"context"
	"mini-product-catalog/internal/auth"
	"mini-product-catalog/internal/response"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type CurrentUser struct {
	ID   uuid.UUID
	Role string
}

type ctxKey int

const currentUserKey ctxKey = iota

func CurrentUserFromContext(ctx context.Context) (CurrentUser, bool) {
	v := ctx.Value(currentUserKey)
	u, ok := v.(CurrentUser)
	return u, ok
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h == "" {
				response.WriteError(w, http.StatusUnauthorized, "missing authorization header", nil)
				return
			}

			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.WriteError(w, http.StatusUnauthorized, "invalid authorization header", nil)
				return
			}

			tokenStr := parts[1]
			claims, err := auth.ParseAccessToken(tokenStr, jwtSecret)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "invalid token", nil)
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "invalid token subject", nil)
				return
			}

			cur := CurrentUser{ID: userID, Role: claims.Role}
			ctx := context.WithValue(r.Context(), currentUserKey, cur)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, ok := CurrentUserFromContext(r.Context())
			if !ok {
				response.WriteError(w, http.StatusUnauthorized, "unauthorized", nil)
				return
			}
			if u.Role != role {
				response.WriteError(w, http.StatusForbidden, "forbidden", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
