package jwt_helper

import (
	"errors"
	"net/http"
	"strings"
)

func ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return "", errors.New("invalid Authorization header format")
		}
		return parts[1], nil
	}

	cookie, err := r.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}
	
	token := r.URL.Query().Get("token")
	if token != "" {
		return token, nil
	}

	return "", errors.New("no token found")
}
