package main

import (
	"encoding/json"
	"net/http"

	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type loginParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	body := loginParams{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 500, "Error decoding params", err)
	}

	// Get user

	user, err := cfg.db.GetUser(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	// Compare password
	err = auth.CheckPasswordHash(user.HashedPassword, body.Password)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	respondWithJSon(w, 200, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
