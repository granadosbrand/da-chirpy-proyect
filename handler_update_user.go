package main

import (
	"encoding/json"
	"net/http"

	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
	"github.com/granadosbrand/da-chirpy-proyect/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	// Get access token

	accesToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error extracting token from headers", err)
		return
	}

	// Validate token

	user, err := auth.ValidateJWT(accesToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token", err)
		return
	}

	// 1. Extract email from body

	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body := params{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 500, "Error decoding json", err)
		return
	}

	if body.Email == "" || body.Password == "" {
		respondWithError(w, 400, "Required fields", nil)
		return
	}

	// Hash the password

	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	// Update user email and password

	newUser, err := cfg.db.UpdateUserData(r.Context(), database.UpdateUserDataParams{
		Email:          body.Email,
		HashedPassword: hashedPassword,
		ID:             user,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user data", err)
	}

	respondWithJSon(w, http.StatusOK, User{
		Id:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed.Bool,
	})
}
