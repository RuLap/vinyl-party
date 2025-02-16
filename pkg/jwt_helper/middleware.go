package jwt_helper

import (
	"context"
	"net/http"
	"strings"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		if isPublicRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := ""

		if cookie, err := r.Cookie("auth_token"); err == nil {
			tokenString = cookie.Value
		}

		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		claims, err := ParseJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isPublicRoute(path string) bool {
	publicRoutes := map[string]bool{
		"/api/login":    true,
		"/api/register": true,
	}
	return publicRoutes[path]
}
