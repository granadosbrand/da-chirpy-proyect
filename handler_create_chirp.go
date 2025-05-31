package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
	"github.com/granadosbrand/da-chirpy-proyect/internal/database"
)

func (apiCfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Body   string    `json:"body"`
		// UserID uuid.UUID `json:"user_id"`
	}

	type respond struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	// Validate JWT

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error extracting token", err)
		return
	}

	userID, err := auth.ValidateJWT( bearerToken, apiCfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating JWT", err)
		return
	}
	

	body := params{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Error decoding params", err)
		return
	}

	chirpyMessage := body.Body

	chirpLength := len(chirpyMessage)

	if chirpLength > 140 {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}

	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	// Validate profane words

	cleanedMessage := badWordReplacement(chirpyMessage, profaneWords)

	// Insert chirp in database

	insertChirp := database.CreateChirpParams{
		Body:   cleanedMessage,
		UserID: userID,
	}

	chirp, err := apiCfg.db.CreateChirp(r.Context(), insertChirp)
	if err != nil {
		respondWithError(w, 500, "Error creating chirp", err)
		return
	}

	respondWithJSon(w, 201, respond{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
