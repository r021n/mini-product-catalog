package middleware

import (
	"context"
	"mini-product-catalog/internal/auth"
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
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			claims, err := auth.ParseAccessToken(tokenStr, jwtSecret)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				http.Error(w, "invalid token subject", http.StatusUnauthorized)
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
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if u.Role != role {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
