package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
	"github.com/granadosbrand/da-chirpy-proyect/internal/database"
)

type User struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// 1. Extract email from body

	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body := params{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 500, "Error decoding json", err)
		return
	}

	if body.Email == "" || body.Password == "" {
		respondWithError(w, 400, "Required fields", nil)
		return
	}

	// 2. Create the user

	// 2.1 Hash the password

	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		respondWithError(w, 500, "Couldn't hash the password", err)
	}

	user, err := apiCfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          body.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 500, "Error creating user: ", err)
		return
	}

	respondWithJSon(w, http.StatusCreated, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})

}
