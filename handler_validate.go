package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func cleanChirpsFromProfanity(chirp string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(chirp, " ")
	for i, word := range words {
		if slices.Contains(profanities, strings.ToLower(word)) {
			words[i] = "****"
		}
	}

	cleanChirp := strings.Join(words, " ")
	return cleanChirp
}

func validateHandler(w http.ResponseWriter, r *http.Request) {

	const chirpMaxLength = 140

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if len(params.Body) > chirpMaxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanChirp := cleanChirpsFromProfanity(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanChirp,
	})

}
