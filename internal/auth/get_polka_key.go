package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, "ApiKey ")
	token = strings.TrimSpace(token)

	if token == "" {
		return "", fmt.Errorf("empty token")
	}
	return token, nil
}
