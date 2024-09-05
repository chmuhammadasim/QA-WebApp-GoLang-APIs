package middleware

import (
	"context"
	"net/http"
	"qa-app/utils"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware is used to verify the JWT and extract user info
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.SendError(w, http.StatusUnauthorized, "Missing token")
			return
		}

		// Validate token and extract claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})
		if err != nil || !token.Valid {
			utils.SendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.SendError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Set user role in context
		role, ok := claims["role"].(string)
		if !ok {
			utils.SendError(w, http.StatusUnauthorized, "Role not found in token")
			return
		}
		ctx := context.WithValue(r.Context(), "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
