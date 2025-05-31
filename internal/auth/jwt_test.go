package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret"

	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("Error al crear token: %v", err)
	}

	casos := []struct {
		nombre     string
		token      string
		secret     string
		esperarErr bool
	}{
		{"Token válido", token, secret, false},
		{"Clave incorrecta", token, "otra-clave", true},
		{"Token mal formado", "no-es-un-token", secret, true},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := ValidateJWT(c.token, c.secret)
			if c.esperarErr && err == nil {
				t.Error("Se esperaba un error, pero no se obtuvo ninguno")
			} else if !c.esperarErr && err != nil {
				t.Errorf("No se esperaba error, pero se obtuvo: %v", err)
			}
		})
	}
}
