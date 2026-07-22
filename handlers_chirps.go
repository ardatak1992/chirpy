package main

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ardatak1992/chirpy/internal/auth"
	"github.com/ardatak1992/chirpy/internal/database"
	"github.com/google/uuid"
)

// Chirp represents a chirp response payload for the API.
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

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

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirpsRes, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, c := range chirpsRes {
		chirps = append(
			chirps,
			Chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
	}

	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirp", err)
		return
	}

	chirpRes, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	chirp := Chirp{
		ID:        chirpRes.ID,
		CreatedAt: chirpRes.CreatedAt,
		UpdatedAt: chirpRes.UpdatedAt,
		Body:      chirpRes.Body,
		UserID:    chirpRes.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)

}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	const chirpMaxLength = 140

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error posting chirp", err)
		return
	}

	userIDFromToken, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error posting chirp", err)
		return
	}

	log.Printf("User id from token: %s\n", userIDFromToken)

	params := parameters{}
	err = decodeRequestJSON(r, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding json", err)
		return
	}

	if len(params.Body) > chirpMaxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleanChirp := cleanChirpsFromProfanity(params.Body)

	chirp, err := cfg.dbQueries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleanChirp,
			UserID: params.UserID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
