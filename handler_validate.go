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
		Valid bool `json:"valid"`
	}


	body := params{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Error decoding params", err)
		return
	}

	chirpLength := len(body.Body)

	if chirpLength > 140 {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}

	respondWithJSon(w, 200,  validChirp{
		Valid: true,
	})


}
