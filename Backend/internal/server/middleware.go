package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["type"] != "access" {
			http.Error(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		uid := claims["uid"].(string)
		ctx := context.WithValue(r.Context(), "uid", uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
