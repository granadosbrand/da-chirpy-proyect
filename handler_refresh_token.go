package main

import (
	"net/http"
	"time"

	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type respond struct {
		Token string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error extracting token", err)
		return
	}

	// Get the token from database

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error retrieving refresh token: ", err)
		return
	}

	// Create a new JWT token

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
	}

	respondWithJSon(w, http.StatusOK, respond{
		Token: token,
	})

}
