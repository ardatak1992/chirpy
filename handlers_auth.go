package main

import (
	"net/http"

	"github.com/ardatak1992/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	params := parameters{ExpiresInSeconds: 60 * 60}
	err := decodeRequestJSON(r, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while logging in", err)
		return
	}

	

	userRes, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while logging in", err)
		return
	}

	passwordCorrect, err := auth.CheckPasswordHash(params.Password, userRes.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while logging in", err)
		return
	}

	if !passwordCorrect {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        userRes.ID,
		CreatedAt: userRes.CreatedAt,
		UpdatedAt: userRes.UpdatedAt,
		Email:     userRes.Email,
	})

}
