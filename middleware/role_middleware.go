package middleware

import (
	"net/http"
	"qa-app/utils"
	"strings"

	"github.com/gorilla/mux"
)

func RoleMiddleware(roles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract user role from request context
			// Assuming user role is stored in the context
			userRole, ok := r.Context().Value("role").(string)
			if !ok || userRole == "" {
				utils.SendError(w, http.StatusForbidden, "User role not found")
				return
			}

			// Check if the user role is in the allowed roles
			roleAllowed := false
			for _, role := range roles {
				if strings.EqualFold(userRole, role) {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				utils.SendError(w, http.StatusForbidden, "User role is not authorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
