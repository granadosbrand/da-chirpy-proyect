package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
	"github.com/granadosbrand/da-chirpy-proyect/internal/database"
)

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpId := r.PathValue("chirpID")

	fmt.Println("chirpID: ", chirpId)
	if chirpId == "" {
		respondWithError(w, http.StatusUnauthorized, "Chirp id needed", fmt.Errorf("chirp id not found"))
		return
	}

	// Validate JWT

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error extracting token", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating JWT", err)
		return
	}

	// Get the chirp
	parsedChirpId, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing chirpy id", err)
		return
	}

	chirp, err := cfg.db.GetOneChirp(r.Context(), parsedChirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Error validating chirp author", err)
		return
	}

	// Delete the chirp
	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirp.ID,
		UserID: chirp.UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	respondWithJSon(w, http.StatusNoContent, nil)

}
