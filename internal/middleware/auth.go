package middleware

import (
	"context"
	"net/http"
	"strings"

	"quadlingo/internal/models"
	"quadlingo/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error": "Требуется токен"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return utils.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Неверный токен"}`, http.StatusUnauthorized)
			return
		}

		userIDFloat, _ := claims["user_id"].(float64)
		roleStr, _ := claims["role"].(string)

		user := models.User{
			ID:   int(userIDFloat),
			Role: models.Role(roleStr),
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCurrentUser(r *http.Request) models.User {
	if user, ok := r.Context().Value(UserContextKey).(models.User); ok {
		return user
	}
	return models.User{}
}

func RequireRole(required ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetCurrentUser(r)
			allowed := false
			for _, r := range required {
				if user.Role == r || user.Role == models.RoleAdmin {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, `{"error": "Недостаточно прав"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
