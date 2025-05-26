package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {

	type respond struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirps, err := apiCfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error creating chirp", err)
		return
	}

	response := []respond{}

	for _, chirp := range chirps {
		response = append(response, respond{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSon(w, 200, response)

}

func (apiCfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid UUID format", err)
	}

	type respond struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirp, err := apiCfg.db.GetOneChirp(r.Context(), parsedID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, 404, "Chirp not found", err)
			return
		}
		respondWithError(w, 500, "Error getting chirp", err)
		return
	}

	response := respond{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSon(w, 200, response)

}
