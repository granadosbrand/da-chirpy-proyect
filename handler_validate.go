package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func handlerValidate(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Body string `json:"body"`
	}

	type validChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	body := params{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&body)
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

	respondWithJSon(w, 200, validChirp{
		CleanedBody: cleanedMessage,
	})

}
