package middleware

import (
	"github.com/ZaharBorisenko/jwt-auth/storage/repositories"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func AdminOnly(userRepo *repositories.UserRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("userClaims").(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		email, exists := claims["email"].(string)
		if !exists {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		user, err := userRepo.GetUserByEmail(r.Context(), email)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		if user.Role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
