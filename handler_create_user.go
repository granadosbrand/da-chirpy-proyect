package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type params struct {
	Email string `json."email"`
}

type respond struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// 1. Extract email from body

	body := params{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 500, "Error decoding json", err)
		return
	}

	// 2. Create the user

	user, err := apiCfg.db.CreateUser(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, 500, "Error creating user: ", err)
	}

	respondWithJSon(w, http.StatusCreated, respond{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
