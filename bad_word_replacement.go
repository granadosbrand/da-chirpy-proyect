package main

import (
	"strings"
)

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

func badWordReplacement(message string, profaneWords map[string]struct{}) string {

	words := strings.Fields(message)

	for i, word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	cleanMessage := strings.Join(words, " ")

	return cleanMessage
}
