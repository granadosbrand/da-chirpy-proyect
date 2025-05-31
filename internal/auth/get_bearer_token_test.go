package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       "Bearer esteEsMiTokenDePrueba",
			expectError: false,
		},
		{
			name:        "Missing Bearer prefix",
			token:       "esteEsMiTokenDePrueba",
			expectError: true,
		},
		{
			name:        "Empty token",
			token:       "Bearer ",
			expectError: true,
		},
		{
			name:        "Empty header",
			token:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Add("Authorization", tt.token)

			token, err := GetBearerToken(headers)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && token == "" {
				t.Error("Expected token but got empty string")
			}
		})
	}
}
