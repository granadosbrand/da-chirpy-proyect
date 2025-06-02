package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/granadosbrand/da-chirpy-proyect/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {

	polka_key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	if polka_key != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Wrong API Key", fmt.Errorf("Wrong API Key"))
		return
	}

	type UserID struct {
		UserId uuid.UUID `json:"user_id"`
	}

	type Params struct {
		Event string `json:"event"`
		Data  UserID `json:"data"`
	}

	body := Params{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 500, "Error decoding json", err)
		return
	}

	if body.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "", fmt.Errorf("wrong event"))
		return
	}

	_, err = cfg.db.UpgradeUser(r.Context(), body.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error upgrading user", err)
		return
	}

	respondWithJSon(w, http.StatusNoContent, nil)

}
