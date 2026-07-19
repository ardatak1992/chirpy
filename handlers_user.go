package main

import (
	"net/http"
	"time"

	"github.com/ardatak1992/chirpy/internal/auth"
	"github.com/ardatak1992/chirpy/internal/database"
	"github.com/google/uuid"
)

// User represents a user response payload for the API.
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	err := decodeRequestJSON(r, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user", err)
		return
	}

	userRes, err := cfg.dbQueries.CreateUser(
		r.Context(),
		database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        userRes.ID,
		CreatedAt: userRes.CreatedAt,
		UpdatedAt: userRes.UpdatedAt,
		Email:     userRes.Email,
	})

}
