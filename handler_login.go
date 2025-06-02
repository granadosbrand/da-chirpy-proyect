package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
	"github.com/granadosbrand/da-chirpy-proyect/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type loginParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type respond struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	// Generate token

	expiresIn := time.Hour

	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token: ", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}
	// Insert refresh token into table

	responseRefresh, err := cfg.db.InsertRefreshToken(r.Context(), database.InsertRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error insert refresh token", err)
	}

	respondWithJSon(w, http.StatusOK, respond{
		User: User{
			Id:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		},
		Token:        token,
		RefreshToken: responseRefresh.Token,
	})
}
